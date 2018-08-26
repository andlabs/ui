// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Button is a Control that represents a button that the user can
// click to perform an action. A Button has a text label that should
// describe what the button does.
type Button struct {
	ControlBase
	b	*C.uiButton
	onClicked		func(*Button)
}

// NewButton creates a new Button with the given text as its label.
func NewButton(text string) *Button {
	b := new(Button)

	ctext := C.CString(text)
	b.b = C.uiNewButton(ctext)
	freestr(ctext)

	C.pkguiButtonOnClicked(b.b)

	b.ControlBase = NewControlBase(b, uintptr(unsafe.Pointer(b.b)))
	return b
}

// Text returns the Button's text.
func (b *Button) Text() string {
	ctext := C.uiButtonText(b.b)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the Button's text to text.
func (b *Button) SetText(text string) {
	ctext := C.CString(text)
	C.uiButtonSetText(b.b, ctext)
	freestr(ctext)
}

// OnClicked registers f to be run when the user clicks the Button.
// Only one function can be registered at a time.
func (b *Button) OnClicked(f func(*Button)) {
	b.onClicked = f
}

//export pkguiDoButtonOnClicked
func pkguiDoButtonOnClicked(bb *C.uiButton, data unsafe.Pointer) {
	b := ControlFromLibui(uintptr(unsafe.Pointer(bb))).(*Button)
	if b.onClicked != nil {
		b.onClicked(b)
	}
}
