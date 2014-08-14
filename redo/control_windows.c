// 17 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

HWND newControl(LPWSTR class, DWORD style, DWORD extstyle)
{
	HWND hwnd;

	hwnd = CreateWindowExW(
		extstyle,
		class, L"",
		style | WS_CHILD | WS_VISIBLE,
		CW_USEDEFAULT, CW_USEDEFAULT,
		CW_USEDEFAULT, CW_USEDEFAULT,
		// the following has the consequence of making the control message-only at first
		// this shouldn't cause any problems... hopefully not
		// but see the msgwndproc() for caveat info
		// also don't use low control IDs as they will conflict with dialog boxes (IDCANCEL, etc.)
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

	// don't generate an event if the control (if there is one) is unparented (a child of the message-only window)
	if (control != NULL && IsChild(msgwin, control) == 0)
		return SendMessageW(control, msgCOMMAND, wParam, lParam);
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

LRESULT forwardNotify(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	HWND control = nmhdr->hwndFrom;

	// don't generate an event if the control (if there is one) is unparented (a child of the message-only window)
	if (control != NULL && IsChild(msgwin, control) == 0)
		return SendMessageW(control, msgNOTIFY, wParam, lParam);
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

void moveWindow(HWND hwnd, int x, int y, int width, int height)
{
	if (MoveWindow(hwnd, x, y, width, height, TRUE) == 0)
		xpanic("error setting window/control rect", GetLastError());
}

LONG controlTextLength(HWND hwnd, LPWSTR text)
{
	HDC dc;
	HFONT prev;
	SIZE size;

	dc = GetDC(hwnd);
	if (dc == NULL)
		xpanic("error getting DC of control for text length", GetLastError());
	prev = SelectObject(dc, controlFont);
	if (prev == NULL)
		xpanic("error setting control font to DC for text length", GetLastError());
	if (GetTextExtentPoint32W(dc, text, wcslen(text), &size) == 0)
		xpanic("error actually getting text length", GetLastError());
	if (SelectObject(dc, prev) != controlFont)
		xpanic("error restoring previous control font to DC for text length", GetLastError());
	if (ReleaseDC(hwnd, dc) == 0)
		xpanic("error releasing DC of control for text length", GetLastError());
	return size.cx;
}
