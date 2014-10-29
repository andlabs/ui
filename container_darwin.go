// 4 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type container struct {
	id			C.id
	resize		func(x int, y int, width int, height int, d *sizing)
	margined		bool
}

type sizing struct {
	sizingbase

	// for size calculations
	// nothing on Mac OS X

	// for the actual resizing
	neighborAlign C.struct_xalignment
}

// containerResized() gets called early so we have to do this in the constructor
func newContainer(resize func(x int, y int, width int, height int, d *sizing)) *container {
	c := new(container)
	c.resize = resize
	c.id = C.newContainerView(unsafe.Pointer(c))
	return c
}

func (c *container) parent() *controlParent {
	return &controlParent{c.id}
}

//export containerResized
func containerResized(data unsafe.Pointer) {
	c := (*container)(data)
	d := beginResize()
	// TODO make this a parameter
	b := C.containerBounds(c.id)
	if c.margined {
		b.x += C.intptr_t(macXMargin)
		b.y += C.intptr_t(macYMargin)
		b.width -= C.intptr_t(macXMargin) * 2
		b.height -= C.intptr_t(macYMargin) * 2
	}
	c.resize(int(b.x), int(b.y), int(b.width), int(b.height), d)
}

// These are based on measurements from Interface Builder.
// TODO reverify these against /layout rects/, not /frame rects/
const (
	macXMargin  = 20
	macYMargin  = 20
	macXPadding = 8
	macYPadding = 8
)

func beginResize() (d *sizing) {
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
