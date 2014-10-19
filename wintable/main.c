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

#define tableWindowClass L"gouitable"

struct table {
	HWND hwnd;
	HFONT defaultFont;
	HFONT font;
	intptr_t selected;
};

static void drawItems(struct table *t, HDC dc)
{
	HFONT prevfont;
	TEXTMETRICW tm;
	LONG y;
	intptr_t i;
	RECT r;

	if (GetClientRect(t->hwnd, &r) == 0)
		abort();
	prevfont = (HFONT) SelectObject(dc, t->font);
	if (prevfont == NULL)
		abort();
	if (GetTextMetricsW(dc, &tm) == 0)
		abort();
	y = 0;
	for (i = 0; i < 100; i++) {
		// TODO check errors
		// TODO verify correct colors
		if (t->selected == i) {
			RECT rsel;

			SetTextColor(dc, GetSysColor(COLOR_HIGHLIGHTTEXT));
			rsel.left = r.left;
			rsel.top = y;
			rsel.right = r.right - r.left;
			rsel.bottom = y + tm.tmHeight;
			FillRect(dc, &rsel, (HBRUSH) (COLOR_HIGHLIGHT + 1));
		} else {
			RECT rsel;

			SetTextColor(dc, GetSysColor(COLOR_WINDOWTEXT));
			rsel.left = r.left;
			rsel.top = y;
			rsel.right = r.right - r.left;
			rsel.bottom = y + tm.tmHeight;
			FillRect(dc, &rsel, (HBRUSH) (COLOR_WINDOW + 1));
		}
		SetBkMode(dc, TRANSPARENT);
		TextOutW(dc, r.left, y, L"Item", 4);
		y += tm.tmHeight;
	}
	if (SelectObject(dc, prevfont) != (HGDIOBJ) (t->font))
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
t->selected = 5;//TODO
		SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) t);
	}
	switch (uMsg) {
	case WM_PAINT:
		dc = BeginPaint(hwnd, &ps);
		if (dc == NULL)
			abort();
		drawItems(t, dc);
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
