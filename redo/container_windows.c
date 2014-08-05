/* 17 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

/*
This could all just be part of Window, but doing so just makes things complex.
In this case, I chose to waste a window handle rather than keep things super complex.
If this is seriously an issue in the future, I can roll it back.
*/

#define containerclass L"gouicontainer"

static LRESULT CALLBACK containerWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	void *data;
	RECT r;

	data = (void *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
	if (data == NULL) {
		/* the lpParam is available during WM_NCCREATE and WM_CREATE */
		if (uMsg == WM_NCCREATE) {
			storelpParam(hwnd, lParam);
			data = (void *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
			storeContainerHWND(data, hwnd);
		}
		/* act as if we're not ready yet, even during WM_NCCREATE (nothing important to the switch statement below happens here anyway) */
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}

	switch (uMsg) {
	case WM_COMMAND:
		return forwardCommand(hwnd, uMsg, wParam, lParam);
	case WM_NOTIFY:
		return forwardNotify(hwnd, uMsg, wParam, lParam);
	case WM_SIZE:
		if (GetClientRect(hwnd, &r) == 0)
			xpanic("error getting client rect for Window in WM_SIZE", GetLastError());
		containerResize(data, &r);
		return 0;
	default:
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("container", "containerWndProc()", uMsg);
	return 0;		/* unreached */
}

DWORD makeContainerWindowClass(char **errmsg)
{
	WNDCLASSW wc;

	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpfnWndProc = containerWndProc;
	wc.hInstance = hInstance;
	wc.hIcon = hDefaultIcon;
	wc.hCursor = hArrowCursor;
	/* TODO won't this override our visual styles? */
	wc.hbrBackground = (HBRUSH) (COLOR_BTNFACE + 1);
	wc.lpszClassName = containerclass;
	if (RegisterClassW(&wc) == 0) {
		*errmsg = "error registering container window class";
		return GetLastError();
	}
	return 0;
}

HWND newContainer(void *data)
{
	HWND hwnd;

	hwnd = CreateWindowExW(
		0,
		containerclass, L"",
		WS_CHILD | WS_VISIBLE,
		CW_USEDEFAULT, CW_USEDEFAULT,
		100, 100,
		msgwin, NULL, hInstance, data);
	if (hwnd == NULL)
		xpanic("container creation failed", GetLastError());
	return hwnd;
}
