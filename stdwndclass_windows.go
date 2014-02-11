// 8 february 2014
package main

import (
	"fmt"
	"syscall"
	"unsafe"
	"sync"
)

const (
	stdWndClass = "gouiwndclass"
)

var (
	defWindowProc = user32.NewProc("DefWindowProcW")
)

func stdWndProc(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	sysData := getSysData(hwnd)
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
			sysData.closing <- struct{}
		}
		return 0
	default:
		r1, _, _ := defWindowProc.Call(
			uintptr(hwnd),
			uintptr(uMsg),
			uintptr(wParam),
			uintptr(lParam))
		return LRESULT(r1)
	}
	panic(fmt.Sprintf("stdWndProc message %d did not return: internal bug in ui library", uMsg))
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
	lpszClassName		*uint16
}

func registerStdWndClass() (err error) {
	const (
		_IDI_APPLICATION = 32512
		_IDC_ARROW = 32512
	)

	icon, err := user32.NewProc("LoadIconW").Call(
		uintptr(_NULL),
		uintptr(_IDI_APPLICATION))
	if err != nil {
		return fmt.Errorf("error getting window icon: %v", err)
	}
	cursor, err := user32.NewProc("LoadCursorW").Call(
		uintptr(_NULL),
		uintptr(_IDC_ARROW))
	if err != nil {
		return fmt.Errorf("error getting window cursor: %v", err)
	}

	wc := &_WNDCLASS{
		lpszClassName:	syscall.StringToUTF16Ptr(stdWndClass),
		lpfnWndProc:		syscall.NewCallback(stdWndProc),
		hInstance:		hInstance,
		hIcon:			icon,
		hCursor:			cursor,
		hbrBackground:	_HBRUSH(_COLOR_BTNFACE + 1),
	}

	r1, _, err := user32.NewProc("RegisterClassW").Call(uintptr(unsafe.Pointer(wc)))
	if r1 == 0 {		// failure
		return fmt.Errorf("error registering class: %v", err)
	}
	return nil
}
