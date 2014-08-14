/* 17 july 2014 */

/* cgo will include this file multiple times */
#ifndef __GO_UI_WINAPI_WINDOWS_H__
#define __GO_UI_WINAPI_WINDOWS_H__

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
#include <uxtheme.h>
#include <string.h>
#include <wchar.h>
#include <windowsx.h>

/* global messages unique to everything */
enum {
	msgRequest = WM_APP + 1,		/* + 1 just to be safe */
	msgCOMMAND,				/* WM_COMMAND proxy; see forwardCommand() in controls_windows.go */
	msgNOTIFY,					/* WM_NOTIFY proxy */
	msgAreaSizeChanged,
	msgAreaRepaintAll,
};

/* uitask_windows.c */
extern void uimsgloop(void);
extern void issue(void *);
extern HWND msgwin;
extern DWORD makemsgwin(char **);

/* comctl32_windows.c */
extern DWORD initCommonControls(char **);
/* these are listed as WINAPI in both Microsoft's and MinGW's headers, but not on MSDN for some reason */
extern BOOL (*WINAPI fv_SetWindowSubclass)(HWND, SUBCLASSPROC, UINT_PTR, DWORD_PTR);
extern BOOL (*WINAPI fv_RemoveWindowSubclass)(HWND, SUBCLASSPROC, UINT_PTR);
extern LRESULT (*WINAPI fv_DefSubclassProc)(HWND, UINT, WPARAM, LPARAM);

/* control_windows.c */
extern HWND newControl(LPWSTR, DWORD, DWORD);
extern void controlSetParent(HWND, HWND);
extern void controlSetControlFont(HWND);
extern LRESULT forwardCommand(HWND, UINT, WPARAM, LPARAM);
extern LRESULT forwardNotify(HWND, UINT, WPARAM, LPARAM);
extern void moveWindow(HWND, int, int, int, int);
extern LONG controlTextLength(HWND, LPWSTR);

/* basicctrls_windows.c */
extern void setButtonSubclass(HWND, void *);
extern void setCheckboxSubclass(HWND, void *);
extern BOOL checkboxChecked(HWND);
extern void checkboxSetChecked(HWND, BOOL);

/* init_windows.c */
extern HINSTANCE hInstance;
extern int nCmdShow;
extern HICON hDefaultIcon;
extern HCURSOR hArrowCursor;
extern HFONT controlFont;
extern HFONT titleFont;
extern HFONT smallTitleFont;
extern HFONT menubarFont;
extern HFONT statusbarFont;
extern HBRUSH hollowBrush;
extern DWORD initWindows(char **);

/* window_windows.c */
extern DWORD makeWindowWindowClass(char **);
extern HWND newWindow(LPWSTR, int, int, void *);
extern void windowClose(HWND);

/* common_windows.c */
extern LRESULT getWindowTextLen(HWND);
extern void getWindowText(HWND, WPARAM, LPWSTR);
extern void setWindowText(HWND, LPWSTR);
extern void updateWindow(HWND);
extern void storelpParam(HWND, LPARAM);

/* tab_windows.go */
extern LPWSTR xWC_TABCONTROL;
extern void setTabSubclass(HWND, void *);
extern void tabAppend(HWND, LPWSTR);
extern void tabGetContentRect(HWND, RECT *);
extern LONG tabGetTabHeight(HWND);

/* table_windows.go */
extern LPWSTR xWC_LISTVIEW;
extern void setTableSubclass(HWND, void *);
extern void tableAppendColumn(HWND, int, LPWSTR);
extern void tableUpdate(HWND, int);
extern void tableAddExtendedStyles(HWND, LPARAM);
extern void tableAutosizeColumns(HWND, int);

/* container_windows.c */
extern DWORD makeContainerWindowClass(char **);
extern HWND newContainer(void *);
extern void calculateBaseUnits(HWND, int *, int *, LONG *);

/* area_window.c */
#define areaWindowClass L"gouiarea"
extern void repaintArea(HWND);
extern DWORD makeAreaWindowClass(char **);
extern HWND newArea(void *);

#endif
