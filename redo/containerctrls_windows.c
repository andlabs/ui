/* 25 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

/* provided for cgo's benefit */
LPWSTR xWC_TABCONTROL = WC_TABCONTROL;

static LRESULT CALLBACK tabSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	LRESULT r;

	switch (uMsg) {
	case msgNOTIFY:
		switch (nmhdr->code) {
		case TCN_SELCHANGING:
			r = SendMessageW(hwnd, TCM_GETCURSEL, 0, 0);
			if (r == (LRESULT) -1)	/* no tab currently selected */
				return FALSE;
			tabChanging((void *) data, r);
			return FALSE;			/* allow change */
		case TCN_SELCHANGE:
			tabChanged((void *) data, SendMessageW(hwnd, TCM_GETCURSEL, 0, 0));
			return 0;
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, tabSubProc, id) == FALSE)
			xpanic("error removing Tab subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Tab", "tabSubProc()", uMsg);
	return 0;		/* unreached */
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
	/* MSDN's example code uses the first invalid index directly for this */
	n = SendMessageW(hwnd, TCM_GETITEMCOUNT, 0, 0);
	if (SendMessageW(hwnd, TCM_INSERTITEM, (WPARAM) n, (LPARAM) (&item)) == (LRESULT) -1)
		xpanic("error adding tab to Tab", GetLastError());
}

void tabGetContentRect(HWND hwnd, RECT *r)
{
	/* not &r; already a pointer (thanks MindChild in irc.efnet.net/#winprog for spotting my failure) */
	SendMessageW(hwnd, TCM_ADJUSTRECT, FALSE, (LPARAM) r);
}

/* TODO this assumes that all inactive tabs have the same height */
LONG tabGetTabHeight(HWND hwnd)
{
	RECT r;
	RECT r2;
	LRESULT n, current, other;

	n = SendMessageW(hwnd, TCM_GETITEMCOUNT, 0, 0);
	/* if there are no tabs, then the control just draws a box over the full window rect, reserving no space for tabs (TODO check on windows xp and 7) */
	if (n == 0)
		return 0;
	/* get the current tab's height */
	/* note that Windows calls the tabs themselves "items" */
	current = SendMessageW(hwnd, TCM_GETCURSEL, 0, 0);
	if (SendMessageW(hwnd, TCM_GETITEMRECT, (WPARAM) current, (LPARAM) (&r)) == FALSE)
		xpanic("error getting current tab's tab height for Tab.preferredSize()", GetLastError());
	/* if there's only one tab, then it's the current one; just get its size and return it */
	if (n == 1)
		goto onlyOne;
	/* otherwise, get an inactive tab's height and return the taller of the two heights */
	other = current + 1;
	if (other >= n)
		other = 0;
	if (SendMessageW(hwnd, TCM_GETITEMRECT, (WPARAM) other, (LPARAM) (&r2)) == FALSE)
		xpanic("error getting other tab's tab height for Tab.preferredSize()", GetLastError());
	if ((r2.bottom - r2.top) > (r.bottom - r.top))
		return r2.bottom - r2.top;
onlyOne:
	return r.bottom - r.top;
}
