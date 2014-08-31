// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type button struct {
	_hwnd	C.HWND
	_textlen	C.LONG
	clicked	*event
}

var buttonclass = toUTF16("BUTTON")

func newButton(text string) *button {
	hwnd := C.newControl(buttonclass,
		C.BS_PUSHBUTTON | C.WS_TABSTOP,
		0)
	b := &button{
		_hwnd:	hwnd,
		clicked:	newEvent(),
	}
	b.SetText(text)
	C.controlSetControlFont(b._hwnd)
	C.setButtonSubclass(b._hwnd, unsafe.Pointer(b))
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return baseText(b)
}

func (b *button) SetText(text string) {
	baseSetText(b, text)
}

//export buttonClicked
func buttonClicked(data unsafe.Pointer) {
	b := (*button)(data)
	b.clicked.fire()
	println("button clicked")
}

func (b *button) hwnd() C.HWND {
	return b._hwnd
}

func (b *button) textlen() C.LONG {
	return b._textlen
}

func (b *button) settextlen(len C.LONG) {
	b._textlen = len
}

func (b *button) setParent(p *controlParent) {
	basesetParent(b, p)
}

func (b *button) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(b, x, y, width, height, d)
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	buttonHeight = 14
)

func (b *button) preferredSize(d *sizing) (width, height int) {
	// comctl32.dll version 6 thankfully provides a method to grab this...
	var size C.SIZE

	size.cx = 0		// explicitly ask for ideal size
	size.cy = 0
	if C.SendMessageW(b._hwnd, C.BCM_GETIDEALSIZE, 0, C.LPARAM(uintptr(unsafe.Pointer(&size)))) != C.FALSE {
		return int(size.cx), int(size.cy)
	}
	// that failed, fall back
println("message failed; falling back")
	// don't worry about the error return from GetSystemMetrics(); there's no way to tell (explicitly documented as such)
	xmargins := 2 * int(C.GetSystemMetrics(C.SM_CXEDGE))
	return xmargins + int(b._textlen), fromdlgunitsY(buttonHeight, d)
}

func (b *button) commitResize(a *allocation, d *sizing) {
	basecommitResize(b, a, d)
}

func (b *button) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(b, d)
}
