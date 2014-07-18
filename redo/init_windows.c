/* 17 july 2014 */

#include "winapi_windows.h"

HINSTANCE hInstnace;
int nCmdShow;

HICON hDefaultIcon;
HCURSOR hArrowCursor;

DWORD initWindows(char **errmsg)
{
	STARTUPINFOW si;

	/* WinMain() parameters */
	hInstance = GetModuleHandleW(NULL);
	if (hInstance == NULL) {
		*errmsg = "error getting hInstance";
		return GetLastError();
	}
	nCmdShow = SW_SHOWDEFAULT;
	GetStartupInfoW(&si);
	if ((si.dwFlags & STARTF_USESHOWWINDOW) != 0)
		nCmdShow = si.wShowWindow;

	/* icons and cursors */
	hDefaultIcon = LoadIconW(NULL, IDI_APPLICATION);
	if (hDefaultIcon == NULL) {
		*errmsg = "error loading default icon";
		return GetLastError();
	}
	hDefaultCursor = LoadCursorW(NULL, IDC_ARROW);
	if (hArrowCursor == NULL) {
		*errmsg = "error loading arrow (default) cursor";
		return GetLastError();
	}

	return 0;
}
