// 11 february 2014
//package ui
package main

import (
	"fmt"
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

//	go msgloop()
	var msg struct {
		Hwnd	_HWND
		Message	uint32
		WParam	_WPARAM
		LParam	_LPARAM
		Time		uint32
		Pt		_POINT
	}
	var _peekMessage = user32.NewProc("PeekMessageW")
	const _PM_REMOVE = 0x0001

	for {
		select {
		case m := <-uitask:
			r1, _, err := m.call.Call(m.p...)
			m.ret <- uiret{
				ret:	r1,
				err:	err,
			}
		default:
			// TODO figure out how to handle errors
			_peekMessage.Call(
				uintptr(unsafe.Pointer(&msg)),
				uintptr(_NULL),
				0,
				0,
				uintptr(_PM_REMOVE))
		}
	}
}

var (
	_dispatchMessage = user32.NewProc("DispatchMessageW")
	_getMessage = user32.NewProc("GetMessageW")
	_postQuitMessage = user32.NewProc("PostQuitMessage")
	_sendMessage = user32.NewProc("SendMessageW")
	_translateMessage = user32.NewProc("TranslateMessage")
)

var getMessageFail = -1		// because Go doesn't let me

func msgloop() {
	runtime.LockOSThread()

	var msg struct {
		Hwnd	_HWND
		Message	uint32
		WParam	_WPARAM
		LParam	_LPARAM
		Time		uint32
		Pt		_POINT
	}

	for {
		r1, _, err := _getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			uintptr(_NULL),
			uintptr(0),
			uintptr(0))
		if r1 == uintptr(getMessageFail) {		// failure
			panic(fmt.Sprintf("GetMessage failed: %v", err))
		} else if r1 == 0 {	// quit
			break
		}
		// TODO handle potential errors in TranslateMessage() and DispatchMessage()
		_translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		_dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
