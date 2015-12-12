// +build !windows,!darwin

// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type tab struct {
	*controlSingleWidget
	container *C.GtkContainer
	notebook  *C.GtkNotebook

	tabs []*container
	children	[]Control
}

func newTab() Tab {
	widget := C.gtk_notebook_new()
	t := &tab{
		controlSingleWidget:	newControlSingleWidget(widget),
		container: (*C.GtkContainer)(unsafe.Pointer(widget)),
		notebook:  (*C.GtkNotebook)(unsafe.Pointer(widget)),
	}
	// there are no scrolling arrows by default; add them in case there are too many tabs
	C.gtk_notebook_set_scrollable(t.notebook, C.TRUE)
	return t
}

func (t *tab) Append(name string, control Control) {
	c := newContainer()
	t.tabs = append(t.tabs, c)
	// this calls gtk_container_add(), which, according to gregier in irc.gimp.net/#gtk+, acts just like gtk_notebook_append_page()
	C.gtk_container_add(t.container, c.widget)
	control.setParent(c.parent())
	c.resize = control.resize
	t.children = append(t.children, control)
	cname := togstr(name)
	defer freegstr(cname)
	C.gtk_notebook_set_tab_label_text(t.notebook,
		// unfortunately there does not seem to be a gtk_notebook_set_nth_tab_label_text()
		C.gtk_notebook_get_nth_page(t.notebook, C.gint(len(t.tabs)-1)),
		cname)
}

// no need to handle resize; the children containers handle that for us
