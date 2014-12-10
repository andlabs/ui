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
	RECT r;

	if (GetClientRect(t->hwnd, &r) == 0)
		panic("error getting Table client rect for recomputeHScroll()");
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_RANGE;
	si.nPage = r.right - r.left;
	si.nMin = 0;
	si.nMax = t->width - 1;		// endpoint inclusive
	SetScrollInfo(t->hwnd, SB_HORZ, &si, TRUE);
	// TODO what happens if the above call renders the current scroll position moot?
}
