// 24 march 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
	"sync"
)

const (
	areastyle = 0 | controlstyle
	areaxstyle = 0 | controlxstyle
)

const (
	areaWndClassFormat = "gouiarea%X"
)

var (
	areaWndClassNum uintptr
	areaWndClassNumLock sync.Mutex
)

func areaWndProc(s *sysData) func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	return func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
		switch uMsg {
		default:
			r1, _, _ := defWindowProc.Call(
				uintptr(hwnd),
				uintptr(uMsg),
				uintptr(wParam),
				uintptr(lParam))
			return _LRESULT(r1)
		}
		panic(fmt.Sprintf("areaWndProc message %d did not return: internal bug in ui library", uMsg))
	}
}

func registerAreaWndClass(s *sysData) (newClassName string, err error) {
	const (
		// from winuser.h
		_CS_DBLCLKS = 0x0008
	)

	areaWndClassNumLock.Lock()
	newClassName = fmt.Sprintf(areaWndClassFormat, areaWndClassNum)
	areaWndClassNum++
	areaWndClassNumLock.Unlock()

	wc := &_WNDCLASS{
		style:			_CS_DBLCLKS,		// needed to be able to register double-clicks
		lpszClassName:	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(newClassName))),
		lpfnWndProc:		syscall.NewCallback(areaWndProc(s)),
		hInstance:		hInstance,
		hIcon:			icon,
		hCursor:			cursor,
		hbrBackground:	_HBRUSH(_COLOR_BTNFACE + 1),
	}

	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_registerClass,
		p:		[]uintptr{uintptr(unsafe.Pointer(wc))},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		return "", r.err
	}
	return newClassName, nil
}
