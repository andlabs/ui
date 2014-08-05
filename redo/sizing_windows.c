/* 17 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

/* TODO figure out where these should go */

void calculateBaseUnits(HWND hwnd, int *baseX, int *baseY)
{
	HDC dc;
	HFONT prevFont;
	TEXTMETRICW tm;

	dc = GetDC(hwnd);
	if (dc == NULL)
		xpanic("error getting DC for preferred size calculations", GetLastError());
	prevFont = (HFONT) SelectObject(dc, controlFont);
	if (prevFont == NULL)
		xpanic("error loading control font into device context for preferred size calculation", GetLastError());
	if (GetTextMetricsW(dc, &tm) == 0)
		xpanic("error getting text metrics for preferred size calculations", GetLastError());
	*baseX = (int) tm.tmAveCharWidth;		/* TODO not optimal; third reference below has better way */
	*baseY = (int) tm.tmHeight;
	if (SelectObject(dc, prevFont) != controlFont)
		xpanic("error restoring previous font into device context after preferred size calculations", GetLastError());
	if (ReleaseDC(hwnd, dc) == 0)
		xpanic("error releasing DC for preferred size calculations", GetLastError());
}

void moveWindow(HWND hwnd, int x, int y, int width, int height)
{
	if (MoveWindow(hwnd, x, y, width, height, TRUE) == 0)
		xpanic("error setting window/control rect", GetLastError());
}

LONG controlTextLength(HWND hwnd, LPWSTR text)
{
	HDC dc;
	HFONT prev;
	SIZE size;

	dc = GetDC(hwnd);
	if (dc == NULL)
		xpanic("error getting DC of control for text length", GetLastError());
	prev = SelectObject(dc, controlFont);
	if (prev == NULL)
		xpanic("error setting control font to DC for text length", GetLastError());
	if (GetTextExtentPoint32W(dc, text, wcslen(text), &size) == 0)
		xpanic("error actually getting text length", GetLastError());
	if (SelectObject(dc, prev) != controlFont)
		xpanic("error restoring previous control font to DC for text length", GetLastError());
	if (ReleaseDC(hwnd, dc) == 0)
		xpanic("error releasing DC of control for text length", GetLastError());
	return size.cx;
}
