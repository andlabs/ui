// +build !windows,!darwin

// 23 february 2014

package ui

import (
	"unsafe"
"fmt"
)

// #include "gtk_unix.h"
// extern void containerResizing(GtkWidget *, GdkRectangle *, gpointer);
import "C"

type container struct {
	containerbase
	layoutwidget		*C.GtkWidget
	layoutcontainer	*C.GtkContainer
	layout			*C.GtkLayout
}

func newContainer(child Control) *container {
	widget := C.gtk_layout_new(nil, nil)
	c := &container{
		layoutwidget:		widget,
		layoutcontainer:	(*C.GtkContainer)(unsafe.Pointer(widget)),
		layout:			(*C.GtkLayout)(unsafe.Pointer(widget)),
	}
	c.child = child
	c.child.setParent(&controlParent{c.layoutcontainer})
	// we connect to the layout's size-allocate, not to the window's configure-event
	// this allows us to handle client-side decoration-based configurations (such as GTK+ on Wayland) properly
	// also see basecommitResize() in control_unix.go for additional notes
	// thanks to many people in irc.gimp.net/#gtk+ for help (including tristan for suggesting g_signal_connect_after())
	g_signal_connect_after(
		C.gpointer(unsafe.Pointer(c.layout)),
		"size-allocate",
		C.GCallback(C.containerResizing),
		C.gpointer(unsafe.Pointer(c)))
	return c
}

//export containerResizing
func containerResizing(wid *C.GtkWidget, r *C.GdkRectangle, data C.gpointer) {
	c := (*container)(unsafe.Pointer(data))
	// the GtkLayout's coordinate system is localized, so the origin is (0, 0)
	c.resize(0, 0, int(r.width), int(r.height))
fmt.Printf("new size %d x %d\n", r.width, r.height)
}

type sizing struct {
	sizingbase

	// for size calculations
	// gtk+ needs nothing

	// for the actual resizing
	shouldVAlignTop	bool
}

const (
	gtkXMargin = 12
	gtkYMargin = 12
	gtkXPadding = 12
	gtkYPadding = 6
)

func (c *container) beginResize() (d *sizing) {
	d = new(sizing)
	if spaced {
		d.xmargin = gtkXMargin
		d.ymargin = gtkYMargin
		d.xpadding = gtkXPadding
		d.ypadding = gtkYPadding
	}
	return d
}

func (c *container) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// no need for coordinate conversion with gtk+
}
