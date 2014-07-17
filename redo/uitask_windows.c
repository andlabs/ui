/* 17 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

void uimsgloop(void)
{
	MSG msg;
	int res;

	for (;;) {
		SetLastError(0);
		res = GetMessage(&msg, NULL, 0, 0);
		if (res < 0)
			xpanic("error calling GetMessage()", GetLastError());
		if (res == 0)		/* WM_QUIT */
			break;
		/* TODO IsDialogMessage() */
		TranslateMessage(&msg);
		DispatchMessage(&msg);
	}
}

void issue(void *request)
{
	SetLastError(0);
	if (PostMessage(msgwin, msgRequested, 0, (LPARAM) request) == 0)
		xpanic("error issuing request", GetLastError());
}

HWND msgwin;

#define msgwinclass L"gouimsgwin"

static LRESULT CALLBACK msgwinproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	switch (uMsg) {
	case WM_COMMAND:
		return forwardCommand(hwnd, uMsg, wParam, lParam);
	case msgRequested:
		xperform((void *) lParam);
		return 0;
	default:
		return DefWindowProc(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("message-only", "msgwinproc()", uMsg);
	return 0;		/* unreachable */
}

DWORD makemsgwin(char **errmsg)
{
	WNDCLASS wc;
	HWND hwnd;

	ZeroMemory(&wc, sizeof (WNDCLASS));
	wc.lpfnWndProc = msgwinproc;
	wc.hInstance = hInstance;
	wc.hIcon = hDefaultIcon;
	wc.hCursor = hArrowCursor;
	wc.hbrBackground = (HBRUSH) (COLOR_BTNFACE + 1);
	wc.lpszClassName = msgwinclass;
	if (RegisterClass(&wc) == 0) {
		*errmsg = "error registering message-only window classs";
		return GetLastError();
	}
	msgwin = CreateWindowEx(
		0,
		msgwinclass, L"package ui message-only window",
		0,
		CW_USEDEFAULT, CW_USEDEFAULT,
		CW_USEDEFAULT, CW_USEDEFAULT,
		HWND_MESSAGE, NULL, hInstance, NULL);
	if (msgwin == NULL) {
		*errmsg = "error creating message-only window";
		return GetLastError();
	}
	return 0;
}
