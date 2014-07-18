/* 17 july 2014 */

#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
/* get Windows version right; right now Windows XP */
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <commctrl.h>
#include <stdint.h>

/* global messages unique to everything */
enum {
	msgRequest = WM_APP + 1,		/* + 1 just to be safe */
	msgCOMMAND,				/* WM_COMMAND proxy; see forwardCommand() in controls_windows.go */
};

/* uitask_windows.c */
extern void uimsgloop(void);
extern void issue(void *);
extern HWND msgwin;
extern DWORD makemsgwin(char **);

/* comctl32_windows.c */
extern DWORD initCommonControls(LPCWSTR, char **);
/* TODO do any of these take WINAPI? */
extern BOOL (*WINAPI fv_SetWindowSubclass)(HWND, SUBCLASSPROC, UINT_PTR, DWORD_PTR);
extern BOOL (*WINAPI fv_RemoveWindowSubclass)(HWND, SUBCLASSPROC, UINT_PTR);
extern LRESULT (*WINAPI fv_DefSubclassProc)(HWND, UINT, WPARAM, LPARAM);

/* controls_windows.c */
extern HWND newWidget(LPCWSTR, DWORD, DWORD);
extern void controlSetParent(HWND, HWND);
extern LRESULT forwardCommand(HWND, UINT, WPARAM, LPARAM);
extern void setButtonSubclass(HWND, void *);

/* init_windows.c */
extern HINSTANCE hInstnace;
extern int nCmdShow;
extern HICON hDefaultIcon;
extern HCURSOR hArrowCursor;
extern DWORD initWindows(char **);

/* sizing_windows.c */
extern HDC getDC(HWND);
extern void releaseDC(HWND, HDC);
extern void getTextMetricsW(HDC, TEXTMETRICW *);
extern void moveWindow(HWND, int, int, int, int);

/* window_windows.c */
extern DWORD makeWindowWindowClass(char **);
extern HWND newWindow(LPCWSTR, int, int, void *);
extern void windowClose(HWND);
