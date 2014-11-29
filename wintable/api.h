// 29 november 2014

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
