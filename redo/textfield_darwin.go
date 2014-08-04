// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type textField struct {
	_id	C.id
}

func newTextField() *textField {
	return &textField{
		_id:		C.newTextField(),
	}
}

func newPasswordField() *textField {
	return &textField{
		_id:		C.newPasswordField(),
	}
}

func (t *textField) Text() string {
	return C.GoString(C.textFieldText(t._id))
}

func (t *textField) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textFieldSetText(t._id, ctext)
}

func (t *textField) id() C.id {
	return t._id
}

func (t *textField) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *textField) containerShow() {
	basecontainerShow(t)
}

func (t *textField) containerHide() {
	basecontainerHide(t)
}

func (t *textField) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *textField) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(t, d)
}

func (t *textField) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *textField) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
