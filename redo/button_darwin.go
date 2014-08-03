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
