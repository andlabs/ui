// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// Label is a Control that represents a line of text that cannot be
// interacted with. TODO rest of documentation.
type Label struct {
	c	*C.uiControl
	l	*C.uiLabel
}

// NewLabel creates a new Label with the given text.
func NewLabel(text string) *Label {
	l := new(Label)

	ctext := C.CString(text)
	l.l = C.uiNewLabel(ctext)
	l.c = (*C.uiControl)(unsafe.Pointer(l.l))
	freestr(ctext)

	return l
}

// Destroy destroys the Label.
func (l *Label) Destroy() {
	C.uiControlDestroy(l.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (l *Label) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(l.c))
}

// Handle returns the OS-level handle associated with this Label.
// On Windows this is an HWND of a standard Windows API STATIC
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkLabel.
// On OS X this is a pointer to a NSTextField.
func (l *Label) Handle() uintptr {
	return uintptr(C.uiControlHandle(l.c))
}

// Show shows the Label.
func (l *Label) Show() {
	C.uiControlShow(l.c)
}

// Hide hides the Label.
func (l *Label) Hide() {
	C.uiControlHide(l.c)
}

// Enable enables the Label.
func (l *Label) Enable() {
	C.uiControlEnable(l.c)
}

// Disable disables the Label.
func (l *Label) Disable() {
	C.uiControlDisable(l.c)
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
