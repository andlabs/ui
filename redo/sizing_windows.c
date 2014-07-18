/* 17 july 2014 */

#include "winapi_windows.h"

HDC getDC(HWND hwnd)
{
	HDC dc;

	dc = GetDC(hwnd);
	if (dc == NULL)
		xpanic("error getting DC for preferred size calculations", GetLastError());
/* TODO */
	/* TODO save for restoring later */
/*
	if (SelectObject(dc, controlFont) == NULL)
		xpanic("error loading control font into device context for preferred size calculation", GetLastError());
*/
	return dc;
}

void releaseDC(HWND hwnd, HDC dc)
{
	if (ReleaseDC(hwnd, dc) == 0)
		xpanic("error releasing DC for preferred size calculations", GetLastError());
}

void getTextMetricsW(HDC dc, TEXTMETRICW *tm)
{
	if (GetTextMetricsW(dc, tm) == 0)
		xpanic("error getting text metrics for preferred size calculations", GetLastError());
}

void moveWindow(HWND hwnd, int x, int y, int width, int height)
{
	if (MoveWindow(hwnd, x, y, width, height, TRUE) == 0)
		xpanic("error setting window/control rect", GetLastError());
}
