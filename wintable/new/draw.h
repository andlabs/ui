// 8 december 2014

static void draw(struct table *t, HDC dc, RECT cliprect, RECT client)
{
	LRESULT i, n;
	RECT r;
	int x = 0;
	HFONT prevfont, newfont;

	n = SendMessageW(t->header, HDM_GETITEMCOUNT, 0, 0);
	for (i = 0; i < n; i++) {
		SendMessage(t->header, HDM_GETITEMRECT, (WPARAM) i, (LPARAM) (&r));
		r.top = client.top;
		r.bottom = client.bottom;
		FillRect(dc, &r, GetSysColorBrush(x));
		x++;
	}

	prevfont = selectFont(t, dc, &newfont);
	TextOutW(dc, 100, 100, L"come on", 7);
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
