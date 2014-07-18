// 12 july 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type window struct {
	hwnd		C.HWND
	shownbefore	bool

	child			Control

	closing		*event

	spaced		bool
}

const windowclassname = ""
var windowclassptr = syscall.StringToUTF16Ptr(windowclassname)

func makeWindowWindowClass() error {
	var errmsg *C.char

	err := C.makeWindowWindowClass(&errmsg)
	if err != 0 || errmsg != nil {
		return fmt.Errorf("%s: %v", C.GoString(errmsg), syscall.Errno(err))
	}
	return nil
}

func newWindow(title string, width int, height int) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			w := &window{
				// hwnd set in WM_CREATE handler
				closing:	newEvent(),
			}
			hwnd := C.newWindow(toUTF16(title), C.int(width), C.int(height), unsafe.Pointer(w))
			if hwnd != w.hwnd {
				panic(fmt.Errorf("inconsistency: hwnd returned by CreateWindowEx() (%p) and hwnd stored in window (%p) differ", hwnd, w.hwnd))
			}
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
			c <- getWindowText(w.hwnd)
		},
		resp:		c,
	}
}

func (w *window) SetTitle(title string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.setWindowText(w.hwnd, toUTF16(title))
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Show() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			if !w.shownbefore {
				C.ShowWindow(w.hwnd, C.nCmdShow)
				C.updateWindow(w.hwnd)
				w.shownbefore = true
			} else {
				C.ShowWindow(w.hwnd, C.SW_SHOW)
			}
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Hide() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.ShowWindow(w.hwnd, C.SW_HIDE)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Close() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.windowClose(w.hwnd)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) OnClosing(e func(Doer) bool) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			w.closing.setbool(e)
			c <- struct{}{}
		},
		resp:		c,
	}
}

//export storeWindowHWND
func storeWindowHWND(data unsafe.Pointer, hwnd C.HWND) {
	w := (*window)(data)
	w.hwnd = hwnd
}

//export windowResize
func windowResize(data unsafe.Pointer, r *C.RECT) {
	w := (*window)(data)
	w.doresize(int(r.right - r.left), int(r.bottom - r.top))
}

//export windowClosing
func windowClosing(data unsafe.Pointer) {
	w := (*window)(data)
	close := w.closing.fire()
	if close {
		C.windowClose(w.hwnd)
	}
}
