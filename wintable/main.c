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

// start at WM_USER + 20 just in case for whatever reason we ever get the various dialog manager messages (see also http://blogs.msdn.com/b/oldnewthing/archive/2003/10/21/55384.aspx)
enum {
	// wParam - one of the type constants
	// lParam - column name as a Unicode string
	tableAddColumn = WM_USER + 20,
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

#define HANDLER(what) static BOOL what ## Handler(struct table *t, UINT uMsg, WPARAM wParam, LPARAM lParam, LRESULT *lResult)
#include "util.h"
#include "hscroll.h"
#include "vscroll.h"
#include "selection.h"
#include "draw.h"
#include "api.h"

typedef BOOL (*handlerfunc)(struct table *, UINT, WPARAM, LPARAM, LRESULT *);

const handlerfunc handlerfuncs[] = {
	vscrollHandler,
	APIHandler,
	NULL,
};

// TODO create a system where each of the above modules provide their own window procedures
static LRESULT CALLBACK tableWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	struct table *t;
	handlerfunc *hf;
	LRESULT lResult;
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
	for (hf = handlerfuncs; *hf != NULL; hf++)
		if ((*hf)(t, uMsg, wParam, lParam, &lResult))
			return lResult;
	switch (uMsg) {
	case WM_PAINT:
		dc = BeginPaint(hwnd, &ps);
		if (dc == NULL)
			abort();
		drawItems(t, dc, ps.rcPaint);
		EndPaint(hwnd, &ps);
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
