// 8 december 2014

HANDLER(apiHandlers)
{
	switch (uMsg) {
	case WM_SETFONT:
		// TODO release old font?
		t->font = (HFONT) wParam;
		SendMessageW(t->header, WM_SETFONT, wParam, lParam);
		// TODO reposition header?
		// TODO how to properly handle LOWORD(lParam) != FALSE?
		*lResult = 0;
		return TRUE;
	case WM_GETFONT:
		*lResult = (LRESULT) (t->font);
		return TRUE;
	case tableAddColumn:
		// TODO
		return FALSE;
	}
	return FALSE;
}
