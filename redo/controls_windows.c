/* 17 july 2014 */

#include "winapi_windows.h"

HWND newWidget(LPCWSTR class, DWORD style, DWORD extstyle)
{
	HWND hwnd;

	hwnd = CreateWindowExW(
		extstyle,
		class, L"",
		style | WS_CHILD | WS_VISIBLE,
		CW_USEDEFAULT, CW_USEDEFAULT,
		CW_USEDEFAULT, CW_USEDEFAULT,
		/*
		the following has the consequence of making the control message-only at first
		this shouldn't cause any problems... hopefully not
		but see the msgwndproc() for caveat info
		also don't use low control IDs as they will conflict with dialog boxes (IDCANCEL, etc.)
		*/
		msgwin, (HMENU) 100, hInstance, NULL);
	if (hwnd == NULL)
		xpanic("error creating control", GetLastError());
	return hwnd;
}

void controlSetParent(HWND control, HWND parent)
{
	if (SetParent(control, parent) == NULL)
		xpanic("error changing control parent", GetLastError());
}

/*
all controls that have events receive the events themselves through subclasses
to do this, all windows (including the message-only window; see http://support.microsoft.com/default.aspx?scid=KB;EN-US;Q104069) forward WM_COMMAND to each control with this function
*/
LRESULT forwardCommand(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	HWND control = (HWND) lParam;

	/* don't generate an event if the control (if there is one) is unparented (a child of the message-only window) */
	if (control != NULL && IsChild(msgwin, control) == 0)
		return SendMessageW(control, msgCOMMAND, wParam, lParam);
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

static LRESULT CALLBACK buttonSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	switch (uMsg) {
	case msgCOMMAND:
		if (HIWORD(wParam) == BN_CLICKED) {
			buttonClicked(data);
			return 0;
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, buttonSubProc, id) == FALSE)
			xpanic("error removing Button subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Button", "buttonSubProc()", uMsg);
	return 0;		/* unreached */
}

void setButtonSubclass(HWND hwnd, void *data)
{
	if ((*fv_SetWindowSubclass)(hwnd, buttonSubProc, 0, (DWORD_PTR) data) == FALSE)
		xpanic("error subclassing Button to give it its own event handler", GetLastError());
}
