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
			tableToggled(gotable, tnm->row, tnm->column);
			return 0;
		case tableNotificationSelectionChanged:
			tableSelectionChanged(gotable);
			return 0;
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

// TODO rename all of these functions to start with gotable, and all the exported ones in Go too
void gotableSetRowCount(HWND hwnd, intptr_t count)
{
	SendMessageW(hwnd, tableSetRowCount, 0, (LPARAM) (&count));
}

void tableAutosizeColumns(HWND hwnd, int nColumns)
{
	int i;

	for (i = 0; i < nColumns; i++)
		if (SendMessageW(hwnd, LVM_SETCOLUMNWIDTH, (WPARAM) i, (LPARAM) LVSCW_AUTOSIZE_USEHEADER) == FALSE)
			xpanic("error resizing columns of results list view", GetLastError());
}

intptr_t tableSelectedItem(HWND hwnd)
{
	intptr_t row;

	SendMessageW(hwnd, tableGetSelection, (WPARAM) (&row), (LPARAM) NULL);
	return row;
}

void tableSelectItem(HWND hwnd, intptr_t index)
{
	SendMessageW(hwnd, tableSetSelection, (WPARAM) (&index), (LPARAM) NULL);
}
