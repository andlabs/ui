// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type label struct {
	*controlSingleWidget
	misc       *C.GtkMisc
	label      *C.GtkLabel
}

func newLabel(text string) Label {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_label_new(ctext)
	l := &label{
		controlSingleWidget:    newControlSingleWidget(widget),
		misc:       (*C.GtkMisc)(unsafe.Pointer(widget)),
		label:      (*C.GtkLabel)(unsafe.Pointer(widget)),
	}
	return l
}

/*TODO
func newStandaloneLabel(text string) Label {
	l := finishNewLabel(text, true)
	// standalone labels are always at the top left
	C.gtk_misc_set_alignment(l.misc, 0, 0)
	return l
}
*/

func (l *label) Text() string {
	return fromgstr(C.gtk_label_get_text(l.label))
}

func (l *label) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_label_set_text(l.label, ctext)
}

/*TODO
func (l *label) commitResize(c *allocation, d *sizing) {
	if !l.standalone && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.shouldVAlignTop {
			// don't bother aligning it to the first line of text in the control; this is harder than it's worth (thanks gregier in irc.gimp.net/#gtk+)
			C.gtk_misc_set_alignment(l.misc, 0, 0)
		} else {
			C.gtk_misc_set_alignment(l.misc, 0, 0.5)
		}
	}
	basecommitResize(l, c, d)
}
*/
