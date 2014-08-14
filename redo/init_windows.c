// 17 july 2014

#include "winapi_windows.h"

HINSTANCE hInstance;
int nCmdShow;

HICON hDefaultIcon;
HCURSOR hArrowCursor;

HFONT controlFont;
HFONT titleFont;
HFONT smallTitleFont;
HFONT menubarFont;
HFONT statusbarFont;

HBRUSH hollowBrush;

DWORD initWindows(char **errmsg)
{
	STARTUPINFOW si;
	NONCLIENTMETRICSW ncm;

	// WinMain() parameters
	hInstance = GetModuleHandleW(NULL);
	if (hInstance == NULL) {
		*errmsg = "error getting hInstance";
		return GetLastError();
	}
	nCmdShow = SW_SHOWDEFAULT;
	GetStartupInfoW(&si);
	if ((si.dwFlags & STARTF_USESHOWWINDOW) != 0)
		nCmdShow = si.wShowWindow;

	// icons and cursors
	hDefaultIcon = LoadIconW(NULL, IDI_APPLICATION);
	if (hDefaultIcon == NULL) {
		*errmsg = "error loading default icon";
		return GetLastError();
	}
	hArrowCursor = LoadCursorW(NULL, IDC_ARROW);
	if (hArrowCursor == NULL) {
		*errmsg = "error loading arrow (default) cursor";
		return GetLastError();
	}

	// standard fonts
#define GETFONT(l, f, n) l = CreateFontIndirectW(&ncm.f); \
	if (l == NULL) { \
		*errmsg = "error loading " n " font"; \
		return GetLastError(); \
	}

	ZeroMemory(&ncm, sizeof (NONCLIENTMETRICSW));
	ncm.cbSize = sizeof (NONCLIENTMETRICSW);
	if (SystemParametersInfoW(SPI_GETNONCLIENTMETRICS, sizeof (NONCLIENTMETRICSW), &ncm, sizeof (NONCLIENTMETRICSW)) == 0) {
		*errmsg = "error getting non-client metrics parameters";
		return GetLastError();
	}
	GETFONT(controlFont, lfMessageFont, "control");
	GETFONT(titleFont, lfCaptionFont, "titlebar");
	GETFONT(smallTitleFont, lfSmCaptionFont, "small title bar");
	GETFONT(menubarFont, lfMenuFont, "menu bar");
	GETFONT(statusbarFont, lfStatusFont, "status bar");

	hollowBrush = GetStockObject(HOLLOW_BRUSH);
	if (hollowBrush == NULL) {
		*errmsg = "error getting hollow brush";
		return GetLastError();
	}

	return 0;
}
