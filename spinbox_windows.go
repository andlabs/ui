// 28 october 2014

package ui

import (
	"strconv"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

// TODO do we have to manually monitor user changes to the edit control?
// TODO WS_EX_CLIENTEDGE on the updown?

type spinbox struct {
	hwndEdit			C.HWND
	hwndUpDown		C.HWND
	changed			*event
	// updown state
	updownVisible		bool
	// keep these here to avoid having to get them out
	value			int
	min				int
	max				int
}

func newSpinbox(min int, max int) Spinbox {
	s := new(spinbox)
	s.hwndEdit = C.newControl(editclass,
		C.textfieldStyle | C.ES_NUMBER,
		C.textfieldExtStyle)
	s.changed = newEvent()
	s.updownVisible = true		// initially shown
	s.min = min
	s.max = max
	s.value = s.min
	s.remakeUpDown()
	C.controlSetControlFont(s.hwndEdit)
	C.setSpinboxEditSubclass(s.hwndEdit, unsafe.Pointer(s))
	return s
}

func (s *spinbox) cap() {
	if s.value < s.min {
		s.value = s.min
	}
	if s.value > s.max {
		s.value = s.max
	}
}

func (s *spinbox) Value() int {
	return s.value
}

func (s *spinbox) SetValue(value int) {
	// UDM_SETPOS32 is documented to do what we want, but since we're keeping a copy of value we need to do it anyway
	s.value = value
	s.cap()
	C.SendMessageW(s.hwndUpDown, C.UDM_SETPOS32, 0, C.LPARAM(s.value))
}

func (s *spinbox) OnChanged(e func()) {
	s.changed.set(e)
}

//export spinboxUpDownClicked
func spinboxUpDownClicked(data unsafe.Pointer, nud *C.NMUPDOWN) {
	// this is where we do custom increments
	s := (*spinbox)(data)
	s.value = int(nud.iPos + nud.iDelta)
	// this can go above or below the bounds (the spinbox only rejects invalid values after the UDN_DELTAPOS notification is processed)
	// because we have a copy of the value, we need to fix that here
	s.cap()
	s.changed.fire()
}

//export spinboxEditChanged
func spinboxEditChanged(data unsafe.Pointer) {
	// we're basically on our own here
	s := (*spinbox)(unsafe.Pointer(data))
	// this basically does what OS X does: values too low get clamped to the minimum, values too high get clamped to the maximum, and deleting everything clamps to the minimum
	value, err := strconv.Atoi(getWindowText(s.hwndEdit))
	if err != nil {
		// best we can do fo rnow in this case :S
		// a partial atoi() like in C would be more optimal
		// it handles the deleting everything case just fine
		value = s.min
	}
	s.value = value
	s.cap()
	C.SendMessageW(s.hwndUpDown, C.UDM_SETPOS32, 0, C.LPARAM(s.value))
	// TODO position the insertion caret at the end (or wherever is appropriate)
	s.changed.fire()
}

func (s *spinbox) setParent(p *controlParent) {
	C.controlSetParent(s.hwndEdit, p.hwnd)
	C.controlSetParent(s.hwndUpDown, p.hwnd)
}

// an up-down control will only properly position itself the first time
// stupidly, there are no messages to force a size calculation, nor can I seem to reset the buddy window to force a new position
// alas, we have to make a new up/down control each time :(
// TODO we'll need to store a copy of the current position and range for this
func (s *spinbox) remakeUpDown() {
	// destroying the previous one, setting the parent properly, and subclassing are handled here
	s.hwndUpDown = C.newUpDown(s.hwndUpDown, unsafe.Pointer(s))
	// for this to work, hwndUpDown needs to have rect [0 0 0 0]
	C.moveWindow(s.hwndUpDown, 0, 0, 0, 0)
	C.SendMessageW(s.hwndUpDown, C.UDM_SETBUDDY, C.WPARAM(uintptr(unsafe.Pointer(s.hwndEdit))), 0)
	C.SendMessageW(s.hwndUpDown, C.UDM_SETRANGE32, C.WPARAM(s.min), C.LPARAM(s.max))
	C.SendMessageW(s.hwndUpDown, C.UDM_SETPOS32, 0, C.LPARAM(s.value))
	if s.updownVisible {
		C.ShowWindow(s.hwndUpDown, C.SW_SHOW)
	}
}

// use the same height as normal text fields
// TODO constrain the width somehow
func (s *spinbox) preferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(textfieldWidth, d), fromdlgunitsY(textfieldHeight, d)
}

func (s *spinbox) resize(x int, y int, width int, height int, d *sizing) {
	C.moveWindow(s.hwndEdit, C.int(x), C.int(y), C.int(width), C.int(height))
	s.remakeUpDown()
}

func (s *spinbox) nTabStops() int {
	// TODO does the up-down control count?
	return 1
}

// TODO be sure to modify this when we add Show()/Hide()
func (s *spinbox) containerShow() {
	C.ShowWindow(s.hwndEdit, C.SW_SHOW)
	C.ShowWindow(s.hwndUpDown, C.SW_SHOW)
	s.updownVisible = true
}

// TODO be sure to modify this when we add Show()/Hide()
func (s *spinbox) containerHide() {
	C.ShowWindow(s.hwndEdit, C.SW_HIDE)
	C.ShowWindow(s.hwndUpDown, C.SW_HIDE)
	s.updownVisible = false
}
