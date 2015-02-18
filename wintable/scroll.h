// 9 december 2014

struct scrollParams {
	intptr_t *pos;
	intptr_t pagesize;
	intptr_t length;
	intptr_t scale;
	void (*post)(struct table *);
	int *wheelCarry;
};

static void scrollto(struct table *t, int which, struct scrollParams *p, intptr_t pos)
{
	RECT scrollArea;
	SCROLLINFO si;
	intptr_t xamount, yamount;

	if (pos < 0)
		pos = 0;
	if (pos > p->length - p->pagesize)
		pos = p->length - p->pagesize;
	// TODO this shouldn't have been necessary but alas
	// TODO the logic is really intended for the whole y origin thing in the scrollbar series; fix that
	if (pos < 0)
		pos = 0;

	// we don't want to scroll the header
	if (GetClientRect(t->hwnd, &scrollArea) == 0)
		panic("error getting Table client rect for scrollto()");
	scrollArea.top += t->headerHeight;

	// negative because ScrollWindowEx() is "backwards"
	xamount = -(pos - *(p->pos)) * p->scale;
	yamount = 0;
	if (which == SB_VERT) {
		yamount = xamount;
		xamount = 0;
	}

	if (ScrollWindowEx(t->hwnd, xamount, yamount,
		&scrollArea, &scrollArea, NULL, NULL,
		SW_ERASE | SW_INVALIDATE) == ERROR)
;// TODO failure case ignored for now; see https://bugs.winehq.org/show_bug.cgi?id=37706
//		panic("error scrolling Table");
	// TODO call UpdateWindow()?

	*(p->pos) = pos;

	// now commit our new scrollbar setup...
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_POS | SIF_RANGE;
	si.nPage = p->pagesize;
	si.nMin = 0;
	si.nMax = p->length - 1;		// endpoint inclusive
	si.nPos = *(p->pos);
	SetScrollInfo(t->hwnd, which, &si, TRUE);

	if (p->post != NULL)
		(*(p->post))(t);

	// EVENT_OBJECT_CONTENTSCROLLED is Vista and up only
	// TODO send state changes for all affected rows/cells?
}

static void scrollby(struct table *t, int which, struct scrollParams *p, intptr_t delta)
{
	scrollto(t, which, p, *(p->pos) + delta);
}

static void scroll(struct table *t, int which, struct scrollParams *p, WPARAM wParam, LPARAM lParam)
{
	intptr_t pos;
	SCROLLINFO si;

	pos = *(p->pos);
	switch (LOWORD(wParam)) {
	case SB_LEFT:			// also SB_TOP
		pos = 0;
		break;
	case SB_RIGHT:		// also SB_BOTTOM
		pos = p->length - p->pagesize;
		break;
	case SB_LINELEFT:		// also SB_LINEUP
		pos--;
		break;
	case SB_LINERIGHT:		// also SB_LINEDOWN
		pos++;
		break;
	case SB_PAGELEFT:		// also SB_PAGEUP
		pos -= p->pagesize;
		break;
	case SB_PAGERIGHT:	// also SB_PAGEDOWN
		pos += p->pagesize;
		break;
	case SB_THUMBPOSITION:
		ZeroMemory(&si, sizeof (SCROLLINFO));
		si.cbSize = sizeof (SCROLLINFO);
		si.fMask = SIF_POS;
		if (GetScrollInfo(t->hwnd, which, &si) == 0)
			panic("error getting thumb position for scroll() in Table");
		pos = si.nPos;
		break;
	case SB_THUMBTRACK:
		ZeroMemory(&si, sizeof (SCROLLINFO));
		si.cbSize = sizeof (SCROLLINFO);
		si.fMask = SIF_TRACKPOS;
		if (GetScrollInfo(t->hwnd, which, &si) == 0)
			panic("error getting thumb track position for scroll() in Table");
		pos = si.nTrackPos;
		break;
	}
	scrollto(t, which, p, pos);
}

static void wheelscroll(struct table *t, int which, struct scrollParams *p, WPARAM wParam, LPARAM lParam)
{
	int delta;
	int lines;
	UINT scrollAmount;

	delta = GET_WHEEL_DELTA_WPARAM(wParam);
	// TODO make a note of what the appropriate hscroll constant is
	if (SystemParametersInfoW(SPI_GETWHEELSCROLLLINES, 0, &scrollAmount, 0) == 0)
		// TODO use scrollAmount == 3 instead?
		panic("error getting wheel scroll amount in wheelscroll()");
	if (scrollAmount == WHEEL_PAGESCROLL)
		scrollAmount = p->pagesize;
	if (scrollAmount == 0)		// no mouse wheel scrolling (or t->pagesize == 0)
		return;
	// the rest of this is basically http://blogs.msdn.com/b/oldnewthing/archive/2003/08/07/54615.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2003/08/11/54624.aspx
	// see those pages for information on subtleties
	delta += *(p->wheelCarry);
	lines = delta * ((int) scrollAmount) / WHEEL_DELTA;
	*(p->wheelCarry) = delta - lines * WHEEL_DELTA / ((int) scrollAmount);
	scrollby(t, which, p, -lines);
}
