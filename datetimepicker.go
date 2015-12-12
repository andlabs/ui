// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// DateTimePicker is a Control that represents a field where the user
// can enter a date and/or a time.
type DateTimePicker struct {
	c	*C.uiControl
	d	*C.uiDateTimePicker
}

// NewDateTimePicker creates a new DateTimePicker that shows
// both a date and a time.
func NewDateTimePicker() *DateTimePicker {
	d := new(DateTimePicker)

	d.d = C.uiNewDateTimePicker()
	d.c = (*C.uiControl)(unsafe.Pointer(d.d))

	return d
}

// NewDatePicker creates a new DateTimePicker that shows
// only a date.
func NewDatePicker() *DateTimePicker {
	d := new(DateTimePicker)

	d.d = C.uiNewDatePicker()
	d.c = (*C.uiControl)(unsafe.Pointer(d.d))

	return d
}

// NewTimePicker creates a new DateTimePicker that shows
// only a time.
func NewTimePicker() *DateTimePicker {
	d := new(DateTimePicker)

	d.d = C.uiNewTimePicker()
	d.c = (*C.uiControl)(unsafe.Pointer(d.d))

	return d
}

// Destroy destroys the DateTimePicker.
func (d *DateTimePicker) Destroy() {
	C.uiControlDestroy(d.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (d *DateTimePicker) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(d.c))
}

// Handle returns the OS-level handle associated with this DateTimePicker.
// On Windows this is an HWND of a standard Windows API
// DATETIMEPICK_CLASS class (as provided by Common Controls
// version 6).
// On GTK+ this is a pointer to a libui-internal class.
// On OS X this is a pointer to a NSDatePicker.
func (d *DateTimePicker) Handle() uintptr {
	return uintptr(C.uiControlHandle(d.c))
}

// Show shows the DateTimePicker.
func (d *DateTimePicker) Show() {
	C.uiControlShow(d.c)
}

// Hide hides the DateTimePicker.
func (d *DateTimePicker) Hide() {
	C.uiControlHide(d.c)
}

// Enable enables the DateTimePicker.
func (d *DateTimePicker) Enable() {
	C.uiControlEnable(d.c)
}

// Disable disables the DateTimePicker.
func (d *DateTimePicker) Disable() {
	C.uiControlDisable(d.c)
}
