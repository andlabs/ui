// 8 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type window struct {
	id		C.id

	closing	*event
}

func newWindow(title string, width int, height int) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			id := C.newWindow(C.intptr_t(width), C.intptr_t(height))
			ctitle := C.CString(title)
			defer C.free(unsafe.Pointer(ctitle))
			C.windowSetTitle(id, ctitle)
			C.windowSetAppDelegate(id)
			c <- &window{
				id:		id,
				closing:	newEvent(),
			}
		},
		resp:		c,
	}
}

func (w *window) SetControl(control Control) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			// TODO unparent
			// TODO reparent
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

// TODO windowClosing

// TODO for testing
func newButton(string) *Request { return nil }
