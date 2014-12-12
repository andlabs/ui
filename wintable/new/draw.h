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

//TODO delete when done
static int current = 0;

static void drawCell(struct table *t, HDC dc, struct drawCellParams *p)
{
	RECT r;
	WCHAR msg[200];
	int n;

	r.left = p->x;//TODO + p->xoff;
	r.right = p->x + p->width;
	r.top = p->y;
	r.bottom = p->y + p->height;
	// TODO fill this rect with the appropriate background color
	// TODO then vertical center content
	n = wsprintf(msg, L"(%d,%d)", p->row, p->column);

	FillRect(dc, &r, (HBRUSH) (current + 1));
	current++;
	if (current >= 31)
		current = 0;

	r.left += p->xoff;
	if (DrawTextExW(dc, msg, n, &r, DT_END_ELLIPSIS | DT_LEFT | DT_NOPREFIX | DT_SINGLELINE, NULL) == 0)
		panic("error drawing Table cell text");
}

static void draw(struct table *t, HDC dc, RECT cliprect, RECT client)
{
	intptr_t i, j;
	RECT r;
	int x = 0;
	HFONT prevfont, newfont;
	struct drawCellParams p;

	prevfont = selectFont(t, dc, &newfont);

current = 0;

	client.top += t->headerHeight;

	ZeroMemory(&p, sizeof (struct drawCellParams));
	p.height = rowHeight(t, dc, FALSE);
	p.xoff = SendMessageW(t->header, HDM_GETBITMAPMARGIN, 0, 0);

	p.y = client.top;
	for (i = 0; i < t->count; i++) {
		p.row = i;
		p.x = client.left - t->hscrollpos;
		for (j = 0; j < t->nColumns; j++) {
			p.column = j;
			// TODO error check
			SendMessage(t->header, HDM_GETITEMRECT, (WPARAM) j, (LPARAM) (&r));
			p.width = r.right - r.left;
			drawCell(t, dc, &p);
			p.x += p.width;
		}
		p.y += p.height;
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
