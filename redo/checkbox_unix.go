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

type checkbox struct {
	// embed button so its methods and events carry over
	*button
	toggle		*C.GtkToggleButton
	checkbox		*C.GtkCheckButton
}

func newCheckbox(text string) *checkbox {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_check_button_new_with_label(ctext)
	return &checkbox{
		button:		finishNewButton(widget, "toggled", C.checkboxToggled),
		toggle:		(*C.GtkToggleButton)(unsafe.Pointer(widget)),
		checkbox:	(*C.GtkCheckButton)(unsafe.Pointer(widget)),
	}
}

//export checkboxToggled
func checkboxToggled(bwid *C.GtkToggleButton, data C.gpointer) {
	// note that the finishNewButton() call uses the embedded *button as data
	// this is fine because we're only deferring to buttonClicked() anyway
	buttonClicked(nil, data)
}

func (c *checkbox) Checked() bool {
	return fromgbool(C.gtk_toggle_button_get_active(c.toggle))
}

func (c *checkbox) SetChecked(checked bool) {
	C.gtk_toggle_button_set_active(c.toggle, togbool(checked))
}

