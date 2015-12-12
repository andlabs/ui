// 7 december 2014

// TODO verify header events (double-clicking on a divider, for example)

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

// to avoid weird bugs, the only functions allowed to call this one are the horizontal scroll functions
// when we need to reposition the header in a situation other than a user-initiated scroll, we use a dummy scroll (hscrollby(t, 0))
// see update() in update.h
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

// TODO is this triggered if we programmatically move headers (for autosizing)?
HANDLER(headerNotifyHandler)
{
	NMHDR *nmhdr = (NMHDR *) lParam;

	if (nmhdr->hwndFrom != t->header)
		return FALSE;
	if (nmhdr->code != HDN_ITEMCHANGED)
		return FALSE;
	update(t, TRUE);
	// TODO make more intelligent
	// to do this, we have to redraw the column to the left of the divider that was dragged and scroll everything to the right normally (leaving the hole that was scrolled invalidated as well)
	// of course, this implies that dragging a divider only resizes the column of which it is the right side of and moves all others
	InvalidateRect(t->hwnd, NULL, TRUE);
	*lResult = 0;
	return TRUE;
}
