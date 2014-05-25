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

// Predefined cursor resource IDs.
const (
	_IDC_APPSTARTING = 32650
	_IDC_ARROW = 32512
	_IDC_CROSS = 32515
	_IDC_HAND = 32649
	_IDC_HELP = 32651
	_IDC_IBEAM = 32513
//	_IDC_ICON = 32641		// [Obsolete for applications marked version 4.0 or later.]
	_IDC_NO = 32648
//	_IDC_SIZE = 32640		// [Obsolete for applications marked version 4.0 or later. Use IDC_SIZEALL.]
	_IDC_SIZEALL = 32646
	_IDC_SIZENESW = 32643
	_IDC_SIZENS = 32645
	_IDC_SIZENWSE = 32642
	_IDC_SIZEWE = 32644
	_IDC_UPARROW = 32516
	_IDC_WAIT = 32514
)

// Predefined icon resource IDs.
const (
	_IDI_APPLICATION = 32512
	_IDI_ASTERISK = 32516
	_IDI_ERROR = 32513
	_IDI_EXCLAMATION = 32515
	_IDI_HAND = 32513
	_IDI_INFORMATION = 32516
	_IDI_QUESTION = 32514
	_IDI_SHIELD = 32518
	_IDI_WARNING = 32515
	_IDI_WINLOGO = 32517
)
