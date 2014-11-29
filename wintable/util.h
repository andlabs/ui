// 29 november 2014

static LONG rowHeight(struct table *t)
{
	HFONT thisfont, prevfont;
	TEXTMETRICW tm;
	HDC dc;
	LONG ret;

	dc = GetDC(t->hwnd);
	if (dc == NULL)
		abort();
	thisfont = t->font;		// in case WM_SETFONT happens before we return
	prevfont = (HFONT) SelectObject(dc, thisfont);
	if (prevfont == NULL)
		abort();
	if (GetTextMetricsW(dc, &tm) == 0)
		abort();
	if (SelectObject(dc, prevfont) != (HGDIOBJ) (thisfont))
		abort();
	if (ReleaseDC(t->hwnd, dc) == 0)
		abort();
	ret = tm.tmHeight;
	if (ret < t->imagelistHeight)
		ret = t->imagelistHeight;
	if (ret < t->checkboxHeight)
		ret = t->checkboxHeight;
	return ret;
}

static void redrawAll(struct table *t)
{
	if (InvalidateRect(t->hwnd, NULL, TRUE) == 0)
		abort();
	if (UpdateWindow(t->hwnd) == 0)
		abort();
}

static RECT realClientRect(struct table *t)
{
	RECT r;

	if (GetClientRect(t->hwnd, &r) == 0)
		abort();
	r.top += t->headerHeight;
	return r;
}

static void repositionHeader(struct table *t)
{
	RECT r;
	HDLAYOUT headerlayout;
	WINDOWPOS headerpos;

	if (GetClientRect(t->hwnd, &r) == 0)		// use the whole client rect
		abort();
	// grow the rectangle to the left to fake scrolling
	r.left -= t->hpos;
	headerlayout.prc = &r;
	headerlayout.pwpos = &headerpos;
	if (SendMessageW(t->header, HDM_LAYOUT, 0, (LPARAM) (&headerlayout)) == FALSE)
		abort();
	if (SetWindowPos(t->header, headerpos.hwndInsertAfter, headerpos.x, headerpos.y, headerpos.cx, headerpos.cy, headerpos.flags | SWP_SHOWWINDOW) == 0)
		abort();
	t->headerHeight = headerpos.cy;
}

// this counts partially visible rows
// for all fully visible rows use t->pagesize
// cliprect and rowHeight must be specified here to avoid recomputing things multiple times
static intptr_t lastVisible(struct table *t, RECT cliprect, LONG rowHeight)
{
	intptr_t last;

	last = ((cliprect.bottom + rowHeight - 1) / rowHeight) + t->firstVisible;
	if (last >= t->count)
		last = t->count;
	return last;
}

static void redrawRow(struct table *t, intptr_t row)
{
	RECT r;
	intptr_t height;

	r = realClientRect(t);
	height = rowHeight(t);
	if (row < t->firstVisible || row > lastVisible(t, r, height))		// not visible; don't bother
		return;
	r.top = (row - t->firstVisible) * height + t->headerHeight;
	r.bottom = r.top + height;
	// keep the width and height the same; it spans the client area anyway
	if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
		abort();
	if (UpdateWindow(t->hwnd) == 0)
		abort();
}

static intptr_t hitTestColumn(struct table *t, int x)
{
	HDITEMW item;
	intptr_t i;

	// TODO count dividers
	for (i = 0; i < t->nColumns; i++) {
		ZeroMemory(&item, sizeof (HDITEMW));
		item.mask = HDI_WIDTH;
		if (SendMessageW(t->header, HDM_GETITEM, (WPARAM) i, (LPARAM) (&item)) == FALSE)
			abort();
		if (x < item.cxy)
			return i;
		x -= item.cxy;		// not yet
	}
	// no column
	return -1;
}

static void lParamToRowColumn(struct table *t, LPARAM lParam, intptr_t *row, intptr_t *column)
{
	int x, y;
	LONG h;

	x = GET_X_LPARAM(lParam);
	y = GET_Y_LPARAM(lParam);
	h = rowHeight(t);
	y += t->firstVisible * h;
	y -= t->headerHeight;
	y /= h;
	if (row != NULL) {
		*row = y;
		if (*row >= t->count)
			*row = -1;
	}
	if (column != NULL)
		*column = hitTestColumn(t, x);
}
