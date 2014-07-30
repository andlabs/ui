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

// TODOs:
// - standalone label on its own: should it be centered or not?

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

type label struct {
	*controlbase
	misc					*C.GtkMisc
	label					*C.GtkLabel
	standalone			bool
	supercommitResize		func(c *allocation, d *sizing)
}

func finishNewLabel(text string, standalone bool) *label {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_label_new(ctext)
	l := &label{
		controlbase:	newControl(widget),
		misc:		(*C.GtkMisc)(unsafe.Pointer(widget)),
		label:		(*C.GtkLabel)(unsafe.Pointer(widget)),
		standalone:	standalone,
	}
	l.supercommitResize = l.fcommitResize
	l.fcommitResize = l.labelcommitResize
	return l
}

func newLabel(text string) Label {
	return finishNewLabel(text, false)
}

func newStandaloneLabel(text string) Label {
	return finishNewLabel(text, true)
}

func (l *label) Text() string {
	return fromgstr(C.gtk_label_get_text(l.label))
}

func (l *label) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_label_set_text(l.label, ctext)
}

func (l *label) labelcommitResize(c *allocation, d *sizing) {
	if !l.standalone && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.shouldVAlignTop {
			// TODO should it be center-aligned to the first line or not
			C.gtk_misc_set_alignment(l.misc, 0, 0)
		} else {
			C.gtk_misc_set_alignment(l.misc, 0, 0.5)
		}
	}
	l.supercommitResize(c, d)
}
