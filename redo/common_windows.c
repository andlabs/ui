// 17 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

LRESULT getWindowTextLen(HWND hwnd)
{
	return SendMessageW(hwnd, WM_GETTEXTLENGTH, 0, 0);
}

void getWindowText(HWND hwnd, WPARAM n, LPWSTR buf)
{
	SetLastError(0);
	if (SendMessageW(hwnd, WM_GETTEXT, n + 1, (LPARAM) buf) != (LRESULT) n)
		xpanic("WM_GETTEXT did not copy the correct number of characters out", GetLastError());
}

void setWindowText(HWND hwnd, LPWSTR text)
{
	switch (SendMessageW(hwnd, WM_SETTEXT, 0, (LPARAM) text)) {
	case FALSE:
		xpanic("WM_SETTEXT failed", GetLastError());
	}
}

void updateWindow(HWND hwnd)
{
	if (UpdateWindow(hwnd) == 0)
		xpanic("error calling UpdateWindow()", GetLastError());
}

void *getWindowData(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, LRESULT *lResult, void (*storeHWND)(void *, HWND))
{
	CREATESTRUCTW *cs = (CREATESTRUCTW *) lParam;
	void *data;

	data = (void *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
	if (data == NULL) {
		// the lpParam is available during WM_NCCREATE and WM_CREATE
		if (uMsg == WM_NCCREATE) {
			SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) (cs->lpCreateParams));
			data = (void *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
			(*storeHWND)(data, hwnd);
		}
		// act as if we're not ready yet, even during WM_NCCREATE (nothing important to the switch statement below happens here anyway)
		*lResult = DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	return data;
}

/*
all container windows (including the message-only window, hence this is not in container_windows.c) have to call the sharedWndProc() to ensure messages go in the right place and control colors are handled properly
*/

/*
all controls that have events receive the events themselves through subclasses
to do this, all container windows (including the message-only window; see http://support.microsoft.com/default.aspx?scid=KB;EN-US;Q104069) forward WM_COMMAND to each control with this function, WM_NOTIFY with forwardNotify, etc.
*/
static LRESULT forwardCommand(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	HWND control = (HWND) lParam;

	// don't generate an event if the control (if there is one) is unparented (a child of the message-only window)
	if (control != NULL && IsChild(msgwin, control) == 0)
		return SendMessageW(control, msgCOMMAND, wParam, lParam);
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

static LRESULT forwardNotify(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	HWND control = nmhdr->hwndFrom;

	// don't generate an event if the control (if there is one) is unparented (a child of the message-only window)
	if (control != NULL && IsChild(msgwin, control) == 0)
		return SendMessageW(control, msgNOTIFY, wParam, lParam);
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

BOOL sharedWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, LRESULT *lResult)
{
	switch (uMsg) {
	case WM_COMMAND:
		*lResult = forwardCommand(hwnd, uMsg, wParam, lParam);
		return TRUE;
	case WM_NOTIFY:
		*lResult = forwardNotify(hwnd, uMsg, wParam, lParam);
		return TRUE;
	case WM_CTLCOLORSTATIC:
	case WM_CTLCOLORBTN:
		if (SetBkMode((HDC) wParam, TRANSPARENT) == 0)
			xpanic("error setting transparent background mode to Labels", GetLastError());
		paintControlBackground((HWND) lParam, (HDC) wParam);
		*lResult = (LRESULT) hollowBrush;
		return TRUE;
	}
	return FALSE;
}

void paintControlBackground(HWND hwnd, HDC dc)
{
	HWND parent;
	RECT r;
	POINT p;
	int saved;
	WCHAR classname[128] = L"";		// more than enough to avoid collisions

	parent = hwnd;
	do {
		parent = GetParent(parent);
		if (parent == NULL)
			xpanic("error getting parent container of control in paintControlBackground()", GetLastError());
		parent = GetParent(parent);
		if (parent == NULL)
			xpanic("error getting parent control of control in paintControlBackground()", GetLastError());
		if (GetClassNameW(parent, classname, 128) == 0)
			xpanic("error getting name of focused window class in paintControlBackground()", GetLastError());
	} while (_wcsicmp(classname, L"button") == 0);		// skip groupboxes
	if (GetWindowRect(hwnd, &r) == 0)
		xpanic("error getting control's window rect in paintControlBackground()", GetLastError());
	// the above is a window rect; convert to client rect
	p.x = r.left;
	p.y = r.top;
	if (ScreenToClient(parent, &p) == 0)
		xpanic("error getting client origin of control in paintControlBackground()", GetLastError());
	saved = SaveDC(dc);
	if (saved == 0)
		xpanic("error saving DC info in paintControlBackground()", GetLastError());
	if (SetWindowOrgEx(dc, p.x, p.y, NULL) == 0)
		xpanic("error moving window origin in paintControlBackground()", GetLastError());
	SendMessageW(parent, WM_PRINTCLIENT, (WPARAM) dc, PRF_CLIENT);
	if (RestoreDC(dc, saved) == 0)
		xpanic("error restoring DC info in paintControlBackground()", GetLastError());
}
