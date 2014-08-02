// 1 march 2014

package ui

// #include "objc_darwin.h"
import "C"

type sizing struct {
	sizingbase

	// for size calculations
	// nothing for mac

	// for the actual resizing
	neighborAlign		C.struct_xalignment
}

// THIS IS A GUESS. TODO.
// The only indication that this is remotely correct is the Auto Layout Guide implying that 12 pixels is the "Aqua space".
const (
	macXMargin = 12
	macYMargin = 12
	macXPadding = 12
	macYPadding = 12
)

func (s *sizer) beginResize() (d *sizing) {
	d = new(sizing)
	if spaced {
		d.xmargin = macXMargin
		d.ymargin = macYMargin
		d.xpadding = macXPadding
		d.ypadding = macYPadding
	}
	return d
}

func (s *sizer) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	for _, a := range allocations {
		// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
		// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
		a.y = (winheight - a.y) - a.height
	}
}
