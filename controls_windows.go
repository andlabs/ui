// 9 february 2014
package main

import (
//	"syscall"
//	"unsafe"
)

// Button styles.
const (
	// from winuser.h
	BS_PUSHBUTTON = 0x00000000
	BS_DEFPUSHBUTTON = 0x00000001
	BS_CHECKBOX = 0x00000002
	BS_AUTOCHECKBOX = 0x00000003
	BS_RADIOBUTTON = 0x00000004
	BS_3STATE = 0x00000005
	BS_AUTO3STATE = 0x00000006
	BS_GROUPBOX = 0x00000007
	BS_USERBUTTON = 0x00000008
	BS_AUTORADIOBUTTON = 0x00000009
	BS_PUSHBOX = 0x0000000A
	BS_OWNERDRAW = 0x0000000B
	BS_TYPEMASK = 0x0000000F
	BS_LEFTTEXT = 0x00000020
	BS_TEXT = 0x00000000
	BS_ICON = 0x00000040
	BS_BITMAP = 0x00000080
	BS_LEFT = 0x00000100
	BS_RIGHT = 0x00000200
	BS_CENTER = 0x00000300
	BS_TOP = 0x00000400
	BS_BOTTOM = 0x00000800
	BS_VCENTER = 0x00000C00
	BS_PUSHLIKE = 0x00001000
	BS_MULTILINE = 0x00002000
	BS_NOTIFY = 0x00004000
	BS_FLAT = 0x00008000
	BS_RIGHTBUTTON = BS_LEFTTEXT
	// from commctrl.h
//	BS_SPLITBUTTON = 0x0000000C		// Windows Vista and newer and(/or?) comctl6 only
//	BS_DEFSPLITBUTTON = 0x0000000D	// Windows Vista and newer and(/or?) comctl6 only
//	BS_COMMANDLINK = 0x0000000E		// Windows Vista and newer and(/or?) comctl6 only
//	BS_DEFCOMMANDLINK = 0x0000000F	// Windows Vista and newer and(/or?) comctl6 only
)

// Button WM_COMMAND notifications.
const (
	// from winuser.h
	BN_CLICKED = 0
	BN_PAINT = 1
	BN_HILITE = 2
	BN_UNHILITE = 3
	BN_DISABLE = 4
	BN_DOUBLECLICKED = 5
	BN_PUSHED = BN_HILITE
	BN_UNPUSHED = BN_UNHILITE
	BN_DBLCLK = BN_DOUBLECLICKED
	BN_SETFOCUS = 6
	BN_KILLFOCUS = 7
)

// Button check states.
const (
	// from winuser.h
	BST_UNCHECKED = 0x0000
	BST_CHECKED = 0x0001
	BST_INDETERMINATE = 0x0002
)

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
