// +build !windows,!darwin

// 23 february 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

// TODO avoid direct access to contents?
type container struct {
	widget		*C.GtkWidget
	container		*C.GtkContainer
	resize		func(x int, y int, width int, height int, d *sizing)
	margined		bool
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
	c.widget = C.newContainer(unsafe.Pointer(c))
	c.container = (*C.GtkContainer)(unsafe.Pointer(c.widget))
	return c
}

func (c *container) parent() *controlParent {
	return &controlParent{c.container}
}

//export containerResize
func containerResize(data unsafe.Pointer, aorig *C.GtkAllocation) {
	c := (*container)(data)
	d := beginResize()
	// copy aorig
	a := *aorig
	if c.margined {
		a.x += C.int(gtkXMargin)
		a.y += C.int(gtkYMargin)
		a.width -= C.int(gtkXMargin) * 2
		a.height -= C.int(gtkYMargin) * 2
	}
	c.resize(int(a.x), int(a.y), int(a.width), int(a.height), d)
}

const (
	gtkXMargin  = 12
	gtkYMargin  = 12
	gtkXPadding = 12
	gtkYPadding = 6
)

func beginResize() (d *sizing) {
	d = new(sizing)
	d.xpadding = gtkXPadding
	d.ypadding = gtkYPadding
	return d
}
