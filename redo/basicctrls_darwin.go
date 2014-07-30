// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type button struct {
	*controlbase
	clicked		*event
}

func finishNewButton(id C.id, text string) *button {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	b := &button{
		controlbase:	newControl(id),
		clicked:		newEvent(),
	}
	C.buttonSetText(b.id, ctext)
	C.buttonSetDelegate(b.id, unsafe.Pointer(b))
	return b
}

func newButton(text string) *button {
	return finishNewButton(C.newButton(), text)
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

//export buttonClicked
func buttonClicked(xb unsafe.Pointer) {
	b := (*button)(unsafe.Pointer(xb))
	b.clicked.fire()
	println("button clicked")
}

func (b *button) Text() string {
	return C.GoString(C.buttonText(b.id))
}

func (b *button) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.buttonSetText(b.id, ctext)
}

type checkbox struct {
	*button
}

func newCheckbox(text string) *checkbox {
	return &checkbox{
		button:	finishNewButton(C.newCheckbox(), text),
	}
}

// we don't need to define our own event here; we can just reuse Button's
// (it's all target-action anyway)

func (c *checkbox) Checked() bool {
	return fromBOOL(C.checkboxChecked(c.id))
}

func (c *checkbox) SetChecked(checked bool) {
	C.checkboxSetChecked(c.id, toBOOL(checked))
}

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

// cheap trick
type label struct {
	*textField
	standalone	bool
}

func finishNewLabel(text string, standalone bool) *label {
	l := &label{
		textField:		finishNewTextField(C.newLabel()),
		standalone:	standalone,
	}
	l.SetText(text)
	return l
}

func newLabel(text string) Label {
	return finishNewLabel(text, false)
}

func newStandaloneLabel(text string) Label {
	return finishNewLabel(text, true)
}

// TODO label commitResize
