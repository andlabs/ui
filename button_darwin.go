// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type button struct {
	_id		C.id
	clicked	*event
}

func newButton(text string) *button {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	b := &button{
		_id:			C.newButton(),
		clicked:		newEvent(),
	}
	C.buttonSetText(b._id, ctext)
	C.buttonSetDelegate(b._id, unsafe.Pointer(b))
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return C.GoString(C.buttonText(b._id))
}

func (b *button) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.buttonSetText(b._id, ctext)
}

//export buttonClicked
func buttonClicked(xb unsafe.Pointer) {
	b := (*button)(unsafe.Pointer(xb))
	b.clicked.fire()
	println("button clicked")
}

func (b *button) id() C.id {
	return b._id
}

func (b *button) setParent(p *controlParent) {
	basesetParent(b, p)
}

func (b *button) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(b, x, y, width, height, d)
}

func (b *button) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(b, d)
}

func (b *button) commitResize(a *allocation, d *sizing) {
	basecommitResize(b, a, d)
}

func (b *button) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(b, d)
}
