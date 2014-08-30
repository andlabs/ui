// 9 february 2014

package ui

import (
//	"syscall"
//	"unsafe"
)

/*
var (
	_checkRadioButton = user32.NewProc("CheckRadioButton")
)

func CheckRadioButton(hDlg HWND, nIDFirstButton int, nIDLastButton int, nIDCheckButton int) (err error) {
	r1, _, err := _checkRadioButton.Call(
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
	_setScrollInfo  = user32.NewProc("SetScrollInfo")
	_scrollWindowEx = user32.NewProc("ScrollWindowEx")
)

type _SCROLLINFO struct {
	cbSize    uint32
	fMask     uint32
	nMin      int32 // originally int
	nMax      int32 // originally int
	nPage     uint32
	nPos      int32 // originally int
	nTrackPos int32 // originally int
}
