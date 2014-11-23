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
	BOOL lastmouse;
	intptr_t lastmouseRow;
	intptr_t lastmouseColumn;
	BOOL mouseDown;			// TRUE if over a checkbox; the next two decide which ones
	intptr_t mouseDownRow;
	intptr_t mouseDownColumn;
};

static LONG rowHeight(struct table *t)
{
	HFONT thisfont, prevfont;
	TEXTMETRICW tm;
	HDC dc;
	LONG ret;

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
	ret = tm.tmHeight;
	if (ret < t->imagelistHeight)
		ret = t->imagelistHeight;
	if (ret < t->checkboxHeight)
		ret = t->checkboxHeight;
	return ret;
}

static void redrawAll(struct table *t)
{
	if (InvalidateRect(t->hwnd, NULL, TRUE) == 0)
		abort();
	if (UpdateWindow(t->hwnd) == 0)
		abort();
}

static RECT realClientRect(struct table *t)
{
	RECT r;

	if (GetClientRect(t->hwnd, &r) == 0)
		abort();
	r.top += t->headerHeight;
	return r;
}

static void repositionHeader(struct table *t)
{
	RECT r;
	HDLAYOUT headerlayout;
	WINDOWPOS headerpos;

	if (GetClientRect(t->hwnd, &r) == 0)		// use the whole client rect
		abort();
	// grow the rectangle to the left to fake scrolling
	r.left -= t->hpos;
	headerlayout.prc = &r;
	headerlayout.pwpos = &headerpos;
	if (SendMessageW(t->header, HDM_LAYOUT, 0, (LPARAM) (&headerlayout)) == FALSE)
		abort();
	if (SetWindowPos(t->header, headerpos.hwndInsertAfter, headerpos.x, headerpos.y, headerpos.cx, headerpos.cy, headerpos.flags | SWP_SHOWWINDOW) == 0)
		abort();
	t->headerHeight = headerpos.cy;
}

// this counts partially visible rows
// for all fully visible rows use t->pagesize
// cliprect and rowHeight must be specified here to avoid recomputing things multiple times
static intptr_t lastVisible(struct table *t, RECT cliprect, LONG rowHeight)
{
	intptr_t last;

	last = ((cliprect.bottom + rowHeight - 1) / rowHeight) + t->firstVisible;
	if (last >= t->count)
		last = t->count;
	return last;
}

static void redrawRow(struct table *t, intptr_t row)
{
	RECT r;
	intptr_t height;

	r = realClientRect(t);
	height = rowHeight(t);
	if (row < t->firstVisible || row > lastVisible(t, r, height))		// not visible; don't bother
		return;
	r.top = (row - t->firstVisible) * height + t->headerHeight;
	r.bottom = r.top + height;
	// keep the width and height the same; it spans the client area anyway
	if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
		abort();
	if (UpdateWindow(t->hwnd) == 0)
		abort();
}

static intptr_t hitTestColumn(struct table *t, int x)
{
	HDITEMW item;
	intptr_t i;

	// TODO count dividers
	for (i = 0; i < t->nColumns; i++) {
		ZeroMemory(&item, sizeof (HDITEMW));
		item.mask = HDI_WIDTH;
		if (SendMessageW(t->header, HDM_GETITEM, (WPARAM) i, (LPARAM) (&item)) == FALSE)
			abort();
		if (x < item.cxy)
			return i;
		x -= item.cxy;		// not yet
	}
	// no column
	return -1;
}

static void lParamToRowColumn(struct table *t, LPARAM lParam, intptr_t *row, intptr_t *column)
{
	int x, y;
	LONG h;

	x = GET_X_LPARAM(lParam);
	y = GET_Y_LPARAM(lParam);
	h = rowHeight(t);
	y += t->firstVisible * h;
	y -= t->headerHeight;
	y /= h;
	if (row != NULL) {
		*row = y;
		if (*row >= t->count)
			*row = -1;
	}
	if (column != NULL)
		*column = hitTestColumn(t, x);
}

static RECT checkboxRect(struct table *t, intptr_t row, intptr_t column, LONG rowHeight)
{
	RECT r;
	HDITEMW item;
	intptr_t i;

	// TODO count dividers
	for (i = 0; i < column; i++) {
		ZeroMemory(&item, sizeof (HDITEMW));
		item.mask = HDI_WIDTH;
		if (SendMessageW(t->header, HDM_GETITEM, (WPARAM) i, (LPARAM) (&item)) == FALSE)
			abort();
		r.left += item.cxy;
	}
	// TODO double-check to see if this takes any parameters
	r.left += SendMessageW(t->header, HDM_GETBITMAPMARGIN, 0, 0);
	r.right = r.left + t->checkboxWidth;
	// TODO vertical center
	r.top = row * rowHeight;
	r.bottom = r.top + t->checkboxHeight;
	return r;
}

// TODO clean up variables
static BOOL lParamInCheckbox(struct table *t, LPARAM lParam, intptr_t *row, intptr_t *column)
{
	int x, y;
	LONG h;
	intptr_t col;
	RECT r;
	POINT pt;

	x = GET_X_LPARAM(lParam);
	y = GET_Y_LPARAM(lParam);
	h = rowHeight(t);
	y += t->firstVisible * h;
	y -= t->headerHeight;
	pt.y = y;		// save actual y coordinate now
	y /= h;		// turn it into a row count
	if (y >= t->count)
		return FALSE;
	col = hitTestColumn(t, x);
	if (col == -1)
		return FALSE;
	if (t->columnTypes[col] != tableColumnCheckbox)
		return FALSE;
	r = checkboxRect(t, y, col, h);
	pt.x = x;
	if (PtInRect(&r, pt) == 0)
		return FALSE;
	if (row != NULL)
		*row = y;
	if (column != NULL)
		*column = col;
	return TRUE;
}

static void retrack(struct table *t)
{
	TRACKMOUSEEVENT tm;

	ZeroMemory(&tm, sizeof (TRACKMOUSEEVENT));
	tm.cbSize = sizeof (TRACKMOUSEEVENT);
	tm.dwFlags = TME_LEAVE;		// TODO also TME_NONCLIENT?
	tm.hwndTrack = t->hwnd;
	if (_TrackMouseEvent(&tm) == 0)
		abort();
}

static void addColumn(struct table *t, WPARAM wParam, LPARAM lParam)
{
	HDITEMW item;

	if (((int) wParam) >= nTableColumnTypes)
		abort();

	t->nColumns++;
	t->columnTypes = (int *) realloc(t->columnTypes, t->nColumns * sizeof (int));
	if (t->columnTypes == NULL)
		abort();
	t->columnTypes[t->nColumns - 1] = (int) wParam;

	ZeroMemory(&item, sizeof (HDITEMW));
	item.mask = HDI_WIDTH | HDI_TEXT | HDI_FORMAT;
	item.cxy = 200;		// TODO
	item.pszText = (WCHAR *) lParam;
	item.fmt = HDF_LEFT | HDF_STRING;
	if (SendMessage(t->header, HDM_INSERTITEM, (WPARAM) (t->nColumns - 1), (LPARAM) (&item)) == (LRESULT) (-1))
		abort();
	// TODO resize(t)?
	redrawAll(t);
}

static void track(struct table *t, LPARAM lParam)
{
	intptr_t row, column;
	BOOL prev;
	intptr_t prevrow, prevcolumn;

	prev = t->lastmouse;
	prevrow = t->lastmouseRow;
	prevcolumn = t->lastmouseColumn;
	t->lastmouse = lParamInCheckbox(t, lParam, &(t->lastmouseRow), &(t->lastmouseColumn));
	if (prev)
		if (prevrow != row || prevcolumn != column)
			redrawRow(t, prevrow);
	redrawRow(t, t->lastmouseRow);
}

static void hscrollto(struct table *t, intptr_t newpos)
{
	SCROLLINFO si;
	RECT scrollArea;

	if (newpos < 0)
		newpos = 0;
	if (newpos > (t->width - t->hpagesize))
		newpos = (t->width - t->hpagesize);

	scrollArea = realClientRect(t);

	// negative because ScrollWindowEx() is "backwards"
	if (ScrollWindowEx(t->hwnd, -(newpos - t->hpos), 0,
		&scrollArea, &scrollArea, NULL, NULL,
		SW_ERASE | SW_INVALIDATE) == ERROR)
		abort();
	t->hpos = newpos;
	// TODO text in header controls doesn't redraw?

	// TODO put this in a separate function? same for vscroll?
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_POS | SIF_RANGE;
	si.nPage = t->hpagesize;
	si.nMin = 0;
	si.nMax = t->width - 1;		// nMax is inclusive
	si.nPos = t->hpos;
	SetScrollInfo(t->hwnd, SB_HORZ, &si, TRUE);

	// and finally reposition the header
	repositionHeader(t);
}

static void hscrollby(struct table *t, intptr_t n)
{
	hscrollto(t, t->hpos + n);
}

// unfortunately horizontal wheel scrolling was only added in Vista

static void hscroll(struct table *t, WPARAM wParam)
{
	SCROLLINFO si;
	intptr_t newpos;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_POS | SIF_TRACKPOS;
	if (GetScrollInfo(t->hwnd, SB_HORZ, &si) == 0)
		abort();

	newpos = t->hpos;
	switch (LOWORD(wParam)) {
	case SB_LEFT:
		newpos = 0;
		break;
	case SB_RIGHT:
		newpos = t->width - t->hpagesize;
		break;
	case SB_LINELEFT:
		newpos--;
		break;
	case SB_LINERIGHT:
		newpos++;
		break;
	case SB_PAGELEFT:
		newpos -= t->hpagesize;
		break;
	case SB_PAGERIGHT:
		newpos += t->hpagesize;
		break;
	case SB_THUMBPOSITION:
		newpos = (intptr_t) (si.nPos);
		break;
	case SB_THUMBTRACK:
		newpos = (intptr_t) (si.nTrackPos);
	}

	hscrollto(t, newpos);
}

static void recomputeHScroll(struct table *t)
{
	HDITEMW item;
	intptr_t i;
	int width = 0;
	RECT r;
	SCROLLINFO si;

	// TODO count dividers
	for (i = 0; i < t->nColumns; i++) {
		ZeroMemory(&item, sizeof (HDITEMW));
		item.mask = HDI_WIDTH;
		if (SendMessageW(t->header, HDM_GETITEM, (WPARAM) i, (LPARAM) (&item)) == FALSE)
			abort();
		width += item.cxy;
	}
	t->width = (intptr_t) width;

	if (GetClientRect(t->hwnd, &r) == 0)
		abort();
	t->hpagesize = r.right - r.left;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_RANGE;
	si.nPage = t->hpagesize;
	si.nMin = 0;
	si.nMax = t->width - 1;			// - 1 because endpoints inclusive
	SetScrollInfo(t->hwnd, SB_HORZ, &si, TRUE);
}

static void vscrollto(struct table *t, intptr_t newpos)
{
	SCROLLINFO si;
	RECT scrollArea;

	if (newpos < 0)
		newpos = 0;
	if (newpos > (t->count - t->pagesize))
		newpos = (t->count - t->pagesize);

	scrollArea = realClientRect(t);

	// negative because ScrollWindowEx() is "backwards"
	if (ScrollWindowEx(t->hwnd, 0, (-(newpos - t->firstVisible)) * rowHeight(t),
		&scrollArea, &scrollArea, NULL, NULL,
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

static void selectItem(struct table *t, WPARAM wParam, LPARAM lParam)
{
	intptr_t prev;

	prev = t->selected;
	lParamToRowColumn(t, lParam, &(t->selected), &(t->focusedColumn));
	// TODO only if inside a checkbox
	t->mouseDown = TRUE;
	t->mouseDownRow = t->selected;
	t->mouseDownColumn = t->focusedColumn;
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
			// TODO replace all this
			rsel.left = headeritem.left + xoff;
			rsel.top = y;
			rsel.right = rsel.left + t->checkboxWidth;
			rsel.bottom = rsel.top + t->checkboxHeight;
			{ COLORREF c;

			c = RGB(255, 0, 0);
			if (t->mouseDown) {
				if (i == t->mouseDownRow && j == t->mouseDownColumn)
					c = RGB(0, 0, 255);
			} else if (t->lastmouse) {
				if (i == t->lastmouseRow && j == t->lastmouseColumn)
					c = RGB(0, 255, 0);
			}
			if (SetDCBrushColor(dc, c) == CLR_INVALID)
				abort();
			}
			if (FillRect(dc, &rsel, GetStockObject(DC_BRUSH)) == 0)
				abort();
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
			retrack(t);
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
	case WM_LBUTTONUP:
		// TODO toggle checkbox
		if (t->mouseDown) {
			t->mouseDown = FALSE;
			redrawRow(t, t->mouseDownRow);
		}
		return 0;
	// TODO other mouse buttons?
	case WM_MOUSEMOVE:
		track(t, lParam);
		return 0;
	case WM_MOUSELEAVE:
		t->lastmouse = FALSE;
		retrack(t);
		// TODO redraw row mouse is currently over
		// TODO split into its own function
		if (t->mouseDown) {
			t->mouseDown = FALSE;
			redrawRow(t, t->mouseDownRow);
		}
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

int main(void)
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
	ShowWindow(mainwin, SW_SHOWDEFAULT);
	if (UpdateWindow(mainwin) == 0)
		abort();
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
	return 0;
}
