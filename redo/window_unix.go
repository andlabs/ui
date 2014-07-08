// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type window struct {
	widget	*C.GtkWidget
	container	*C.GtkContainer
	bin		*C.GtkBin
	window	*C.GtkWindow
}

func newWindow(title string, width int, height int) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			widget := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
			ctitle := togstr(title)
			defer freegstr(ctitle)
			w := &window{
				widget:		widget,
				container:		(*C.GtkContainer)(unsafe.Pointer(widget)),
				bin:			(*C.GtkBin)(unsafe.Pointer(widget)),
				window:		(*C.GtkWindow)(unsafe.Pointer(widget)),
			}
			C.gtk_window_set_title(w.window, ctitle)
			// TODO size
			// TODO content
			c <- w
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
			c <- fromgstr(C.gtk_window_get_title(w.window))
		},
		resp:		c,
	}
}

func (w *window) SetTitle(title string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			ctitle := togstr(title)
			defer freegstr(ctitle)
			C.gtk_window_set_title(w.window, ctitle)
			c <- struct{}{}
		},
		resp:		c,
	}
}


func (w *window) Show() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.gtk_widget_show_all(w.widget)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Hide() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.gtk_widget_hide(w.widget)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Close() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.gtk_widget_destroy(w.widget)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) OnClosing(e func(c Doer) bool) *Request {
	// TODO
	return nil
}
