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
