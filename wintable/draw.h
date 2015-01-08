// 8 december 2014

struct drawCellParams {
	intptr_t row;
	intptr_t column;
	LONG x;
	LONG y;
	LONG width;		// of column
	LONG height;		// rowHeight()
	LRESULT xoff;		// result of HDM_GETBITMAPMARGIN
};

static void drawCell(struct table *t, HDC dc, struct drawCellParams *p)
{
	RECT r;
	WCHAR *text;
	HBRUSH background;
	int textColor;
	POINT pt;
	int cbState;
	RECT cellrect;
	HDC idc;
	HBITMAP previbitmap;
	BLENDFUNCTION bf;

	// TODO verify these two
	background = (HBRUSH) (COLOR_WINDOW + 1);
	textColor = COLOR_WINDOWTEXT;
	if (t->selectedRow == p->row) {
		// these are the colors wine uses (http://source.winehq.org/source/dlls/comctl32/listview.c)
		// the two for unfocused are also suggested by http://stackoverflow.com/questions/10428710/windows-forms-inactive-highlight-color
		background = (HBRUSH) (COLOR_HIGHLIGHT + 1);
		textColor = COLOR_HIGHLIGHTTEXT;
		if (GetFocus() != t->hwnd) {
			background = (HBRUSH) (COLOR_BTNFACE + 1);
			textColor = COLOR_BTNTEXT;
		}
		// TODO disabled
	}

	r.left = p->x;
	r.right = p->x + p->width;
	r.top = p->y;
	r.bottom = p->y + p->height;
	if (FillRect(dc, &r, background) == 0)
		panic("error filling Table cell background");
	cellrect = r;		// save for drawing the focus rect

	switch (t->columnTypes[p->column]) {
	case tableColumnText:
		toCellContentRect(t, &r, p->xoff, 0, 0);		// TODO get the text height
		if (SetTextColor(dc, GetSysColor(textColor)) == CLR_INVALID)
			panic("error setting Table cell text color");
		if (SetBkMode(dc, TRANSPARENT) == 0)
			panic("error setting transparent text drawing mode for Table cell");
		text = (WCHAR *) notify(t, tableNotificationGetCellData, p->row, p->column, 0);
		if (DrawTextExW(dc, text, -1, &r, DT_END_ELLIPSIS | DT_LEFT | DT_NOPREFIX | DT_SINGLELINE, NULL) == 0)
			panic("error drawing Table cell text");
		notify(t, tableNotificationFinishedWithCellData, p->row, p->column, (uintptr_t) text);
		break;
	case tableColumnImage:
		toCellContentRect(t, &r, p->xoff, tableImageWidth(), tableImageHeight());
		idc = CreateCompatibleDC(dc);
		if (idc == NULL)
			panic("error creating compatible DC for Table image cell drawing");
		previbitmap = SelectObject(idc, testbitmap);
		if (previbitmap == NULL)
			panic("error selecting Table cell image into compatible DC for image drawing");
		ZeroMemory(&bf, sizeof (BLENDFUNCTION));
		bf.BlendOp = AC_SRC_OVER;
		bf.BlendFlags = 0;
		bf.SourceConstantAlpha = 255;			// per-pixel alpha values
		bf.AlphaFormat = AC_SRC_ALPHA;
		// TODO 16 and 16 are the width and height of the image; we would need to get that out somehow
		if (AlphaBlend(dc, r.left, r.top, r.right - r.left, r.bottom - r.top,
			idc, 0, 0, 16, 16, bf) == FALSE)
			panic("error drawing image into Table cell");
		if (SelectObject(idc, previbitmap) != testbitmap)
			panic("error deselecting Table cell image for drawing image");
		if (DeleteDC(idc) == 0)
			panic("error deleting Table compatible DC for image cell drawing");
		break;
	case tableColumnCheckbox:
		toCheckboxRect(t, &r, p->xoff);
		cbState = 0;
		if (p->row == lastCheckbox.row && p->column == lastCheckbox.column)
			cbState |= checkboxStateChecked;
		if (t->checkboxMouseDown)
			if (p->row == t->checkboxMouseDownRow && p->column == t->checkboxMouseDownColumn)
				cbState |= checkboxStatePushed;
		if (t->checkboxMouseOverLast) {
			pt.x = GET_X_LPARAM(t->checkboxMouseOverLastPoint);
			pt.y = GET_Y_LPARAM(t->checkboxMouseOverLastPoint);
			if (PtInRect(&r, pt) != 0)
				cbState |= checkboxStateHot;
		}
		drawCheckbox(t, dc, &r, cbState);
		break;
	}

	// TODO in front of or behind the cell contents?
	if (t->selectedRow == p->row && t->selectedColumn == p->column)
		if (DrawFocusRect(dc, &cellrect) == 0)
			panic("error drawing focus rect on current Table cell");
}

// TODO use cliprect
static void draw(struct table *t, HDC dc, RECT cliprect, RECT client)
{
	intptr_t i, j;
	int x = 0;
	HFONT prevfont, newfont;
	struct drawCellParams p;

	prevfont = selectFont(t, dc, &newfont);

	client.top += t->headerHeight;

	ZeroMemory(&p, sizeof (struct drawCellParams));
	p.height = rowHeight(t, dc, FALSE);
	p.xoff = SendMessageW(t->header, HDM_GETBITMAPMARGIN, 0, 0);

	p.y = client.top;
	for (i = t->vscrollpos; i < t->count; i++) {
		p.row = i;
		p.x = client.left - t->hscrollpos;
		for (j = 0; j < t->nColumns; j++) {
			p.column = j;
			p.width = columnWidth(t, p.column);
			drawCell(t, dc, &p);
			p.x += p.width;
		}
		p.y += p.height;
		if (p.y >= client.bottom)			// >= because RECT.bottom is the first pixel outside the rect
			break;
	}

	deselectFont(dc, prevfont, newfont);
}

HANDLER(drawHandlers)
{
	HDC dc;
	PAINTSTRUCT ps;
	RECT client;
	RECT r;

	if (uMsg != WM_PAINT && uMsg != WM_PRINTCLIENT)
		return FALSE;
	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting client rect for Table painting");
	if (uMsg == WM_PAINT) {
		dc = BeginPaint(t->hwnd, &ps);
		if (dc == NULL)
			panic("error beginning Table painting");
		r = ps.rcPaint;
	} else {
		dc = (HDC) wParam;
		r = client;
	}
	draw(t, dc, r, client);
	if (uMsg == WM_PAINT)
		EndPaint(t->hwnd, &ps);
	// this is correct for WM_PRINTCLIENT; see http://stackoverflow.com/a/27362258/3408572
	*lResult = 0;
	return TRUE;
}

// TODO redraw selected row on focus change
// TODO here or in select.h?
