// 30 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type controlPrivate interface {
	// TODO
	Control
}

type controlbase struct {
	hwnd	C.HWND
	parent	C.HWND		// for Tab and Group
	textlen	C.LONG
}

type controlParent struct {
	hwnd	C.HWND
}

func newControl(class C.LPWSTR, style C.DWORD, extstyle C.DWORD) *controlbase {
	c := new(controlbase)
	// TODO rename to newWidget
	c.hwnd = C.newWidget(class, style, extstyle)
	return c
}

// TODO for maximum correctness these shouldn't take controlbases... but then the amount of duplicated code would skyrocket

func basesetParent(c *controlbase, p *controlParent) {
	C.controlSetParent(c.hwnd, p.hwnd)
	c.parent = p.hwnd
}

func basecontainerShow(c *controlbase) {
	C.ShowWindow(c.hwnd, C.SW_SHOW)
}

func basecontainerHide(c *controlbase) {
	C.ShowWindow(c.hwnd, C.SW_HIDE)
}

// don't specify basepreferredSize; it is custom on ALL controls

func basecommitResize(c *controlbase, a *allocation, d *sizing) {
	C.moveWindow(c.hwnd, C.int(a.x), C.int(a.y), C.int(a.width), C.int(a.height))
}

func basegetAuxResizeInfo(c controlPrivate, d *sizing) {
	// do nothing
}

// these are provided for convenience

func (c *controlbase) text() string {
	return getWindowText(c.hwnd)
}

func (c *controlbase) setText(text string) {
	t := toUTF16(text)
	C.setWindowText(c.hwnd, t)
	c.textlen = C.controlTextLength(c.hwnd, t)
}
