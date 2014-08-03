// +build !windows,!darwin

// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type tab struct {
	*controlbase
	notebook		*C.GtkNotebook

	tabs			[]*layout
}

func newTab() Tab {
	widget := C.gtk_notebook_new()
	t := &tab{
		controlbase:	newControl(widget),
		notebook:		(*C.GtkNotebook)(unsafe.Pointer(widget)),
	}
	// there are no scrolling arrows by default; add them in case there are too many tabs
	C.gtk_notebook_set_scrollable(t.notebook, C.TRUE)
	return t
}

func (t *tab) Append(name string, control Control) {
	tl := newLayout(control)
	t.tabs = append(t.tabs, tl)
	cname := togstr(name)
	defer freegstr(cname)
	tab := C.gtk_notebook_append_page(t.notebook,
		tl.layoutwidget,
		C.gtk_label_new(cname))
	if tab == -1 {
		panic("gtk_notebook_append_page() failed")
	}
}

// no need to override Control.commitResize() as only prepared the tabbed control; its children will be reallocated when that one is resized
