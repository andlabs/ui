// 7 february 2014

package ui

import (
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	comctl32 *syscall.LazyDLL // comctl32 not defined here; see comctl_windows.go
	msimg32  = syscall.NewLazyDLL("msimg32.dll")
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

func (l _LPARAM) X() int32 {
	// according to windowsx.h
	loword := uint16(l & 0xFFFF)
	short := int16(loword) // convert to signed...
	return int32(short)    // ...and sign extend
}

func (l _LPARAM) Y() int32 {
	// according to windowsx.h
	hiword := uint16((l & 0xFFFF0000) >> 16)
	short := int16(hiword) // convert to signed...
	return int32(short)    // ...and sign extend
}

type _POINT struct {
	x int32
	y int32
}

type _RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

// Go doesn't allow negative constants to be forced into unsigned types at compile-time; this will do it at runtime.
// This is safe; see http://stackoverflow.com/questions/24022225/what-are-the-sign-extension-rules-for-calling-windows-api-functions-stdcall-t
func negConst(c int) uintptr {
	return uintptr(c)
}

// the next two are convenience wrappers
// the intention is not to say utf16ToArg(toUTF16(s)) - even though it appears to work fine, that's just because the garbage collector scheduling makes it run long after we're finished; if we store these uintptrs globally instead, then things will break
// instead, call them separately - s := toUTF16(str); winapifunc.Call(utf16ToArg(s))

func toUTF16(s string) *uint16 {
	return syscall.StringToUTF16Ptr(s)
}

func utf16ToArg(s *uint16) uintptr {
	return uintptr(unsafe.Pointer(s))
}

func utf16ToLPARAM(s *uint16) uintptr {
	return uintptr(_LPARAM(unsafe.Pointer(s)))
}

var (
	_adjustWindowRectEx = user32.NewProc("AdjustWindowRectEx")
	_createWindowEx     = user32.NewProc("CreateWindowExW")
	_getClientRect      = user32.NewProc("GetClientRect")
	_moveWindow         = user32.NewProc("MoveWindow")
	_setWindowPos       = user32.NewProc("SetWindowPos")
	_setWindowText      = user32.NewProc("SetWindowTextW")
	_showWindow         = user32.NewProc("ShowWindow")
)

type _MINMAXINFO struct {
	ptReserved     _POINT
	ptMaxSize      _POINT
	ptMaxPosition  _POINT
	ptMinTrackSize _POINT
	ptMaxTrackSize _POINT
}

func (l _LPARAM) MINMAXINFO() *_MINMAXINFO {
	return (*_MINMAXINFO)(unsafe.Pointer(l))
}
