// 4 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type container struct {
	*controlSingleObject
}

type sizing struct {
	sizingbase

	// for size calculations
	// nothing for mac

	// for the actual resizing
	neighborAlign C.struct_xalignment
}

func newContainer() *container {
	c := new(container)
	c.controlSingleObject = newControlSingleObject(C.newContainerView(unsafe.Pointer(c)))
	return c
}

func (c *container) parent() *controlParent {
	return &controlParent{c.id}
}

func (c *container) allocation(margined bool) C.struct_xrect {
	b := C.containerBounds(c.id)
	if margined {
		b.x += C.intptr_t(macXMargin)
		b.y += C.intptr_t(macYMargin)
		b.width -= C.intptr_t(macXMargin) * 2
		b.height -= C.intptr_t(macYMargin) * 2
	}
	return b
}

// we can just return these values as is
func (c *container) bounds(d *sizing) (int, int, int, int) {
	b := C.containerBounds(c.id)
	return int(b.x), int(b.y), int(b.width), int(b.height)
}

// These are based on measurements from Interface Builder.
const (
	macXMargin  = 20
	macYMargin  = 20
	macXPadding = 8
	macYPadding = 8
)

func (w *window) beginResize() (d *sizing) {
	d = new(sizing)
	d.xpadding = macXPadding
	d.ypadding = macYPadding
	return d
}

/*TODO
func (c *container) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	for _, a := range allocations {
		// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
		// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
		a.y = (winheight - a.y) - a.height
	}
}
*/
