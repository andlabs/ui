// 8 february 2014

//
package ui

import (
//	"syscall"
//	"unsafe"
)

// Predefined cursor resource IDs.
const (
	IDC_APPSTARTING = 32650
	IDC_ARROW       = 32512
	IDC_CROSS       = 32515
	IDC_HAND        = 32649
	IDC_HELP        = 32651
	IDC_IBEAM       = 32513
	//	IDC_ICON = 32641		// [Obsolete for applications marked version 4.0 or later.]
	IDC_NO = 32648
	//	IDC_SIZE = 32640		// [Obsolete for applications marked version 4.0 or later. Use IDC_SIZEALL.]
	IDC_SIZEALL  = 32646
	IDC_SIZENESW = 32643
	IDC_SIZENS   = 32645
	IDC_SIZENWSE = 32642
	IDC_SIZEWE   = 32644
	IDC_UPARROW  = 32516
	IDC_WAIT     = 32514
)

var (
	loadCursor = user32.NewProc("LoadCursorW")
)

func LoadCursor_ResourceID(hInstance HANDLE, lpCursorName uint16) (cursor HANDLE, err error) {
	r1, _, err := loadCursor.Call(
		uintptr(hInstance),
		MAKEINTRESOURCE(lpCursorName))
	if r1 == 0 { // failure
		return NULL, err
	}
	return HANDLE(r1), nil
}
