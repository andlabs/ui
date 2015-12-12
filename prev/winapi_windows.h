// 17 july 2014

// cgo will include this file multiple times
#ifndef __GO_UI_WINAPI_WINDOWS_H__
#define __GO_UI_WINAPI_WINDOWS_H__

#include "wininclude_windows.h"

// if by some stroke of luck Go ever supports compiling with MSVC, this will need to change
// noe that this has to come after the headers above because it's not predefined
#ifndef __MINGW64_VERSION_MAJOR
#error Sorry, you must use MinGW-w64 (http://mingw-w64.sourceforge.net/) to build package ui, as vanilla MinGW does not support Windows XP features (in 2014!).
#endif

// global messages unique to everything
enum {
	msgRequest = WM_APP + 1,		// + 1 just to be safe
	msgCOMMAND,				// WM_COMMAND proxy; see forwardCommand() in controls_windows.go
	msgNOTIFY,					// WM_NOTIFY proxy
	msgAreaSizeChanged,
	msgAreaGetScroll,
	msgAreaRepaint,
	msgAreaRepaintAll,
	msgTabCurrentTabHasChildren,
	msgAreaKeyDown,
	msgAreaKeyUp,
	msgOpenFileDone,
};

// uitask_windows.c
extern void uimsgloop(void);
extern void issue(void *);
extern HWND msgwin;
extern DWORD makemsgwin(char **);

// comctl32_windows.c
extern DWORD initCommonControls(char **);
// these are listed as WINAPI in both Microsoft's and MinGW's headers, but not on MSDN for some reason
extern BOOL (*WINAPI fv_SetWindowSubclass)(HWND, SUBCLASSPROC, UINT_PTR, DWORD_PTR);
extern BOOL (*WINAPI fv_RemoveWindowSubclass)(HWND, SUBCLASSPROC, UINT_PTR);
extern LRESULT (*WINAPI fv_DefSubclassProc)(HWND, UINT, WPARAM, LPARAM);
// these are listed as WINAPI on MSDN
extern BOOL (*WINAPI fv__TrackMouseEvent)(LPTRACKMOUSEEVENT);

// control_windows.c
extern HWND newControl(LPWSTR, DWORD, DWORD);
extern void controlSetParent(HWND, HWND);
extern void controlSetControlFont(HWND);
extern void moveWindow(HWND, int, int, int, int);
extern LONG controlTextLength(HWND, LPWSTR);

// basicctrls_windows.c
extern void setButtonSubclass(HWND, void *);
extern void setCheckboxSubclass(HWND, void *);
extern BOOL checkboxChecked(HWND);
extern void checkboxSetChecked(HWND, BOOL);
#define textfieldStyle (ES_AUTOHSCROLL | ES_LEFT | ES_NOHIDESEL | WS_TABSTOP)
#define textfieldExtStyle (WS_EX_CLIENTEDGE)
extern void setTextFieldSubclass(HWND, void *);
extern void textfieldSetAndShowInvalidBalloonTip(HWND, WCHAR *);
extern void textfieldHideInvalidBalloonTip(HWND);
extern int textfieldReadOnly(HWND);
extern void textfieldSetReadOnly(HWND, BOOL);
extern void setGroupSubclass(HWND, void *);
extern HWND newUpDown(HWND, void *);
extern void setSpinboxEditSubclass(HWND, void *);
extern LPWSTR xPROGRESS_CLASS;

// init_windows.c
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

// window_windows.c
extern DWORD makeWindowWindowClass(char **);
extern HWND newWindow(LPWSTR, int, int, void *);
extern void windowClose(HWND);

// common_windows.c
extern LRESULT getWindowTextLen(HWND);
extern void getWindowText(HWND, WPARAM, LPWSTR);
extern void setWindowText(HWND, LPWSTR);
extern void updateWindow(HWND);
extern void *getWindowData(HWND, UINT, WPARAM, LPARAM, LRESULT *);
extern int windowClassOf(HWND, ...);
extern BOOL sharedWndProc(HWND, UINT, WPARAM, LPARAM, LRESULT *);
extern void paintControlBackground(HWND, HDC);

// tab_windows.go
extern LPWSTR xWC_TABCONTROL;
extern void setTabSubclass(HWND, void *);
extern void tabAppend(HWND, LPWSTR);
extern void tabGetContentRect(HWND, RECT *);
extern LONG tabGetTabHeight(HWND);
extern void tabEnterChildren(HWND);
extern void tabLeaveChildren(HWND);

// table_windows.go
#include "wintable/includethis.h"
extern LPWSTR xtableWindowClass;
extern void doInitTable(void);
extern void setTableSubclass(HWND, void *);
extern void gotableSetRowCount(HWND, intptr_t);
/* TODO
extern void tableAutosizeColumns(HWND, int);
*/
extern intptr_t tableSelectedItem(HWND);
extern void tableSelectItem(HWND, intptr_t);

// container_windows.c
extern RECT containerBounds(HWND);
extern void calculateBaseUnits(HWND, int *, int *, LONG *);

// area_windows.c
#define areaWindowClass L"gouiarea"
extern void repaintArea(HWND, RECT *);
extern DWORD makeAreaWindowClass(char **);
extern HWND newArea(void *);
extern HWND newAreaTextField(HWND, void *);
extern void areaOpenTextField(HWND, HWND, int, int, int, int);
extern void areaMarkTextFieldDone(HWND);

// image_windows.c
extern HBITMAP toBitmap(void *, intptr_t, intptr_t);
extern void freeBitmap(uintptr_t);

// dialog_windows.c
extern void openFile(HWND, void *);

#endif
