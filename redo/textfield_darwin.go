// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type textfield struct {
	_id	C.id
}

func newTextField() *textfield {
	return &textfield{
		_id:		C.newTextField(),
	}
}

func newPasswordField() *textfield {
	return &textfield{
		_id:		C.newPasswordField(),
	}
}

func (t *textfield) Text() string {
	return C.GoString(C.textFieldText(t._id))
}

func (t *textfield) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textFieldSetText(t._id, ctext)
}

func (t *textfield) id() C.id {
	return t._id
}

func (t *textfield) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *textfield) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *textfield) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(t, d)
}

func (t *textfield) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *textfield) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
