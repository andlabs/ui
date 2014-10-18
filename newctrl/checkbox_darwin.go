// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type checkbox struct {
	*controlSingleObject
	toggled *event
}

func newCheckbox(text string) *checkbox {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	c := &checkbox{
		controlSingleObject:		newControlSingleObject(C.newCheckbox()),
		toggled: newEvent(),
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
