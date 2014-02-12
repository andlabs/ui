// 8 february 2014
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	stdWndClass = "gouiwndclass"
)

var (
	defWindowProc = user32.NewProc("DefWindowProcW")
)

func stdWndProc(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	sysData := getSysData(hwnd)
	if sysData == nil {	// not ready for events yet
		goto defwndproc
	}
	switch uMsg {
	case _WM_COMMAND:
		id := wParam.LOWORD()
		// ... member events
		_ = id
		return 0
	case _WM_GETMINMAXINFO:
		mm := lParam.MINMAXINFO()
		// ... minimum size
		_ = mm
		return 0
	case _WM_SIZE:
		// TODO
		return 0
	case _WM_CLOSE:
		if sysData.closing != nil {
			sysData.closing <- struct{}{}
		}
		return 0
	default:
		goto defwndproc
	}
	panic(fmt.Sprintf("stdWndProc message %d did not return: internal bug in ui library", uMsg))
defwndproc:
	r1, _, _ := defWindowProc.Call(
		uintptr(hwnd),
		uintptr(uMsg),
		uintptr(wParam),
		uintptr(lParam))
	return _LRESULT(r1)
}

type _WNDCLASS struct {
	style				uint32
	lpfnWndProc		uintptr
	cbClsExtra		int
	cbWndExtra		int
	hInstance			_HANDLE
	hIcon			_HANDLE
	hCursor			_HANDLE
	hbrBackground	_HBRUSH
	lpszMenuName	*uint16
	lpszClassName		uintptr
}

func registerStdWndClass() (err error) {
	const (
		_IDI_APPLICATION = 32512
		_IDC_ARROW = 32512
	)

	r1, _, err := user32.NewProc("LoadIconW").Call(
		uintptr(_NULL),
		uintptr(_IDI_APPLICATION))
	if r1 == 0 {		// failure
		return fmt.Errorf("error getting window icon: %v", err)
	}
	icon := _HANDLE(r1)

	r1, _, err = user32.NewProc("LoadCursorW").Call(
		uintptr(_NULL),
		uintptr(_IDC_ARROW))
	if r1 == 0 {		// failure
		return fmt.Errorf("error getting window cursor: %v", err)
	}
	cursor := _HANDLE(r1)

	wc := &_WNDCLASS{
		lpszClassName:	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(stdWndClass))),
		lpfnWndProc:		syscall.NewCallback(stdWndProc),
		hInstance:		hInstance,
		hIcon:			icon,
		hCursor:			cursor,
		hbrBackground:	_HBRUSH(_COLOR_BTNFACE + 1),
	}

	r1, _, err = user32.NewProc("RegisterClassW").Call(uintptr(unsafe.Pointer(wc)))
	if r1 == 0 {		// failure
		return fmt.Errorf("error registering class: %v", err)
	}
	return nil
}
