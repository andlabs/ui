// 30 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type controlParent struct {
	hwnd	C.HWND
}
type sizing struct {
	// TODO
}

// don't specify preferredSize in any of these; they're per-control

type controlSingleHWND struct {
	*controlbase
	hwnd	C.HWND
}

func newControlSingleHWND(hwnd C.HWND) *controlSingleHWND {
	c := new(controlSingleHWND)
	c.controlbase = &controlbase{
		fsetParent:	c.setParent,
		fresize:		c.resize,
	}
	c.hwnd = hwnd
	return c
}

func (c *controlSingleHWND) setParent(p *controlParent) {
	C.controlSetParent(c.hwnd, p.hwnd)
}

func (c *controlSingleHWND) resize(x int, y int, width int, height int, d *sizing) {
	C.moveWindow(c.hwnd, C.int(x), C.int(y), C.int(width), C.int(height))
}

// these are provided for convenience

type controlSingleHWNDWithText struct {
	*controlSingleHWND
	textlen	C.LONG
}

func newControlSingleHWNDWithText(h C.HWND) *controlSingleHWNDWithText {
	return &controlSingleHWNDWithText{
		controlSingleHWND:		newControlSingleHWND(h),
	}
}

func (c *controlSingleHWNDWithText) text() string {
	return getWindowText(c.hwnd)
}

func (c *controlSingleHWNDWithText) setText(text string) {
	t := toUTF16(text)
	C.setWindowText(c.hwnd, t)
	c.textlen = C.controlTextLength(c.hwnd, t)
}
