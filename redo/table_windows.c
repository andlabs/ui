/* 28 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

/* provided for cgo's benefit */
LPWSTR xWC_LISTVIEW = WC_LISTVIEW;

static LRESULT CALLBACK tableSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	NMLVDISPINFOW *fill = (NMLVDISPINFO *) lParam;

	switch (uMsg) {
	case msgNOTIFY:
		switch (nmhdr->code) {
		case LVN_GETDISPINFO:
			tableGetCellText((void *) data, fill->item.iItem, fill->item.iSubItem, &(fill->item.pszText));
			return 0;
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	/* the autosize behavior is simple: always autosize until the user manually sizes, then never autosize again (this is my best guess as to how GTK+ behaves) */
	case WM_NOTIFY:		/* from the contained header control */
		if (nmhdr->code == HDN_BEGINTRACK)
			tableStopColumnAutosize((void *) data);
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_SIZE:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
		/* TODO this causes weird issues with regards to item positioning on Windows XP */
		if (tableAutosizeColumns((void *) data)) {
			int i, nColumns;

			nColumns = tableColumnCount((void *) data);
			for (i = 0; i < nColumns; i++)
				if (SendMessageW(hwnd, LVM_SETCOLUMNWIDTH, (WPARAM) i, (LPARAM) LVSCW_AUTOSIZE_USEHEADER) == FALSE)
					xpanic("error resizing columns of results list view", GetLastError());
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, tableSubProc, id) == FALSE)
			xpanic("error removing Table subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Button", "tableSubProc()", uMsg);
	return 0;		/* unreached */
}

void setTableSubclass(HWND hwnd, void *data)
{
	if ((*fv_SetWindowSubclass)(hwnd, tableSubProc, 0, (DWORD_PTR) data) == FALSE)
		xpanic("error subclassing Table to give it its own event handler", GetLastError());
}

void tableAppendColumn(HWND hwnd, int index, LPWSTR name)
{
	LVCOLUMNW col;

	ZeroMemory(&col, sizeof (LVCOLUMNW));
	col.mask = LVCF_FMT | LVCF_TEXT | LVCF_SUBITEM | LVCF_ORDER;
	col.fmt = LVCFMT_LEFT;
	col.pszText = name;
	col.iSubItem = index;
	col.iOrder = index;
	if (SendMessageW(hwnd, LVM_INSERTCOLUMN, (WPARAM) index, (LPARAM) (&col)) == (LRESULT) -1)
		xpanic("error adding column to Table", GetLastError());
}

void tableUpdate(HWND hwnd, int nItems)
{
	if (SendMessageW(hwnd, LVM_SETITEMCOUNT, (WPARAM) nItems, 0) == 0)
		xpanic("error setting number of items in Table", GetLastError());
}

void tableAddExtendedStyles(HWND hwnd, LPARAM styles)
{
	/* the bits of WPARAM specify which bits of LPARAM to look for; having WPARAM == LPARAM ensures that only the bits we want to add are affected */
	SendMessageW(hwnd, LVM_SETEXTENDEDLISTVIEWSTYLE, (WPARAM) styles, styles);
}
