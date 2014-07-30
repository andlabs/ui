// 30 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type controlbase struct {
	*controldefs
	hwnd	C.HWND
	parent	C.HWND		// for Tab and Group
}

type controlParent struct {
	hwnd	C.HWND
}

func newControl(class C.LPCWSTR, style C.DWORD, extstyle C.DWORD) *controlbase {
	c := new(controlbase)
	c.hwnd = C.newWidget(class, style, extstyle)
	c.controldefs = new(controldefs)
	c.fsetParent = func(p *controlParent) {
		C.controlSetParent(c.hwnd, p.hwnd)
		c.parent = p.hwnd
	}
	c.fcontainerShow = func() {
		C.ShowWindow(c.hwnd, C.SW_SHOW)
	}
	c.fcontainerHide = func() {
		C.ShowWindow(c.hwnd, C.SW_HIDE)
	}
	c.fallocate = func(x int, y int, width int, height int, d *sizing) []*allocation {
		// TODO split into its own function
		return []*allocation{&allocation{
			x:		x,
			y:		y,
			width:	width,
			height:	height,
			this:		c,
		}}
	}
	c.fpreferredSize = func(d *sizing) (int, int) {
		// TODO
		return 75, 23
	}
	c.fcommitResize = func(a *allocation, d *sizing) {
		C.moveWindow(c.hwnd, C.int(a.x), C.int(a.y), C.int(a.width), C.int(a.height))
	}
	c.fgetAuxResizeInfo = func(d *sizing) {
		// do nothing
	}
	return c
}

// these are provided for convenience

func (c *controlbase) text() string {
	return getWindowText(c.hwnd)
}

func (c *controlbase) setText(text string) {
	C.setWindowText(c.hwnd, toUTF16(text))
}
