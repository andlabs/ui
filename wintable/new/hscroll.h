// 9 december 2014

// forward declaration needed here
static void repositionHeader(struct table *);

static void hscrollto(struct table *t, intptr_t pos)
{
	RECT scrollArea;
	SCROLLINFO si;

	if (pos < 0)
		pos = 0;
	if (pos > t->width - t->hpagesize)
		pos = t->width - t->hpagesize;

	// we don't want to scroll the header
	if (GetClientRect(t->hwnd, &scrollArea) == 0)
		panic("error getting Table client rect for hscrollto()");
	scrollArea.top += t->headerHeight;

	// negative because ScrollWindowEx() is "backwards"
	if (ScrollWindowEx(t->hwnd, -(pos - t->hscrollpos), 0,
		&scrollArea, &scrollArea, NULL, NULL,
		SW_ERASE | SW_INVALIDATE) == ERROR)
		panic("error horizontally scrolling Table");
	// TODO call UpdateWindow()?

	t->hscrollpos = pos;

	// now commit our new scrollbar setup...
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_PAGE | SIF_POS | SIF_RANGE;
	// the width of scrollArea is unchanged here; use it
	t->hpagesize = scrollArea.right - scrollArea.left;
	si.nPage = t->hpagesize;
	si.nMin = 0;
	si.nMax = t->width - 1;		// endpoint inclusive
	si.nPos = t->hscrollpos;
	SetScrollInfo(t->hwnd, SB_HORZ, &si, TRUE);

	// and finally move the header
	repositionHeader(t);
}

static void hscrollby(struct table *t, intptr_t delta)
{
	hscrollto(t, t->hscrollpos + delta);
}

static void hscroll(struct table *t, WPARAM wParam, LPARAM lParam)
{
	intptr_t pos;
	SCROLLINFO si;

	pos = t->hscrollpos;
	switch (LOWORD(wParam)) {
	case SB_LEFT:
		pos = 0;
		break;
	case SB_RIGHT:
		pos = t->width - t->hpagesize;
		break;
	case SB_LINELEFT:
		pos--;
		break;
	case SB_LINERIGHT:
		pos++;
		break;
	case SB_PAGELEFT:
		pos -= t->hpagesize;
		break;
	case SB_PAGERIGHT:
		pos += t->hpagesize;
		break;
	case SB_THUMBPOSITION:
		ZeroMemory(&si, sizeof (SCROLLINFO));
		si.cbSize = sizeof (SCROLLINFO);
		si.fMask = SIF_POS;
		if (GetScrollInfo(t->hwnd, SB_HORZ, &si) == 0)
			panic("error getting thumb position for WM_HSCROLL in Table");
		pos = si.nPos;
		break;
	case SB_THUMBTRACK:
		ZeroMemory(&si, sizeof (SCROLLINFO));
		si.cbSize = sizeof (SCROLLINFO);
		si.fMask = SIF_TRACKPOS;
		if (GetScrollInfo(t->hwnd, SB_HORZ, &si) == 0)
			panic("error getting thumb track position for WM_HSCROLL in Table");
		pos = si.nTrackPos;
		break;
	}
	hscrollto(t, pos);
}

// TODO find out if we can indicriminately check for WM_WHEELHSCROLL
HANDLER(hscrollHandler)
{
	if (uMsg != WM_HSCROLL)
		return FALSE;
	hscroll(t, wParam, lParam);
	*lResult = 0;
	return TRUE;
}

// TODO when we write vscroll.h, see just /what/ is common so we can isolate it
