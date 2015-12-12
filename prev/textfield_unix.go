// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern void textfieldChanged(GtkEditable *, gpointer);
// /* because cgo doesn't like GTK_STOCK_DIALOG_ERROR */
// static inline void setErrorIcon(GtkEntry *entry)
// {
// 	gtk_entry_set_icon_from_stock(entry, GTK_ENTRY_ICON_SECONDARY, GTK_STOCK_DIALOG_ERROR);
// }
import "C"

type textfield struct {
	*controlSingleWidget
	editable	*C.GtkEditable
	entry   *C.GtkEntry
	changed *event
}

func startNewTextField() *textfield {
	widget := C.gtk_entry_new()
	t := &textfield{
		controlSingleWidget: newControlSingleWidget(widget),
		editable:	(*C.GtkEditable)(unsafe.Pointer(widget)),
		entry:   (*C.GtkEntry)(unsafe.Pointer(widget)),
		changed: newEvent(),
	}
	g_signal_connect(
		C.gpointer(unsafe.Pointer(t.widget)),
		"changed",
		C.GCallback(C.textfieldChanged),
		C.gpointer(unsafe.Pointer(t)))
	return t
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

func (t *textfield) OnChanged(f func()) {
	t.changed.set(f)
}

func (t *textfield) Invalid(reason string) {
	if reason == "" {
		C.gtk_entry_set_icon_from_stock(t.entry, C.GTK_ENTRY_ICON_SECONDARY, nil)
		return
	}
	C.setErrorIcon(t.entry)
	creason := togstr(reason)
	defer freegstr(creason)
	C.gtk_entry_set_icon_tooltip_text(t.entry, C.GTK_ENTRY_ICON_SECONDARY, creason)
	C.gtk_widget_error_bell(t.widget)
}

// note that the property here is editable, which is the opposite of read-only

func (t *textfield) ReadOnly() bool {
	return !fromgbool(C.gtk_editable_get_editable(t.editable))
}

func (t *textfield) SetReadOnly(readonly bool) {
	C.gtk_editable_set_editable(t.editable, togbool(!readonly))
}

//export textfieldChanged
func textfieldChanged(editable *C.GtkEditable, data C.gpointer) {
	t := (*textfield)(unsafe.Pointer(data))
	t.changed.fire()
}
