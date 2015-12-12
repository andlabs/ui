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
	// TODO only redraw the part of the client area where the new client went, if any
	// (TODO when — if — adding autoresize, figure this one out)

	// TODO send a notification for all rows?
}

// TODO what happens if the currently selected row is lost?
static void setRowCount(struct table *t, intptr_t rc)
{
	intptr_t old, i;

	old = t->count;
	t->count = rc;
	// we DO redraw everything because we don't want any rows that should no longer be there to remain on screen!
	updateAll(t);			// DONE
	// TODO reset checkbox and selection logic if the current row for both no longer exists
	// TODO really send all these notifications?
	if (old < t->count)		// rows added
		for (i = old; i < t->count; i++)
			NotifyWinEvent(EVENT_OBJECT_CREATE, t->hwnd, OBJID_CLIENT, i);
	else if (old > t->count)	// rows removed
		for (i = old; i > t->count; i++)
			NotifyWinEvent(EVENT_OBJECT_DESTROY, t->hwnd, OBJID_CLIENT, i);
	// TODO update existing rows?
}

HANDLER(apiHandlers)
{
	intptr_t *rcp;
	intptr_t row;

	switch (uMsg) {
	case WM_SETFONT:
		// don't free the old font; see http://blogs.msdn.com/b/oldnewthing/archive/2008/09/12/8945692.aspx
		t->font = (HFONT) wParam;
		SendMessageW(t->header, WM_SETFONT, wParam, lParam);
		// if we redraw, we have to redraw ALL of it; after all, the font changed!
		if (LOWORD(lParam) != FALSE)
			updateAll(t);			// DONE
		else
			update(t, FALSE);		// DONE
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
		setRowCount(t, *rcp);
		*lResult = 0;
		return TRUE;
	case tableGetSelection:
		rcp = (intptr_t *) wParam;
		if (rcp != NULL)
			*rcp = t->selectedRow;
		rcp = (intptr_t *) lParam;
		if (rcp != NULL)
			*rcp = t->selectedColumn;
		*lResult = 0;
		return TRUE;
	case tableSetSelection:
		// TODO does doselect() do validation?
		rcp = (intptr_t *) wParam;
		row = *rcp;
		rcp = (intptr_t *) lParam;
		if (rcp == NULL)
			if (row == -1)
				doselect(t, -1, -1);
			else		// select column 0, just like keyboard selections; TODO what if there aren't any columns?
				doselect(t, row, 0);
		else
			doselect(t, row, *rcp);
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
