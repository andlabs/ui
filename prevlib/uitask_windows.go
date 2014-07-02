// 11 february 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

/*
TODO rewrite this comment block

problem: messages have to be dispatched on the same thread as system calls, and we can't mux GetMessage() with select, and PeekMessage() every iteration is wasteful (and leads to lag for me (only) with the concurrent garbage collector sweep)
possible: solution: use PostThreadMessage() to send uimsgs out to the message loop, which runs on its own goroutine
(I had come up with this first but wanted to try other things before doing it (and wasn't really sure if user-defined messages were safe, not quite understanding the system); nsf came up with it independently and explained that this was really the only right way to do it, so thanks to him)

problem: if the thread isn't in its main message pump, the thread message is simply lost (see, for example, http://blogs.msdn.com/b/oldnewthing/archive/2005/04/26/412116.aspx)
this happened when scrolling Areas (as scrolling is modal; see http://blogs.msdn.com/b/oldnewthing/archive/2005/04/27/412565.aspx)

the only recourse, and the one both Microsoft (http://support.microsoft.com/kb/183116) and Raymond Chen (http://blogs.msdn.com/b/oldnewthing/archive/2008/12/23/9248851.aspx) suggest (and Treeki/Ninjifox confirmed), is to create an invisible window to dispatch messages instead.

yay.
*/

var msghwnd _HWND

const (
	msgQuit = _WM_APP + iota + 1			// + 1 just to be safe
	msgSetAreaSize
	msgRepaintAll
	msgCreateWindow
)

type uitaskParams struct {
	window	*Window		// createWindow
	control	Control		// createWindow
	show	bool			// createWindow
}

// SendMessage() won't return unti lthe deed is done, even if the deed is on another thread
// SendMessage() does a thread switch if necessary
// this also means we don't have to worry about the uitaskParams object being garbage collected

func (_uitask) createWindow(w *Window, c Control, s bool) {
	uc := &uitaskParams{
		window:	w,
		control:	c,
		show:	s,
	}
	_sendMessage.Call(
		uintptr(msghwnd),
		msgCreateWindow,
		uintptr(0),
		uintptr(unsafe.Pointer(uc)))
}

func uiinit() error {
	err := doWindowsInit()
	if err != nil {
		return fmt.Errorf("error doing general Windows initialization: %v", err)
	}

	msghwnd, err = makeMessageHandler()
	if err != nil {
		return fmt.Errorf("error making invisible window for handling events: %v", err)
	}

	return nil
}

var (
	_postMessage = user32.NewProc("PostMessageW")
)

func ui() {
	go func() {
		<-Stop
		// PostMessage() so it gets handled after any events currently being processed complete
		r1, _, err := _postMessage.Call(
			uintptr(msghwnd),
			msgQuit,
			uintptr(0),
			uintptr(0))
		if r1 == 0 { // failure
			panic("error sending quit message to message loop: " + err.Error())
		}
	}()

	msgloop()
}

var (
	_dispatchMessage  = user32.NewProc("DispatchMessageW")
	_getActiveWindow		= user32.NewProc("GetActiveWindow")
	_getMessage       = user32.NewProc("GetMessageW")
	_isDialogMessage		= user32.NewProc("IsDialogMessageW")
	_postQuitMessage  = user32.NewProc("PostQuitMessage")
	_sendMessage      = user32.NewProc("SendMessageW")
	_translateMessage = user32.NewProc("TranslateMessage")
)

func msgloop() {
	var msg struct {
		hwnd    _HWND
		message uint32
		wParam  _WPARAM
		lParam  _LPARAM
		time    uint32
		pt      _POINT
	}

	for {
		r1, _, err := _getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			uintptr(_NULL),
			uintptr(0),
			uintptr(0))
		if r1 == negConst(-1) { // error
			panic("error getting message in message loop: " + err.Error())
		}
		if r1 == 0 { // WM_QUIT message
			return
		}
		// this next bit handles tab stops
		r1, _, _ = _getActiveWindow.Call()
		r1, _, _ = _isDialogMessage.Call(
			r1,		// active window
			uintptr(unsafe.Pointer(&msg)))
		if r1 != 0 {
			continue
		}
		_translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		_dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

var (
	msghandlerclass = toUTF16("gomsghandler")
	msghandlertitle = toUTF16("ui package message window")
)

func makeMessageHandler() (hwnd _HWND, err error) {
	wc := &_WNDCLASS{
		lpszClassName: utf16ToArg(msghandlerclass),
		lpfnWndProc:   syscall.NewCallback(messageHandlerWndProc),
		hInstance:     hInstance,
		hIcon:         icon,
		hCursor:       cursor,
		hbrBackground: _HBRUSH(_COLOR_BTNFACE + 1),
	}

	r1, _, err := _registerClass.Call(uintptr(unsafe.Pointer(wc)))
	if r1 == 0 { // failure
		return _HWND(_NULL), fmt.Errorf("error registering the class of the invisible window for handling events: %v", err)
	}

	r1, _, err = _createWindowEx.Call(
		uintptr(0),
		utf16ToArg(msghandlerclass),
		utf16ToArg(msghandlertitle),
		uintptr(0),
		negConst(_CW_USEDEFAULT),
		negConst(_CW_USEDEFAULT),
		negConst(_CW_USEDEFAULT),
		negConst(_CW_USEDEFAULT),
		// don't negConst() HWND_MESSAGE; windowsconstgen was given a pointer by windows.h, and pointers are unsigned, so converting it back to signed doesn't work
		uintptr(_HWND_MESSAGE),
		uintptr(_NULL),
		uintptr(hInstance),
		uintptr(_NULL))
	if r1 == 0 { // failure
		return _HWND(_NULL), fmt.Errorf("error actually creating invisible window for handling events: %v", err)
	}

	return _HWND(r1), nil
}

func messageHandlerWndProc(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	switch uMsg {
	case msgQuit:
		// does not return a value according to MSDN
		_postQuitMessage.Call(0)
		return 0
	case msgCreateWindow:
		uc := (*uitaskParams)(unsafe.Pointer(lParam))
		uc.window.create(uc.control, uc.show)
		return 0
	}
	return defWindowProc(hwnd, uMsg, wParam, lParam)
}
