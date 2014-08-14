// 4 august 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type container struct {
	containerbase
	hwnd		C.HWND
	nchildren		int
}

type sizing struct {
	sizingbase

	// for size calculations
	baseX			C.int
	baseY			C.int
	internalLeading	C.LONG		// for Label; see Label.commitResize() for details

	// for the actual resizing
	// possibly the HDWP
}

func makeContainerWindowClass() error {
	var errmsg *C.char

	err := C.makeContainerWindowClass(&errmsg)
	if err != 0 || errmsg != nil {
		return fmt.Errorf("%s: %v", C.GoString(errmsg), syscall.Errno(err))
	}
	return nil
}

func newContainer(control Control) *container {
	c := new(container)
	hwnd := C.newContainer(unsafe.Pointer(c))
	if hwnd != c.hwnd {
		panic(fmt.Errorf("inconsistency: hwnd returned by CreateWindowEx() (%p) and hwnd stored in container (%p) differ", hwnd, c.hwnd))
	}
	c.child = control
	c.child.setParent(&controlParent{c})
	return c
}

func (c *container) setParent(hwnd C.HWND) {
	C.controlSetParent(c.hwnd, hwnd)
}

// this is needed because Windows won't move/resize a child window for us
func (c *container) move(r *C.RECT) {
	C.moveWindow(c.hwnd, C.int(r.left), C.int(r.top), C.int(r.right - r.left), C.int(r.bottom - r.top))
}

func (c *container) show() {
	C.ShowWindow(c.hwnd, C.SW_SHOW)
}

func (c *container) hide() {
	C.ShowWindow(c.hwnd, C.SW_HIDE)
}

//export storeContainerHWND
func storeContainerHWND(data unsafe.Pointer, hwnd C.HWND) {
	c := (*container)(data)
	c.hwnd = hwnd
}

//export containerResize
func containerResize(data unsafe.Pointer, r *C.RECT) {
	c := (*container)(data)
	// the origin of any window's content area is always (0, 0), but let's use the values from the RECT just to be safe
	c.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top))
}

// For Windows, Microsoft just hands you a list of preferred control sizes as part of the MSDN documentation and tells you to roll with it.
// These sizes are given in "dialog units", which are independent of the font in use.
// We need to convert these into standard pixels, which requires we get the device context of the OS window.
// References:
// - http://msdn.microsoft.com/en-us/library/ms645502%28VS.85%29.aspx - the calculation needed
// - http://support.microsoft.com/kb/125681 - to get the base X and Y
// (thanks to http://stackoverflow.com/questions/58620/default-button-size)
// In my tests (see https://github.com/andlabs/windlgunits), the GetTextExtentPoint32() option for getting the base X produces much more accurate results than the tmAveCharWidth option when tested against the sample values given in http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing, but can be off by a pixel in either direction (probably due to rounding errors).

// note on MulDiv():
// div will not be 0 in the usages below
// we also ignore overflow; that isn't likely to happen for our use case anytime soon

func fromdlgunitsX(du int, d *sizing) int {
	return int(C.MulDiv(C.int(du), d.baseX, 4))
}

func fromdlgunitsY(du int, d *sizing) int {
	return int(C.MulDiv(C.int(du), d.baseY, 8))
}

const (
	marginDialogUnits = 7
	paddingDialogUnits = 4
)

func (c *container) beginResize() (d *sizing) {
	var baseX, baseY C.int
	var internalLeading C.LONG

	d = new(sizing)

	C.calculateBaseUnits(c.hwnd, &baseX, &baseY, &internalLeading)
	d.baseX = baseX
	d.baseY = baseY
	d.internalLeading = internalLeading

	if spaced {
		d.xmargin = fromdlgunitsX(marginDialogUnits, d)
		d.ymargin = fromdlgunitsY(marginDialogUnits, d)
		d.xpadding = fromdlgunitsX(paddingDialogUnits, d)
		d.ypadding = fromdlgunitsY(paddingDialogUnits, d)
	}

	return d
}

func (c *container) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// no translation needed on windows
}
