// +build !windows,!darwin,!plan9

// 23 february 2014

package ui

type sysSizeData struct {
	cSysSizeData

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

func (s *sysData) beginResize() (d *sysSizeData) {
	d = new(sysSizeData)
	if s.spaced {
		d.xmargin = gtkXMargin
		d.ymargin = gtkYMargin
		d.xpadding = gtkXPadding
		d.ypadding = gtkYPadding
	}
	return d
}

func (s *sysData) endResize(d *sysSizeData) {
	// redraw
}

func (s *sysData) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// no need for coordinate conversion with gtk+
}

func (s *sysData) commitResize(c *allocation, d *sysSizeData) {
	if s.ctype == c_label && !s.alternate && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.shouldVAlignTop {
			// TODO should it be center-aligned to the first line or not
			gtk_misc_set_alignment(s.widget, 0, 0)
		} else {
			gtk_misc_set_alignment(s.widget, 0, 0.5)
		}
	}
	// TODO merge this here
	s.setRect(c.x, c.y, c.width, c.height, 0)
}

func (s *sysData) getAuxResizeInfo(d *sysSizeData) {
	d.shouldVAlignTop = (s.ctype == c_listbox) || (s.ctype == c_area)
}

// GTK+ 3 makes this easy: controls can tell us what their preferred size is!
// ...actually, it tells us two things: the "minimum size" and the "natural size".
// The "minimum size" is the smallest size we /can/ display /anything/. The "natural size" is the smallest size we would /prefer/ to display.
// The difference? Minimum size takes into account things like truncation with ellipses: the minimum size of a label can allot just the ellipses!
// So we use the natural size instead.
// There is a warning about height-for-width controls, but in my tests this isn't an issue.
// For Areas, we manually save the Area size and use that, just to be safe.

// We don't need to worry about y-offset because label alignment is "vertically center", which GtkLabel does for us.

func (s *sysData) preferredSize(d *sysSizeData) (width int, height int) {
	if s.ctype == c_area {
		return s.areawidth, s.areaheight
	}

	_, _, width, height = gtk_widget_get_preferred_size(s.widget)
	return width, height
}
