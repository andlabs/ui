// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type widgetbase struct {
	hwnd	C.HWND
}

func newWidget(class C.LPCWSTR, style C.DWORD, extstyle C.DWORD) *widgetbase {
	return &widgetbase{
		hwnd:	C.newWidget(class, style, extstyle),
	}
}

// these few methods are embedded by all the various Controls since they all will do the same thing

func (w *widgetbase) unparent() {
	C.controlSetParent(w.hwnd, C.msgwin)
}

func (w *widgetbase) parent(win *window) {
	C.controlSetParent(w.hwnd, win.hwnd)
}

// don't embed these as exported; let each Control decide if it should

func (w *widgetbase) text() string {
	return getWindowText(w.hwnd)
}

func (w *widgetbase) settext(text string) {
	C.setWindowText(w.hwnd, toUTF16(text))
}

type button struct {
	*widgetbase
	clicked		*event
}

var buttonclass = toUTF16("BUTTON")

func newButton(text string) *button {
	w := newWidget(buttonclass,
		C.BS_PUSHBUTTON | C.WS_TABSTOP,
		0)
	C.setWindowText(w.hwnd, toUTF16(text))
	C.controlSetControlFont(w.hwnd)
	b := &button{
		widgetbase:	w,
		clicked:		newEvent(),
	}
	C.setButtonSubclass(w.hwnd, unsafe.Pointer(b))
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return b.text()
}

func (b *button) SetText(text string) {
	b.settext(text)
}

//export buttonClicked
func buttonClicked(data unsafe.Pointer) {
	b := (*button)(data)
	b.clicked.fire()
	println("button clicked")
}
