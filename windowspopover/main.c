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

#define ARROWHEIGHT 8
#define ARROWWIDTH 8		/* should be the same for smooth lines */

LRESULT CALLBACK popoverproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	PAINTSTRUCT ps;
	HDC dc;
	HRGN region;
	POINT pt;
	RECT r;
	LONG width;
	LONG height;

	switch (uMsg) {
	case WM_NCPAINT:
		GetWindowRect(hwnd, &r);
		width = r.right - r.left;
		height = r.bottom - r.top;
		dc = GetDCEx(hwnd, (HRGN) wParam, DCX_WINDOW | DCX_INTERSECTRGN);
		if (dc == NULL) abort();
		BeginPath(dc);
		r.left = 0; r.top = 0;		// everything's in device coordinates
		pt.x = r.left;
		pt.y = r.top + ARROWHEIGHT;
		if (MoveToEx(dc, pt.x, pt.y, NULL) == 0) abort();
		pt.y += height - ARROWHEIGHT;
		if (LineTo(dc, pt.x, pt.y) == 0) abort();
		pt.x += width;
		LineTo(dc, pt.x, pt.y);
		pt.y -= height - ARROWHEIGHT;
		LineTo(dc, pt.x, pt.y);
		pt.x -= (width / 2) - ARROWWIDTH;
		LineTo(dc, pt.x, pt.y);
		pt.x -= ARROWWIDTH;
		pt.y -= ARROWHEIGHT;
		LineTo(dc, pt.x, pt.y);
		pt.x -= ARROWWIDTH;
		pt.y += ARROWHEIGHT;
		LineTo(dc, pt.x, pt.y);
		pt.x = 0;
		LineTo(dc, pt.x, pt.y);
		EndPath(dc);
		SetDCBrushColor(dc, RGB(255, 0, 0));
		region = PathToRegion(dc);
		FrameRgn(dc, region, GetStockObject(DC_BRUSH), 1, 1);
		SetWindowRgn(hwnd, region, TRUE);
		ReleaseDC(hwnd, dc);
		return 0;
	case WM_NCCALCSIZE:
		break;
	case WM_ERASEBKGND:
		return (LRESULT) NULL;
	case WM_PAINT:
/*		dc = BeginPaint(hwnd, &ps);
		GetClientRect(hwnd, &r);
		FillRect(dc, &r, GetSysColorBrush(COLOR_ACTIVECAPTION));
		EndPaint(hwnd, &ps);
*/		return 0;
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
