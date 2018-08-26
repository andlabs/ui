// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doSliderOnChanged(uiSlider *, void *);
// // see golang/go#19835
// typedef void (*sliderCallback)(uiSlider *, void *);
import "C"

// Slider is a Control that represents a horizontal bar that represents
// a range of integers. The user can drag a pointer on the bar to
// select an integer.
type Slider struct {
	ControlBase
	s	*C.uiSlider
	onChanged		func(*Slider)
}

// NewSlider creates a new Slider. If min >= max, they are swapped.
func NewSlider(min int, max int) *Slider {
	s := new(Slider)

	s.s = C.uiNewSlider(C.int(min), C.int(max))

	C.uiSliderOnChanged(s.s, C.sliderCallback(C.doSliderOnChanged), nil)

	s.ControlBase = NewControlBase(s, uintptr(unsafe.Pointer(s.s)))
	return s
}

// Value returns the Slider's current value.
func (s *Slider) Value() int {
	return int(C.uiSliderValue(s.s))
}

// SetValue sets the Slider's current value to value.
func (s *Slider) SetValue(value int) {
	C.uiSliderSetValue(s.s, C.int(value))
}

// OnChanged registers f to be run when the user changes the value
// of the Slider. Only one function can be registered at a time.
func (s *Slider) OnChanged(f func(*Slider)) {
	s.onChanged = f
}

//export doSliderOnChanged
func doSliderOnChanged(ss *C.uiSlider, data unsafe.Pointer) {
	s := ControlFromLibui(uintptr(unsafe.Pointer(ss))).(*Slider)
	if s.onChanged != nil {
		s.onChanged(s)
	}
}
