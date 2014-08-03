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
	*controlbase
	entry		*C.GtkEntry
}

func startNewTextField() *textField {
	w := C.gtk_entry_new()
	return &textField{
		controlbase:	newControl(w),
		entry:		(*C.GtkEntry)(unsafe.Pointer(w)),
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
