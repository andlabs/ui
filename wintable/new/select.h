// 13 december 2014

// damn winsock
static void doselect(struct table *t, intptr_t row, intptr_t column)
{
	t->selectedRow = row;
	t->selectedColumn = column;
	// TODO scroll to ensure the full cell is visible
	// TODO redraw only the old and new columns /if there was no scrolling/
	InvalidateRect(t->hwnd, NULL, TRUE);
}

// TODO which WM_xBUTTONDOWNs?
HANDLER(mouseDownSelectHandler)
{
	struct rowcol rc;

	rc = lParamToRowColumn(t, lParam);
	// don't check if lParamToRowColumn() returned row -1 or column -1; we want deselection behavior
	doselect(t, rc.row, rc.column);
	*lResult = 0;
	return TRUE;
}
