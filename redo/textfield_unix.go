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

type textField struct {
	_widget		*C.GtkWidget
	entry		*C.GtkEntry
}

func startNewTextField() *textField {
	widget := C.gtk_entry_new()
	return &textField{
		_widget:		widget,
		entry:		(*C.GtkEntry)(unsafe.Pointer(widget)),
	}
}

func newTextField() *textField {
	return startNewTextField()
}

func newPasswordField() *textField {
	t := startNewTextField()
	C.gtk_entry_set_visibility(t.entry, C.FALSE)
	return t
}

func (t *textField) Text() string {
	return fromgstr(C.gtk_entry_get_text(t.entry))
}

func (t *textField) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_entry_set_text(t.entry, ctext)
}

func (t *textField) widget() *C.GtkWidget {
	return t._widget
}

func (t *textField) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *textField) containerShow() {
	basecontainerShow(t)
}

func (t *textField) containerHide() {
	basecontainerHide(t)
}

func (t *textField) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *textField) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(t, d)
}

func (t *textField) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *textField) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(d)
}
