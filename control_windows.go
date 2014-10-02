// 30 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type controlPrivate interface {
	hwnd() C.HWND
	Control
}

type controlParent struct {
	c *container
}

func basesetParent(c controlPrivate, p *controlParent) {
	C.controlSetParent(c.hwnd(), p.c.hwnd)
	p.c.nchildren++
}

// don't specify basepreferredSize; it is custom on ALL controls

func basecommitResize(c controlPrivate, a *allocation, d *sizing) {
	C.moveWindow(c.hwnd(), C.int(a.x), C.int(a.y), C.int(a.width), C.int(a.height))
}

func basegetAuxResizeInfo(c controlPrivate, d *sizing) {
	// do nothing
}

// these are provided for convenience

type textableControl interface {
	controlPrivate
	textlen() C.LONG
	settextlen(C.LONG)
}

func baseText(c textableControl) string {
	return getWindowText(c.hwnd())
}

func baseSetText(c textableControl, text string) {
	hwnd := c.hwnd()
	t := toUTF16(text)
	C.setWindowText(hwnd, t)
	c.settextlen(C.controlTextLength(hwnd, t))
}
