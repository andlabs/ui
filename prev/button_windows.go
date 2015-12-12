// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type button struct {
	*controlSingleHWNDWithText
	clicked  *event
}

var buttonclass = toUTF16("BUTTON")

func newButton(text string) *button {
	hwnd := C.newControl(buttonclass,
		C.BS_PUSHBUTTON|C.WS_TABSTOP,
		0)
	b := &button{
		controlSingleHWNDWithText:		newControlSingleHWNDWithText(hwnd),
		clicked: newEvent(),
	}
	b.fpreferredSize = b.xpreferredSize
	b.SetText(text)
	C.controlSetControlFont(b.hwnd)
	C.setButtonSubclass(b.hwnd, unsafe.Pointer(b))
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return b.text()
}

func (b *button) SetText(text string) {
	b.setText(text)
}

//export buttonClicked
func buttonClicked(data unsafe.Pointer) {
	b := (*button)(data)
	b.clicked.fire()
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	buttonHeight = 14
)

func (b *button) xpreferredSize(d *sizing) (width, height int) {
	// comctl32.dll version 6 thankfully provides a method to grab this...
	var size C.SIZE

	size.cx = 0 // explicitly ask for ideal size
	size.cy = 0
	if C.SendMessageW(b.hwnd, C.BCM_GETIDEALSIZE, 0, C.LPARAM(uintptr(unsafe.Pointer(&size)))) != C.FALSE {
		return int(size.cx), int(size.cy)
	}
	// that failed, fall back
	println("message failed; falling back")
	// don't worry about the error return from GetSystemMetrics(); there's no way to tell (explicitly documented as such)
	xmargins := 2 * int(C.GetSystemMetrics(C.SM_CXEDGE))
	return xmargins + int(b.textlen), fromdlgunitsY(buttonHeight, d)
}
