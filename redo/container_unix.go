// +build !windows,!darwin

// 23 february 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type container struct {
	containerbase
	layoutwidget		*C.GtkWidget
	layoutcontainer	*C.GtkContainer
}

type sizing struct {
	sizingbase

	// for size calculations
	// gtk+ needs nothing

	// for the actual resizing
	shouldVAlignTop	bool
}

func newContainer(child Control) *container {
	c := new(container)
	widget := C.newContainer(unsafe.Pointer(c))
	c.layoutwidget = widget
	c.layoutcontainer = (*C.GtkContainer)(unsafe.Pointer(widget))
	c.child = child
	c.child.setParent(&controlParent{c.layoutcontainer})
	return c
}

func (c *container) setParent(p *controlParent) {
	C.gtk_container_add(p.c, c.layoutwidget)
}

//export containerResizing
func containerResizing(data unsafe.Pointer, r *C.GtkAllocation) {
	c := (*container)(data)
	c.resize(int(r.x), int(r.y), int(r.width), int(r.height))
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
