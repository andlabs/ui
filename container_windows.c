// 17 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

RECT containerBounds(HWND hwnd)
{
	RECT r;

	if (GetClientRect(hwnd, &r) == 0)
		xpanic("error getting container client rect for container.bounds()", GetLastError());
	return r;
}

void calculateBaseUnits(HWND hwnd, int *baseX, int *baseY, LONG *internalLeading)
{
	HDC dc;
	HFONT prevFont;
	TEXTMETRICW tm;
	SIZE size;

	dc = GetDC(hwnd);
	if (dc == NULL)
		xpanic("error getting DC for preferred size calculations", GetLastError());
	prevFont = (HFONT) SelectObject(dc, controlFont);
	if (prevFont == NULL)
		xpanic("error loading control font into device context for preferred size calculation", GetLastError());
	if (GetTextMetricsW(dc, &tm) == 0)
		xpanic("error getting text metrics for preferred size calculations", GetLastError());
	if (GetTextExtentPoint32W(dc, L"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 52, &size) == 0)
		xpanic("error getting text extent point for preferred size calculations", GetLastError());
	*baseX = (int) ((size.cx / 26 + 1) / 2);
	*baseY = (int) tm.tmHeight;
	*internalLeading = tm.tmInternalLeading;
	if (SelectObject(dc, prevFont) != controlFont)
		xpanic("error restoring previous font into device context after preferred size calculations", GetLastError());
	if (ReleaseDC(hwnd, dc) == 0)
		xpanic("error releasing DC for preferred size calculations", GetLastError());
}
