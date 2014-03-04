// 8 february 2014

//
package ui

import (
	//	"syscall"
	"unsafe"
)

// Window styles.
const (
	_WS_BORDER           = 0x00800000
	_WS_CAPTION          = 0x00C00000
	_WS_CHILD            = 0x40000000
	_WS_CHILDWINDOW      = 0x40000000
	_WS_CLIPCHILDREN     = 0x02000000
	_WS_CLIPSIBLINGS     = 0x04000000
	_WS_DISABLED         = 0x08000000
	_WS_DLGFRAME         = 0x00400000
	_WS_GROUP            = 0x00020000
	_WS_HSCROLL          = 0x00100000
	_WS_ICONIC           = 0x20000000
	_WS_MAXIMIZE         = 0x01000000
	_WS_MAXIMIZEBOX      = 0x00010000
	_WS_MINIMIZE         = 0x20000000
	_WS_MINIMIZEBOX      = 0x00020000
	_WS_OVERLAPPED       = 0x00000000
	_WS_OVERLAPPEDWINDOW = (_WS_OVERLAPPED | _WS_CAPTION | _WS_SYSMENU | _WS_THICKFRAME | _WS_MINIMIZEBOX | _WS_MAXIMIZEBOX)
	_WS_POPUP            = 0x80000000
	_WS_POPUPWINDOW      = (_WS_POPUP | _WS_BORDER | _WS_SYSMENU)
	_WS_SIZEBOX          = 0x00040000
	_WS_SYSMENU          = 0x00080000
	_WS_TABSTOP          = 0x00010000
	_WS_THICKFRAME       = 0x00040000
	_WS_TILED            = 0x00000000
	_WS_TILEDWINDOW      = (_WS_OVERLAPPED | _WS_CAPTION | _WS_SYSMENU | _WS_THICKFRAME | _WS_MINIMIZEBOX | _WS_MAXIMIZEBOX)
	_WS_VISIBLE          = 0x10000000
	_WS_VSCROLL          = 0x00200000
)

// Extended window styles.
const (
	_WS_EX_ACCEPTFILES = 0x00000010
	_WS_EX_APPWINDOW   = 0x00040000
	_WS_EX_CLIENTEDGE  = 0x00000200
	//	_WS_EX_COMPOSITED = 0x02000000	// [Windows 2000:This style is not supported.]
	_WS_EX_CONTEXTHELP      = 0x00000400
	_WS_EX_CONTROLPARENT    = 0x00010000
	_WS_EX_DLGMODALFRAME    = 0x00000001
	_WS_EX_LAYERED          = 0x00080000
	_WS_EX_LAYOUTRTL        = 0x00400000
	_WS_EX_LEFT             = 0x00000000
	_WS_EX_LEFTSCROLLBAR    = 0x00004000
	_WS_EX_LTRREADING       = 0x00000000
	_WS_EX_MDICHILD         = 0x00000040
	_WS_EX_NOACTIVATE       = 0x08000000
	_WS_EX_NOINHERITLAYOUT  = 0x00100000
	_WS_EX_NOPARENTNOTIFY   = 0x00000004
	_WS_EX_OVERLAPPEDWINDOW = (_WS_EX_WINDOWEDGE | _WS_EX_CLIENTEDGE)
	_WS_EX_PALETTEWINDOW    = (_WS_EX_WINDOWEDGE | _WS_EX_TOOLWINDOW | _WS_EX_TOPMOST)
	_WS_EX_RIGHT            = 0x00001000
	_WS_EX_RIGHTSCROLLBAR   = 0x00000000
	_WS_EX_RTLREADING       = 0x00002000
	_WS_EX_STATICEDGE       = 0x00020000
	_WS_EX_TOOLWINDOW       = 0x00000080
	_WS_EX_TOPMOST          = 0x00000008
	_WS_EX_TRANSPARENT      = 0x00000020
	_WS_EX_WINDOWEDGE       = 0x00000100
)

// bizarrely, this value is given on the page for CreateMDIWindow, but not CreateWindow or CreateWindowEx
// I do it this way because Go won't let me shove the exact value into an int
var (
	__CW_USEDEFAULT uint = 0x80000000
	_CW_USEDEFAULT       = int(__CW_USEDEFAULT)
)

// GetSysColor values. These can be cast to HBRUSH (after adding 1) for WNDCLASS as well.
const (
	_COLOR_3DDKSHADOW              = 21
	_COLOR_3DFACE                  = 15
	_COLOR_3DHIGHLIGHT             = 20
	_COLOR_3DHILIGHT               = 20
	_COLOR_3DLIGHT                 = 22
	_COLOR_3DSHADOW                = 16
	_COLOR_ACTIVEBORDER            = 10
	_COLOR_ACTIVECAPTION           = 2
	_COLOR_APPWORKSPACE            = 12
	_COLOR_BACKGROUND              = 1
	_COLOR_BTNFACE                 = 15
	_COLOR_BTNHIGHLIGHT            = 20
	_COLOR_BTNHILIGHT              = 20
	_COLOR_BTNSHADOW               = 16
	_COLOR_BTNTEXT                 = 18
	_COLOR_CAPTIONTEXT             = 9
	_COLOR_DESKTOP                 = 1
	_COLOR_GRADIENTACTIVECAPTION   = 27
	_COLOR_GRADIENTINACTIVECAPTION = 28
	_COLOR_GRAYTEXT                = 17
	_COLOR_HIGHLIGHT               = 13
	_COLOR_HIGHLIGHTTEXT           = 14
	_COLOR_HOTLIGHT                = 26
	_COLOR_INACTIVEBORDER          = 11
	_COLOR_INACTIVECAPTION         = 3
	_COLOR_INACTIVECAPTIONTEXT     = 19
	_COLOR_INFOBK                  = 24
	_COLOR_INFOTEXT                = 23
	_COLOR_MENU                    = 4
	//	COLOR_MENUHILIGHT = 29	// [Windows 2000:This value is not supported.]
	//	COLOR_MENUBAR = 30		// [Windows 2000:This value is not supported.]
	_COLOR_MENUTEXT    = 7
	_COLOR_SCROLLBAR   = 0
	_COLOR_WINDOW      = 5
	_COLOR_WINDOWFRAME = 6
	_COLOR_WINDOWTEXT  = 8
)

// SetWindowPos hWndInsertAfter values.
const (
	_HWND_BOTTOM = _HWND(1)
	_HWND_TOP    = _HWND(0)
)

// SetWindowPos hWndInsertAfter values that Go won't allow as constants.
var (
	__HWND_NOTOPMOST = -2
	_HWND_NOTOPMOST  = _HWND(__HWND_NOTOPMOST)
	__HWND_TOPMOST   = -1
	_HWND_TOPMOST    = _HWND(__HWND_TOPMOST)
)

// SetWindowPos uFlags values.
const (
	_SWP_DRAWFRAME      = 0x0020
	_SWP_FRAMECHANGED   = 0x0020
	_SWP_HIDEWINDOW     = 0x0080
	_SWP_NOACTIVATE     = 0x0010
	_SWP_NOCOPYBITS     = 0x0100
	_SWP_NOMOVE         = 0x0002
	_SWP_NOOWNERZORDER  = 0x0200
	_SWP_NOREDRAW       = 0x0008
	_SWP_NOREPOSITION   = 0x0200
	_SWP_NOSENDCHANGING = 0x0400
	_SWP_NOSIZE         = 0x0001
	_SWP_NOZORDER       = 0x0004
	_SWP_SHOWWINDOW     = 0x0040
	_SWP_ASYNCWINDOWPOS = 0x4000
	_SWP_DEFERERASE     = 0x2000
)

// ShowWindow settings.
const (
	_SW_FORCEMINIMIZE   = 11
	_SW_HIDE            = 0
	_SW_MAXIMIZE        = 3
	_SW_MINIMIZE        = 6
	_SW_RESTORE         = 9
	_SW_SHOW            = 5
	_SW_SHOWDEFAULT     = 10
	_SW_SHOWMAXIMIZED   = 3
	_SW_SHOWMINIMIZED   = 2
	_SW_SHOWMINNOACTIVE = 7
	_SW_SHOWNA          = 8
	_SW_SHOWNOACTIVATE  = 4
	_SW_SHOWNORMAL      = 1
)

var (
	_createWindowEx = user32.NewProc("CreateWindowExW")
	_getClientRect  = user32.NewProc("GetClientRect")
	_moveWindow     = user32.NewProc("MoveWindow")
	_setWindowPos   = user32.NewProc("SetWindowPos")
	_setWindowText  = user32.NewProc("SetWindowTextW")
	_showWindow     = user32.NewProc("ShowWindow")
)

// WM_SETICON and WM_GETICON values.
const (
	_ICON_BIG    = 1
	_ICON_SMALL  = 0
	_ICON_SMALL2 = 2 // WM_GETICON only?
)

// Window messages.
const (
	_MN_GETHMENU      = 0x01E1
	_WM_ERASEBKGND    = 0x0014
	_WM_GETFONT       = 0x0031
	_WM_GETTEXT       = 0x000D
	_WM_GETTEXTLENGTH = 0x000E
	_WM_SETFONT       = 0x0030
	_WM_SETICON       = 0x0080
	_WM_SETTEXT       = 0x000C
)

// WM_INPUTLANGCHANGEREQUEST values.
const (
	_INPUTLANGCHANGE_BACKWARD   = 0x0004
	_INPUTLANGCHANGE_FORWARD    = 0x0002
	_INPUTLANGCHANGE_SYSCHARSET = 0x0001
)

// WM_NCCALCSIZE return values.
const (
	_WVR_ALIGNTOP    = 0x0010
	_WVR_ALIGNRIGHT  = 0x0080
	_WVR_ALIGNLEFT   = 0x0020
	_WVR_ALIGNBOTTOM = 0x0040
	_WVR_HREDRAW     = 0x0100
	_WVR_VREDRAW     = 0x0200
	_WVR_REDRAW      = 0x0300
	_WVR_VALIDRECTS  = 0x0400
)

// WM_SHOWWINDOW reasons (lParam).
const (
	_SW_OTHERUNZOOM   = 4
	_SW_OTHERZOOM     = 2
	_SW_PARENTCLOSING = 1
	_SW_PARENTOPENING = 3
)

// WM_SIZE values.
const (
	_SIZE_MAXHIDE   = 4
	_SIZE_MAXIMIZED = 2
	_SIZE_MAXSHOW   = 3
	_SIZE_MINIMIZED = 1
	_SIZE_RESTORED  = 0
)

// WM_SIZING edge values (wParam).
const (
	_WMSZ_BOTTOM      = 6
	_WMSZ_BOTTOMLEFT  = 7
	_WMSZ_BOTTOMRIGHT = 8
	_WMSZ_LEFT        = 1
	_WMSZ_RIGHT       = 2
	_WMSZ_TOP         = 3
	_WMSZ_TOPLEFT     = 4
	_WMSZ_TOPRIGHT    = 5
)

// WM_STYLECHANGED and WM_STYLECHANGING values (wParam).
const (
	_GWL_EXSTYLE = -20
	_GWL_STYLE   = -16
)

// Window notifications.
const (
	_WM_ACTIVATEAPP   = 0x001C
	_WM_CANCELMODE    = 0x001F
	_WM_CHILDACTIVATE = 0x0022
	_WM_CLOSE         = 0x0010
	_WM_COMPACTING    = 0x0041
	_WM_CREATE        = 0x0001
	_WM_DESTROY       = 0x0002
	//	_WM_DPICHANGED = 0x02E0		// Windows 8.1 and newer only
	_WM_ENABLE                 = 0x000A
	_WM_ENTERSIZEMOVE          = 0x0231
	_WM_EXITSIZEMOVE           = 0x0232
	_WM_GETICON                = 0x007F
	_WM_GETMINMAXINFO          = 0x0024
	_WM_INPUTLANGCHANGE        = 0x0051
	_WM_INPUTLANGCHANGEREQUEST = 0x0050
	_WM_MOVE                   = 0x0003
	_WM_MOVING                 = 0x0216
	_WM_NCACTIVATE             = 0x0086
	_WM_NCCALCSIZE             = 0x0083
	_WM_NCCREATE               = 0x0081
	_WM_NCDESTROY              = 0x0082
	_WM_NULL                   = 0x0000
	_WM_QUERYDRAGICON          = 0x0037
	_WM_QUERYOPEN              = 0x0013
	_WM_QUIT                   = 0x0012
	_WM_SHOWWINDOW             = 0x0018
	_WM_SIZE                   = 0x0005
	_WM_SIZING                 = 0x0214
	_WM_STYLECHANGED           = 0x007D
	_WM_STYLECHANGING          = 0x007C
	//	_WM_THEMECHANGED = 0x031A		// Windows XP and newer only
	//	_WM_USERCHANGED = 0x0054			// Windows XP only: [Note  This message is not supported as of Windows Vista.; also listed as not supported by server Windows]
	_WM_WINDOWPOSCHANGED  = 0x0047
	_WM_WINDOWPOSCHANGING = 0x0046
)

type _MINMAXINFO struct {
	PtReserved     _POINT
	PtMaxSize      _POINT
	PtMaxPosition  _POINT
	PtMinTrackSize _POINT
	PtMaxTrackSize _POINT
}

func (l _LPARAM) MINMAXINFO() *_MINMAXINFO {
	return (*_MINMAXINFO)(unsafe.Pointer(l))
}
