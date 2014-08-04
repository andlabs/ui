// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type checkbox struct {
	_id		C.id
	toggled	*event
}

func newCheckbox(text string) *checkbox {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	c := &checkbox{
		_id:			C.newCheckbox(),
		toggled:		newEvent(),
	}
	C.buttonSetText(c._id, ctext)
	C.checkboxSetDelegate(c._id, unsafe.Pointer(c))
	return c
}

func (c *checkbox) OnToggled(e func()) {
	c.toggled.set(e)
}

func (c *checkbox) Text() string {
	return C.GoString(C.buttonText(c._id))
}

func (c *checkbox) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.buttonSetText(c._id, ctext)
}

func (c *checkbox) Checked() bool {
	return fromBOOL(C.checkboxChecked(c._id))
}

func (c *checkbox) SetChecked(checked bool) {
	C.checkboxSetChecked(c._id, toBOOL(checked))
}

//export checkboxToggled
func checkboxToggled(xc unsafe.Pointer) {
	c := (*checkbox)(unsafe.Pointer(xc))
	c.toggled.fire()
}

func (c *checkbox) id() C.id {
	return c._id
}

func (c *checkbox) setParent(p *controlParent) {
	basesetParent(c, p)
}

func (c *checkbox) containerShow() {
	basecontainerShow(c)
}

func (c *checkbox) containerHide() {
	basecontainerHide(c)
}

func (c *checkbox) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(c, x, y, width, height, d)
}

func (c *checkbox) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(c, d)
}

func (c *checkbox) commitResize(a *allocation, d *sizing) {
	basecommitResize(c, a, d)
}

func (c *checkbox) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(c, d)
}
