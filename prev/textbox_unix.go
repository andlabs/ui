// +build !windows,!darwin

// 23 october 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type textbox struct {
	*scroller
	textview		*C.GtkTextView
}

func newTextbox() Textbox {
	widget := C.gtk_text_view_new()
	t := &textbox{
		scroller:		newScroller(widget, true, true, false),		// natively scrollable, has a border, no overlay
		textview:		(*C.GtkTextView)(unsafe.Pointer(widget)),
	}
	return t
}

func (t *textbox) Text() string {
	var start, end C.GtkTextIter

	buf := C.gtk_text_view_get_buffer(t.textview)
	C.gtk_text_buffer_get_bounds(buf, &start, &end)
	// include hidden chars even though there can't be one since Textbox is explicitly unformatted just to be safe
	// don't worry about embedded pixbufs or widgets; those aren't allowed either
	ctext := C.gtk_text_buffer_get_text(buf, &start, &end, C.TRUE)
	// not explicitly documented: have to manually free this (thanks ste in irc.gimp.net/#gtk+)
	defer C.g_free(C.gpointer(unsafe.Pointer(ctext)))
	return fromgstr(ctext)
}

func (t *textbox) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	buf := C.gtk_text_view_get_buffer(t.textview)
	C.gtk_text_buffer_set_text(buf, ctext, -1)		// null-terminated
}
