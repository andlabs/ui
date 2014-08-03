// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern void buttonClicked(GtkButton *, gpointer);
// extern void checkboxToggled(GtkToggleButton *, gpointer);
import "C"

type button struct {
	*controlbase
	button		*C.GtkButton
	clicked		*event
}

// shared code for setting up buttons, check boxes, etc.
func finishNewButton(widget *C.GtkWidget, event string, handler unsafe.Pointer) *button {
	b := &button{
		controlbase:	newControl(widget),
		button:		(*C.GtkButton)(unsafe.Pointer(widget)),
		clicked:		newEvent(),
	}
	g_signal_connect(
		C.gpointer(unsafe.Pointer(b.button)),
		event,
		C.GCallback(handler),
		C.gpointer(unsafe.Pointer(b)))
	return b
}

func newButton(text string) *button {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_button_new_with_label(ctext)
	return finishNewButton(widget, "clicked", C.buttonClicked)
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

//export buttonClicked
func buttonClicked(bwid *C.GtkButton, data C.gpointer) {
	b := (*button)(unsafe.Pointer(data))
	b.clicked.fire()
	println("button clicked")
}

func (b *button) Text() string {
	return fromgstr(C.gtk_button_get_label(b.button))
}

func (b *button) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_button_set_label(b.button, ctext)
}
