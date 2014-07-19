// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

// TODO move to common_darwin.go
func toBOOL(b bool) C.BOOL {
	if b == true {
		return C.YES
	}
	return C.NO
}

type widgetbase struct {
	id		C.id
	parentw	*window
	floating	bool
}

func newWidget(id C.id) *widgetbase {
	return &widgetbase{
		id:	id,
	}
}

// these few methods are embedded by all the various Controls since they all will do the same thing

func (w *widgetbase) unparent() {
	if w.parentw != nil {
		C.unparent(w.id)
		w.floating = true
		w.parentw = nil
	}
}

func (w *widgetbase) parent(win *window) {
	C.parent(w.id, win.id, toBOOL(w.floating))
	w.floating = false
	w.parentw = win
}

type button struct {
	*widgetbase
	clicked		*event
}

func newButton(text string) *button {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	b := &button{
		widgetbase:	newWidget(C.newButton(ctext)),
		clicked:		newEvent(),
	}
	C.buttonSetDelegate(b.id, unsafe.Pointer(b))
	return b
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
