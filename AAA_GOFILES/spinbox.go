// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doSpinboxOnChanged(uiSpinbox *, void *);
// static inline void realuiSpinboxOnChanged(uiSpinbox *b)
// {
// 	uiSpinboxOnChanged(b, doSpinboxOnChanged, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var spinboxes = make(map[*C.uiSpinbox]*Spinbox)

// Spinbox is a Control that represents a space where the user can
// enter integers. The space also comes with buttons to add or
// subtract 1 from the integer.
type Spinbox struct {
	c	*C.uiControl
	s	*C.uiSpinbox

	onChanged		func(*Spinbox)
}

// NewSpinbox creates a new Spinbox. If min >= max, they are swapped.
func NewSpinbox(min int, max int) *Spinbox {
	s := new(Spinbox)

	s.s = C.uiNewSpinbox(C.intmax_t(min), C.intmax_t(max))
	s.c = (*C.uiControl)(unsafe.Pointer(s.s))

	C.realuiSpinboxOnChanged(s.s)
	spinboxes[s.s] = s

	return s
}

// Destroy destroys the Spinbox.
func (s *Spinbox) Destroy() {
	delete(spinboxes, s.s)
	C.uiControlDestroy(s.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (s *Spinbox) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(s.c))
}

// Handle returns the OS-level handle associated with this Spinbox.
// On Windows this is an HWND of a standard Windows API EDIT
// class (as provided by Common Controls version 6). Due to
// various limitations which affect the lifetime of the associated
// Common Controls version 6 UPDOWN_CLASS window that
// provides the buttons, there is no way to access it.
// On GTK+ this is a pointer to a GtkSpinButton.
// On OS X this is a pointer to a NSView that contains a NSTextField
// and a NSStepper as subviews.
func (s *Spinbox) Handle() uintptr {
	return uintptr(C.uiControlHandle(s.c))
}

// Show shows the Spinbox.
func (s *Spinbox) Show() {
	C.uiControlShow(s.c)
}

// Hide hides the Spinbox.
func (s *Spinbox) Hide() {
	C.uiControlHide(s.c)
}

// Enable enables the Spinbox.
func (s *Spinbox) Enable() {
	C.uiControlEnable(s.c)
}

// Disable disables the Spinbox.
func (s *Spinbox) Disable() {
	C.uiControlDisable(s.c)
}

// Value returns the Spinbox's current value.
func (s *Spinbox) Value() int {
	return int(C.uiSpinboxValue(s.s))
}

// SetText sets the Spinbox's current value to value.
func (s *Spinbox) SetValue(value int) {
	C.uiSpinboxSetValue(s.s, C.intmax_t(value))
}

// OnChanged registers f to be run when the user changes the value
// of the Spinbox. Only one function can be registered at a time.
func (s *Spinbox) OnChanged(f func(*Spinbox)) {
	s.onChanged = f
}

//export doSpinboxOnChanged
func doSpinboxOnChanged(ss *C.uiSpinbox, data unsafe.Pointer) {
	s := spinboxes[ss]
	if s.onChanged != nil {
		s.onChanged(s)
	}
}
