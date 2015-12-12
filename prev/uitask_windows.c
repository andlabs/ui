// 17 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

void uimsgloop_area(HWND active, HWND focus, MSG *msg)
{
	MSG copy;

	copy = *msg;
	switch (copy.message) {
	case WM_KEYDOWN:
	case WM_SYSKEYDOWN:			// Alt+[anything] and F10 send these instead
		copy.message = msgAreaKeyDown;
		break;
	case WM_KEYUP:
	case WM_SYSKEYUP:
		copy.message = msgAreaKeyUp;
		break;
	default:
		goto notkey;
	}
	// if we handled the key, don't do the default behavior
	// don't call TranslateMessage(); we do our own keyboard handling
	if (DispatchMessage(&copy) != FALSE)
		return;
notkey:
	if (IsDialogMessage(active, msg) != 0)
		return;
	DispatchMessage(msg);
}

void uimsgloop_tab(HWND active, HWND focus, MSG *msg)
{
	BOOL hasChildren;
	BOOL idm;

	// THIS BIT IS IMPORTANT: if the current tab has no children, then there will be no children left in the dialog to tab to, and IsDialogMessageW() will loop forever
	hasChildren = SendMessageW(focus, msgTabCurrentTabHasChildren, 0, 0);
	if (hasChildren)
		tabEnterChildren(focus);
	idm = IsDialogMessageW(active, msg);
	if (hasChildren)
		tabLeaveChildren(focus);
	if (idm != 0)
		return;
	TranslateMessage(msg);
	DispatchMessage(msg);
}

void uimsgloop_else(MSG *msg)
{
	TranslateMessage(msg);
	DispatchMessage(msg);
}

void uimsgloop(void)
{
	MSG msg;
	int res;
	HWND active, focus;
	BOOL dodlgmessage;

	for (;;) {
		SetLastError(0);
		res = GetMessageW(&msg, NULL, 0, 0);
		if (res < 0)
			xpanic("error calling GetMessage()", GetLastError());
		if (res == 0)		// WM_QUIT
			break;
		active = GetActiveWindow();
		if (active == NULL) {
			uimsgloop_else(&msg);
			continue;
		}

		// bit of logic involved here:
		// we don't want dialog messages passed into Areas, so we don't call IsDialogMessageW() there
		// as for Tabs, we can't have both WS_TABSTOP and WS_EX_CONTROLPARENT set at the same time, so we hotswap the two styles to get the behavior we want
		focus = GetFocus();
		if (focus != NULL) {
			switch (windowClassOf(focus, areaWindowClass, WC_TABCONTROLW, NULL)) {
			case 0:		// areaWindowClass
				uimsgloop_area(active, focus, &msg);
				continue;
			case 1:		// WC_TABCONTROLW
				uimsgloop_tab(active, focus, &msg);
				continue;
			}
			// else fall through
		}

		if (IsDialogMessage(active, &msg) != 0)
			continue;
		uimsgloop_else(&msg);
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
	LRESULT shared;
	size_t i;

	if (sharedWndProc(hwnd, uMsg, wParam, lParam, &shared))
		return shared;
	switch (uMsg) {
	case msgRequest:
		doissue((void *) lParam);
		return 0;
	case msgOpenFileDone:
		finishOpenFile((WCHAR *) wParam, (void *) lParam);
		return 0;
	default:
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("message-only", "msgwinproc()", uMsg);
	return 0;		// unreachable
}

DWORD makemsgwin(char **errmsg)
{
	WNDCLASSW wc;

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
