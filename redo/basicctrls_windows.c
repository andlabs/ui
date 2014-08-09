/* 17 july 2014 */

#include "winapi_windows.h"
#include "_cgo_export.h"

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
