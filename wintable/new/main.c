// 19 october 2014
#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
#define CINTERFACE
// get Windows version right; right now Windows XP
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <commctrl.h>
#include <stdint.h>
#include <uxtheme.h>
#include <string.h>
#include <wchar.h>
#include <windowsx.h>
#include <vsstyle.h>
#include <vssym32.h>
#include <oleacc.h>

// #qo LIBS: user32 kernel32 gdi32 comctl32 uxtheme

#define tableWindowClass L"gouitable"

// start at WM_USER + 20 just in case for whatever reason we ever get the various dialog manager messages (see also http://blogs.msdn.com/b/oldnewthing/archive/2003/10/21/55384.aspx)
enum {
	// wParam - one of the type constants
	// lParam - column name as a Unicode string
	tableAddColumn = WM_USER + 20,
};

enum {
	tableColumnText,
	tableColumnImage,
	tableColumnCheckbox,
	nTableColumnTypes,
};

struct table {
	HWND hwnd;
};

#include "util.h"
#include "coord.h"
#include "events.h"

static const handlerfunc handlers[] = {
	eventHandlers,
	NULL,
};

static LRESULT CALLBACK tableWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	struct table *t;
	LRESULT lResult;

	t = (struct table *) GetWindowLongPtrW(hwnd, GWLP_USERDATA);
	if (t == NULL) {
		// we have to do things this way because creating the header control will fail mysteriously if we create it first thing
		// (which is fine; we can get the parent hInstance this way too)
		if (uMsg == WM_NCCREATE) {
			CREATESTRUCTW *cs = (CREATESTRUCTW *) lParam;

			t = (struct table *) malloc(sizeof (struct table));
			if (t == NULL)
				abort();
			ZeroMemory(t, sizeof (struct table));
			t->hwnd = hwnd;
			SetWindowLongPtrW(hwnd, GWLP_USERDATA, (LONG_PTR) t);
		}
		// even if we did the above, fall through
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	// TODO decide if we do this in WM_DESTROY or WM_NCDESTROY (need to see how the header control responds) so this can be in sync witht he above
	if (uMsg == WM_DESTROY) {
		// TODO free appropriate (after figuring this part out) components of t
		// TODO send DESTROY events to accessibility listeners (when appropriate)
		free(t);
		return DefWindowProcW(hwnd, uMsg, wParam, lParam);
	}
	if (runHandlers(handlers, t, uMsg, wParam, lParam, &lResult))
		return lResult;
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

void makeTableWindowClass(void)
{
	WNDCLASSW wc;

	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpszClassName = tableWindowClass;
	wc.lpfnWndProc = tableWndProc;
	wc.hCursor = LoadCursorW(NULL, IDC_ARROW);
	wc.hIcon = LoadIconW(NULL, IDI_APPLICATION);
	wc.hbrBackground = (HBRUSH) (COLOR_WINDOW + 1);		// TODO correct?
	wc.style = CS_HREDRAW | CS_VREDRAW;
	wc.hInstance = GetModuleHandle(NULL);
	if (RegisterClassW(&wc) == 0)
		abort();
}

int main(int argc, char *argv[])
{
	HWND mainwin;
	MSG msg;
	INITCOMMONCONTROLSEX icc;

	ZeroMemory(&icc, sizeof (INITCOMMONCONTROLSEX));
	icc.dwSize = sizeof (INITCOMMONCONTROLSEX);
	icc.dwICC = ICC_LISTVIEW_CLASSES;
	if (InitCommonControlsEx(&icc) == 0)
		abort();
	makeTableWindowClass();
	mainwin = CreateWindowExW(0,
		tableWindowClass, L"Main Window",
		WS_OVERLAPPEDWINDOW | WS_HSCROLL | WS_VSCROLL,
		CW_USEDEFAULT, CW_USEDEFAULT,
		400, 400,
		NULL, NULL, GetModuleHandle(NULL), NULL);
	if (mainwin == NULL)
		abort();
	SendMessageW(mainwin, tableAddColumn, tableColumnText, (LPARAM) L"Column");
	SendMessageW(mainwin, tableAddColumn, tableColumnImage, (LPARAM) L"Column 2");
	SendMessageW(mainwin, tableAddColumn, tableColumnCheckbox, (LPARAM) L"Column 3");
	if (argc > 1) {
		NONCLIENTMETRICSW ncm;
		HFONT font;

		ZeroMemory(&ncm, sizeof (NONCLIENTMETRICSW));
		ncm.cbSize = sizeof (NONCLIENTMETRICSW);
		if (SystemParametersInfoW(SPI_GETNONCLIENTMETRICS, sizeof (NONCLIENTMETRICSW), &ncm, sizeof (NONCLIENTMETRICSW)) == 0)
			abort();
		font = CreateFontIndirectW(&ncm.lfMessageFont);
		if (font == NULL)
			abort();
		SendMessageW(mainwin, WM_SETFONT, (WPARAM) font, TRUE);
	}
	ShowWindow(mainwin, SW_SHOWDEFAULT);
	if (UpdateWindow(mainwin) == 0)
		abort();
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
	return 0;
}
