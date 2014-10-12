// 4 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type container struct {
	containerbase
	id C.id
}

type sizing struct {
	sizingbase

	// for size calculations
	// nothing for mac

	// for the actual resizing
	neighborAlign C.struct_xalignment
}

func newContainer(child Control) *container {
	c := new(container)
	c.id = C.newContainerView(unsafe.Pointer(c))
	c.child = child
	c.child.setParent(&controlParent{c.id})
	return c
}

//export containerResized
func containerResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
	c := (*container)(unsafe.Pointer(data))
	// the origin of a view's content area is always (0, 0)
	c.resize(0, 0, int(width), int(height))
}

// These are based on measurements from Interface Builder.
const (
	macXMargin  = 20
	macYMargin  = 20
	macXPadding = 8
	macYPadding = 8
)

func (c *container) beginResize() (d *sizing) {
	d = new(sizing)
	if c.spaced {
		d.xmargin = macXMargin
		d.ymargintop = macYMargin
		d.ymarginbottom = d.ymargintop
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
