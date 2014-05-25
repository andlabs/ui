// 9 february 2014

package ui

import (
//	"syscall"
//	"unsafe"
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
*/

var (
	_getScrollInfo = user32.NewProc("GetScrollInfo")
	_setScrollInfo = user32.NewProc("SetScrollInfo")
	_scrollWindowEx = user32.NewProc("ScrollWindowEx")
)

type _SCROLLINFO struct {
	cbSize		uint32
	fMask		uint32
	nMin			int32		// originally int
	nMax		int32		// originally int
	nPage		uint32
	nPos			int32		// originally int
	nTrackPos	int32		// originally int
}
