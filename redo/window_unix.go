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
			ctext := togstr(text)
			defer freegstr(ctext)
			w := &window{
				widget:		widget,
				container:		(*C.GtkContainer)(unsafe.Pointer(widget)),
				bin:			(*C.GtkBin)(unsafe.Pointer(widget)),
				window:		(*C.GtkWindow)(unsafe.Pointer(widget)),
			}
			C.gtk_window_set_title(w.window, ctext)
			// TODO size
			// TODO content
			c <-  w
		},
		resp:		c,
	}
}

func (w *window) SetControl(c Control) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			// TODO unparent
			// TODO reparent
			c <- struct{}{}
		},
		done:	c,
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
			ctext := togstr(text)
			defer freegstr(ctext)
			C.gtk_window_set_title(w.window, ctext)
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

func (w *window) OnClosing(func e(c Doer) bool) *Request {
	// TODO
	return nil
}
