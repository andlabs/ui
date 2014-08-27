// 28 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

// provided for cgo's benefit
LPWSTR xWC_LISTVIEW = WC_LISTVIEW;

static void handle(HWND hwnd, LPARAM lParam, void (*handler)(void *, int, int), void *data)
{
	LVHITTESTINFO ht;

	ZeroMemory(&ht, sizeof (LVHITTESTINFO));
	ht.pt.x = GET_X_LPARAM(lParam);
	ht.pt.y = GET_Y_LPARAM(lParam);
	if (SendMessageW(hwnd, LVM_SUBITEMHITTEST, 0, (LPARAM) (&ht)) == (LRESULT) -1) {
		(*handler)(data, -1, -1);
		return;		// no item
	}
	if (ht.flags != LVHT_ONITEMSTATEICON) {
		(*handler)(data, -1, -1);
		return;		// not on a checkbox
	}
	(*handler)(data, ht.iItem, ht.iSubItem);
}

struct tableData {
	void *gotable;
	HIMAGELIST imagelist;
	HTHEME theme;
	HIMAGELIST checkboxImageList;
};

static void tableLoadImageList(HWND hwnd, struct tableData *t, HIMAGELIST new)
{
	HIMAGELIST old;

	old = t->imagelist;
	t->imagelist = new;
	applyImageList(hwnd, LVM_SETIMAGELIST, LVSIL_SMALL, t->imagelist, old);
}

static void tableSetCheckboxImageList(HWND hwnd, struct tableData *t)
{
	HIMAGELIST old;

	old = t->checkboxImageList;
	t->checkboxImageList = makeCheckboxImageList(hwnd, &t->theme);
	applyImageList(hwnd, LVM_SETIMAGELIST, LVSIL_STATE, t->checkboxImageList, old);
	// thanks to Jonathan Potter (http://stackoverflow.com/questions/25354448/why-do-my-owner-data-list-view-state-images-come-up-as-blank-on-windows-xp)
	if (SendMessageW(hwnd, LVM_SETCALLBACKMASK, LVIS_STATEIMAGEMASK, 0) == FALSE)
		xpanic("error marking state image list as application-managed", GetLastError());
}

static LRESULT CALLBACK tableSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	NMLVDISPINFOW *fill = (NMLVDISPINFO *) lParam;
	NMLISTVIEW *nlv = (NMLISTVIEW *) lParam;
	struct tableData *t = (struct tableData *) data;

	switch (uMsg) {
	case msgNOTIFY:
		switch (nmhdr->code) {
		case LVN_GETDISPINFO:
			tableGetCell(t->gotable, &(fill->item));
			return 0;
		case LVN_ITEMCHANGED:
			if ((nlv->uChanged & LVIF_STATE) == 0)
				break;
			// if both old and new states have the same value for the selected bit, then the selection state did not change, regardless of selected or deselected
			if ((nlv->uOldState & LVIS_SELECTED) == (nlv->uNewState & LVIS_SELECTED))
				break;
			tableSelectionChanged(t->gotable);
			return 0;
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_MOUSEMOVE:
		handle(hwnd, lParam, tableSetHot, t->gotable);
		// and let the list view do its thing
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_LBUTTONDOWN:
	case WM_LBUTTONDBLCLK:			// listviews have CS_DBLCICKS; check this to better mimic the behavior of a real checkbox
		handle(hwnd, lParam, tablePushed, t->gotable);
		// and let the list view do its thing
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_LBUTTONUP:
		handle(hwnd, lParam, tableToggled, t->gotable);
		// and let the list view do its thing
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_MOUSELEAVE:
		// TODO doesn't work
		tablePushed(t->gotable, -1, -1);			// in case button held as drag out
		// and let the list view do its thing
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case msgLoadImageList:
		tableLoadImageList(hwnd, t, (HIMAGELIST) lParam);
		return 0;
	case msgTableMakeInitialCheckboxImageList:
		tableSetCheckboxImageList(hwnd, t);
		return 0;
	case WM_THEMECHANGED:
		tableSetCheckboxImageList(hwnd, t);
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	// see table.autoresize() in table_windows.go for the column autosize policy
	case WM_NOTIFY:		// from the contained header control
		if (nmhdr->code == HDN_BEGINTRACK)
			tableStopColumnAutosize(t->gotable);
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		free(t);
		if ((*fv_RemoveWindowSubclass)(hwnd, tableSubProc, id) == FALSE)
			xpanic("error removing Table subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Button", "tableSubProc()", uMsg);
	return 0;		// unreached
}

void setTableSubclass(HWND hwnd, void *data)
{
	struct tableData *t;

	t = (struct tableData *) malloc(sizeof (struct tableData));
	if (t == NULL)
		xpanic("error allocating structure for Table extra data", GetLastError());
	ZeroMemory(t, sizeof (struct tableData));
	t->gotable = data;
	if ((*fv_SetWindowSubclass)(hwnd, tableSubProc, 0, (DWORD_PTR) t) == FALSE)
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
	// the bits of WPARAM specify which bits of LPARAM to look for; having WPARAM == LPARAM ensures that only the bits we want to add are affected
	SendMessageW(hwnd, LVM_SETEXTENDEDLISTVIEWSTYLE, (WPARAM) styles, styles);
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
