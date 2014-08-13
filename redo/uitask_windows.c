/* 17 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

/* note that this includes the terminating '\0' */
#define NAREACLASS (sizeof areaWindowClass / sizeof areaWindowClass[0])

void uimsgloop(void)
{
	MSG msg;
	int res;
	HWND active;
	WCHAR classchk[NAREACLASS];
	BOOL dodlgmessage;

	for (;;) {
		SetLastError(0);
		res = GetMessageW(&msg, NULL, 0, 0);
		if (res < 0)
			xpanic("error calling GetMessage()", GetLastError());
		if (res == 0)		/* WM_QUIT */
			break;
		active = GetActiveWindow();
		if (active != NULL) {
			HWND focus;

			dodlgmessage = TRUE;
			focus = GetFocus();
			if (focus != NULL) {
				// theoretically we could use the class atom to avoid a wcscmp()
				// however, raymond chen advises against this - http://blogs.msdn.com/b/oldnewthing/archive/2004/10/11/240744.aspx
				// while I have no idea what could deregister *my* window class from under *me*, better safe than sorry
				// we could also theoretically just send msgAreaDefocuses directly, but what DefWindowProc() does to a WM_APP message is undocumented
				if (GetClassNameW(focus, classchk, NAREACLASS) == 0)
					xpanic("error getting name of focused window class for Area check", GetLastError());
				if (wcscmp(classchk, areaWindowClass) == 0)
					dodlgmessage = FALSE;
			}
			if (dodlgmessage && IsDialogMessageW(active, &msg) != 0)
				continue;
		}
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
}

void issue(void *request)
{
	SetLastError(0);
	if (PostMessageW(msgwin, msgRequest, 0, (LPARAM) request) == 0)
		xpanic("error issuing request", GetLastError());
}

HWND msgwin;

#define msgwinclass L"gouimsgwin"

static LRESULT CALLBACK msgwinproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	switch (uMsg) {
	case WM_COMMAND:
		return forwardCommand(hwnd, uMsg, wParam, lParam);
	case WM_NOTIFY:
		return forwardNotify(hwnd, uMsg, wParam, lParam);
	case msgRequest:
		doissue((void *) lParam);
		return 0;
	default:
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("message-only", "msgwinproc()", uMsg);
	return 0;		/* unreachable */
}

DWORD makemsgwin(char **errmsg)
{
	WNDCLASSW wc;
	HWND hwnd;

	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpfnWndProc = msgwinproc;
	wc.hInstance = hInstance;
	wc.hIcon = hDefaultIcon;
	wc.hCursor = hArrowCursor;
	wc.hbrBackground = (HBRUSH) (COLOR_BTNFACE + 1);
	wc.lpszClassName = msgwinclass;
	if (RegisterClassW(&wc) == 0) {
		*errmsg = "error registering message-only window classs";
		return GetLastError();
	}
	msgwin = CreateWindowExW(
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
