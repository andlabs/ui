// 29 november 2014

static void vscrollto(struct table *t, intptr_t newpos)
{
	SCROLLINFO si;
	RECT scrollArea;

	if (newpos < 0)
		newpos = 0;
	if (newpos > (t->count - t->pagesize))
		newpos = (t->count - t->pagesize);

	scrollArea = realClientRect(t);

	// negative because ScrollWindowEx() is "backwards"
	if (ScrollWindowEx(t->hwnd, 0, (-(newpos - t->firstVisible)) * rowHeight(t),
		&scrollArea, &scrollArea, NULL, NULL,
		SW_ERASE | SW_INVALIDATE) == ERROR)
		abort();
	t->firstVisible = newpos;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_POS | SIF_RANGE;
	si.nPage = t->pagesize;
	si.nMin = 0;
	si.nMax = t->count - 1;		// nMax is inclusive
	si.nPos = t->firstVisible;
	SetScrollInfo(t->hwnd, SB_VERT, &si, TRUE);
}

static void vscrollby(struct table *t, intptr_t n)
{
	vscrollto(t, t->firstVisible + n);
}

static void wheelscroll(struct table *t, WPARAM wParam)
{
	int delta;
	int lines;
	UINT scrollAmount;

	delta = GET_WHEEL_DELTA_WPARAM(wParam);
	if (SystemParametersInfoW(SPI_GETWHEELSCROLLLINES, 0, &scrollAmount, 0) == 0)
		abort();
	if (scrollAmount == WHEEL_PAGESCROLL)
		scrollAmount = t->pagesize;
	if (scrollAmount == 0)		// no mouse wheel scrolling (or t->pagesize == 0)
		return;
	// the rest of this is basically http://blogs.msdn.com/b/oldnewthing/archive/2003/08/07/54615.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2003/08/11/54624.aspx
	// see those pages for information on subtleties
	delta += t->wheelCarry;
	lines = delta * ((int) scrollAmount) / WHEEL_DELTA;
	t->wheelCarry = delta - lines * WHEEL_DELTA / ((int) scrollAmount);
	vscrollby(t, -lines);
}

static void vscroll(struct table *t, WPARAM wParam)
{
	SCROLLINFO si;
	intptr_t newpos;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_POS | SIF_TRACKPOS;
	if (GetScrollInfo(t->hwnd, SB_VERT, &si) == 0)
		abort();

	newpos = t->firstVisible;
	switch (LOWORD(wParam)) {
	case SB_TOP:
		newpos = 0;
		break;
	case SB_BOTTOM:
		newpos = t->count - t->pagesize;
		break;
	case SB_LINEUP:
		newpos--;
		break;
	case SB_LINEDOWN:
		newpos++;
		break;
	case SB_PAGEUP:
		newpos -= t->pagesize;
		break;
	case SB_PAGEDOWN:
		newpos += t->pagesize;
		break;
	case SB_THUMBPOSITION:
		newpos = (intptr_t) (si.nPos);
		break;
	case SB_THUMBTRACK:
		newpos = (intptr_t) (si.nTrackPos);
	}

	vscrollto(t, newpos);
}
