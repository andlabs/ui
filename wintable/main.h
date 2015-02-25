// 7 january 2015

// TODO remove
#include <stdio.h>

// TODO
// - should tablePanic be CALLBACK or some other equivalent macro? and definitely export initTable somehow, but which alias macro to use?
// - make panic messages grammatically correct ("Table error: adding...")
// - make access to column widths consistent; see whether HDITEMW.cxy == (ITEMRECT.right - ITEMRECT.left)
// - make sure all uses of t->headerHeight are ADDED to RECT.top
// - WM_THEMECHANGED, etc.
// - see if vertical centering is really what we want or if we just want to offset by a few pixels or so
// - going right from column 0 to column 2 with the right arrow key deselects
// - make sure all error messages involving InvalidateRect() are consistent with regards to "redrawing" and "queueing for redraw"
// - collect all resize-related tasks in a single function (so things like adding columns will refresh everything, not just horizontal scrolls; also would fix initial coordinates)
// - checkbox columns don't clip to the column width
// - send standard notification codes
// - seriously figure out how we're going to update everything about the table in one place
// - redraw on focus change
// - draw themed WS_EX_CLIENTEDGE (and WS_EX_WINDOWEDGE?)
// - in-place tooltips for text columns
// - checkboxes: use capture instead of manual tracking logic? http://blogs.msdn.com/b/oldnewthing/archive/2015/02/25/10595729.aspx

#include "includethis.h"

static void (*tablePanic)(const char *, DWORD) = NULL;
#define panic(...) (*tablePanic)(__VA_ARGS__, GetLastError())
#define abort $$$$		// prevent accidental use of abort()

static BOOL (*WINAPI tableTrackMouseEvent)(LPTRACKMOUSEEVENT);

// forward declaration
struct tableAcc;

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
	BOOL checkboxMouseDown;
	intptr_t checkboxMouseDownRow;
	intptr_t checkboxMouseDownColumn;
	struct tableAcc *firstAcc;
};

// forward declaration (TODO needed?)
static LRESULT notify(struct table *, UINT, intptr_t, intptr_t, uintptr_t);

// necessary forward declarations
static void update(struct table *, BOOL);
static void updateAll(struct table *);

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
#include "accessibility.h"
#include "update.h"

static const handlerfunc handlers[] = {
	eventHandlers,
	childrenHandlers,
	resizeHandler,
	drawHandlers,
	apiHandlers,
	hscrollHandler,
	vscrollHandler,
	accessibilityHandler,
	NULL,
};

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
			SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) t);
		}
		// even if we did the above, fall through
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	if (uMsg == WM_DESTROY) {
printf("destroy\n");
		// TODO free appropriate (after figuring this part out) components of t
		invalidateTableAccs(t);
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
	fprintf(stderr, "Table error: %s (last error %I32u)\n", msg, lastError);
	fprintf(stderr, "This is the default Table error handler function; programs that use Table should provide their own instead.\nThe program will now break into the debugger.\n");
	DebugBreak();
}

// TODO have hInstance passed in
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
