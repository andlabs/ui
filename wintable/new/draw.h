// 8 december 2014

static void draw(struct table *t, HDC dc, RECT cliprect, RECT client)
{
	LRESULT i, n;
	RECT r;
	int x = 0;

	n = SendMessageW(t->header, HDM_GETITEMCOUNT, 0, 0);
	for (i = 0; i < n; i++) {
		SendMessage(t->header, HDM_GETITEMRECT, (WPARAM) i, (LPARAM) (&r));
		r.top = client.top;
		r.bottom = client.bottom;
		FillRect(dc, &r, GetSysColorBrush(x));
		x++;
	}
}

// TODO handle WM_PRINTCLIENT flags?

HANDLER(drawHandlers)
{
	HDC dc;
	PAINTSTRUCT ps;
	RECT client;
	RECT r;
	BOOL wmpaint;

	if (uMsg != WM_PAINT && uMsg != WM_PRINTCLIENT)
		return FALSE;
	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting client rect for Table painting");
	// let's be nice: some controls don't support WM_PRINTCLIENT but do allow you to pass a HDC as the WPARAM to WM_PAINT, so let's support that too as an option
	// TODO find out how susch controls handle LPARAM
	wmpaint = uMsg == WM_PAINT && ((HDC) wParam) == NULL;
	if (wmpaint) {
		dc = BeginPaint(t->hwnd, &ps);
		if (dc == NULL)
			panic("error beginning Table painting");
		r = ps.rcPaint;
	} else {
		dc = (HDC) wParam;
		r = client;
	}
	draw(t, dc, r, client);
	if (wmpaint)
		EndPaint(t->hwnd, &ps);
	// TODO is this correct for WM_PRINTCLIENT? MSDN doesn't say
	*lResult = 0;
	return TRUE;
}
