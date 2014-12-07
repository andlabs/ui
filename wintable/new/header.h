// 7 december 2014

static void makeHeader(struct table *t, HINSTANCE hInstance)
{
	t->header = CreateWindowExW(0,
		WC_HEADERW, L"",
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
