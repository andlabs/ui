// 7 february 2014
package main

import (
	"syscall"
	"unsafe"
)

// TODO filter out commctrl.h stuff because the combobox and listbox stuff has values that differ between windows versions and bleh
// either that or switch to ComboBoxEx and ListView because they might not have that problem???

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32 = syscall.NewLazyDLL("gdi32.dll")
)

type _HANDLE uintptr
type _HWND _HANDLE
type _HBRUSH _HANDLE
type _HMENU _HANDLE

const (
	_NULL = 0
	_FALSE = 0		// from windef.h
	_TRUE = 1			// from windef.h
)

// TODO pull the thanks for these three from the old wingo source
// TODO put these in windows.go
type _WPARAM uintptr
type _LPARAM uintptr
type _LRESULT uintptr

func (w _WPARAM) LOWORD() uint16 {
	// according to windef.h
	return uint16(w & 0xFFFF)
}

func (w _WPARAM) HIWORD() uint16 {
	// according to windef.h
	return uint16((w >> 16) & 0xFFFF)
}

func _LPARAMFromString(str string) _LPARAM {
	return _LPARAM(unsafe.Pointer(syscall.StringToUTF16Ptr(str)))
}

// microsoft's header files do this
func _MAKEINTRESOURCE(what uint16) uintptr {
	return uintptr(what)
}

type _POINT struct {
	X	int32
	Y	int32
}

type _RECT struct {
	Left		int32
	Top		int32
	Right	int32
	Bottom	int32
}
