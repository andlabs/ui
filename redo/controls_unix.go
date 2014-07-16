// +build !darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type widgetbase struct {
	widget	*C.GtkWidget
	parentw	*window
	floating	bool
}

func newWidget(w *C.GtkWidget) *widgetbase {
	return &widgetbase{
		widget:	w,
	}
}

// these few methods are embedded by all the various Controls since they all will do the same thing

func (w *widgetbase) unparent() {
	if w.parentw != nil {
		// add another reference so it doesn't get removed by accident
		C.g_object_ref(C.gpointer(unsafe.Pointer(w.widget)))
		// we unref this in parent() below
		w.floating = true
		C.gtk_container_remove(w.parentw.layoutc, w.widget)
		w.parentw = nil
	}
}

func (w *widgetbase) parent(win *window) {
	C.gtk_container_add(win.layoutc, w.widget)
	w.parentw = win
	// was previously parented; unref our saved ref
	if w.floating {
		C.g_object_unref(C.gpointer(unsafe.Pointer(w.widget)))
		w.floating = false
	}
}

type button struct {
	*widgetbase
	button		*C.GtkButton
}

func newButton(text string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			ctext := togstr(text)
			defer freegstr(ctext)
			widget := C.gtk_button_new_with_label(ctext)
			c <- &button{
				widgetbase:	newWidget(widget),
				button:		(*C.GtkButton)(unsafe.Pointer(widget)),
			}
		},
		resp:		c,
	}
}

func (b *button) OnClicked(e func(c Doer)) *Request {
	// TODO
	return nil
}

func (b *button) Text() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			c <- fromgstr(C.gtk_button_get_label(b.button))
		},
		resp:		c,
	}
}

func (b *button) SetText(text string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			ctext := togstr(text)
			defer freegstr(ctext)
			C.gtk_button_set_label(b.button, ctext)
			c <- struct{}{}
		},
		resp:		c,
	}
}
