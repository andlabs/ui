// 8 february 2014
package main

import (
	"syscall"
	"unsafe"
)

type WNDCLASS struct {
	Style				uint32
	LpfnWndProc		WNDPROC
	CbClsExtra		int		// TODO exact Go type for C int? MSDN says C int
	CbWndExtra		int		// TODO exact Go type for C int? MSDN says C int
	HInstance			HANDLE	// actually HINSTANCE
	HIcon			HANDLE	// actually HICON
	HCursor			HANDLE	// actually HCURSOR
	HbrBackground	HBRUSH
	LpszMenuName	*string	// TODO this should probably just be a regular string with "" indicating no name but MSDN doesn't say if that's legal or not
	LpszClassName	string
}

type _WNDCLASSW struct {
	style				uint32
	lpfnWndProc		WNDPROC
	cbClsExtra		int
	cbWndExtra		int
	hInstance			HANDLE
	hIcon			HANDLE
	hCursor			HANDLE
	hbrBackground	HBRUSH
	lpszMenuName	*uint16
	lpszClassName		*uint16
}

func (w *WNDCLASS) toNative() *_WNDCLASSW {
	menuName := (*uint16)(nil)
	if w.LpszMenuName != nil {
		menuName = syscall.StringToUTF16Ptr(*w.LpszMenuName)
	}
	return &_WNDCLASSW{
		style:			w.Style,
		lpfnWndProc:		w.LpfnWndProc,
		cbClsExtra:		w.CbClsExtra,
		cbWndExtra:		w.CbWndExtra,
		hInstance:		w.HInstance,
		hIcon:			w.HIcon,
		hCursor:			w.HCursor,
		hbrBackground:	w.HbrBackground,
		lpszMenuName:	menuName,
		lpszClassName:	syscall.StringToUTF16Ptr(w.LpszClassName),
	}
}

var (
	registerClass = user32.NewProc("RegisterClassW")
)

func RegisterClass(lpWndClass *WNDCLASS) (class ATOM, err error) {
	r1, _, err := registerClass.Call(uintptr(unsafe.Pointer(lpWndClass.toNative())))
	if r1 == 0 {		// failure
		return 0, err
	}
	return ATOM(r1), nil
}
