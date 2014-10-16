// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type textfield struct {
	*controlSingleHWNDWithText
	changed  *event
}

var editclass = toUTF16("EDIT")

func startNewTextField(style C.DWORD) *textfield {
	hwnd := C.newControl(editclass,
		style|C.textfieldStyle,
		C.textfieldExtStyle) // WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
	t := &textfield{
		controlSingleHWNDWithText:		newControlSingleHWNDWithText(hwnd),
		changed: newEvent(),
	}
	t.fpreferredSize = t.preferredSize
	C.controlSetControlFont(t.hwnd)
	C.setTextFieldSubclass(t.hwnd, unsafe.Pointer(t))
	return t
}

func newTextField() *textfield {
	return startNewTextField(0)
}

func newPasswordField() *textfield {
	return startNewTextField(C.ES_PASSWORD)
}

func (t *textfield) Text() string {
	return t.text()
}

func (t *textfield) SetText(text string) {
	t.setText(text)
}

func (t *textfield) OnChanged(f func()) {
	t.changed.set(f)
}

func (t *textfield) Invalid(reason string) {
	if reason == "" {
		C.textfieldHideInvalidBalloonTip(t.hwnd)
		return
	}
	C.textfieldSetAndShowInvalidBalloonTip(t.hwnd, toUTF16(reason))
}

//export textfieldChanged
func textfieldChanged(data unsafe.Pointer) {
	t := (*textfield)(data)
	t.changed.fire()
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	textfieldWidth  = 107 // this is actually the shorter progress bar width, but Microsoft only indicates as wide as necessary
	textfieldHeight = 14
)

func (t *textfield) preferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(textfieldWidth, d), fromdlgunitsY(textfieldHeight, d)
}
