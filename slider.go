// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doSliderOnChanged(uiSlider *, void *);
// static inline void realuiSliderOnChanged(uiSlider *b)
// {
// 	uiSliderOnChanged(b, doSliderOnChanged, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var sliders = make(map[*C.uiSlider]*Slider)

// Slider is a Control that represents a horizontal bar that represents
// a range of integers. The user can drag a pointer on the bar to
// select an integer.
type Slider struct {
	c	*C.uiControl
	s	*C.uiSlider

	onChanged		func(*Slider)
}

// NewSlider creates a new Slider. If min >= max, they are swapped.
func NewSlider(min int, max int) *Slider {
	s := new(Slider)

	s.s = C.uiNewSlider(C.intmax_t(min), C.intmax_t(max))
	s.c = (*C.uiControl)(unsafe.Pointer(s.s))

	C.realuiSliderOnChanged(s.s)
	sliders[s.s] = s

	return s
}

// Destroy destroys the Slider.
func (s *Slider) Destroy() {
	delete(sliders, s.s)
	C.uiControlDestroy(s.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (s *Slider) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(s.c))
}

// Handle returns the OS-level handle associated with this Slider.
// On Windows this is an HWND of a standard Windows API
// TRACKBAR_CLASS class (as provided by Common Controls
// version 6).
// On GTK+ this is a pointer to a GtkScale.
// On OS X this is a pointer to a NSSlider.
func (s *Slider) Handle() uintptr {
	return uintptr(C.uiControlHandle(s.c))
}

// Show shows the Slider.
func (s *Slider) Show() {
	C.uiControlShow(s.c)
}

// Hide hides the Slider.
func (s *Slider) Hide() {
	C.uiControlHide(s.c)
}

// Enable enables the Slider.
func (s *Slider) Enable() {
	C.uiControlEnable(s.c)
}

// Disable disables the Slider.
func (s *Slider) Disable() {
	C.uiControlDisable(s.c)
}

// Value returns the Slider's current value.
func (s *Slider) Value() int {
	return int(C.uiSliderValue(s.s))
}

// SetText sets the Slider's current value to value.
func (s *Slider) SetValue(value int) {
	C.uiSliderSetValue(s.s, C.intmax_t(value))
}

// OnChanged registers f to be run when the user changes the value
// of the Slider. Only one function can be registered at a time.
func (s *Slider) OnChanged(f func(*Slider)) {
	s.onChanged = f
}

//export doSliderOnChanged
func doSliderOnChanged(ss *C.uiSlider, data unsafe.Pointer) {
	s := sliders[ss]
	if s.onChanged != nil {
		s.onChanged(s)
	}
}
