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

type textfield struct {
	_widget		*C.GtkWidget
	entry		*C.GtkEntry
}

func startNewTextField() *textfield {
	widget := C.gtk_entry_new()
	return &textfield{
		_widget:		widget,
		entry:		(*C.GtkEntry)(unsafe.Pointer(widget)),
	}
}

func newTextField() *textfield {
	return startNewTextField()
}

func newPasswordField() *textfield {
	t := startNewTextField()
	C.gtk_entry_set_visibility(t.entry, C.FALSE)
	return t
}

func (t *textfield) Text() string {
	return fromgstr(C.gtk_entry_get_text(t.entry))
}

func (t *textfield) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_entry_set_text(t.entry, ctext)
}

func (t *textfield) widget() *C.GtkWidget {
	return t._widget
}

func (t *textfield) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *textfield) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *textfield) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(t, d)
}

func (t *textfield) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *textfield) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
