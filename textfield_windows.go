// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type textfield struct {
	_hwnd    C.HWND
	_textlen C.LONG
	changed  *event
}

var editclass = toUTF16("EDIT")

func startNewTextField(style C.DWORD) *textfield {
	hwnd := C.newControl(editclass,
		style|C.textfieldStyle,
		C.textfieldExtStyle) // WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
	t := &textfield{
		_hwnd:   hwnd,
		changed: newEvent(),
	}
	C.controlSetControlFont(t._hwnd)
	C.setTextFieldSubclass(t._hwnd, unsafe.Pointer(t))
	return t
}

func newTextField() *textfield {
	return startNewTextField(0)
}

func newPasswordField() *textfield {
	return startNewTextField(C.ES_PASSWORD)
}

func (t *textfield) Text() string {
	return baseText(t)
}

func (t *textfield) SetText(text string) {
	baseSetText(t, text)
}

func (t *textfield) OnChanged(f func()) {
	t.changed.set(f)
}

func (t *textfield) Invalid(reason string) {
	if reason == "" {
		C.textfieldHideInvalidBalloonTip(t._hwnd)
		return
	}
	C.textfieldSetAndShowInvalidBalloonTip(t._hwnd, toUTF16(reason))
}

//export textfieldChanged
func textfieldChanged(data unsafe.Pointer) {
	t := (*textfield)(data)
	t.changed.fire()
}

func (t *textfield) hwnd() C.HWND {
	return t._hwnd
}

func (t *textfield) textlen() C.LONG {
	return t._textlen
}

func (t *textfield) settextlen(len C.LONG) {
	t._textlen = len
}

func (t *textfield) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *textfield) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	textfieldWidth  = 107 // this is actually the shorter progress bar width, but Microsoft only indicates as wide as necessary
	textfieldHeight = 14
)

func (t *textfield) preferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(textfieldWidth, d), fromdlgunitsY(textfieldHeight, d)
}

func (t *textfield) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *textfield) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
