// 9 february 2014

//
package ui

import (
	//	"syscall"
	"unsafe"
)

// SendMessage constants.
const (
	HWND_BROADCAST = HWND(0xFFFF)
)

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  WPARAM
	LParam  LPARAM
	Time    uint32
	Pt      POINT
}

var (
	dispatchMessage  = user32.NewProc("DispatchMessageW")
	getMessage       = user32.NewProc("GetMessageW")
	postQuitMessage  = user32.NewProc("PostQuitMessage")
	sendMessage      = user32.NewProc("SendMessageW")
	translateMessage = user32.NewProc("TranslateMessage")
)

// TODO handle errors
func DispatchMessage(lpmsg *MSG) (result LRESULT, err error) {
	r1, _, _ := dispatchMessage.Call(uintptr(unsafe.Pointer(lpmsg)))
	return LRESULT(r1), nil
}

var getMessageFail = -1 // because Go doesn't let me

func GetMessage(hWnd HWND, wMsgFilterMin uint32, wMsgFilterMax uint32) (lpMsg *MSG, quit bool, err error) {
	lpMsg = new(MSG)
	r1, _, err := getMessage.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr(hWnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax))
	if r1 == uintptr(getMessageFail) { // failure
		return nil, false, err
	}
	return lpMsg, r1 == 0, nil
}

// TODO handle errors
func PostQuitMessage(nExitCode int) (err error) {
	postQuitMessage.Call(uintptr(nExitCode))
	return nil
}

// TODO handle errors
func SendMessage(hWnd HWND, Msg uint32, wParam WPARAM, lParam LPARAM) (result LRESULT, err error) {
	r1, _, _ := sendMessage.Call(
		uintptr(hWnd),
		uintptr(Msg),
		uintptr(wParam),
		uintptr(lParam))
	return LRESULT(r1), nil
}

// TODO handle errors
func TranslateMessage(lpMsg *MSG) (translated bool, err error) {
	r1, _, _ := translateMessage.Call(uintptr(unsafe.Pointer(lpMsg)))
	return r1 != 0, nil
}
