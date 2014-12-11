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
	WCHAR msg[200];
	int n;

	r.left = p->x + p->xoff;
	r.right = p->x + p->width;
	r.top = p->y;
	r.bottom = p->y + p->height;
	// TODO fill this rect with the appropriate background color
	// TODO then vertical center content
	n = wsprintf(msg, L"(%d,%d)", p->row, p->column);
	if (DrawTextExW(dc, msg, n, &r, DT_END_ELLIPSIS | DT_LEFT | DT_NOPREFIX | DT_SINGLELINE, NULL) == 0)
		panic("error drawing Table cell text");
}

static void draw(struct table *t, HDC dc, RECT cliprect, RECT client)
{
	intptr_t i;
	RECT r;
	int x = 0;
	HFONT prevfont, newfont;
	struct drawCellParams p;

	for (i = 0; i < t->nColumns; i++) {
		SendMessage(t->header, HDM_GETITEMRECT, (WPARAM) i, (LPARAM) (&r));
		r.left -= t->hscrollpos;
		r.right -= t->hscrollpos;
		r.top = client.top;
		r.bottom = client.bottom;
		FillRect(dc, &r, GetSysColorBrush(x));
		x++;
	}

	prevfont = selectFont(t, dc, &newfont);
	ZeroMemory(&p, sizeof (struct drawCellParams));
	p.row = 0;
	p.column = 0;
	p.x = r.left - t->hscrollpos;
	p.y = 100;
	p.width = r.right - r.left;
	p.height = rowHeight(t, dc, FALSE);
	p.xoff = SendMessageW(t->header, HDM_GETBITMAPMARGIN, 0, 0);
	drawCell(t, dc, &p);
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
