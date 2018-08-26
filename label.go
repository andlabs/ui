// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Label is a Control that represents a line of text that cannot be
// interacted with.
type Label struct {
	ControlBase
	l	*C.uiLabel
}

// NewLabel creates a new Label with the given text.
func NewLabel(text string) *Label {
	l := new(Label)

	ctext := C.CString(text)
	l.l = C.uiNewLabel(ctext)
	freestr(ctext)

	l.ControlBase = NewControlBase(l, uintptr(unsafe.Pointer(l.l)))
	return l
}

// Text returns the Label's text.
func (l *Label) Text() string {
	ctext := C.uiLabelText(l.l)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the Label's text to text.
func (l *Label) SetText(text string) {
	ctext := C.CString(text)
	C.uiLabelSetText(l.l, ctext)
	freestr(ctext)
}
