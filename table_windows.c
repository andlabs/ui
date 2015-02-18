// 28 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

#include "wintable/main.h"

// provided for cgo's benefit
LPWSTR xtableWindowClass = tableWindowClass;

void doInitTable(void)
{
	initTable(xpanic, fv__TrackMouseEvent);
}

static LRESULT CALLBACK tableSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	tableNM *tnm = (tableNM *) lParam;
	void *gotable = (void *) data;

	switch (uMsg) {
	case msgNOTIFY:
		switch (nmhdr->code) {
		case tableNotificationGetCellData:
			return tableGetCell(gotable, tnm);
		case tableNotificationFinishedWithCellData:
			tableFreeCellData(gotable, tnm->data);
			return 0;
		case tableNotificationCellCheckboxToggled:
			// TODO
		// TODO selection changed
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
/* TODO
	// see table.autoresize() in table_windows.go for the column autosize policy
	case WM_NOTIFY:		// from the contained header control
		if (nmhdr->code == HDN_BEGINTRACK)
			tableStopColumnAutosize(t->gotable);
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
*/
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, tableSubProc, id) == FALSE)
			xpanic("error removing Table subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Table", "tableSubProc()", uMsg);
	return 0;		// unreached
}

void setTableSubclass(HWND hwnd, void *data)
{
	if ((*fv_SetWindowSubclass)(hwnd, tableSubProc, 0, (DWORD_PTR) data) == FALSE)
		xpanic("error subclassing Table to give it its own event handler", GetLastError());
}

void tableAutosizeColumns(HWND hwnd, int nColumns)
{
	int i;

	for (i = 0; i < nColumns; i++)
		if (SendMessageW(hwnd, LVM_SETCOLUMNWIDTH, (WPARAM) i, (LPARAM) LVSCW_AUTOSIZE_USEHEADER) == FALSE)
			xpanic("error resizing columns of results list view", GetLastError());
}

// because Go won't let me do C.WPARAM(-1)
intptr_t tableSelectedItem(HWND hwnd)
{
	return (intptr_t) SendMessageW(hwnd, LVM_GETNEXTITEM, (WPARAM) -1, LVNI_SELECTED);
}

/*
TODO
void tableSelectItem(HWND hwnd, intptr_t index)
{
	LVITEMW item;
	LRESULT current;

	// via http://support.microsoft.com/kb/131284
	// we don't need to clear the other bits; Tables don't support cutting or drag/drop
	current = SendMessageW(hwnd, LVM_GETNEXTITEM, (WPARAM) -1, LVNI_SELECTED);
	if (current != (LRESULT) -1) {
		ZeroMemory(&item, sizeof (LVITEMW));
		item.mask = LVIF_STATE;
		item.state = 0;
		item.stateMask = LVIS_FOCUSED | LVIS_SELECTED;
		if (SendMessageW(hwnd, LVM_SETITEMSTATE, (WPARAM) current, (LPARAM) (&item)) == FALSE)
			xpanic("error deselecting current Table item", GetLastError());
	}
	if (index == -1)			// select nothing
		return;
	ZeroMemory(&item, sizeof (LVITEMW));
	item.mask = LVIF_STATE;
	item.state = LVIS_FOCUSED | LVIS_SELECTED;
	item.stateMask = LVIS_FOCUSED | LVIS_SELECTED;
	if (SendMessageW(hwnd, LVM_SETITEMSTATE, (WPARAM) index, (LPARAM) (&item)) == FALSE)
		xpanic("error selecting new Table item", GetLastError());
}
*/
