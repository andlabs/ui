// 30 november 2014

static void resize(struct table *t)
{
	RECT r;
	SCROLLINFO si;

	// do this first so our scrollbar calculations can be correct
	repositionHeader(t);

	// now adjust the scrollbars
	r = realClientRect(t);
	t->pagesize = (r.bottom - r.top) / rowHeight(t);
	ZeroMemory(&si, sizeof (SCROLLINFO));
	si.cbSize = sizeof (SCROLLINFO);
	si.fMask = SIF_RANGE | SIF_PAGE;
	si.nMin = 0;
	si.nMax = t->count - 1;
	si.nPage = t->pagesize;
	SetScrollInfo(t->hwnd, SB_VERT, &si, TRUE);

	recomputeHScroll(t);
}

// TODO alter this so that only the visible columns are redrawn
// TODO this means rename controlSize to clientRect
static void drawItem(struct table *t, HDC dc, intptr_t i, LONG y, LONG height, RECT controlSize)
{
	RECT rsel;
	HBRUSH background;
	int textColor;
	WCHAR msg[100];
	RECT headeritem;
	intptr_t j;
	LRESULT xoff;
	IMAGELISTDRAWPARAMS ip;
	POINT pt;

	// TODO verify these two
	background = (HBRUSH) (COLOR_WINDOW + 1);
	textColor = COLOR_WINDOWTEXT;
	if (t->selected == i) {
		// these are the colors wine uses (http://source.winehq.org/source/dlls/comctl32/listview.c)
		// the two for unfocused are also suggested by http://stackoverflow.com/questions/10428710/windows-forms-inactive-highlight-color
		background = (HBRUSH) (COLOR_HIGHLIGHT + 1);
		textColor = COLOR_HIGHLIGHTTEXT;
		if (GetFocus() != t->hwnd) {
			background = (HBRUSH) (COLOR_BTNFACE + 1);
			textColor = COLOR_BTNTEXT;
		}
	}

	// first fill the selection rect
	// note that this already only draws the visible area
	rsel.left = controlSize.left;
	rsel.top = y;
	rsel.right = controlSize.right - controlSize.left;
	rsel.bottom = y + height;
	if (FillRect(dc, &rsel, background) == 0)
		abort();

	// TODO double-check to see if this takes any parameters
	xoff = SendMessageW(t->header, HDM_GETBITMAPMARGIN, 0, 0);
	// now adjust for horizontal scrolling
	xoff -= t->hpos;

	// now draw the cells
	if (SetTextColor(dc, GetSysColor(textColor)) == CLR_INVALID)
		abort();
	if (SetBkMode(dc, TRANSPARENT) == 0)
		abort();
	for (j = 0; j < t->nColumns; j++) {
		if (SendMessageW(t->header, HDM_GETITEMRECT, (WPARAM) j, (LPARAM) (&headeritem)) == 0)
			abort();
		switch (t->columnTypes[j]) {
		case tableColumnText:
			rsel.left = headeritem.left + xoff;
			rsel.top = y;
			rsel.right = headeritem.right;
			rsel.bottom = y + height;
			// TODO vertical center in case the height is less than the icon height?
			if (DrawTextExW(dc, msg, wsprintf(msg, L"Item %d", i), &rsel, DT_END_ELLIPSIS | DT_LEFT | DT_NOPREFIX | DT_SINGLELINE, NULL) == 0)
				abort();
			break;
		case tableColumnImage:
			// TODO vertically center if image is smaller than text height
			// TODO same for checkboxes
			ZeroMemory(&ip, sizeof (IMAGELISTDRAWPARAMS));
			ip.cbSize = sizeof (IMAGELISTDRAWPARAMS);
			ip.himl = t->checkboxes;//t->imagelist;
			ip.i = (i%8);//0;
			ip.hdcDst = dc;
			ip.x = headeritem.left + xoff;
			ip.y = y;
			ip.cx = 0;		// draw whole image
			ip.cy = 0;
			ip.xBitmap = 0;
			ip.yBitmap = 0;
			ip.rgbBk = CLR_NONE;
			ip.fStyle = ILD_NORMAL | ILD_SCALE;		// TODO alpha-blend; ILD_DPISCALE?
			// TODO ILS_ALPHA?
			if (ImageList_DrawIndirect(&ip) == 0)
				abort();
			break;
		case tableColumnCheckbox:
			// TODO
			break;
		}
		if (t->selected == i && t->focusedColumn == j) {
			rsel.left = headeritem.left;
			rsel.top = y;
			rsel.right = headeritem.right;
			rsel.bottom = y + height;
			if (DrawFocusRect(dc, &rsel) == 0)
				abort();
		}
	}
}

static void drawItems(struct table *t, HDC dc, RECT cliprect)
{
	HFONT thisfont, prevfont;
	LONG height;
	LONG y;
	intptr_t i;
	RECT controlSize;		// for filling the entire selected row
	intptr_t first, last;

	if (GetClientRect(t->hwnd, &controlSize) == 0)
		abort();

	height = rowHeight(t);

	thisfont = t->font;		// in case WM_SETFONT happens before we return
	prevfont = (HFONT) SelectObject(dc, thisfont);
	if (prevfont == NULL)
		abort();

	// ignore anything beneath the header
	if (cliprect.top < t->headerHeight)
		cliprect.top = t->headerHeight;
	// now let's pretend the header isn't there
	// we only need it in (or rather, before) the drawItem() calls below
	cliprect.top -= t->headerHeight;
	cliprect.bottom -= t->headerHeight;

	// see http://blogs.msdn.com/b/oldnewthing/archive/2003/07/29/54591.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2003/07/30/54600.aspx
	// we need to add t->firstVisible here because cliprect is relative to the visible area
	first = (cliprect.top / height) + t->firstVisible;
	if (first < 0)
		first = 0;
	last = lastVisible(t, cliprect, height);

	// now for the first y, discount firstVisible
	y = (first - t->firstVisible) * height;
	// and offset by the header height
	y += t->headerHeight;
	for (i = first; i < last; i++) {
		drawItem(t, dc, i, y, height, controlSize);
		y += height;
	}

	// reset everything
	if (SelectObject(dc, prevfont) != (HGDIOBJ) (thisfont))
		abort();
}
