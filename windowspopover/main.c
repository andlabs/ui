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

HWND popover;

// #qo LIBS: user32 kernel32 gdi32

#define ARROWHEIGHT 6
#define ARROWWIDTH 8

LRESULT CALLBACK popoverproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	PAINTSTRUCT ps;
	HDC dc;
	HRGN region;
	POINT pt;

	switch (uMsg) {
	case WM_PAINT:
		dc = BeginPaint(hwnd, &ps);
		if (dc == NULL) abort();
		BeginPath(dc);
		pt.x = 0;
		pt.y = ARROWHEIGHT;
		if (MoveToEx(dc, pt.x, pt.y, NULL) == 0) abort();
		pt.y = 100;
		if (LineTo(dc, pt.x, pt.y) == 0) abort();
		pt.x = 100;
		LineTo(dc, pt.x, pt.y);
		pt.y = ARROWHEIGHT;
		LineTo(dc, pt.x, pt.y);
		pt.x = 50 + ARROWWIDTH;
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
		region = PathToRegion(dc);
		FrameRgn(dc, region, GetStockObject(BLACK_PEN), 1, 1);
		SetWindowRgn(hwnd, region, TRUE);
		EndPaint(hwnd, &ps);
		return 0;
	case WM_NCCALCSIZE:
		break;
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
