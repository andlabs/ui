// 17 july 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

LRESULT getWindowTextLen(HWND hwnd)
{
	return SendMessageW(hwnd, WM_GETTEXTLENGTH, 0, 0);
}

void getWindowText(HWND hwnd, WPARAM n, LPCWSTR buf)
{
	SetLastError(0);
	if (SendMessageW(hwnd, WM_GETTEXT, n + 1, (LPARAM) buf) != n)
		xpanic("WM_GETTEXT did not copy the correct number of characters out", GetLastError());
}

void setWindowText(HWND hwnd, LPCWSTR text)
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

void storelpParam(HWND hwnd, LPARAM lParam)
{
	CREATESTRUCTW *cs = (CREATESTRUCTW *) lParam;

	SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) (cs->lpCreateParams));
}
