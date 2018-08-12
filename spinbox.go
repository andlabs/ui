// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doSpinboxOnChanged(uiSpinbox *, void *);
import "C"

// Spinbox is a Control that represents a space where the user can
// enter integers. The space also comes with buttons to add or
// subtract 1 from the integer.
type Spinbox struct {
	ControlBase
	s	*C.uiSpinbox
	onChanged		func(*Spinbox)
}

// NewSpinbox creates a new Spinbox. If min >= max, they are swapped.
func NewSpinbox(min int, max int) *Spinbox {
	s := new(Spinbox)

	s.s = C.uiNewSpinbox(C.int(min), C.int(max))

	C.uiSpinboxOnChanged(s.s, C.doSpinboxOnChanged, nil)

	s.ControlBase = NewControlBase(s, uintptr(unsafe.Pointer(s.s)))
	return s
}

// Value returns the Spinbox's current value.
func (s *Spinbox) Value() int {
	return int(C.uiSpinboxValue(s.s))
}

// SetValue sets the Spinbox's current value to value.
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
	s := ControlFromLibui(uintptr(unsafe.Pointer(ss))).(*Spinbox)
	if s.onChanged != nil {
		s.onChanged(s)
	}
}
