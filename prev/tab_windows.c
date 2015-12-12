// 25 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

// provided for cgo's benefit
LPWSTR xWC_TABCONTROL = WC_TABCONTROL;

static LRESULT CALLBACK tabSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	LRESULT lResult, r;
	RECT resizeRect;
	WINDOWPOS *wp;

	if (sharedWndProc(hwnd, uMsg, wParam, lParam, &lResult))
		return lResult;
	switch (uMsg) {
	case msgNOTIFY:
		switch (nmhdr->code) {
		case TCN_SELCHANGING:
			r = SendMessageW(hwnd, TCM_GETCURSEL, 0, 0);
			if (r == (LRESULT) -1)	// no tab currently selected
				return FALSE;
			tabChanging((void *) data, r);
			return FALSE;			// allow change
		case TCN_SELCHANGE:
			tabChanged((void *) data, SendMessageW(hwnd, TCM_GETCURSEL, 0, 0));
			return 0;
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case msgTabCurrentTabHasChildren:
		return (LRESULT) tabTabHasChildren((void *) data, SendMessageW(hwnd, TCM_GETCURSEL, 0, 0));
	// don't do this on WM_WINDOWPOSCHANGING; weird redraw issues will happen
	case WM_WINDOWPOSCHANGED:
		wp = (WINDOWPOS *) lParam;
		resizeRect.left = wp->x;
		resizeRect.top = wp->y;
		resizeRect.right = wp->x + wp->cx;
		resizeRect.bottom = wp->y + wp->cy;
		tabGetContentRect(hwnd, &resizeRect);
		tabResized((void *) data, resizeRect);
		// and chain up
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, tabSubProc, id) == FALSE)
			xpanic("error removing Tab subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Tab", "tabSubProc()", uMsg);
	return 0;		// unreached
}

void setTabSubclass(HWND hwnd, void *data)
{
	if ((*fv_SetWindowSubclass)(hwnd, tabSubProc, 0, (DWORD_PTR) data) == FALSE)
		xpanic("error subclassing Tab to give it its own event handler", GetLastError());
}

void tabAppend(HWND hwnd, LPWSTR name)
{
	TCITEM item;
	LRESULT n;

	ZeroMemory(&item, sizeof (TCITEM));
	item.mask = TCIF_TEXT;
	item.pszText = name;
	// MSDN's example code uses the first invalid index directly for this
	n = SendMessageW(hwnd, TCM_GETITEMCOUNT, 0, 0);
	if (SendMessageW(hwnd, TCM_INSERTITEM, (WPARAM) n, (LPARAM) (&item)) == (LRESULT) -1)
		xpanic("error adding tab to Tab", GetLastError());
}

void tabGetContentRect(HWND hwnd, RECT *r)
{
	// not &r; already a pointer (thanks MindChild in irc.efnet.net/#winprog for spotting my failure)
	SendMessageW(hwnd, TCM_ADJUSTRECT, FALSE, (LPARAM) r);
}

// theoretically we don't need to iterate over every tab for this, but let's do it just to be safe
LONG tabGetTabHeight(HWND hwnd)
{
	RECT r;
	LRESULT i, n;
	LONG tallest;

	n = SendMessageW(hwnd, TCM_GETITEMCOUNT, 0, 0);
	// if there are no tabs, then the control just draws a box over the full window rect, reserving no space for tabs; this is handled with the next line
	tallest = 0;
	for (i = 0; i < n; i++) {
		if (SendMessageW(hwnd, TCM_GETITEMRECT, (WPARAM) i, (LPARAM) (&r)) == FALSE)
			xpanic("error getting tab height for Tab.preferredSize()", GetLastError());
		if (tallest < (r.bottom - r.top))
			tallest = r.bottom - r.top;
	}
	return tallest;
}

void tabEnterChildren(HWND hwnd)
{
	DWORD style, xstyle;

	style = (DWORD) GetWindowLongPtrW(hwnd, GWL_STYLE);
	xstyle = (DWORD) GetWindowLongPtrW(hwnd, GWL_EXSTYLE);
	style &= ~((DWORD) WS_TABSTOP);
	xstyle |= WS_EX_CONTROLPARENT;
	SetWindowLongPtrW(hwnd, GWL_STYLE, (LONG_PTR) style);
	SetWindowLongPtrW(hwnd, GWL_EXSTYLE, (LONG_PTR) xstyle);
}

void tabLeaveChildren(HWND hwnd)
{
	DWORD style, xstyle;

	style = (DWORD) GetWindowLongPtrW(hwnd, GWL_STYLE);
	xstyle = (DWORD) GetWindowLongPtrW(hwnd, GWL_EXSTYLE);
	style |= WS_TABSTOP;
	xstyle &= ~((DWORD) WS_EX_CONTROLPARENT);
	SetWindowLongPtrW(hwnd, GWL_STYLE, (LONG_PTR) style);
	SetWindowLongPtrW(hwnd, GWL_EXSTYLE, (LONG_PTR) xstyle);
}
