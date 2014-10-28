// 28 october 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

// TODO do we have to manually monitor user changes to the edit control?

type spinbox struct {
	hwndEdit			C.HWND
	hwndUpDown		C.HWND
}

func newSpinbox() Spinbox {
	s := new(spinbox)
	s.hwndEdit = C.newControl(editclass,
		C.textfieldStyle | C.ES_NUMBER,
		C.textfieldExtStyle)
	s.hwndUpDown = C.newControl(C.xUPDOWN_CLASSW,
		C.UDS_ALIGNRIGHT | C.UDS_ARROWKEYS | C.UDS_HOTTRACK | C.UDS_NOTHOUSANDS | C.UDS_SETBUDDYINT,
		0)
	C.SendMessageW(s.hwndUpDown, C.UDM_SETBUDDY, C.WPARAM(uintptr(unsafe.Pointer(s.hwndEdit))), 0)
	C.SendMessageW(s.hwndUpDown, C.UDM_SETRANGE32, 0, 100)
	C.SendMessageW(s.hwndUpDown, C.UDM_SETPOS32, 0, 0)
	return s
}

func (s *spinbox) setParent(p *controlParent) {
	C.controlSetParent(s.hwndEdit, p.hwnd)
	C.controlSetParent(s.hwndUpDown, p.hwnd)
}

func (s *spinbox) preferredSize(d *sizing) (width, height int) {
	// TODO
	return 20, 20
}

func (s *spinbox) resize(x int, y int, width int, height int, d *sizing) {
	C.moveWindow(s.hwndEdit, C.int(x), C.int(y), C.int(width), C.int(height))
}

func (s *spinbox) nTabStops() int {
	// TODO does the up-down control count?
	return 1
}

// TODO be sure to modify this when we add Show()/Hide()
func (s *spinbox) containerShow() {
	C.ShowWindow(s.hwndEdit, C.SW_SHOW)
	C.ShowWindow(s.hwndUpDown, C.SW_SHOW)
}

// TODO be sure to modify this when we add Show()/Hide()
func (s *spinbox) containerHide() {
	C.ShowWindow(s.hwndEdit, C.SW_HIDE)
	C.ShowWindow(s.hwndUpDown, C.SW_HIDE)
}
