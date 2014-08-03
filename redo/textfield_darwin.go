// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type textField struct {
	*controlbase
}

func finishNewTextField(id C.id) *textField {
	return &textField{
		controlbase:	newControl(id),
	}
}

func newTextField() *textField {
	return finishNewTextField(C.newTextField())
}

func newPasswordField() *textField {
	return finishNewTextField(C.newPasswordField())
}

func (t *textField) Text() string {
	return C.GoString(C.textFieldText(t.id))
}

func (t *textField) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textFieldSetText(t.id, ctext)
}
