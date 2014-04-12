// 7 february 2014

package ui

import (
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32 = syscall.NewLazyDLL("gdi32.dll")
	comctl32 *syscall.LazyDLL		// comctl32 not defined here; see comctl_windows.go
	msimg32 = syscall.NewLazyDLL("msimg32.dll")
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

// In MSDN, _LPARAM and _LRESULT are listed as signed pointers, however their interpretation is message-specific. Ergo, just cast them yourself; it'll be the same. (Thanks to Tv` in #go-nuts for helping me realize this.)
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

func (l _LPARAM) _X() int32 {
	// according to windowsx.h
	loword := uint16(l & 0xFFFF)
	short := int16(loword)	// convert to signed...
	return int32(short)		// ...and sign extend
}

func (l _LPARAM) _Y() int32 {
	// according to windowsx.h
	hiword := uint16((l & 0xFFFF0000) >> 16)
	short := int16(hiword)	// convert to signed...
	return int32(short)		// ...and sign extend
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
