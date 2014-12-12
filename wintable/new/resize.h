// 7 december 2014

// TODO why doesn't this trigger on first show?
// TODO last few bits of the scrollbar series, that talk about WM_WINDOWPOSCHANGING and the metaphor
// TODO rename this to boot

HANDLER(resizeHandler)
{
	WINDOWPOS *wp;
	RECT client;
	intptr_t height;

	if (uMsg != WM_WINDOWPOSCHANGED)
		return FALSE;
	wp = (WINDOWPOS *) lParam;
	if ((wp->flags & SWP_NOSIZE) != 0)
		return FALSE;

	// TODO does wp store the window rect or the client rect?
	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting Table client rect in resizeHandler()");
	// TODO do this before calling updateTableWidth() (which calls repositionHeader()?)?
	client.top -= t->headerHeight;

	// update the width...
	// this will call repositionHeader(); there's a good reason... (see comments)
	// TODO when I clean that mess up, remove this comment
	updateTableWidth(t);

	// ...and the height
	// TODO find out if order matters
	height = client.bottom - client.top;
	t->vpagesize = height / rowht(t);
	// do a dummy scroll to reflect those changes
	vscrollby(t, 0);

	*lResult = 0;
	return TRUE;
}
