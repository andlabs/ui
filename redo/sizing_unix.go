// +build !windows,!darwin

// 23 february 2014

package ui

// #include "gtk_unix.h"
import "C"

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
