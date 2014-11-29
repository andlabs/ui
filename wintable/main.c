// 19 october 2014
#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
#define CINTERFACE
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
extern HIMAGELIST makeCheckboxImageList(HWND hwnddc, HTHEME *theme, int *, int *);
enum {
        checkboxStateChecked = 1 << 0,
        checkboxStateHot = 1 << 1,
        checkboxStatePushed = 1 << 2,
        checkboxnStates = 1 << 3,
};
#include <windowsx.h>
#include <vsstyle.h>
#include <vssym32.h>
#include <oleacc.h>

// #qo LIBS: user32 kernel32 gdi32 comctl32 uxtheme

// TODO
// - http://blogs.msdn.com/b/oldnewthing/archive/2003/09/09/54826.aspx (relies on the integrality parts? IDK)
// 	- might want to http://blogs.msdn.com/b/oldnewthing/archive/2003/09/17/54944.aspx instead
// - http://msdn.microsoft.com/en-us/library/windows/desktop/bb775574%28v=vs.85%29.aspx
// - hscroll
// 	- keyboard navigation
// 		- how will this affect hot-tracking?
// 	- automatic hscroll when scrolling columns
// - accessibility
// 	- must use MSAA as UI Automation is not included by default on Windows XP (and apparently requires SP3?)
// - try horizontally scrolling the initail window and watch the selection rect corrupt itself *sometimes*
// - preallocate t->columnTypes instead of keeping it at exactly the right size
// - checkbox events
// 	- space to toggle (TODO); + or = to set; - to clear (see http://msdn.microsoft.com/en-us/library/windows/desktop/bb775941%28v=vs.85%29.aspx)
// 	- TODO figure out which notification is needed
// - http://blogs.msdn.com/b/oldnewthing/archive/2006/01/03/508694.aspx
// - free all allocated resources on WM_DESTROY
// - rename lastmouse
// 	- or perhaps do a general cleanup of the checkbox and mouse event code...
// - figure out why initial draw pretends there is no header
// - find places where the top-left corner of the client rect is assumed to be (0, 0)
// - figure out how we can split this into multiple files to make this easier to manage

#define tableWindowClass L"gouitable"

enum {
	// wParam - one of the type constants
	// lParam - column name as a Unicode string
	tableAddColumn = WM_USER,
};

enum {
	tableColumnText,
	tableColumnImage,
	tableColumnCheckbox,
	nTableColumnTypes,
};

struct table {
	HWND hwnd;
	HFONT defaultFont;
	HFONT font;
	intptr_t selected;
	intptr_t count;
	intptr_t firstVisible;
	intptr_t pagesize;		// in rows
	int wheelCarry;
	HWND header;
	int headerHeight;
	intptr_t nColumns;
	HIMAGELIST imagelist;
	int imagelistHeight;
	intptr_t width;
	intptr_t hpagesize;
	intptr_t hpos;
	HIMAGELIST checkboxes;
	HTHEME theme;
	int *columnTypes;
	intptr_t focusedColumn;
	int checkboxWidth;
	int checkboxHeight;
};

#include "util.h"
#include "api.h"
#include "hscroll.h"
#include "vscroll.h"

static void finishSelect(struct table *t, intptr_t prev)
{
	if (t->selected < 0)
		t->selected = 0;
	if (t->selected >= t->count)
		t->selected = t->count - 1;

	// always redraw the old and new rows to avoid artifacts when scrolling, even if they are the same (since the focused column may have changed)
	redrawRow(t, prev);
	if (prev != t->selected)
		redrawRow(t, t->selected);

	// if we need to scroll, the scrolling will force a redraw, so we don't have to worry about doing so ourselves
	if (t->selected < t->firstVisible)
		vscrollto(t, t->selected);
	// note that this is not lastVisible(t) because the last visible row may only be partially visible and we want selections to make them fully visible
	else if (t->selected >= (t->firstVisible + t->pagesize))
		vscrollto(t, t->selected - t->pagesize + 1);
}

static void keySelect(struct table *t, WPARAM wParam, LPARAM lParam)
{
	intptr_t prev;

	// TODO figure out correct behavior with nothing selected
	if (t->count == 0)		// don't try to do anything if there's nothing to do
		return;
	prev = t->selected;
	switch (wParam) {
	case VK_UP:
		t->selected--;
		break;
	case VK_DOWN:
		t->selected++;
		break;
	case VK_PRIOR:
		t->selected -= t->pagesize;
		break;
	case VK_NEXT:
		t->selected += t->pagesize;
		break;
	case VK_HOME:
		t->selected = 0;
		break;
	case VK_END:
		t->selected = t->count - 1;
		break;
	case VK_LEFT:
		t->focusedColumn--;
		if (t->focusedColumn < 0)
			if (t->nColumns == 0)		// peg at -1
				t->focusedColumn = -1;
			else
				t->focusedColumn = 0;
		break;
	case VK_RIGHT:
		t->focusedColumn++;
		if (t->focusedColumn >= t->nColumns)
			if (t->nColumns == 0)		// peg at -1
				t->focusedColumn = -1;
			else
				t->focusedColumn = t->nColumns - 1;
		break;
	// TODO keyboard shortcuts for going to the first/last column?
	default:
		// don't touch anything
		return;
	}
	finishSelect(t, prev);
}

// TODO rename
static void selectItem(struct table *t, WPARAM wParam, LPARAM lParam)
{
	intptr_t prev;

	prev = t->selected;
	lParamToRowColumn(t, lParam, &(t->selected), &(t->focusedColumn));
	finishSelect(t, prev);
}

static void resize(struct table *t)
{
	RECT r;
	SCROLLINFO si;

	// do this first so our scrollbar calculations can be correct
	repositionHeader(t);

	// now adjust the scrollbars
	r = realClientRect(t);
	t->pagesize = (r.bottom - r.top) / rowHeight(t);
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_RANGE | SIF_PAGE;
	si.nMin = 0;
	si.nMax = t->count - 1;
	si.nPage = t->pagesize;
	SetScrollInfo(t->hwnd, SB_VERT, &si, TRUE);

	recomputeHScroll(t);
}

// TODO alter this so that only the visible columns are redrawn
// TODO this means rename controlSize to clientRect
static void drawItem(struct table *t, HDC dc, intptr_t i, LONG y, LONG height, RECT controlSize)
{
	RECT rsel;
	HBRUSH background;
	int textColor;
	WCHAR msg[100];
	RECT headeritem;
	intptr_t j;
	LRESULT xoff;
	IMAGELISTDRAWPARAMS ip;
	POINT pt;

	// TODO verify these two
	background = (HBRUSH) (COLOR_WINDOW + 1);
	textColor = COLOR_WINDOWTEXT;
	if (t->selected == i) {
		// these are the colors wine uses (http://source.winehq.org/source/dlls/comctl32/listview.c)
		// the two for unfocused are also suggested by http://stackoverflow.com/questions/10428710/windows-forms-inactive-highlight-color
		background = (HBRUSH) (COLOR_HIGHLIGHT + 1);
		textColor = COLOR_HIGHLIGHTTEXT;
		if (GetFocus() != t->hwnd) {
			background = (HBRUSH) (COLOR_BTNFACE + 1);
			textColor = COLOR_BTNTEXT;
		}
	}

	// first fill the selection rect
	// note that this already only draws the visible area
	rsel.left = controlSize.left;
	rsel.top = y;
	rsel.right = controlSize.right - controlSize.left;
	rsel.bottom = y + height;
	if (FillRect(dc, &rsel, background) == 0)
		abort();

	// TODO double-check to see if this takes any parameters
	xoff = SendMessageW(t->header, HDM_GETBITMAPMARGIN, 0, 0);
	// now adjust for horizontal scrolling
	xoff -= t->hpos;

	// now draw the cells
	if (SetTextColor(dc, GetSysColor(textColor)) == CLR_INVALID)
		abort();
	if (SetBkMode(dc, TRANSPARENT) == 0)
		abort();
	for (j = 0; j < t->nColumns; j++) {
		if (SendMessageW(t->header, HDM_GETITEMRECT, (WPARAM) j, (LPARAM) (&headeritem)) == 0)
			abort();
		switch (t->columnTypes[j]) {
		case tableColumnText:
			rsel.left = headeritem.left + xoff;
			rsel.top = y;
			rsel.right = headeritem.right;
			rsel.bottom = y + height;
			// TODO vertical center in case the height is less than the icon height?
			if (DrawTextExW(dc, msg, wsprintf(msg, L"Item %d", i), &rsel, DT_END_ELLIPSIS | DT_LEFT | DT_NOPREFIX | DT_SINGLELINE, NULL) == 0)
				abort();
			break;
		case tableColumnImage:
			// TODO vertically center if image is smaller than text height
			// TODO same for checkboxes
			ZeroMemory(&ip, sizeof (IMAGELISTDRAWPARAMS));
			ip.cbSize = sizeof (IMAGELISTDRAWPARAMS);
			ip.himl = t->checkboxes;//t->imagelist;
			ip.i = (i%8);//0;
			ip.hdcDst = dc;
			ip.x = headeritem.left + xoff;
			ip.y = y;
			ip.cx = 0;		// draw whole image
			ip.cy = 0;
			ip.xBitmap = 0;
			ip.yBitmap = 0;
			ip.rgbBk = CLR_NONE;
			ip.fStyle = ILD_NORMAL | ILD_SCALE;		// TODO alpha-blend; ILD_DPISCALE?
			// TODO ILS_ALPHA?
			if (ImageList_DrawIndirect(&ip) == 0)
				abort();
			break;
		case tableColumnCheckbox:
			// TODO
			break;
		}
		if (t->selected == i && t->focusedColumn == j) {
			rsel.left = headeritem.left;
			rsel.top = y;
			rsel.right = headeritem.right;
			rsel.bottom = y + height;
			if (DrawFocusRect(dc, &rsel) == 0)
				abort();
		}
	}
}

static void drawItems(struct table *t, HDC dc, RECT cliprect)
{
	HFONT thisfont, prevfont;
	LONG height;
	LONG y;
	intptr_t i;
	RECT controlSize;		// for filling the entire selected row
	intptr_t first, last;

	if (GetClientRect(t->hwnd, &controlSize) == 0)
		abort();

	height = rowHeight(t);

	thisfont = t->font;		// in case WM_SETFONT happens before we return
	prevfont = (HFONT) SelectObject(dc, thisfont);
	if (prevfont == NULL)
		abort();

	// ignore anything beneath the header
	if (cliprect.top < t->headerHeight)
		cliprect.top = t->headerHeight;
	// now let's pretend the header isn't there
	// we only need it in (or rather, before) the drawItem() calls below
	cliprect.top -= t->headerHeight;
	cliprect.bottom -= t->headerHeight;

	// see http://blogs.msdn.com/b/oldnewthing/archive/2003/07/29/54591.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2003/07/30/54600.aspx
	// we need to add t->firstVisible here because cliprect is relative to the visible area
	first = (cliprect.top / height) + t->firstVisible;
	if (first < 0)
		first = 0;
	last = lastVisible(t, cliprect, height);

	// now for the first y, discount firstVisible
	y = (first - t->firstVisible) * height;
	// and offset by the header height
	y += t->headerHeight;
	for (i = first; i < last; i++) {
		drawItem(t, dc, i, y, height, controlSize);
		y += height;
	}

	// reset everything
	if (SelectObject(dc, prevfont) != (HGDIOBJ) (thisfont))
		abort();
}

static LRESULT CALLBACK tableWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	struct table *t;
	HDC dc;
	PAINTSTRUCT ps;
	NMHDR *nmhdr = (NMHDR *) lParam;
	NMHEADERW *nm = (NMHEADERW *) lParam;

	t = (struct table *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
	if (t == NULL) {
		// we have to do things this way because creating the header control will fail mysteriously if we create it first thing
		// (which is fine; we can get the parent hInstance this way too)
		if (uMsg == WM_NCCREATE) {
			CREATESTRUCTW *cs = (CREATESTRUCTW *) lParam;

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
			t->header = CreateWindowExW(0,
				WC_HEADERW, L"",
				// TODO is HOTTRACK needed?
				WS_CHILD | HDS_FULLDRAG | HDS_HORZ | HDS_HOTTRACK,
				0, 0, 0, 0,
				t->hwnd, (HMENU) 100, cs->hInstance, NULL);
			if (t->header == NULL)
				abort();
{t->imagelist = ImageList_Create(GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CYSMICON), ILC_COLOR32, 1, 1);
if(t->imagelist==NULL)abort();
{
HICON icon;
int unused;
icon = LoadIconW(NULL, IDI_ERROR);
if(icon == NULL)abort();
if (ImageList_AddIcon(t->imagelist, icon) == -1)abort();
if (ImageList_GetIconSize(t->imagelist, &unused, &(t->imagelistHeight)) == 0)abort();
}
}
			t->checkboxes = makeCheckboxImageList(t->hwnd, &(t->theme), &(t->checkboxWidth), &(t->checkboxHeight));
			t->focusedColumn = -1;
//TODO			retrack(t);
			SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) t);
		}
		// even if we did the above, fall through
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
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
		// also set the header font
		SendMessageW(t->header, WM_SETFONT, wParam, lParam);
		if (LOWORD(lParam) != FALSE) {
			// the scrollbar page size will change so redraw that too
			// also recalculate the header height
			// TODO do that when this is FALSE too somehow
			resize(t);
			redrawAll(t);
		}
		return 0;
	case WM_GETFONT:
		return (LRESULT) t->font;
	case WM_VSCROLL:
		vscroll(t, wParam);
		return 0;
	case WM_MOUSEWHEEL:
		wheelscroll(t, wParam);
		return 0;
	case WM_HSCROLL:
		hscroll(t, wParam);
		return 0;
	case WM_SIZE:
		resize(t);
		return 0;
	case WM_LBUTTONDOWN:
		selectItem(t, wParam, lParam);
		return 0;
	case WM_SETFOCUS:
	case WM_KILLFOCUS:
		// all we need to do here is redraw the highlight
		// TODO ensure giving focus works right
		redrawRow(t, t->selected);
		return 0;
	case WM_KEYDOWN:
		keySelect(t, wParam, lParam);
		return 0;
	// TODO header double-click
	case WM_NOTIFY:
		if (nmhdr->hwndFrom == t->header)
			switch (nmhdr->code) {
			// I could use HDN_TRACK but wine doesn't emit that
			case HDN_ITEMCHANGING:
			case HDN_ITEMCHANGED:		// TODO needed?
				recomputeHScroll(t);
				redrawAll(t);
				return FALSE;
			}
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	// TODO others?
	case WM_WININICHANGE:
	case WM_SYSCOLORCHANGE:
	case WM_THEMECHANGED:
		if (ImageList_Destroy(t->checkboxes) == 0)
			abort();
		t->checkboxes = makeCheckboxImageList(t->hwnd, &(t->theme), &(t->checkboxWidth), &(t->checkboxHeight));
		resize(t);		// TODO needed?
		redrawAll(t);
		// now defer back to DefWindowProc() in case other things are needed
		// TODO needed?
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	case tableAddColumn:
		addColumn(t, wParam, lParam);
		return 0;
	case WM_GETOBJECT:		// accessibility
/*
		if (((DWORD) lParam) == OBJID_CLIENT) {
			TODO *server;
			LRESULT lResult;

			// TODO create the server object
			lResult = LresultFromObject(IID_IAccessible, wParam, server);
			if (/* TODO failure *|/)
				abort();
			// TODO release object
			return lResult;
		}
*/
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
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

int main(int argc, char *argv[])
{
	HWND mainwin;
	MSG msg;
	INITCOMMONCONTROLSEX icc;

	ZeroMemory(&icc, sizeof (INITCOMMONCONTROLSEX));
	icc.dwSize = sizeof (INITCOMMONCONTROLSEX);
	icc.dwICC = ICC_LISTVIEW_CLASSES;
	if (InitCommonControlsEx(&icc) == 0)
		abort();
	makeTableWindowClass();
	mainwin = CreateWindowExW(0,
		tableWindowClass, L"Main Window",
		WS_OVERLAPPEDWINDOW | WS_HSCROLL | WS_VSCROLL,
		CW_USEDEFAULT, CW_USEDEFAULT,
		400, 400,
		NULL, NULL, GetModuleHandle(NULL), NULL);
	if (mainwin == NULL)
		abort();
	SendMessageW(mainwin, tableAddColumn, tableColumnText, (LPARAM) L"Column");
	SendMessageW(mainwin, tableAddColumn, tableColumnImage, (LPARAM) L"Column 2");
	SendMessageW(mainwin, tableAddColumn, tableColumnCheckbox, (LPARAM) L"Column 3");
	if (argc > 1) {
		NONCLIENTMETRICSW ncm;
		HFONT font;

		ZeroMemory(&ncm, sizeof (NONCLIENTMETRICSW));
		ncm.cbSize = sizeof (NONCLIENTMETRICSW);
		if (SystemParametersInfoW(SPI_GETNONCLIENTMETRICS, sizeof (NONCLIENTMETRICSW), &ncm, sizeof (NONCLIENTMETRICSW)) == 0)
			abort();
		font = CreateFontIndirectW(&ncm.lfMessageFont);
		if (font == NULL)
			abort();
		SendMessageW(mainwin, WM_SETFONT, (WPARAM) font, TRUE);
	}
	ShowWindow(mainwin, SW_SHOWDEFAULT);
	if (UpdateWindow(mainwin) == 0)
		abort();
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
	return 0;
}
