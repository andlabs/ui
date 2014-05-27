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

func (l _LPARAM) X() int32 {
	// according to windowsx.h
	loword := uint16(l & 0xFFFF)
	short := int16(loword)	// convert to signed...
	return int32(short)		// ...and sign extend
}

func (l _LPARAM) Y() int32 {
	// according to windowsx.h
	hiword := uint16((l & 0xFFFF0000) >> 16)
	short := int16(hiword)	// convert to signed...
	return int32(short)		// ...and sign extend
}

type _POINT struct {
	x	int32
	y	int32
}

type _RECT struct {
	left		int32
	top		int32
	right		int32
	bottom	int32
}

// Go doesn't allow negative constants to be forced into unsigned types at compile-time; this will do it at runtime.
// TODO make sure sign extension works fine here (check Go's rules and ABI sign extension rules)
func negConst(c int) uintptr {
	return uintptr(c)
}

var (
	_adjustWindowRectEx = user32.NewProc("AdjustWindowRectEx")
	_createWindowEx = user32.NewProc("CreateWindowExW")
	_getClientRect = user32.NewProc("GetClientRect")
	_moveWindow = user32.NewProc("MoveWindow")
	_setWindowLong = user32.NewProc("SetWindowLongW")
	_setWindowPos = user32.NewProc("SetWindowPos")
	_setWindowText = user32.NewProc("SetWindowTextW")
	_showWindow = user32.NewProc("ShowWindow")
)

type _MINMAXINFO struct {
	ptReserved		_POINT
	ptMaxSize		_POINT
	ptMaxPosition		_POINT
	ptMinTrackSize		_POINT
	ptMaxTrackSize	_POINT
}

func (l _LPARAM) MINMAXINFO() *_MINMAXINFO {
	return (*_MINMAXINFO)(unsafe.Pointer(l))
}
