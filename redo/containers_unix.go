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
	return t
}

func (t *tab) Append(name string, control Control) {
	// TODO isolate and standardize
	layout := C.gtk_layout_new(nil, nil)
	t.layoutws = append(t.layoutws, layout)
	t.layoutcs = append(t.layoutcs, (*C.GtkContainer)(unsafe.Pointer(layout)))
	t.layouts = append(t.layouts, (*C.GtkLayout)(unsafe.Pointer(layout)))
	c := new(container)
	// don't set beginResize; this container's resize() will be a recursive call
	t.containers = append(t.containers, c)
	c.child = control
	c.child.setParent((*C.GtkContainer)(unsafe.Pointer(layout)))
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
	// set up the recursive calls
	for _, c := range t.containers {
		c.d = d
	}
	// and prepare the tabbed control itself
	return t.widgetbase.allocate(x, y, width, height, d)
}

//export layoutResizing
func layoutResizing(wid *C.GtkWidget, r *C.GdkRectangle, data C.gpointer) {
	c := (*container)(unsafe.Pointer(data))
	c.resize(int(r.width), int(r.height))
}
