// 17 july 2014

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
	return 0;		// unreached
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

			// we didn't use BS_AUTOCHECKBOX (see controls_windows.go) so we have to manage the check state ourselves
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
	return 0;		// unreached
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

static LRESULT CALLBACK textfieldSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	switch (uMsg) {
	case msgCOMMAND:
		if (HIWORD(wParam) == EN_CHANGE) {
			textfieldChanged((void *) data);
			return 0;
		}
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, textfieldSubProc, id) == FALSE)
			xpanic("error removing TextField subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("TextField", "textfieldSubProc()", uMsg);
	return 0;		// unreached
}

void setTextFieldSubclass(HWND hwnd, void *data)
{
	if ((*fv_SetWindowSubclass)(hwnd, textfieldSubProc, 0, (DWORD_PTR) data) == FALSE)
		xpanic("error subclassing TextField to give it its own event handler", GetLastError());
}

void textfieldSetAndShowInvalidBalloonTip(HWND hwnd, WCHAR *text)
{
	EDITBALLOONTIP ti;

	ZeroMemory(&ti, sizeof (EDITBALLOONTIP));
	ti.cbStruct = sizeof (EDITBALLOONTIP);
	// this is required to show the error icon
	// this probably should be localized...
	ti.pszTitle = L"Invalid Input";
	ti.pszText = text;
	ti.ttiIcon = TTI_ERROR;
	if (SendMessageW(hwnd, EM_SHOWBALLOONTIP, 0, (LPARAM) (&ti)) == FALSE)
		xpanic("error showing TextField.Invalid() balloon tip", GetLastError());
	if (MessageBeep(0xFFFFFFFF) == 0)
		xpanic("error beeping in response to TextField.Invalid()", GetLastError());
}

void textfieldHideInvalidBalloonTip(HWND hwnd)
{
	if (SendMessageW(hwnd, EM_HIDEBALLOONTIP, 0, 0) == FALSE)
		xpanic("error hiding TextField.Invalid() balloon tip", GetLastError());
}

static LRESULT CALLBACK groupSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	LRESULT lResult;
	RECT r;

	if (sharedWndProc(hwnd, uMsg, wParam, lParam, &lResult))
		return lResult;
	switch (uMsg) {
	// don't do this on WM_WINDOWPOSCHANGING; weird redraw issues will happen
	case WM_WINDOWPOSCHANGED:
		// don't use the WINDOWPOS rect here; the coordinates of the controls have to be in real client coordinates
		if (GetClientRect(hwnd, &r) == 0)
			xpanic("error getting client rect of Group for resizing its child Control", GetLastError());
		groupResized((void *) data, r);
		// and chain up
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, groupSubProc, id) == FALSE)
			xpanic("error removing Group subclass (which was for its own event handler)", GetLastError());
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	default:
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
	}
	xmissedmsg("Group", "groupSubProc()", uMsg);
	return 0;		// unreached
}

void setGroupSubclass(HWND hwnd, void *data)
{
	if ((*fv_SetWindowSubclass)(hwnd, groupSubProc, 0, (DWORD_PTR) data) == FALSE)
		xpanic("error subclassing Group to give it its own event handler", GetLastError());
}

// provided for cgo's benefit
LPWSTR xUPDOWN_CLASSW = UPDOWN_CLASSW;
