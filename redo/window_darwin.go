// 8 july 2014

package ui

import (
	"unsafe"
"fmt"
)

// #include "objc_darwin.h"
import "C"

type window struct {
	id		C.id

	child		Control

	closing	*event

	spaced	bool
}

func newWindow(title string, width int, height int) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			id := C.newWindow(C.intptr_t(width), C.intptr_t(height))
			ctitle := C.CString(title)
			defer C.free(unsafe.Pointer(ctitle))
			C.windowSetTitle(id, ctitle)
			w := &window{
				id:		id,
				closing:	newEvent(),
			}
			C.windowSetDelegate(id, unsafe.Pointer(w))
			c <- w
		},
		resp:		c,
	}
}

func (w *window) SetControl(control Control) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			if w.child != nil {		// unparent existing control
				w.child.unparent()
			}
			control.unparent()
			control.parent(w)
			w.child = control
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Title() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			c <- C.GoString(C.windowTitle(w.id))
		},
		resp:		c,
	}
}

func (w *window) SetTitle(title string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			ctitle := C.CString(title)
			defer C.free(unsafe.Pointer(ctitle))
			C.windowSetTitle(w.id, ctitle)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Show() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.windowShow(w.id)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Hide() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.windowHide(w.id)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Close() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.windowClose(w.id)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) OnClosing(e func(c Doer) bool) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			w.closing.setbool(e)
			c <- struct{}{}
		},
		resp:		c,
	}
}

//export windowClosing
func windowClosing(xw unsafe.Pointer) C.BOOL {
	w := (*window)(unsafe.Pointer(xw))
	close := w.closing.fire()
	if close {
		// TODO make sure this actually closes the window the way we want
		return C.YES
	}
	return C.NO
}

//export windowResized
func windowResized(xw unsafe.Pointer, width C.uintptr_t, height C.uintptr_t) {
	// TODO this isn't called when the window first opens up
	w := (*window)(unsafe.Pointer(xw))
	w.doresize(int(width), int(height))
	fmt.Printf("new size %d x %d\n", width, height)
}
