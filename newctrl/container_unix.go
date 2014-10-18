// +build !windows,!darwin

// 23 february 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type container struct {
	*controlSingleWidget
	container *C.GtkContainer
}

type sizing struct {
	sizingbase

	// for size calculations
	// gtk+ needs nothing

	// for the actual resizing
	// gtk+ needs nothing
}

func newContainer() *container {
	c := new(container)
	c.controlSingleWidget = newControlSingleWidget(C.newContainer(unsafe.Pointer(c)))
	c.container = (*C.GtkContainer)(unsafe.Pointer(c.widget))
	return c
}

func (c *container) parent() *controlParent {
	return &controlParent{c.container}
}

//export containerResizing
func containerResizing(data unsafe.Pointer, r *C.GtkAllocation) {
	c := (*container)(data)
	c.resize(int(r.x), int(r.y), int(r.width), int(r.height))
}

const (
	gtkXMargin  = 12
	gtkYMargin  = 12
	gtkXPadding = 12
	gtkYPadding = 6
)

func (w *window) beginResize() (d *sizing) {
	d = new(sizing)
	if spaced {
		d.xmargin = gtkXMargin
		d.ymargintop = gtkYMargin
		d.ymarginbottom = d.ymargintop
		d.xpadding = gtkXPadding
		d.ypadding = gtkYPadding
	}
	return d
}
