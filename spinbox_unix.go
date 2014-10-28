// +build !windows,!darwin

// 28 october 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

// TODO preferred width may be too wide

type spinbox struct {
	*controlSingleWidget
	spinbutton	*C.GtkSpinButton
}

func newSpinbox() Spinbox {
	widget := C.gtk_spin_button_new_with_range(0, 100, 1)
	s := &spinbox{
		controlSingleWidget:	newControlSingleWidget(widget),
		spinbutton:			(*C.GtkSpinButton)(unsafe.Pointer(widget)),
	}
	C.gtk_spin_button_set_digits(s.spinbutton, 0)				// integers
	C.gtk_spin_button_set_numeric(s.spinbutton, C.TRUE)		// digits only
	return s
}
