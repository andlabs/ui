// +build !windows,!darwin

// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern void layoutResizing(GtkWidget *, GdkRectangle *, gpointer);
import "C"

type tab struct {
	*widgetbase
	notebook		*C.GtkNotebook

	containers	[]*container
	layoutws		[]*C.GtkWidget
	layoutcs		[]*C.GtkContainer
	layouts		[]*C.GtkLayout
}

func newTab() Tab {
	widget := C.gtk_notebook_new()
	t := &tab{
		widgetbase:	newWidget(widget),
		notebook:		(*C.GtkNotebook)(unsafe.Pointer(widget)),
	}
	// there are no scrolling arrows by default; add them in case there are too many tabs
	C.gtk_notebook_set_scrollable(t.notebook, C.TRUE)
	return t
}

func (t *tab) Append(name string, control Control) {
	// TODO isolate and standardize
	layout := C.gtk_layout_new(nil, nil)
	t.layoutws = append(t.layoutws, layout)
	t.layoutcs = append(t.layoutcs, (*C.GtkContainer)(unsafe.Pointer(layout)))
	t.layouts = append(t.layouts, (*C.GtkLayout)(unsafe.Pointer(layout)))
	c := new(container)
	t.containers = append(t.containers, c)
	c.child = control
	c.child.setParent(&controlParent{(*C.GtkContainer)(unsafe.Pointer(layout))})
	g_signal_connect_after(
		C.gpointer(unsafe.Pointer(layout)),
		"size-allocate",
		C.GCallback(C.layoutResizing),
		C.gpointer(unsafe.Pointer(c)))
	cname := togstr(name)
	defer freegstr(cname)
	tab := C.gtk_notebook_append_page(t.notebook,
		layout,
		C.gtk_label_new(cname))
	if tab == -1 {
		panic("gtk_notebook_append_page() failed")
	}
}

func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	// only prepared the tabbed control; its children will be reallocated when that one is resized
	return t.widgetbase.allocate(x, y, width, height, d)
}

//export layoutResizing
func layoutResizing(wid *C.GtkWidget, r *C.GdkRectangle, data C.gpointer) {
	c := (*container)(unsafe.Pointer(data))
	// the layout's coordinate system is localized, so the origin is (0, 0)
	c.resize(0, 0, int(r.width), int(r.height))
}
