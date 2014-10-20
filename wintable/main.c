// 19 october 2014
#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
// get Windows version right; right now Windows XP
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <commctrl.h>
#include <stdint.h>
#include <uxtheme.h>
#include <string.h>
#include <wchar.h>
#include <windowsx.h>
#include <vsstyle.h>
#include <vssym32.h>

// #qo LIBS: user32 kernel32 gdi32

// TODO
// - http://blogs.msdn.com/b/oldnewthing/archive/2003/09/09/54826.aspx (relies on the integrality parts? IDK)
// 	- might want to http://blogs.msdn.com/b/oldnewthing/archive/2003/09/17/54944.aspx instead

#define tableWindowClass L"gouitable"

struct table {
	HWND hwnd;
	HFONT defaultFont;
	HFONT font;
	intptr_t selected;
	intptr_t count;
	intptr_t firstVisible;
	intptr_t pagesize;
	int wheelCarry;
};

static LONG rowHeight(struct table *t)
{
	HFONT thisfont, prevfont;
	TEXTMETRICW tm;
	HDC dc;

	dc = GetDC(t->hwnd);
	if (dc == NULL)
		abort();
	thisfont = t->font;		// in case WM_SETFONT happens before we return
	prevfont = (HFONT) SelectObject(dc, thisfont);
	if (prevfont == NULL)
		abort();
	if (GetTextMetricsW(dc, &tm) == 0)
		abort();
	if (SelectObject(dc, prevfont) != (HGDIOBJ) (thisfont))
		abort();
	if (ReleaseDC(t->hwnd, dc) == 0)
		abort();
	return tm.tmHeight;
}

static void vscrollto(struct table *t, intptr_t newpos)
{
	SCROLLINFO si;

	if (newpos < 0)
		newpos = 0;
	if (newpos > (t->count - t->pagesize))
		newpos = (t->count - t->pagesize);

	// negative because ScrollWindowEx() is "backwards"
	if (ScrollWindowEx(t->hwnd, 0, (-(newpos - t->firstVisible)) * rowHeight(t),
		NULL, NULL, NULL, NULL,
		SW_ERASE | SW_INVALIDATE) == ERROR)
		abort();
	t->firstVisible = newpos;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_POS | SIF_RANGE;
	si.nPage = t->pagesize;
	si.nMin = 0;
	si.nMax = t->count - 1;		// nMax is inclusive
	si.nPos = t->firstVisible;
	SetScrollInfo(t->hwnd, SB_VERT, &si, TRUE);
}

static void vscrollby(struct table *t, intptr_t n)
{
	vscrollto(t, t->firstVisible + n);
}

static void wheelscroll(struct table *t, WPARAM wParam)
{
	int delta;
	int lines;
	UINT scrollAmount;

	delta = GET_WHEEL_DELTA_WPARAM(wParam);
	if (SystemParametersInfoW(SPI_GETWHEELSCROLLLINES, 0, &scrollAmount, 0) == 0)
		abort();
	if (scrollAmount == WHEEL_PAGESCROLL)
		scrollAmount = t->pagesize;
	if (scrollAmount == 0)		// no mouse wheel scrolling (or t->pagesize == 0)
		return;
	// the rest of this is basically http://blogs.msdn.com/b/oldnewthing/archive/2003/08/07/54615.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2003/08/11/54624.aspx
	// see those pages for information on subtleties
	delta += t->wheelCarry;
	lines = delta * ((int) scrollAmount) / WHEEL_DELTA;
	t->wheelCarry = delta - lines * WHEEL_DELTA / ((int) scrollAmount);
	vscrollby(t, -lines);
}

static void vscroll(struct table *t, WPARAM wParam)
{
	SCROLLINFO si;
	intptr_t newpos;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_POS | SIF_TRACKPOS;
	if (GetScrollInfo(t->hwnd, SB_VERT, &si) == 0)
		abort();

	newpos = t->firstVisible;
	switch (LOWORD(wParam)) {
	case SB_TOP:
		newpos = 0;
		break;
	case SB_BOTTOM:
		newpos = t->count - t->pagesize;
		break;
	case SB_LINEUP:
		newpos--;
		break;
	case SB_LINEDOWN:
		newpos++;
		break;
	case SB_PAGEUP:
		newpos -= t->pagesize;
		break;
	case SB_PAGEDOWN:
		newpos += t->pagesize;
		break;
	case SB_THUMBPOSITION:
		newpos = (intptr_t) (si.nPos);
		break;
	case SB_THUMBTRACK:
		newpos = (intptr_t) (si.nTrackPos);
	}

	vscrollto(t, newpos);
}

static void resize(struct table *t)
{
	RECT r;
	SCROLLINFO si;

	if (GetClientRect(t->hwnd, &r) == 0)
		abort();
	t->pagesize = (r.bottom - r.top) / rowHeight(t);
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_RANGE | SIF_PAGE;
	si.nMin = 0;
	si.nMax = t->count - 1;
	si.nPage = t->pagesize;
	SetScrollInfo(t->hwnd, SB_VERT, &si, TRUE);
}

static void drawItems(struct table *t, HDC dc, RECT cliprect)
{
	HFONT thisfont, prevfont;
	TEXTMETRICW tm;
	LONG y;
	intptr_t i;
	RECT r;
	intptr_t first, last;
	POINT prevOrigin;

	// TODO eliminate the need (only use cliprect)
	if (GetClientRect(t->hwnd, &r) == 0)
		abort();

	thisfont = t->font;		// in case WM_SETFONT happens before we return
	prevfont = (HFONT) SelectObject(dc, thisfont);
	if (prevfont == NULL)
		abort();
	if (GetTextMetricsW(dc, &tm) == 0)
		abort();

	// adjust the clip rect and the viewport so that (0, 0) is always the first item
	if (OffsetRect(&cliprect, 0, t->firstVisible * tm.tmHeight) == 0)
		abort();
	if (GetWindowOrgEx(dc, &prevOrigin) == 0)
		abort();
	if (SetWindowOrgEx(dc, prevOrigin.x, prevOrigin.y + (t->firstVisible * tm.tmHeight), NULL) == 0)
		abort();

	// see http://blogs.msdn.com/b/oldnewthing/archive/2003/07/29/54591.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2003/07/30/54600.aspx
	first = cliprect.top / tm.tmHeight;
	if (first < 0)
		first = 0;
	last = (cliprect.bottom + tm.tmHeight - 1) / tm.tmHeight;
	if (last >= t->count)
		last = t->count;

	y = first * tm.tmHeight;
	for (i = first; i < last; i++) {
		RECT rsel;
		HBRUSH background;
		WCHAR msg[100];

		// TODO check errors
		// TODO verify correct colors
		rsel.left = r.left;
		rsel.top = y;
		rsel.right = r.right - r.left;
		rsel.bottom = y + tm.tmHeight;
		background = (HBRUSH) (COLOR_WINDOW + 1);
		if (t->selected == i) {
			background = (HBRUSH) (COLOR_HIGHLIGHT + 1);
			SetTextColor(dc, GetSysColor(COLOR_HIGHLIGHTTEXT));
		} else
			SetTextColor(dc, GetSysColor(COLOR_WINDOWTEXT));
		FillRect(dc, &rsel, background);
		SetBkMode(dc, TRANSPARENT);
		TextOutW(dc, r.left, y, msg, wsprintf(msg, L"Item %d", i));
		y += tm.tmHeight;
	}

	// reset everything
	if (SetWindowOrgEx(dc, prevOrigin.x, prevOrigin.y, NULL) == 0)
		abort();
	if (SelectObject(dc, prevfont) != (HGDIOBJ) (thisfont))
		abort();
}

static LRESULT CALLBACK tableWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	struct table *t;
	HDC dc;
	PAINTSTRUCT ps;

	t = (struct table *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
	if (t == NULL) {
		t = (struct table *) malloc(sizeof (struct table));
		if (t == NULL)
			abort();
		ZeroMemory(t, sizeof (struct table));
		t->hwnd = hwnd;
		// TODO this should be a global
		t->defaultFont = (HFONT) GetStockObject(SYSTEM_FONT);
		if (t->defaultFont == NULL)
			abort();
		t->font = t->defaultFont;
t->selected = 5;t->count=100;//TODO
		SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) t);
	}
	switch (uMsg) {
	case WM_PAINT:
		dc = BeginPaint(hwnd, &ps);
		if (dc == NULL)
			abort();
		drawItems(t, dc, ps.rcPaint);
		EndPaint(hwnd, &ps);
		return 0;
	case WM_SETFONT:
		t->font = (HFONT) wParam;
		if (t->font == NULL)
			t->font = t->defaultFont;
		if (LOWORD(lParam) != FALSE)
			;	// TODO
		return 0;
	case WM_GETFONT:
		return (LRESULT) t->font;
	case WM_VSCROLL:
		vscroll(t, wParam);
		return 0;
	case WM_MOUSEWHEEL:
		wheelscroll(t, wParam);
		return 0;
	case WM_SIZE:
		resize(t);
		return 0;
	default:
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	abort();
	return 0;		// unreached
}

void makeTableWindowClass(void)
{
	WNDCLASSW wc;

	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpszClassName = tableWindowClass;
	wc.lpfnWndProc = tableWndProc;
	wc.hCursor = LoadCursorW(NULL, IDC_ARROW);
	wc.hIcon = LoadIconW(NULL, IDI_APPLICATION);
	wc.hbrBackground = (HBRUSH) (COLOR_WINDOW + 1);		// TODO correct?
	wc.style = CS_HREDRAW | CS_VREDRAW;
	wc.hInstance = GetModuleHandle(NULL);
	if (RegisterClassW(&wc) == 0)
		abort();
}

int main(void)
{
	HWND mainwin;
	MSG msg;

	makeTableWindowClass();
	mainwin = CreateWindowExW(0,
		tableWindowClass, L"Main Window",
		WS_OVERLAPPEDWINDOW | WS_HSCROLL | WS_VSCROLL,
		CW_USEDEFAULT, CW_USEDEFAULT,
		400, 400,
		NULL, NULL, GetModuleHandle(NULL), NULL);
	if (mainwin == NULL)
		abort();
	ShowWindow(mainwin, SW_SHOWDEFAULT);
	if (UpdateWindow(mainwin) == 0)
		abort();
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
	return 0;
}
