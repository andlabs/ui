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
#include <windowsx.h>
#include <vsstyle.h>
#include <vssym32.h>
#include <oleacc.h>

// #qo LIBS: user32 kernel32 gdi32 comctl32 uxtheme

// TODO
// - should tablePanic be CALLBACK or some other equivalent macro? and definitely export initTable somehow, but which alias macro to use?
// - make panic messages grammatically correct ("Table error: adding...")
// - make access to column widths consistent; see whether HDITEMW.cxy == (ITEMRECT.right - ITEMRECT.left)
// - make sure all uses of t->headerHeight are ADDED to RECT.top
// - do we maintain the scrolling metaphor?
// - WM_THEMECHANGED, etc.
// - see if vertical centering is really what we want or if we just want to offset by a few pixels or so
// - going right from column 0 to column 2 with the right arrow key deselects

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

static void (*tablePanic)(const char *, DWORD) = NULL;
#define panic(...) (*tablePanic)(__VA_ARGS__, GetLastError())
#define abort $$$$		// prevent accidental use of abort()

static BOOL (*WINAPI tableTrackMouseEvent)(LPTRACKMOUSEEVENT);

struct table {
	HWND hwnd;
	HWND header;
	HFONT font;
	intptr_t nColumns;
	int *columnTypes;
	intptr_t width;
	intptr_t headerHeight;
	intptr_t hscrollpos;		// in logical units
	intptr_t hpagesize;		// in logical units
	intptr_t count;
	intptr_t vscrollpos;		// in rows
	intptr_t vpagesize;		// in rows
	int hwheelCarry;
	int vwheelCarry;
	intptr_t selectedRow;
	intptr_t selectedColumn;
	HTHEME theme;
	int checkboxWidth;
	int checkboxHeight;
	BOOL checkboxMouseOverLast;
	LPARAM checkboxMouseOverLastPoint;
};

#include "util.h"
#include "coord.h"
#include "scroll.h"
#include "hscroll.h"
#include "vscroll.h"
#include "select.h"
#include "checkboxes.h"
#include "events.h"
#include "header.h"
#include "children.h"
#include "resize.h"
#include "draw.h"
#include "api.h"

static const handlerfunc handlers[] = {
	eventHandlers,
	childrenHandlers,
	resizeHandler,
	drawHandlers,
	apiHandlers,
	hscrollHandler,
	vscrollHandler,
	NULL,
};

static void initDummyTableStuff(struct table *t)
{
	t->count = 100;
	t->selectedRow = 2;
	t->selectedColumn = 1;
}

static LRESULT CALLBACK tableWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	struct table *t;
	LRESULT lResult;

	t = (struct table *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
	if (t == NULL) {
		// we have to do things this way because creating the header control will fail mysteriously if we create it first thing
		// (which is fine; we can get the parent hInstance this way too)
		// we use WM_CREATE because we have to use WM_DESTROY to destroy the header; we can't do it in WM_NCDESTROY because Windows will have destroyed it for us by then, and let's match message pairs to be safe
		if (uMsg == WM_CREATE) {
			CREATESTRUCTW *cs = (CREATESTRUCTW *) lParam;

			t = (struct table *) tableAlloc(sizeof (struct table), "error allocating internal Table data structure");
			t->hwnd = hwnd;
			makeHeader(t, cs->hInstance);
			t->selectedRow = -1;
			t->selectedColumn = -1;
			loadCheckboxThemeData(t);
initDummyTableStuff(t);
			SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) t);
		}
		// even if we did the above, fall through
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	if (uMsg == WM_DESTROY) {
printf("destroy\n");
		// TODO free appropriate (after figuring this part out) components of t
		// TODO send EVENT_OBJECT_DESTROY events to accessibility listeners (when appropriate); see the note on proxy objects as well
		freeCheckboxThemeData(t);
		destroyHeader(t);
		tableFree(t, "error allocating internal Table data structure");
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	if (runHandlers(handlers, t, uMsg, wParam, lParam, &lResult))
		return lResult;
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

static void deftablePanic(const char *msg, DWORD lastError)
{
	fprintf(stderr, "Table error: %s (last error %d)\n", msg, lastError);
	fprintf(stderr, "This is the default Table error handler function; programs that use Table should provide their own instead.\nThe program will now break into the debugger.\n");
	DebugBreak();
}

void initTable(void (*panicfunc)(const char *msg, DWORD lastError), BOOL (*WINAPI tme)(LPTRACKMOUSEEVENT))
{
	WNDCLASSW wc;

	tablePanic = panicfunc;
	if (tablePanic == NULL)
		tablePanic = deftablePanic;
	if (tme == NULL)
		// TODO errorless version
		panic("must provide a TrackMouseEvent() to initTable()");
	tableTrackMouseEvent = tme;
	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpszClassName = tableWindowClass;
	wc.lpfnWndProc = tableWndProc;
	wc.hCursor = LoadCursorW(NULL, IDC_ARROW);
	wc.hIcon = LoadIconW(NULL, IDI_APPLICATION);
	wc.hbrBackground = (HBRUSH) (COLOR_WINDOW + 1);		// TODO correct?
	wc.style = CS_HREDRAW | CS_VREDRAW;
	wc.hInstance = GetModuleHandle(NULL);
	if (RegisterClassW(&wc) == 0)
		panic("error registering Table window class");
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
		panic("(test program) error initializing comctl32.dll");
	initTable(NULL, _TrackMouseEvent);
	mainwin = CreateWindowExW(0,
		tableWindowClass, L"Main Window",
		WS_OVERLAPPEDWINDOW | WS_HSCROLL | WS_VSCROLL,
		CW_USEDEFAULT, CW_USEDEFAULT,
		400, 400,
		NULL, NULL, GetModuleHandle(NULL), NULL);
	if (mainwin == NULL)
		panic("(test program) error creating Table");
	SendMessageW(mainwin, tableAddColumn, tableColumnText, (LPARAM) L"Column");
	SendMessageW(mainwin, tableAddColumn, tableColumnImage, (LPARAM) L"Column 2");
	SendMessageW(mainwin, tableAddColumn, tableColumnCheckbox, (LPARAM) L"Column 3");
	if (argc > 1) {
		NONCLIENTMETRICSW ncm;
		HFONT font;

		ZeroMemory(&ncm, sizeof (NONCLIENTMETRICSW));
		ncm.cbSize = sizeof (NONCLIENTMETRICSW);
		if (SystemParametersInfoW(SPI_GETNONCLIENTMETRICS, sizeof (NONCLIENTMETRICSW), &ncm, sizeof (NONCLIENTMETRICSW)) == 0)
			panic("(test program) error getting non-client metrics");
		font = CreateFontIndirectW(&ncm.lfMessageFont);
		if (font == NULL)
			panic("(test program) error creating lfMessageFont HFONT");
		SendMessageW(mainwin, WM_SETFONT, (WPARAM) font, TRUE);
	}
	ShowWindow(mainwin, SW_SHOWDEFAULT);
	if (UpdateWindow(mainwin) == 0)
		panic("(test program) error updating window");
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
	return 0;
}
