// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern void checkboxToggled(GtkToggleButton *, gpointer);
import "C"

type checkbox struct {
	*controlSingleWidget
	button   *C.GtkButton
	toggle   *C.GtkToggleButton
	checkbox *C.GtkCheckButton
	toggled  *event
}

func newCheckbox(text string) *checkbox {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_check_button_new_with_label(ctext)
	c := &checkbox{
		controlSingleWidget:  newControlSingleWidget(widget),
		button:   (*C.GtkButton)(unsafe.Pointer(widget)),
		toggle:   (*C.GtkToggleButton)(unsafe.Pointer(widget)),
		checkbox: (*C.GtkCheckButton)(unsafe.Pointer(widget)),
		toggled:  newEvent(),
	}
	g_signal_connect(
		C.gpointer(unsafe.Pointer(c.checkbox)),
		"toggled",
		C.GCallback(C.checkboxToggled),
		C.gpointer(unsafe.Pointer(c)))
	return c
}

func (c *checkbox) OnToggled(e func()) {
	c.toggled.set(e)
}

func (c *checkbox) Text() string {
	return fromgstr(C.gtk_button_get_label(c.button))
}

func (c *checkbox) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_button_set_label(c.button, ctext)
}

func (c *checkbox) Checked() bool {
	return fromgbool(C.gtk_toggle_button_get_active(c.toggle))
}

func (c *checkbox) SetChecked(checked bool) {
	C.gtk_toggle_button_set_active(c.toggle, togbool(checked))
}

//export checkboxToggled
func checkboxToggled(bwid *C.GtkToggleButton, data C.gpointer) {
	c := (*checkbox)(unsafe.Pointer(data))
	c.toggled.fire()
}
