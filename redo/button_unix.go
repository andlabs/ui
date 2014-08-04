// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern void buttonClicked(GtkButton *, gpointer);
import "C"

type button struct {
	_widget		*C.GtkWidget
	button		*C.GtkButton
	clicked		*event
}

// shared code for setting up buttons, check boxes, etc.
func newButton(text string) *button {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_button_new_with_label(ctext)
	b := &button{
		_widget:		widget,
		button:		(*C.GtkButton)(unsafe.Pointer(widget)),
		clicked:		newEvent(),
	}
	g_signal_connect(
		C.gpointer(unsafe.Pointer(b.button)),
		"clicked",
		C.GCallback(C.buttonClicked),
		C.gpointer(unsafe.Pointer(b)))
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return fromgstr(C.gtk_button_get_label(b.button))
}

func (b *button) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_button_set_label(b.button, ctext)
}

//export buttonClicked
func buttonClicked(bwid *C.GtkButton, data C.gpointer) {
	b := (*button)(unsafe.Pointer(data))
	b.clicked.fire()
	println("button clicked")
}

func (b *button) widget() *C.GtkWidget {
	return b._widget
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
