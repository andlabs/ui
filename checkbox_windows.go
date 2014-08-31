// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type checkbox struct {
	_hwnd	C.HWND
	_textlen	C.LONG
	toggled	*event
}

func newCheckbox(text string) *checkbox {
	// don't use BS_AUTOCHECKBOX here because it creates problems when refocusing (see http://blogs.msdn.com/b/oldnewthing/archive/2014/05/22/10527522.aspx)
	// we'll handle actually toggling the check state ourselves (see controls_windows.c)
	hwnd := C.newControl(buttonclass,
		C.BS_CHECKBOX | C.WS_TABSTOP,
		0)
	c := &checkbox{
		_hwnd:		hwnd,
		toggled:		newEvent(),
	}
	c.SetText(text)
	C.controlSetControlFont(c._hwnd)
	C.setCheckboxSubclass(c._hwnd, unsafe.Pointer(c))
	return c
}

func (c *checkbox) OnToggled(e func()) {
	c.toggled.set(e)
}

func (c *checkbox) Text() string {
	return baseText(c)
}

func (c *checkbox) SetText(text string) {
	baseSetText(c, text)
}

func (c *checkbox) Checked() bool {
	return C.checkboxChecked(c._hwnd) != C.FALSE
}

func (c *checkbox) SetChecked(checked bool) {
	if checked {
		C.checkboxSetChecked(c._hwnd, C.TRUE)
		return
	}
	C.checkboxSetChecked(c._hwnd, C.FALSE)
}

//export checkboxToggled
func checkboxToggled(data unsafe.Pointer) {
	c := (*checkbox)(data)
	c.toggled.fire()
	println("checkbox toggled")
}

func (c *checkbox) hwnd() C.HWND {
	return c._hwnd
}

func (c *checkbox) textlen() C.LONG {
	return c._textlen
}

func (c *checkbox) settextlen(len C.LONG) {
	c._textlen = len
}

func (c *checkbox) setParent(p *controlParent) {
	basesetParent(c, p)
}

func (c *checkbox) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(c, x, y, width, height, d)
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	checkboxHeight = 10
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/bb226818%28v=vs.85%29.aspx
	checkboxXFromLeftOfBoxToLeftOfLabel = 12
)

func (c *checkbox) preferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(checkboxXFromLeftOfBoxToLeftOfLabel, d) + int(c._textlen),
		fromdlgunitsY(checkboxHeight, d)
}

func (c *checkbox) commitResize(a *allocation, d *sizing) {
	basecommitResize(c, a, d)
}

func (c *checkbox) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(c, d)
}
