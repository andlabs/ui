// 11 february 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
	"runtime"
)

/*
problem: messages have to be dispatched on the same thread as system calls, and we can't mux GetMessage() with select, and PeekMessage() every iteration is wasteful (and leads to lag for me (only) with the concurrent garbage collector sweep)
possible: solution: use PostThreadMessage() to send uimsgs out to the message loop, which runs on its own goroutine
(I had come up with this first but wanted to try other things before doing it (and wasn't really sure if user-defined messages were safe, not quite understanding the system); nsf came up with it independently and explained that this was really the only right way to do it, so thanks to him)

problem: if the thread isn't in its main message pump, the thread message is simply lost (see, for example, http://blogs.msdn.com/b/oldnewthing/archive/2005/04/26/412116.aspx)
this happened when scrolling Areas (as scrolling is modal; see http://blogs.msdn.com/b/oldnewthing/archive/2005/04/27/412565.aspx)

the only recourse, and the one both Microsoft (http://support.microsoft.com/kb/183116) and Raymond Chen (http://blogs.msdn.com/b/oldnewthing/archive/2008/12/23/9248851.aspx) suggest (and Treeki/Ninjifox confirmed), is to create an invisible window to dispatch messages instead.

yay.
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
	msgQuit
	msgSetAreaSize
)

var (
	_postMessage = user32.NewProc("PostMessageW")
)

func ui(main func()) error {
	runtime.LockOSThread()

	uitask = make(chan *uimsg)
	err := doWindowsInit()
	if err != nil {
		return fmt.Errorf("error doing general Windows initialization: %v", err)
	}

	hwnd, err := makeMessageHandler()
	if err != nil {
		return fmt.Errorf("error making invisible window for handling events: %v", err)
	}

	go func() {
		for m := range uitask {
			// TODO use _sendMessage instead?
			r1, _, err := _postMessage.Call(
				uintptr(hwnd),
				msgRequested,
				uintptr(0),
				uintptr(unsafe.Pointer(m)))
			if r1 == 0 {		// failure
				panic("error sending message to message loop to call function: " + err.Error())
			}
		}
	}()

	go func() {
		main()
		r1, _, err := _postMessage.Call(
			uintptr(hwnd),
			msgQuit,
			uintptr(0),
			uintptr(0))
		if r1 == 0 {		// failure
			panic("error sending quit message to message loop: " + err.Error())
		}
	}()

	msgloop()
	return nil
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
	var msg struct {
		hwnd	_HWND
		message	uint32
		wParam	_WPARAM
		lParam	_LPARAM
		time		uint32
		pt		_POINT
	}

	for {
		r1, _, err := _getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			uintptr(_NULL),
			uintptr(0),
			uintptr(0))
		if r1 == uintptr(getMessageFail) {		// error
			panic("error getting message in message loop: " + err.Error())
		}
		if r1 == 0 {		// WM_QUIT message
			return
		}
		_translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		_dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// TODO move to init?

const (
	msghandlerclass = "gomsghandler"
)

var (
	// fron winuser.h; var because Go won't let me
	_HWND_MESSAGE = -3
)

func makeMessageHandler() (hwnd _HWND, err error) {
	wc := &_WNDCLASS{
		lpszClassName:	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msghandlerclass))),
		lpfnWndProc:		syscall.NewCallback(messageHandlerWndProc),
		hInstance:		hInstance,
		hIcon:			icon,
		hCursor:			cursor,
		hbrBackground:	_HBRUSH(_COLOR_BTNFACE + 1),
	}

	r1, _, err := _registerClass.Call(uintptr(unsafe.Pointer(wc)))
	if r1 == 0 {		// failure
		return _HWND(_NULL), fmt.Errorf("error registering the class of the invisible window for handling events: %v", err)
	}

	r1, _, err = _createWindowEx.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msghandlerclass))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("ui package message window"))),
		uintptr(0),
		uintptr(_CW_USEDEFAULT),
		uintptr(_CW_USEDEFAULT),
		uintptr(_CW_USEDEFAULT),
		uintptr(_CW_USEDEFAULT),
		uintptr(_HWND_MESSAGE),
		uintptr(_NULL),
		uintptr(hInstance),
		uintptr(_NULL))
	if r1 == 0 {		// failure
		return _HWND(_NULL), fmt.Errorf("error actually creating invisible window for handling events: %v", err)
	}

	return _HWND(r1), nil
}

func messageHandlerWndProc(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	switch uMsg {
	case msgRequested:
		m := (*uimsg)(unsafe.Pointer(lParam))
		r1, _, err := m.call.Call(m.p...)
		m.ret <- uiret{
			ret:	r1,
			err:	err,
		}
		return 0
	case msgQuit:
		// does not return a value according to MSDN
		_postQuitMessage.Call(0)
		return 0
	}
	r1, _, _ := defWindowProc.Call(
		uintptr(hwnd),
		uintptr(uMsg),
		uintptr(wParam),
		uintptr(lParam))
	return _LRESULT(r1)
}
