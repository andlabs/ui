// 29 november 2014

static void hscrollto(struct table *t, intptr_t newpos)
{
	SCROLLINFO si;
	RECT scrollArea;

	if (newpos < 0)
		newpos = 0;
	if (newpos > (t->width - t->hpagesize))
		newpos = (t->width - t->hpagesize);

	scrollArea = realClientRect(t);

	// negative because ScrollWindowEx() is "backwards"
	if (ScrollWindowEx(t->hwnd, -(newpos - t->hpos), 0,
		&scrollArea, &scrollArea, NULL, NULL,
		SW_ERASE | SW_INVALIDATE) == ERROR)
		abort();
	t->hpos = newpos;
	// TODO text in header controls doesn't redraw?

	// TODO put this in a separate function? same for vscroll?
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_POS | SIF_RANGE;
	si.nPage = t->hpagesize;
	si.nMin = 0;
	si.nMax = t->width - 1;		// nMax is inclusive
	si.nPos = t->hpos;
	SetScrollInfo(t->hwnd, SB_HORZ, &si, TRUE);

	// and finally reposition the header
	repositionHeader(t);
}

static void hscrollby(struct table *t, intptr_t n)
{
	hscrollto(t, t->hpos + n);
}

// unfortunately horizontal wheel scrolling was only added in Vista

static void hscroll(struct table *t, WPARAM wParam)
{
	SCROLLINFO si;
	intptr_t newpos;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_POS | SIF_TRACKPOS;
	if (GetScrollInfo(t->hwnd, SB_HORZ, &si) == 0)
		abort();

	newpos = t->hpos;
	switch (LOWORD(wParam)) {
	case SB_LEFT:
		newpos = 0;
		break;
	case SB_RIGHT:
		newpos = t->width - t->hpagesize;
		break;
	case SB_LINELEFT:
		newpos--;
		break;
	case SB_LINERIGHT:
		newpos++;
		break;
	case SB_PAGELEFT:
		newpos -= t->hpagesize;
		break;
	case SB_PAGERIGHT:
		newpos += t->hpagesize;
		break;
	case SB_THUMBPOSITION:
		newpos = (intptr_t) (si.nPos);
		break;
	case SB_THUMBTRACK:
		newpos = (intptr_t) (si.nTrackPos);
	}

	hscrollto(t, newpos);
}

static void recomputeHScroll(struct table *t)
{
	HDITEMW item;
	intptr_t i;
	int width = 0;
	RECT r;
	SCROLLINFO si;

	// TODO count dividers
	for (i = 0; i < t->nColumns; i++) {
		ZeroMemory(&item, sizeof (HDITEMW));
		item.mask = HDI_WIDTH;
		if (SendMessageW(t->header, HDM_GETITEM, (WPARAM) i, (LPARAM) (&item)) == FALSE)
			abort();
		width += item.cxy;
	}
	t->width = (intptr_t) width;

	if (GetClientRect(t->hwnd, &r) == 0)
		abort();
	t->hpagesize = r.right - r.left;

	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_RANGE;
	si.nPage = t->hpagesize;
	si.nMin = 0;
	si.nMax = t->width - 1;			// - 1 because endpoints inclusive
	SetScrollInfo(t->hwnd, SB_HORZ, &si, TRUE);
}
