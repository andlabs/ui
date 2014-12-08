// 7 december 2014

// TODO why doesn't this trigger on first show?

HANDLER(resizeHandler)
{
	WINDOWPOS *wp;

	if (uMsg != WM_WINDOWPOSCHANGED)
		return FALSE;
	wp = (WINDOWPOS *) lParam;
	if ((wp->flags & SWP_NOSIZE) != 0)
		return FALSE;
	repositionHeader(t);
	*lResult = 0;
	return TRUE;
}
