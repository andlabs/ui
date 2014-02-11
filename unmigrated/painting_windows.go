// 9 february 2014
package main

import (
//	"syscall"
//	"unsafe"
)

var (
	updateWindow = user32.NewProc("UpdateWindow")
)

// TODO is error handling valid here? MSDN just says zero on failure; syscall.LazyProc.Call() always returns non-nil
func UpdateWindow(hWnd HWND) (err error) {
	r1, _, err := updateWindow.Call(uintptr(hWnd))
	if r1 == 0 {		// failure
		return err
	}
	return nil
}
