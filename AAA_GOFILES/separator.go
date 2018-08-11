// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// Separator is a Control that represents a horizontal line that
// visually separates controls.
type Separator struct {
	c	*C.uiControl
	s	*C.uiSeparator
}

// NewSeparator creates a new horizontal Separator.
func NewHorizontalSeparator() *Separator {
	s := new(Separator)

	s.s = C.uiNewHorizontalSeparator()
	s.c = (*C.uiControl)(unsafe.Pointer(s.s))

	return s
}

// Destroy destroys the Separator.
func (s *Separator) Destroy() {
	C.uiControlDestroy(s.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (s *Separator) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(s.c))
}

// Handle returns the OS-level handle associated with this Separator.
// On Windows this is an HWND of a standard Windows API STATIC
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkSeparator.
// On OS X this is a pointer to a NSBox.
func (s *Separator) Handle() uintptr {
	return uintptr(C.uiControlHandle(s.c))
}

// Show shows the Separator.
func (s *Separator) Show() {
	C.uiControlShow(s.c)
}

// Hide hides the Separator.
func (s *Separator) Hide() {
	C.uiControlHide(s.c)
}

// Enable enables the Separator.
func (s *Separator) Enable() {
	C.uiControlEnable(s.c)
}

// Disable disables the Separator.
func (s *Separator) Disable() {
	C.uiControlDisable(s.c)
}
