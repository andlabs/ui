// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type widgetbase struct {
	widget	*C.GtkWidget
}

func newWidget(w *C.GtkWidget) *widget {
	return &widget{
		widget:	w,
	}
}

type button struct {
	*widgetbase
	button	*C.GtkButton
}

func newButton(text string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			ctext := togstr(text)
			defer freegstr(ctext)
			widget := C.gtk_button_new_with_label(ctext)
			c <- &button{
				widget:	newWidget(widget),
				button:	(*C.GtkButton)(unsafe.Pointer(widget)),
			}
		},
		resp:		c,
	}
}

// TODO OnClicked

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
