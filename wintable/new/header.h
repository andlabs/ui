// 7 december 2014

static void makeHeader(struct table *t, HINSTANCE hInstance)
{
	t->header = CreateWindowExW(0,
		WC_HEADERW, L"",
		// don't set WS_VISIBLE; according to MSDN we create the header hidden as part of setting the initial position (http://msdn.microsoft.com/en-us/library/windows/desktop/ff485935%28v=vs.85%29.aspx)
		// TODO WS_BORDER?
		// TODO is HDS_HOTTRACK needed?
		WS_CHILD | HDS_FULLDRAG | HDS_HORZ | HDS_HOTTRACK,
		0, 0, 0, 0,		// no initial size
		t->hwnd, (HMENU) 100, hInstance, NULL);
	if (t->header == NULL)
		panic("error creating Table header");
}

static void destroyHeader(struct table *t)
{
	if (DestroyWindow(t->header) == 0)
		panic("error destroying Table header");
}

static void repositionHeader(struct table *t)
{
	RECT r;
	WINDOWPOS wp;
	HDLAYOUT l;

	if (GetClientRect(t->hwnd, &r) == 0)
		panic("error getting client rect for Table header repositioning");
	l.prc = &r;
	l.pwpos = &wp;
	if (SendMessageW(t->header, HDM_LAYOUT, 0, (LPARAM) (&l)) == FALSE)
		panic("error getting new Table header position");
	if (SetWindowPos(t->header, wp.hwndInsertAfter,
		wp.x, wp.y, wp.cx, wp.cy,
		// see above on showing the header here instead of in the CreateWindowExW() call
		wp.flags | SWP_SHOWWINDOW) == 0)
		panic("error repositioning Table header");
}

static void headerAddColumn(struct table *t, WCHAR *name)
{
	HDITEMW item;

	ZeroMemory(&item, sizeof (HDITEMW));
	item.mask = HDI_WIDTH | HDI_TEXT | HDI_FORMAT;
	item.cxy = 200;		// TODO
	item.pszText = name;
	item.fmt = HDF_LEFT | HDF_STRING;
	// TODO replace 100 with (t->nColumns - 1)
	if (SendMessage(t->header, HDM_INSERTITEM, (WPARAM) (100), (LPARAM) (&item)) == (LRESULT) (-1))
		panic("error adding column to Table header");
}

HANDLER(headerNotifyHandler)
{
	NMHDR *nmhdr = (NMHDR *) lParam;

	if (nmhdr->hwndFrom != t->header)
		return FALSE;
	if (nmhdr->code != HDN_ITEMCHANGED)
		return FALSE;
	// TODO should this be last?
	recomputeHScroll(t);
	// TODO make more intelligent
	InvalidateRect(t->hwnd, NULL, TRUE);
	// TODO UpdateWindow()?
	*lResult = 0;
	return TRUE;
}
