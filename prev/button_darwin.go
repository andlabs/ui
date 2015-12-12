// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type button struct {
	*controlSingleObject
	clicked *event
}

func newButton(text string) *button {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	b := &button{
		controlSingleObject:		newControlSingleObject(C.newButton()),
		clicked: newEvent(),
	}
	C.buttonSetText(b.id, ctext)
	C.buttonSetDelegate(b.id, unsafe.Pointer(b))
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return C.GoString(C.buttonText(b.id))
}

func (b *button) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.buttonSetText(b.id, ctext)
}

//export buttonClicked
func buttonClicked(xb unsafe.Pointer) {
	b := (*button)(unsafe.Pointer(xb))
	b.clicked.fire()
}
