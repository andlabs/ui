// 9 october 2014
#include "../wininclude_windows.h"
#include "popover.h"

// #qo LIBS: user32 kernel32 gdi32

// TODO
// - should the parent window appear deactivated?

HWND popoverWindow;

void xpanic(char *msg, DWORD err)
{
	printf("%d | %s\n", err, msg);
	abort();
}

popover *p;

HRGN makePopoverRegion(HDC dc, LONG width, LONG height)
{
	popoverPoint ppt[20];
	POINT pt[20];
	int i, n;
	HRGN region;

	n = popoverMakeFramePoints(p, (intptr_t) width, (intptr_t) height, ppt);
	for (i = 0; i < n; i++) {
		pt[i].x = (LONG) (ppt[i].x);
		pt[i].y = (LONG) (ppt[i].y);
	}

	if (BeginPath(dc) == 0)
		xpanic("error beginning path for Popover shape", GetLastError());
	if (Polyline(dc, pt, n) == 0)
		xpanic("error drawing lines in Popover shape", GetLastError());
	if (EndPath(dc) == 0)
		xpanic("error ending path for Popover shape", GetLastError());
	region = PathToRegion(dc);
	if (region == NULL)
		xpanic("error converting Popover shape path to region", GetLastError());
	return region;
}

#define msgPopoverPrepareLeftRight (WM_APP+50)
#define msgPopoverPrepareTopBottom (WM_APP+51)

LRESULT CALLBACK popoverproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	PAINTSTRUCT ps;
	HDC dc;
	HRGN region;
	RECT r;
	LONG width;
	LONG height;
	WINDOWPOS *wp;
	HBRUSH brush;

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
		// don't call FillRgn(); WM_ERASEBKGND seems to do this to the non-client area for us already :S (TODO confirm)
		// TODO arrow is black in wine
		brush = (HBRUSH) GetStockObject(BLACK_BRUSH);
		if (brush == NULL)
			xpanic("error getting Popover border brush", GetLastError());
		if (FrameRgn(dc, region, brush, 1, 1) == 0)
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
			popoverRect pr;

			if (wParam != FALSE)
				r = &np->rgrc[0];
			pr.left = (intptr_t) (r->left);
			pr.top = (intptr_t) (r->top);
			pr.right = (intptr_t) (r->right);
			pr.bottom = (intptr_t) (r->bottom);
			popoverWindowSizeToClientSize(p, &pr);
			r->left = (LONG) (pr.left);
			r->top = (LONG) (pr.top);
			r->right = (LONG) (pr.right);
			r->bottom = (LONG) (pr.bottom);
			return 0;
		}
	case WM_PAINT:
		dc = BeginPaint(hwnd, &ps);
		GetClientRect(hwnd, &r);
		FillRect(dc, &r, GetSysColorBrush(COLOR_ACTIVECAPTION));
		FrameRect(dc, &r, GetStockPen(WHITE_BRUSH));
		EndPaint(hwnd, &ps);
		return 0;
	case msgPopoverPrepareLeftRight:
	case msgPopoverPrepareTopBottom:
		// TODO window edge detection
		{
			RECT r;
			LONG width = 200, height = 200;
			popoverRect control;
			uintptr_t side;
			popoverRect out;

			if (GetWindowRect((HWND) wParam, &r) == 0)
				xpanic("error getting window rect of Popover target", GetLastError());
			control.left = (intptr_t) (r.left);
			control.top = (intptr_t) (r.top);
			control.right = (intptr_t) (r.right);
			control.bottom = (intptr_t) (r.bottom);
			switch (uMsg) {
			case msgPopoverPrepareLeftRight:
				side = popoverPointLeft;
				break;
			case msgPopoverPrepareTopBottom:
				side = popoverPointTop;
				break;
			}
			out = popoverPointAt(p, control, (intptr_t) width, (intptr_t) height, side);
			if (MoveWindow(hwnd, out.left, out.top, out.right - out.left, out.bottom - out.top, TRUE) == 0)
				xpanic("error repositioning Popover", GetLastError());
		}
		return 0;
	}
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

HWND button;

LRESULT CALLBACK wndproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	switch (uMsg) {
	case WM_COMMAND:
		if (HIWORD(wParam) == BN_CLICKED && LOWORD(wParam) == 100) {
			SendMessageW(popoverWindow, msgPopoverPrepareLeftRight, (WPARAM) button, 0);
			ShowWindow(popoverWindow, SW_SHOW);
			UpdateWindow(popoverWindow);
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
	HWND mainwin;
	MSG msg;

	p = popoverDataNew(NULL);
	// TODO null check

	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpszClassName = L"popover";
	wc.lpfnWndProc = popoverproc;
	wc.hbrBackground = (HBRUSH) (COLOR_BTNFACE + 1);
	wc.style = CS_DROPSHADOW | CS_NOCLOSE;
	if (RegisterClassW(&wc) == 0)
		abort();
	popoverWindow = CreateWindowExW(WS_EX_TOPMOST,
		L"popover", L"",
		WS_POPUP,
		0, 0, 150, 100,
		NULL, NULL, NULL, NULL);
	if (popoverWindow == NULL)
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
