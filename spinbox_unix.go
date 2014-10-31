// +build !windows,!darwin

// 28 october 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern void spinboxChanged(GtkSpinButton *, gpointer);
import "C"

// TODO preferred width may be too wide

type spinbox struct {
	*controlSingleWidget
	spinbutton	*C.GtkSpinButton
	changed		*event
}

func newSpinbox(min int, max int) Spinbox {
	// gtk_spin_button_new_with_range() initially sets its value to the minimum value
	widget := C.gtk_spin_button_new_with_range(C.gdouble(min), C.gdouble(max), 1)
	s := &spinbox{
		controlSingleWidget:	newControlSingleWidget(widget),
		spinbutton:			(*C.GtkSpinButton)(unsafe.Pointer(widget)),
		changed:				newEvent(),
	}
	C.gtk_spin_button_set_digits(s.spinbutton, 0)				// integers
	C.gtk_spin_button_set_numeric(s.spinbutton, C.TRUE)		// digits only
	// this isn't specifically documented as the signal to connect to until 3.14
	// it has existed as far back as 3.4, though, if not earlier
	// there's also ::change-value which is for keyboard changing
	g_signal_connect(
		C.gpointer(unsafe.Pointer(s.spinbutton)),
		"value-changed",
		C.GCallback(C.spinboxChanged),
		C.gpointer(unsafe.Pointer(s)))
	return s
}

func (s *spinbox) Value() int {
	return int(C.gtk_spin_button_get_value(s.spinbutton))
}

func (s *spinbox) SetValue(value int) {
	var min, max C.gdouble

	C.gtk_spin_button_get_range(s.spinbutton, &min, &max)
	if value < int(min) {
		value = int(min)
	}
	if value > int(max) {
		value = int(max)
	}
	C.gtk_spin_button_set_value(s.spinbutton, C.gdouble(value))
}

func (s *spinbox) OnChanged(e func()) {
	s.changed.set(e)
}

//export spinboxChanged
func spinboxChanged(swid *C.GtkSpinButton, data C.gpointer) {
	s := (*spinbox)(unsafe.Pointer(data))
	s.changed.fire()
}
