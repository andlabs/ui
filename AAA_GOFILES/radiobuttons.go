// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// RadioButtons is a Control that represents a set of checkable
// buttons from which exactly one may be chosen by the user.
type RadioButtons struct {
	c	*C.uiControl
	r	*C.uiRadioButtons
}

// NewRadioButtons creates a new RadioButtons.
func NewRadioButtons() *RadioButtons {
	r := new(RadioButtons)

	r.r = C.uiNewRadioButtons()
	r.c = (*C.uiControl)(unsafe.Pointer(r.r))

	return r
}

// Destroy destroys the RadioButtons.
func (r *RadioButtons) Destroy() {
	C.uiControlDestroy(r.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (r *RadioButtons) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(r.c))
}

// Handle returns the OS-level handle associated with this RadioButtons.
// On Windows this is an HWND of a libui-internal class; its
// child windows are instances of the standard Windows API
// BUTTON class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkBox containing GtkRadioButtons.
// On OS X this is a pointer to a NSView with each radio button as a NSButton subview.
func (r *RadioButtons) Handle() uintptr {
	return uintptr(C.uiControlHandle(r.c))
}

// Show shows the RadioButtons.
func (r *RadioButtons) Show() {
	C.uiControlShow(r.c)
}

// Hide hides the RadioButtons.
func (r *RadioButtons) Hide() {
	C.uiControlHide(r.c)
}

// Enable enables the RadioButtons.
func (r *RadioButtons) Enable() {
	C.uiControlEnable(r.c)
}

// Disable disables the RadioButtons.
func (r *RadioButtons) Disable() {
	C.uiControlDisable(r.c)
}

// Append adds the named button to the end of the RadioButtons.
// If this button is the first button, it is automatically selected.
func (r *RadioButtons) Append(text string) {
	ctext := C.CString(text)
	C.uiRadioButtonsAppend(r.r, ctext)
	freestr(ctext)
}
