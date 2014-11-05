// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type textfield struct {
	*controlSingleObject
	changed *event
	invalid C.id
	chainpreferredSize	func(d *sizing) (int, int)
}

func finishNewTextField(id C.id) *textfield {
	t := &textfield{
		controlSingleObject:		newControlSingleObject(id),
		changed: newEvent(),
	}
	C.textfieldSetDelegate(t.id, unsafe.Pointer(t))
	t.chainpreferredSize = t.fpreferredSize
	t.fpreferredSize = t.xpreferredSize
	return t
}

func newTextField() *textfield {
	return finishNewTextField(C.newTextField())
}

func newPasswordField() *textfield {
	return finishNewTextField(C.newPasswordField())
}

func (t *textfield) Text() string {
	return C.GoString(C.textfieldText(t.id))
}

func (t *textfield) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textfieldSetText(t.id, ctext)
}

func (t *textfield) OnChanged(f func()) {
	t.changed.set(f)
}

func (t *textfield) Invalid(reason string) {
	if t.invalid != nil {
		C.textfieldCloseInvalidPopover(t.invalid)
		t.invalid = nil
	}
	if reason == "" {
		return
	}
	creason := C.CString(reason)
	defer C.free(unsafe.Pointer(creason))
	t.invalid = C.textfieldOpenInvalidPopover(t.id, creason)
}

// note that the property here is editable, which is the opposite of read-only

func (t *textfield) ReadOnly() bool {
	return !fromBOOL(C.textfieldEditable(t.id))
}

func (t *textfield) SetReadOnly(readonly bool) {
	C.textfieldSetEditable(t.id, toBOOL(!readonly))
}

//export textfieldChanged
func textfieldChanged(data unsafe.Pointer) {
	t := (*textfield)(data)
	t.changed.fire()
}

func (t *textfield) xpreferredSize(d *sizing) (width, height int) {
	_, height = t.chainpreferredSize(d)
	// the returned width is based on the contents; use this instead
	return C.textfieldWidth, height
}
