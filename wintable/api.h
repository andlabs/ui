// 8 december 2014

static void addColumn(struct table *t, WPARAM wParam, LPARAM lParam)
{
	t->nColumns++;
	t->columnTypes = (int *) tableRealloc(t->columnTypes, t->nColumns * sizeof (int), "adding the new column type to the current Table's list of column types");
	t->columnTypes[t->nColumns - 1] = (int) wParam;
	// TODO make a panicNoErrCode() or panicArg() for this
	if (t->columnTypes[t->nColumns - 1] >= nTableColumnTypes)
		panic("invalid column type passed to tableAddColumn");
	headerAddColumn(t, (WCHAR *) lParam);
	update(t, TRUE);
}

HANDLER(apiHandlers)
{
	intptr_t *rcp;

	switch (uMsg) {
	case WM_SETFONT:
		// don't free the old font; see http://blogs.msdn.com/b/oldnewthing/archive/2008/09/12/8945692.aspx
		t->font = (HFONT) wParam;
		SendMessageW(t->header, WM_SETFONT, wParam, lParam);
		update(t, LOWORD(lParam) != FALSE);
		// TODO is this needed?
		if (LOWORD(lParam) != FALSE)
			// TODO check error
			InvalidateRect(t->hwnd, NULL, TRUE);
		*lResult = 0;
		return TRUE;
	case WM_GETFONT:
		*lResult = (LRESULT) (t->font);
		return TRUE;
	case tableAddColumn:
		addColumn(t, wParam, lParam);
		*lResult = 0;
		return TRUE;
	case tableSetRowCount:
		rcp = (intptr_t *) lParam;
		t->count = *rcp;
		// TODO shouldn't we just redraw everything?
		update(t, TRUE);
		*lResult = 0;
		return TRUE;
	}
	return FALSE;
}

static LRESULT notify(struct table *t, UINT code, intptr_t row, intptr_t column, uintptr_t data)
{
	tableNM nm;

	ZeroMemory(&nm, sizeof (tableNM));
	nm.nmhdr.hwndFrom = t->hwnd;
	// TODO check for error from here? 0 is a valid ID (IDCANCEL)
	nm.nmhdr.idFrom = GetDlgCtrlID(t->hwnd);
	nm.nmhdr.code = code;
	nm.row = row;
	nm.column = column;
	nm.columnType = t->columnTypes[nm.column];
	nm.data = data;
	// TODO check for error from GetParent()?
	return SendMessageW(GetParent(t->hwnd), WM_NOTIFY, (WPARAM) (nm.nmhdr.idFrom), (LPARAM) (&nm));
}
