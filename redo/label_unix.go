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
			// don't bother aligning it to the first line of text in the control; this is harder than it's worth (thanks gregier in irc.gimp.net/#gtk+)
			C.gtk_misc_set_alignment(l.misc, 0, 0)
		} else {
			C.gtk_misc_set_alignment(l.misc, 0, 0.5)
		}
	}
	l.supercommitResize(c, d)
}
