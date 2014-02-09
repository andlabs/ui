// 7 february 2014
package main

import (
	"syscall"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

type HANDLE uintptr
type HWND HANDLE
type HBRUSH HANDLE

const (
	NULL = 0
)

type ATOM uint16

// TODO pull the thanks for these three from the old wingo source
type WPARAM uintptr
type LPARAM uintptr
type LRESULT uintptr

// microsoft's header files do this
func MAKEINTRESOURCE(what uint16) uintptr {
	return uintptr(what)
}
