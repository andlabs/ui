// 4 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type container struct {
	containerbase
	// TODO rename to id
	view			C.id
}

type sizing struct {
	sizingbase

	// for size calculations
	// nothing for mac

	// for the actual resizing
	neighborAlign		C.struct_xalignment
}

func newContainer(child Control) *container {
	c := new(container)
	c.view = C.newContainerView(unsafe.Pointer(c))
	c.child = child
	c.child.setParent(&controlParent{c.view})
	return c
}

//export containerResized
func containerResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
	c := (*container)(unsafe.Pointer(data))
	// the origin of a view's content area is always (0, 0)
	c.resize(0, 0, int(width), int(height))
}

// THIS IS A GUESS. TODO.
// The only indication that this is remotely correct is the Auto Layout Guide implying that 12 pixels is the "Aqua space".
const (
	macXMargin = 12
	macYMargin = 12
	macXPadding = 12
	macYPadding = 12
)

func (c *container) beginResize() (d *sizing) {
	d = new(sizing)
	if spaced {
		d.xmargin = macXMargin
		d.ymargin = macYMargin
		d.xpadding = macXPadding
		d.ypadding = macYPadding
	}
	return d
}

func (c *container) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	for _, a := range allocations {
		// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
		// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
		a.y = (winheight - a.y) - a.height
	}
}
