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
	// we fake horizontal scrolling here by extending the client rect to the left by the scroll position
	r.left -= t->hscrollpos;
	l.prc = &r;
	l.pwpos = &wp;
	if (SendMessageW(t->header, HDM_LAYOUT, 0, (LPARAM) (&l)) == FALSE)
		panic("error getting new Table header position");
	if (SetWindowPos(t->header, wp.hwndInsertAfter,
		wp.x, wp.y, wp.cx, wp.cy,
		// see above on showing the header here instead of in the CreateWindowExW() call
		wp.flags | SWP_SHOWWINDOW) == 0)
		panic("error repositioning Table header");
	t->headerHeight = wp.cy;
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

// TODO make a better name for this?
// TODO move to hscroll.h?
// TODO organize this in general...
// TODO because of this function's new extended functionality only hscrollto() is allowed to call repositionHeader()
static void updateTableWidth(struct table *t)
{
	HDITEMW item;
	intptr_t i;
	RECT client;

	t->width = 0;
	// TODO count dividers?
	// TODO use columnWidth()
	for (i = 0; i < t->nColumns; i++) {
		ZeroMemory(&item, sizeof (HDITEMW));
		item.mask = HDI_WIDTH;
		if (SendMessageW(t->header, HDM_GETITEM, (WPARAM) i, (LPARAM) (&item)) == FALSE)
			panic("error getting Table column width for updateTableWidth()");
		t->width += item.cxy;
	}

	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting Table client rect in updateTableWidth()");
	t->hpagesize = client.right - client.left;

	// this part is critical: if we resize the columns to less than the client area width, then the following hscrollby() will make t->hscrollpos negative, which does very bad things
	// note to self: do this regardless of whether the table width or the client area width was changed
	if (t->hpagesize > t->width)
		t->hpagesize = t->width;

	// do a dummy scroll to update the horizontal scrollbar to use the new width
	hscrollby(t, 0);
}

HANDLER(headerNotifyHandler)
{
	NMHDR *nmhdr = (NMHDR *) lParam;

	if (nmhdr->hwndFrom != t->header)
		return FALSE;
	if (nmhdr->code != HDN_ITEMCHANGED)
		return FALSE;
	updateTableWidth(t);
	// TODO make more intelligent
	InvalidateRect(t->hwnd, NULL, TRUE);
	// TODO UpdateWindow()?
	*lResult = 0;
	return TRUE;
}
