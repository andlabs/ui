/* 17 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

BOOL baseUnitsCalculated = FALSE;
int baseX;
int baseY;

/* called by newWindow() so we can calculate base units when we have a window */
void calculateBaseUnits(HWND hwnd)
{
	HDC dc;
	HFONT prevFont;
	TEXTMETRICW tm;

	if (baseUnitsCalculated)
		return;
	dc = GetDC(hwnd);
	if (dc == NULL)
		xpanic("error getting DC for preferred size calculations", GetLastError());
	prevFont = (HFONT) SelectObject(dc, controlFont);
	if (prevFont == NULL)
		xpanic("error loading control font into device context for preferred size calculation", GetLastError());
	if (GetTextMetricsW(dc, &tm) == 0)
		xpanic("error getting text metrics for preferred size calculations", GetLastError());
	baseX = (int) tm.tmAveCharWidth;		/* TODO not optimal; third reference below has better way */
	baseY = (int) tm.tmHeight;
	if (SelectObject(dc, prevFont) != controlFont)
		xpanic("error restoring previous font into device context after preferred size calculations", GetLastError());
	if (ReleaseDC(hwnd, dc) == 0)
		xpanic("error releasing DC for preferred size calculations", GetLastError());
	baseUnitsCalculated = TRUE;
}

void moveWindow(HWND hwnd, int x, int y, int width, int height)
{
	if (MoveWindow(hwnd, x, y, width, height, TRUE) == 0)
		xpanic("error setting window/control rect", GetLastError());
}
