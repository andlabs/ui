/* 17 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

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

void controlSetControlFont(HWND which)
{
	SendMessageW(which, WM_SETFONT, (WPARAM) controlFont, TRUE);
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
			buttonClicked((void *) data);
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

static LRESULT CALLBACK checkboxSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	switch (uMsg) {
	case msgCOMMAND:
		if (HIWORD(wParam) == BN_CLICKED) {
			WPARAM check;

			/* we didn't use BS_AUTOCHECKBOX (see controls_windows.go) so we have to manage the check state ourselves */
			check = BST_CHECKED;
			if (SendMessage(hwnd, BM_GETCHECK, 0, 0) == BST_CHECKED)
				check = BST_UNCHECKED;
			SendMessage(hwnd, BM_SETCHECK, check, 0);
			checkboxToggled((void *) data);
			return 0;
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, checkboxSubProc, id) == FALSE)
			xpanic("error removing Checkbox subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Checkbox", "checkboxSubProc()", uMsg);
	return 0;		/* unreached */
}

void setCheckboxSubclass(HWND hwnd, void *data)
{
	if ((*fv_SetWindowSubclass)(hwnd, checkboxSubProc, 0, (DWORD_PTR) data) == FALSE)
		xpanic("error subclassing Checkbox to give it its own event handler", GetLastError());
}

BOOL checkboxChecked(HWND hwnd)
{
	if (SendMessage(hwnd, BM_GETCHECK, 0, 0) == BST_UNCHECKED)
		return FALSE;
	return TRUE;
}

void checkboxSetChecked(HWND hwnd, BOOL c)
{
	WPARAM check;

	check = BST_CHECKED;
	if (c == FALSE)
		check = BST_UNCHECKED;
	SendMessage(hwnd, BM_SETCHECK, check, 0);
}
