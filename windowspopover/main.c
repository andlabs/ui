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
// - should the parent window appear deactivated?

HWND popover;

void xpanic(char *msg, DWORD err)
{
	printf("%d | %s\n", err, msg);
	abort();
}

#define ARROWHEIGHT 8
#define ARROWWIDTH 8		/* should be the same for smooth lines */

struct popover {
	void *gopopover;

	// a nice consequence of this design is that it allows four arrowheads to jut out at once; in practice only one will ever be used, but hey — simple implementation!
	LONG arrowLeft;
	LONG arrowRight;
	LONG arrowTop;
	LONG arrowBottom;
};

struct popover _p = { NULL, -1, -1, 20, -1 };
struct popover *p = &_p;

HRGN makePopoverRegion(HDC dc, LONG width, LONG height)
{
	POINT pt[20];
	int n;
	HRGN region;
	LONG xmax, ymax;

	if (BeginPath(dc) == 0)
		xpanic("error beginning path for Popover shape", GetLastError());
	n = 0;

	// figure out the xmax and ymax of the box
	xmax = width;
	if (p->arrowRight >= 0)
		xmax -= ARROWWIDTH;
	ymax = height;
	if (p->arrowBottom >= 0)
		ymax -= ARROWHEIGHT;

	// the first point is either at (0,0), (0,arrowHeight), (arrowWidth,0), or (arrowWidth,arrowHeight)
	pt[n].x = 0;
	if (p->arrowLeft >= 0)
		pt[n].x = ARROWWIDTH;
	pt[n].y = 0;
	if (p->arrowTop >= 0)
		pt[n].y = ARROWHEIGHT;
	n++;

	// the left side
	pt[n].x = pt[n - 1].x;
	if (p->arrowLeft >= 0) {
		pt[n].y = pt[n - 1].y + p->arrowLeft;
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x;
	}
	pt[n].y = ymax;
	n++;

	// the bottom side
	pt[n].y = pt[n - 1].y;
	if (p->arrowBottom >= 0) {
		pt[n].x = pt[n - 1].x + p->arrowBottom;
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].y = pt[n - 1].y;
	}
	pt[n].x = xmax;
	n++;

	// the right side
	pt[n].x = pt[n - 1].x;
	if (p->arrowRight >= 0) {
		pt[n].y = pt[0].y + p->arrowRight + (ARROWHEIGHT * 2);
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x;
	}
	pt[n].y = pt[0].y;
	n++;

	// the top side
	pt[n].y = pt[n - 1].y;
	if (p->arrowTop >= 0) {
		pt[n].x = pt[0].x + p->arrowTop + (ARROWWIDTH * 2);
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].y = pt[n - 1].y;
	}
	pt[n].x = pt[0].x;
	n++;

	if (Polyline(dc, pt, n) == 0)
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
		// TODO isolate the brush name to a constant
		// unfortunately FillRgn() doesn't document the COLOR+1 trick as working there
		brush = GetSysColorBrush(COLOR_BTNFACE);
		if (brush == NULL)
			xpanic("error getting Popover background brush", GetLastError());
		if (FillRgn(dc, region, brush) == 0)
			xpanic("error drawing Popover background", GetLastError());
		// TODO use a system color brush?
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
		FrameRect(dc, &r, GetStockPen(WHITE_BRUSH));
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
