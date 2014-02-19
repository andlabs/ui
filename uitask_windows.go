// 11 february 2014
package ui

import (
	"syscall"
	"unsafe"
	"runtime"
)

var uitask chan *uimsg

type uimsg struct {
	call		*syscall.LazyProc
	p		[]uintptr
	ret		chan uiret
}

type uiret struct {
	ret		uintptr
	err		error
}

func ui(initDone chan error) {
	runtime.LockOSThread()

	uitask = make(chan *uimsg)
	initDone <- doWindowsInit()

	quit := false
	for !quit {
		select {
		case m := <-uitask:
			r1, _, err := m.call.Call(m.p...)
			m.ret <- uiret{
				ret:	r1,
				err:	err,
			}
		default:
			quit = msgloopstep()
		}
	}
}

const (
	_PM_REMOVE = 0x0001
)

var (
	_dispatchMessage = user32.NewProc("DispatchMessageW")
	_getMessage = user32.NewProc("GetMessageW")
	_peekMessage = user32.NewProc("PeekMessageW")
	_postQuitMessage = user32.NewProc("PostQuitMessage")
	_sendMessage = user32.NewProc("SendMessageW")
	_translateMessage = user32.NewProc("TranslateMessage")
)

var getMessageFail = -1		// because Go doesn't let me

func msgloopstep() (quit bool) {
	var msg struct {
		Hwnd	_HWND
		Message	uint32
		WParam	_WPARAM
		LParam	_LPARAM
		Time		uint32
		Pt		_POINT
	}

	r1, _, _ := _peekMessage.Call(
		uintptr(unsafe.Pointer(&msg)),
		uintptr(_NULL),
		uintptr(0),
		uintptr(0),
		uintptr(_PM_REMOVE))
	if r1 == 0 {		// no message available
		return false
	}
	if msg.Message == _WM_QUIT {
		return true
	}
	_translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
	_dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	return false
}
