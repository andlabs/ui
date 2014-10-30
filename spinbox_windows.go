// 28 october 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

// TODO do we have to manually monitor user changes to the edit control?
// TODO WS_EX_CLIENTEDGE on the updown?

type spinbox struct {
	hwndEdit			C.HWND
	hwndUpDown		C.HWND
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
	s.updownVisible = true		// initially shown
	s.min = min
	s.max = max
	s.value = s.min
	s.remakeUpDown()
	return s
}

func (s *spinbox) Value() int {
	// TODO TODO TODO TODO TODO
	// this CAN error out!!!
	// we need to update s.value but we need to implement events first
	return int(C.SendMessageW(s.hwndUpDown, C.UDM_GETPOS32, 0, 0))
}

func (s *spinbox) SetValue(value int) {
	// UDM_SETPOS32 is documented to do what we want, but since we're keeping a copy of value we need to do it anyway
	if value < s.min {
		value = s.min
	}
	if value > s.max {
		value = s.max
	}
	s.value = value
	C.SendMessageW(s.hwndUpDown, C.UDM_SETPOS32, 0, C.LPARAM(value))
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
	// destroying the previous one and setting the parent properly is handled here
	s.hwndUpDown = C.newUpDown(s.hwndUpDown)
	// for this to work, hwndUpDown needs to have rect [0 0 0 0]
	C.moveWindow(s.hwndUpDown, 0, 0, 0, 0)
	C.SendMessageW(s.hwndUpDown, C.UDM_SETBUDDY, C.WPARAM(uintptr(unsafe.Pointer(s.hwndEdit))), 0)
	C.SendMessageW(s.hwndUpDown, C.UDM_SETRANGE32, C.WPARAM(s.min), C.LPARAM(s.max))
	C.SendMessageW(s.hwndUpDown, C.UDM_SETPOS32, 0, C.LPARAM(s.value))
	if s.updownVisible {
		C.ShowWindow(s.hwndUpDown, C.SW_SHOW)
	}
}

func (s *spinbox) preferredSize(d *sizing) (width, height int) {
	// TODO
	return 20, 20
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
