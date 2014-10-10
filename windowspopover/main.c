// 9 october 2014
#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
// get Windows version right; right now Windows XP
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <commctrl.h>
#include <stdint.h>
#include <uxtheme.h>
#include <string.h>
#include <wchar.h>
#include <windowsx.h>
#include <vsstyle.h>
#include <vssym32.h>

// #qo LIBS: user32 kernel32 gdi32

// TODO
// - investigate visual styles
// - put the client and non-client areas in the right place
// - make sure redrawing is correct (especially for backgrounds)
// - wine: BLACK_PEN draws a white line? (might change later so eh)
// - should the parent window appear deactivated?

HWND popover;

void xpanic(char *msg, DWORD err)
{
	printf("%d | %s\n", err, msg);
	abort();
}

#define ARROWHEIGHT 8
#define ARROWWIDTH 8		/* should be the same for smooth lines */

HRGN makePopoverRegion(HDC dc, LONG width, LONG height)
{
	POINT pt[8];
	HRGN region;

	if (BeginPath(dc) == 0)
		xpanic("error beginning path for Popover shape", GetLastError());
	pt[0].x = 0;
	pt[0].y = ARROWHEIGHT;
	pt[1].x = pt[0].x;
	pt[1].y = pt[0].y + (height - ARROWHEIGHT);
	pt[2].x = pt[1].x + width;
	pt[2].y = pt[1].y;
	pt[3].x = pt[2].x;
	pt[3].y = pt[2].y - (height - ARROWHEIGHT);
	pt[4].x = pt[3].x - ((width / 2) - ARROWWIDTH);
	pt[4].y = pt[3].y;
	pt[5].x = pt[4].x - ARROWWIDTH;
	pt[5].y = pt[4].y - ARROWHEIGHT;
	pt[6].x = pt[5].x - ARROWWIDTH;
	pt[6].y = pt[5].y + ARROWHEIGHT;
	pt[7].x = pt[0].x;
	pt[7].y = pt[0].y;
	if (Polyline(dc, pt, 8) == 0)
		xpanic("error drawing lines in Popover shape", GetLastError());
	if (EndPath(dc) == 0)
		xpanic("error ending path for Popover shape", GetLastError());
	region = PathToRegion(dc);
	if (region == NULL)
		xpanic("error converting Popover shape path to region", GetLastError());
	return region;
}

LRESULT CALLBACK popoverproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	PAINTSTRUCT ps;
	HDC dc;
	HRGN region;
	RECT r;
	LONG width;
	LONG height;
	WINDOWPOS *wp;

	switch (uMsg) {
	case WM_NCPAINT:
		if (GetWindowRect(hwnd, &r) == 0)
			xpanic("error getting Popover window rect for shape redraw", GetLastError());
		width = r.right - r.left;
		height = r.bottom - r.top;
		dc = GetWindowDC(hwnd);
		if (dc == NULL)
			xpanic("error getting Popover window DC for drawing border", GetLastError());
		region = makePopoverRegion(dc, width, height);
		// TODO isolate the class brush to a constant
		// TODO object errors (or switch to path functions?)
		if (FillRgn(dc, region, GetSysColorBrush(COLOR_BTNFACE)) == 0)
			xpanic("error drawing Popover background", GetLastError());
		if (FrameRgn(dc, region, GetStockObject(BLACK_PEN), 1, 1) == 0)
			xpanic("error drawing Popover border", GetLastError());
		if (DeleteObject(region) == 0)
			xpanic("error deleting Popover shape region", GetLastError());
		if (ReleaseDC(hwnd, dc) == 0)
			xpanic("error releasing Popover window DC for shape drawing", GetLastError());
		return 0;
	case WM_WINDOWPOSCHANGED:
		// this must be here; if it's in WM_NCPAINT weird things happen (see http://stackoverflow.com/questions/26288303/why-is-my-client-rectangle-drawing-behaving-bizarrely-pictures-provided-if-i-t)
		wp = (WINDOWPOS *) lParam;
		if ((wp->flags & SWP_NOSIZE) == 0) {
			dc = GetWindowDC(hwnd);
			if (dc == NULL)
				xpanic("error getting Popover window DC for reshaping", GetLastError());
			region = makePopoverRegion(dc, wp->cx, wp->cy);
			if (SetWindowRgn(hwnd, region, TRUE) == 0)
				xpanic("error setting Popover shape", GetLastError());
			// don't delete the region; the window manager owns it now
			if (ReleaseDC(hwnd, dc) == 0)
				xpanic("error releasing Popover window DC for reshaping", GetLastError());
		}
		break;		// defer to DefWindowProc()
	case WM_NCCALCSIZE:
		{
			RECT *r = (RECT *) lParam;
			NCCALCSIZE_PARAMS *np = (NCCALCSIZE_PARAMS *) lParam;

			if (wParam != FALSE)
				r = &np->rgrc[0];
			r->left++;
			r->top++;
			r->right--;
			r->bottom--;
			r->top += ARROWHEIGHT;
			return 0;
		}
	case WM_PAINT:
		dc = BeginPaint(hwnd, &ps);
		GetClientRect(hwnd, &r);
		FillRect(dc, &r, GetSysColorBrush(COLOR_ACTIVECAPTION));
		FrameRect(dc, &r, GetStockPen(WHITE_PEN));
		EndPaint(hwnd, &ps);
		return 0;
	}
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

LRESULT CALLBACK wndproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	switch (uMsg) {
	case WM_COMMAND:
		if (HIWORD(wParam) == BN_CLICKED && LOWORD(wParam) == 100) {
			MoveWindow(popover, 50, 50,  200, 200, TRUE);
			ShowWindow(popover, SW_SHOW);
			UpdateWindow(popover);
			return 0;
		}
		break;
	case WM_CLOSE:
		PostQuitMessage(0);
		return 0;
	}
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

int main(int argc, char *argv[])
{
	WNDCLASSW wc;
	HWND mainwin, button;
	MSG msg;

	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpszClassName = L"popover";
	wc.lpfnWndProc = popoverproc;
	wc.hbrBackground = (HBRUSH) (COLOR_BTNFACE + 1);
	wc.style = CS_DROPSHADOW | CS_NOCLOSE;
	if (RegisterClassW(&wc) == 0)
		abort();
	popover = CreateWindowExW(WS_EX_TOPMOST,
		L"popover", L"",
		WS_POPUP,
		0, 0, 150, 100,
		NULL, NULL, NULL, NULL);
	if (popover == NULL)
		abort();

	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpszClassName = L"mainwin";
	wc.lpfnWndProc = wndproc;
	wc.hbrBackground = (HBRUSH) (COLOR_BTNFACE + 1);
	if (RegisterClassW(&wc) == 0)
		abort();
	mainwin = CreateWindowExW(0,
		L"mainwin", L"Main Window",
		WS_OVERLAPPEDWINDOW,
		0, 0, 150, 100,
		NULL, NULL, NULL, NULL);
	if (mainwin == NULL)
		abort();
	button = CreateWindowExW(0,
		L"button", L"Click Me",
		BS_PUSHBUTTON | WS_CHILD | WS_VISIBLE,
		20, 20, 100, 40,
		mainwin, (HMENU) 100, NULL, NULL);
	if (button == NULL)
		abort();
	ShowWindow(mainwin, SW_SHOWDEFAULT);
	if (UpdateWindow(mainwin) == 0)
		abort();
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
	return 0;
}
