// 11 february 2014
package ui

import (
	"syscall"
	"unsafe"
	"runtime"
)

/*
problem: messages have to be dispatched on the same thread as system calls, and we can't mux GetMessage() with select, and PeekMessage() every iteration is wasteful (and leads to lag for me (only) with the concurrent garbage collector sweep)
solution: use PostThreadMessage() to send uimsgs out to the message loop, which runs on its own goroutine
I had come up with this first but wanted to try other things before doing it (and wasn't really sure if user-defined messages were safe, not quite understanding the system); nsf came up with it independently and explained that this was really the only right way to do it, so thanks to him
*/

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

const (
	_WM_APP = 0x8000 + iota
	msgRequested
)

var (
	_postThreadMessage = user32.NewProc("PostThreadMessageW")
)

func ui(initDone chan error) {
	runtime.LockOSThread()

	uitask = make(chan *uimsg)
	initDone <- doWindowsInit()

	threadIDReq := make(chan uintptr)
	msglooperrs := make(chan error)
	go msgloop(threadIDReq, msglooperrs)
	threadID := <-threadIDReq

	quit := false
	for !quit {
		select {
		case m := <-uitask:
			r1, _, err := _postThreadMessage.Call(
				threadID,
				msgRequested,
				uintptr(0),
				uintptr(unsafe.Pointer(m)))
			if r1 == 0 {		// failure
				panic("error sending message to message loop to call function: " + err.Error())		// TODO
			}
		case err := <-msglooperrs:
			if err == nil {		// WM_QUIT; no error
				quit = true
			} else {
				panic("unexpected return from message loop: " + err.Error())		// TODO
			}
		}
	}
}

var (
	_dispatchMessage = user32.NewProc("DispatchMessageW")
	_getMessage = user32.NewProc("GetMessageW")
	_getCurrentThreadID = kernel32.NewProc("GetCurrentThreadId")
	_postQuitMessage = user32.NewProc("PostQuitMessage")
	_sendMessage = user32.NewProc("SendMessageW")
	_translateMessage = user32.NewProc("TranslateMessage")
)

var getMessageFail = -1		// because Go doesn't let me

func msgloop(threadID chan uintptr, errors chan error) {
	runtime.LockOSThread()

	var msg struct {
		Hwnd	_HWND
		Message	uint32
		WParam	_WPARAM
		LParam	_LPARAM
		Time		uint32
		Pt		_POINT
	}

	r1, _, _ := _getCurrentThreadID.Call()
	threadID <- r1
	for {
		r1, _, err := _getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			uintptr(_NULL),
			uintptr(0),
			uintptr(0))
		if r1 == uintptr(getMessageFail) {		// error
			errors <- err
			return
		}
		if r1 == 0 {		// WM_QUIT message
			errors <- nil
			return
		}
		if msg.Message == msgRequested {
			m := (*uimsg)(unsafe.Pointer(msg.LParam))
			r1, _, err := m.call.Call(m.p...)
			m.ret <- uiret{
				ret:	r1,
				err:	err,
			}
			continue
		}
		_translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		_dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
