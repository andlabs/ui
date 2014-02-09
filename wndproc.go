// 8 february 2014
package main

import (
//	"syscall"
//	"unsafe"
)

// TODO error handling
type WNDPROC func(hwnd HWND, uMsg uint32, wParam WPARAM, lParam LPARAM) LRESULT

var (
	defWindowProc = user32.NewProc("DefWindowProcW")
)

// TODO error handling
func DefWindowProc(hwnd HWND, uMsg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	r1, _, _ := defWindowProc.Call(
		uintptr(hwnd),
		uintptr(uMsg),
		uintptr(wParam),
		uintptr(lParam))
	return LRESULT(r1)
}
