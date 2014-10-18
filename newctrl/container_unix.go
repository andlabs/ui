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

func (c *container) allocation(margined bool) C.GtkAllocation {
	var a C.GtkAllocation

	C.gtk_widget_get_allocation(c.widget, &a)
	if margined {
		a.x += C.gint(gtkXMargin)
		a.y += C.gint(gtkYMargin)
		a.width -= C.gint(gtkXMargin) * 2
		a.height -= C.gint(gtkYMargin) * 2
	}
	return a
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
		d.xpadding = gtkXPadding
		d.ypadding = gtkYPadding
	}
	return d
}
