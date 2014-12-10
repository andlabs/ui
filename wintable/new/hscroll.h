// 9 december 2014

static void hscrollto(struct table *t, intptr_t pos)
{
	// TODO
}

static void hscrollby(struct table *t, intptr_t delta)
{
	// TODO
}

static void hscroll(struct table *t, WPARAM wParam, LPARAM lParam)
{
	// TODO
}

static void recomputeHScroll(struct table *t)
{
	SCROLLINFO si;
	HDITEMW item;
	intptr_t i, width;
	RECT r;

	width = 0;
	// TODO count dividers?
	for (i = 0; i < t->nColumns; i++) {
		ZeroMemory(&item, sizeof (HDITEMW));
		item.mask = HDI_WIDTH;
		if (SendMessageW(t->header, HDM_GETITEM, (WPARAM) i, (LPARAM) (&item)) == FALSE)
			panic("error getting Table column width for recomputeHScroll()");
		width += item.cxy;
	}
	if (GetClientRect(t->hwnd, &r) == 0)
		panic("error getting Table client rect for recomputeHScroll()");
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_RANGE;
	si.nPage = r.right - r.left;
	si.nMin = 0;
	si.nMax = width - 1;		// endpoint inclusive
	SetScrollInfo(t->hwnd, SB_HORZ, &si, TRUE);
	// TODO what happens if the above call renders the current scroll position moot?
}
