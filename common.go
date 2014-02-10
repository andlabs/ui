// 7 february 2014
package main

import (
	"syscall"
	"unsafe"
)

// TODO filter out commctrl.h stuff because the combobox stuff has values that differ between windows versions and bleh

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

type HANDLE uintptr
type HWND HANDLE
type HBRUSH HANDLE
type HMENU HANDLE

const (
	NULL = 0
)

type ATOM uint16

// TODO pull the thanks for these three from the old wingo source
// TODO put these in windows.go
type WPARAM uintptr
type LPARAM uintptr
type LRESULT uintptr

func (w WPARAM) LOWORD() uint16 {
	// according to windef.h
	return uint16(w & 0xFFFF)
}

func (w WPARAM) HIWORD() uint16 {
	// according to windef.h
	return uint16((w >> 16) & 0xFFFF)
}

func LPARAMFromString(str string) LPARAM {
	return LPARAM(unsafe.Pointer(syscall.StringToUTF16Ptr(str)))
}

// microsoft's header files do this
func MAKEINTRESOURCE(what uint16) uintptr {
	return uintptr(what)
}

// TODO adorn error messages with which step failed?
func getText(hwnd HWND) (text string, err error) {
	var tc []uint16

	length, err := SendMessage(hwnd, WM_GETTEXTLENGTH, 0, 0)
	if err != nil {
		return "", err
	}
	tc = make([]uint16, length + 1)
	_, err = SendMessage(hwnd,
		WM_GETTEXT,
		WPARAM(length + 1),
		LPARAM(unsafe.Pointer(&tc[0])))
	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(tc), nil
}
