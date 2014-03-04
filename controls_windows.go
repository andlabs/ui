// 9 february 2014

//
package ui

import (
//	"syscall"
//	"unsafe"
)

// Button styles.
const (
	// from winuser.h
	_BS_PUSHBUTTON      = 0x00000000
	_BS_DEFPUSHBUTTON   = 0x00000001
	_BS_CHECKBOX        = 0x00000002
	_BS_AUTOCHECKBOX    = 0x00000003
	_BS_RADIOBUTTON     = 0x00000004
	_BS_3STATE          = 0x00000005
	_BS_AUTO3STATE      = 0x00000006
	_BS_GROUPBOX        = 0x00000007
	_BS_USERBUTTON      = 0x00000008
	_BS_AUTORADIOBUTTON = 0x00000009
	_BS_PUSHBOX         = 0x0000000A
	_BS_OWNERDRAW       = 0x0000000B
	_BS_TYPEMASK        = 0x0000000F
	_BS_LEFTTEXT        = 0x00000020
	_BS_TEXT            = 0x00000000
	_BS_ICON            = 0x00000040
	_BS_BITMAP          = 0x00000080
	_BS_LEFT            = 0x00000100
	_BS_RIGHT           = 0x00000200
	_BS_CENTER          = 0x00000300
	_BS_TOP             = 0x00000400
	_BS_BOTTOM          = 0x00000800
	_BS_VCENTER         = 0x00000C00
	_BS_PUSHLIKE        = 0x00001000
	_BS_MULTILINE       = 0x00002000
	_BS_NOTIFY          = 0x00004000
	_BS_FLAT            = 0x00008000
	_BS_RIGHTBUTTON     = _BS_LEFTTEXT
	// from commctrl.h
//	_BS_SPLITBUTTON = 0x0000000C			// Windows Vista and newer and(/or?) comctl6 only
//	_BS_DEFSPLITBUTTON = 0x0000000D		// Windows Vista and newer and(/or?) comctl6 only
//	_BS_COMMANDLINK = 0x0000000E		// Windows Vista and newer and(/or?) comctl6 only
//	_BS_DEFCOMMANDLINK = 0x0000000F		// Windows Vista and newer and(/or?) comctl6 only
)

// Button messages.
// TODO check if any are not defined on Windows 2000
const (
	// from winuser.h
	_BM_GETCHECK     = 0x00F0
	_BM_SETCHECK     = 0x00F1
	_BM_GETSTATE     = 0x00F2
	_BM_SETSTATE     = 0x00F3
	_BM_SETSTYLE     = 0x00F4
	_BM_CLICK        = 0x00F5
	_BM_GETIMAGE     = 0x00F6
	_BM_SETIMAGE     = 0x00F7
	_BM_SETDONTCLICK = 0x00F8
)

// Button WM_COMMAND notifications.
const (
	// from winuser.h
	_BN_CLICKED       = 0
	_BN_PAINT         = 1
	_BN_HILITE        = 2
	_BN_UNHILITE      = 3
	_BN_DISABLE       = 4
	_BN_DOUBLECLICKED = 5
	_BN_PUSHED        = _BN_HILITE
	_BN_UNPUSHED      = _BN_UNHILITE
	_BN_DBLCLK        = _BN_DOUBLECLICKED
	_BN_SETFOCUS      = 6
	_BN_KILLFOCUS     = 7
)

// Button check states.
const (
	// from winuser.h
	_BST_UNCHECKED     = 0x0000
	_BST_CHECKED       = 0x0001
	_BST_INDETERMINATE = 0x0002
)

/*
var (
	checkDlgButton = user32.NewProc("CheckDlgButton")
	checkRadioButton = user32.NewProc("CheckRadioButton")
	isDlgButtonChecked = user32.NewProc("IsDlgButtonChecked")
)

func CheckDlgButton(hDlg HWND, nIDButton int, uCheck uint32) (err error) {
	r1, _, err := checkDlgButton.Call(
		uintptr(hDlg),
		uintptr(nIDButton),
		uintptr(uCheck))
	if r1 == 0 {		// failure
		return err
	}
	return nil
}

func CheckRadioButton(hDlg HWND, nIDFirstButton int, nIDLastButton int, nIDCheckButton int) (err error) {
	r1, _, err := checkRadioButton.Call(
		uintptr(hDlg),
		uintptr(nIDFirstButton),
		uintptr(nIDLastButton),
		uintptr(nIDCheckButton))
	if r1 == 0 {		// failure
		return err
	}
	return nil
}

// TODO handle errors
func IsDlgButtonChecked(hDlg HWND, nIDButton int) (state uint32, err error) {
	r1, _, _ := isDlgButtonChecked.Call(
		uintptr(hDlg),
		uintptr(nIDButton))
	return uint32(r1), nil
}
*/

// Combobox styles.
const (
	// from winuser.h
	_CBS_SIMPLE            = 0x0001
	_CBS_DROPDOWN          = 0x0002
	_CBS_DROPDOWNLIST      = 0x0003
	_CBS_OWNERDRAWFIXED    = 0x0010
	_CBS_OWNERDRAWVARIABLE = 0x0020
	_CBS_AUTOHSCROLL       = 0x0040
	_CBS_OEMCONVERT        = 0x0080
	_CBS_SORT              = 0x0100
	_CBS_HASSTRINGS        = 0x0200
	_CBS_NOINTEGRALHEIGHT  = 0x0400
	_CBS_DISABLENOSCROLL   = 0x0800
	_CBS_UPPERCASE         = 0x2000
	_CBS_LOWERCASE         = 0x4000
)

// Combobox messages.
// TODO filter out messages not provided in windows 2000
const (
	// from winuser.h
	_CB_GETEDITSEL            = 0x0140
	_CB_LIMITTEXT             = 0x0141
	_CB_SETEDITSEL            = 0x0142
	_CB_ADDSTRING             = 0x0143
	_CB_DELETESTRING          = 0x0144
	_CB_DIR                   = 0x0145
	_CB_GETCOUNT              = 0x0146
	_CB_GETCURSEL             = 0x0147
	_CB_GETLBTEXT             = 0x0148
	_CB_GETLBTEXTLEN          = 0x0149
	_CB_INSERTSTRING          = 0x014A
	_CB_RESETCONTENT          = 0x014B
	_CB_FINDSTRING            = 0x014C
	_CB_SELECTSTRING          = 0x014D
	_CB_SETCURSEL             = 0x014E
	_CB_SHOWDROPDOWN          = 0x014F
	_CB_GETITEMDATA           = 0x0150
	_CB_SETITEMDATA           = 0x0151
	_CB_GETDROPPEDCONTROLRECT = 0x0152
	_CB_SETITEMHEIGHT         = 0x0153
	_CB_GETITEMHEIGHT         = 0x0154
	_CB_SETEXTENDEDUI         = 0x0155
	_CB_GETEXTENDEDUI         = 0x0156
	_CB_GETDROPPEDSTATE       = 0x0157
	_CB_FINDSTRINGEXACT       = 0x0158
	_CB_SETLOCALE             = 0x0159
	_CB_GETLOCALE             = 0x015A
	_CB_GETTOPINDEX           = 0x015B
	_CB_SETTOPINDEX           = 0x015C
	_CB_GETHORIZONTALEXTENT   = 0x015D
	_CB_SETHORIZONTALEXTENT   = 0x015E
	_CB_GETDROPPEDWIDTH       = 0x015F
	_CB_SETDROPPEDWIDTH       = 0x0160
	_CB_INITSTORAGE           = 0x0161
	_CB_MULTIPLEADDSTRING     = 0x0163
	_CB_GETCOMBOBOXINFO       = 0x0164
)

// Combobox errors.
var ( // var so they can be cast to uintptr
	// from winuser.h
	_CB_ERR       = (-1)
	_CB_ERRSPACE  = (-2)
	_CBN_ERRSPACE = (-1)
)

// Combobox WM_COMMAND notificaitons.
// TODO filter out notifications not provided in windows 2000
const (
	// from winuser.h
	_CBN_SELCHANGE    = 1
	_CBN_DBLCLK       = 2
	_CBN_SETFOCUS     = 3
	_CBN_KILLFOCUS    = 4
	_CBN_EDITCHANGE   = 5
	_CBN_EDITUPDATE   = 6
	_CBN_DROPDOWN     = 7
	_CBN_CLOSEUP      = 8
	_CBN_SELENDOK     = 9
	_CBN_SELENDCANCEL = 10
)

// Edit control styles.
const (
	// from winuser.h
	_ES_LEFT        = 0x0000
	_ES_CENTER      = 0x0001
	_ES_RIGHT       = 0x0002
	_ES_MULTILINE   = 0x0004
	_ES_UPPERCASE   = 0x0008
	_ES_LOWERCASE   = 0x0010
	_ES_PASSWORD    = 0x0020
	_ES_AUTOVSCROLL = 0x0040
	_ES_AUTOHSCROLL = 0x0080
	_ES_NOHIDESEL   = 0x0100
	_ES_OEMCONVERT  = 0x0400
	_ES_READONLY    = 0x0800
	_ES_WANTRETURN  = 0x1000
	_ES_NUMBER      = 0x2000
)

// Edit control messages.
// TODO filter out messages not provided in windows 2000
const (
	// from winuser.h
	_EM_GETSEL              = 0x00B0
	_EM_SETSEL              = 0x00B1
	_EM_GETRECT             = 0x00B2
	_EM_SETRECT             = 0x00B3
	_EM_SETRECTNP           = 0x00B4
	_EM_SCROLL              = 0x00B5
	_EM_LINESCROLL          = 0x00B6
	_EM_SCROLLCARET         = 0x00B7
	_EM_GETMODIFY           = 0x00B8
	_EM_SETMODIFY           = 0x00B9
	_EM_GETLINECOUNT        = 0x00BA
	_EM_LINEINDEX           = 0x00BB
	_EM_SETHANDLE           = 0x00BC
	_EM_GETHANDLE           = 0x00BD
	_EM_GETTHUMB            = 0x00BE
	_EM_LINELENGTH          = 0x00C1
	_EM_REPLACESEL          = 0x00C2
	_EM_GETLINE             = 0x00C4
	_EM_LIMITTEXT           = 0x00C5
	_EM_CANUNDO             = 0x00C6
	_EM_UNDO                = 0x00C7
	_EM_FMTLINES            = 0x00C8
	_EM_LINEFROMCHAR        = 0x00C9
	_EM_SETTABSTOPS         = 0x00CB
	_EM_SETPASSWORDCHAR     = 0x00CC
	_EM_EMPTYUNDOBUFFER     = 0x00CD
	_EM_GETFIRSTVISIBLELINE = 0x00CE
	_EM_SETREADONLY         = 0x00CF
	_EM_SETWORDBREAKPROC    = 0x00D0
	_EM_GETWORDBREAKPROC    = 0x00D1
	_EM_GETPASSWORDCHAR     = 0x00D2
	_EM_SETMARGINS          = 0x00D3
	_EM_GETMARGINS          = 0x00D4
	_EM_SETLIMITTEXT        = _EM_LIMITTEXT // [;win40 Name change]
	_EM_GETLIMITTEXT        = 0x00D5
	_EM_POSFROMCHAR         = 0x00D6
	_EM_CHARFROMPOS         = 0x00D7
	_EM_SETIMESTATUS        = 0x00D8
	_EM_GETIMESTATUS        = 0x00D9
)

// Edit control WM_COMMAND notifications.
// TODO filter out notifications not provided in windows 2000
const (
	// from winuser.h
	_EN_SETFOCUS                    = 0x0100
	_EN_KILLFOCUS                   = 0x0200
	_EN_CHANGE                      = 0x0300
	_EN_UPDATE                      = 0x0400
	_EN_ERRSPACE                    = 0x0500
	_EN_MAXTEXT                     = 0x0501
	_EN_HSCROLL                     = 0x0601
	_EN_VSCROLL                     = 0x0602
	_EN_ALIGN_LTR_EC                = 0x0700
	_EN_ALIGN_RTL_EC                = 0x0701
	_EC_LEFTMARGIN                  = 0x0001
	_EC_RIGHTMARGIN                 = 0x0002
	_EC_USEFONTINFO                 = 0xFFFF
	_EMSIS_COMPOSITIONSTRING        = 0x0001
	_EIMES_GETCOMPSTRATONCE         = 0x0001
	_EIMES_CANCELCOMPSTRINFOCUS     = 0x0002
	_EIMES_COMPLETECOMPSTRKILLFOCUS = 0x0004
)

// Listbox styles.
const (
	// from winuser.h
	_LBS_NOTIFY            = 0x0001
	_LBS_SORT              = 0x0002
	_LBS_NOREDRAW          = 0x0004
	_LBS_MULTIPLESEL       = 0x0008
	_LBS_OWNERDRAWFIXED    = 0x0010
	_LBS_OWNERDRAWVARIABLE = 0x0020
	_LBS_HASSTRINGS        = 0x0040
	_LBS_USETABSTOPS       = 0x0080
	_LBS_NOINTEGRALHEIGHT  = 0x0100
	_LBS_MULTICOLUMN       = 0x0200
	_LBS_WANTKEYBOARDINPUT = 0x0400
	_LBS_EXTENDEDSEL       = 0x0800
	_LBS_DISABLENOSCROLL   = 0x1000
	_LBS_NODATA            = 0x2000
	_LBS_NOSEL             = 0x4000
	_LBS_COMBOBOX          = 0x8000
	_LBS_STANDARD          = (_LBS_NOTIFY | _LBS_SORT | _WS_VSCROLL | _WS_BORDER)
)

// Listbox messages.
// TODO filter out messages not provided in windows 2000
const (
	// from winuser.h
	_LB_ADDSTRING           = 0x0180
	_LB_INSERTSTRING        = 0x0181
	_LB_DELETESTRING        = 0x0182
	_LB_SELITEMRANGEEX      = 0x0183
	_LB_RESETCONTENT        = 0x0184
	_LB_SETSEL              = 0x0185
	_LB_SETCURSEL           = 0x0186
	_LB_GETSEL              = 0x0187
	_LB_GETCURSEL           = 0x0188
	_LB_GETTEXT             = 0x0189
	_LB_GETTEXTLEN          = 0x018A
	_LB_GETCOUNT            = 0x018B
	_LB_SELECTSTRING        = 0x018C
	_LB_DIR                 = 0x018D
	_LB_GETTOPINDEX         = 0x018E
	_LB_FINDSTRING          = 0x018F
	_LB_GETSELCOUNT         = 0x0190
	_LB_GETSELITEMS         = 0x0191
	_LB_SETTABSTOPS         = 0x0192
	_LB_GETHORIZONTALEXTENT = 0x0193
	_LB_SETHORIZONTALEXTENT = 0x0194
	_LB_SETCOLUMNWIDTH      = 0x0195
	_LB_ADDFILE             = 0x0196
	_LB_SETTOPINDEX         = 0x0197
	_LB_GETITEMRECT         = 0x0198
	_LB_GETITEMDATA         = 0x0199
	_LB_SETITEMDATA         = 0x019A
	_LB_SELITEMRANGE        = 0x019B
	_LB_SETANCHORINDEX      = 0x019C
	_LB_GETANCHORINDEX      = 0x019D
	_LB_SETCARETINDEX       = 0x019E
	_LB_GETCARETINDEX       = 0x019F
	_LB_SETITEMHEIGHT       = 0x01A0
	_LB_GETITEMHEIGHT       = 0x01A1
	_LB_FINDSTRINGEXACT     = 0x01A2
	_LB_SETLOCALE           = 0x01A5
	_LB_GETLOCALE           = 0x01A6
	_LB_SETCOUNT            = 0x01A7
	_LB_INITSTORAGE         = 0x01A8
	_LB_ITEMFROMPOINT       = 0x01A9
	_LB_MULTIPLEADDSTRING   = 0x01B1
	_LB_GETLISTBOXINFO      = 0x01B2
)

// Listbox errors.
var ( // var so they can be cast to uintptr
	// from winuser.h
	_LB_OKAY      = 0
	_LB_ERR       = (-1)
	_LB_ERRSPACE  = (-2)
	_LBN_ERRSPACE = (-2)
)

// Listbox WM_COMMAND notifications and message returns.
// TODO filter out notifications not provided in windows 2000
const (
	// from winuser.h
	_LBN_SELCHANGE = 1
	_LBN_DBLCLK    = 2
	_LBN_SELCANCEL = 3
	_LBN_SETFOCUS  = 4
	_LBN_KILLFOCUS = 5
)

// Static control styles.
const (
	// from winuser.h
	_SS_LEFT            = 0x00000000
	_SS_CENTER          = 0x00000001
	_SS_RIGHT           = 0x00000002
	_SS_ICON            = 0x00000003
	_SS_BLACKRECT       = 0x00000004
	_SS_GRAYRECT        = 0x00000005
	_SS_WHITERECT       = 0x00000006
	_SS_BLACKFRAME      = 0x00000007
	_SS_GRAYFRAME       = 0x00000008
	_SS_WHITEFRAME      = 0x00000009
	_SS_USERITEM        = 0x0000000A
	_SS_SIMPLE          = 0x0000000B
	_SS_LEFTNOWORDWRAP  = 0x0000000C
	_SS_OWNERDRAW       = 0x0000000D
	_SS_BITMAP          = 0x0000000E
	_SS_ENHMETAFILE     = 0x0000000F
	_SS_ETCHEDHORZ      = 0x00000010
	_SS_ETCHEDVERT      = 0x00000011
	_SS_ETCHEDFRAME     = 0x00000012
	_SS_TYPEMASK        = 0x0000001F
	_SS_REALSIZECONTROL = 0x00000040
	_SS_NOPREFIX        = 0x00000080
	_SS_NOTIFY          = 0x00000100
	_SS_CENTERIMAGE     = 0x00000200
	_SS_RIGHTJUST       = 0x00000400
	_SS_REALSIZEIMAGE   = 0x00000800
	_SS_SUNKEN          = 0x00001000
	_SS_EDITCONTROL     = 0x00002000
	_SS_ENDELLIPSIS     = 0x00004000
	_SS_PATHELLIPSIS    = 0x00008000
	_SS_WORDELLIPSIS    = 0x0000C000
	_SS_ELLIPSISMASK    = 0x0000C000
)

// Static control messages and WM_COMMAND notifications.
const (
	// from winuser.h
	_STM_SETICON  = 0x0170
	_STM_GETICON  = 0x0171
	_STM_SETIMAGE = 0x0172
	_STM_GETIMAGE = 0x0173
	_STN_CLICKED  = 0
	_STN_DBLCLK   = 1
	_STN_ENABLE   = 2
	_STN_DISABLE  = 3
)
