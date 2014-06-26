// +build SKIP

// 25 june 2014

package ui

type sysSizeData struct {
	// for size calculations
	// all platforms
	margin	int	// windows: calculated
	spacing	int	// gtk+, cocoa: constants
	// windows
	baseX	int
	baseY	int
	// gtk, mac os x: nothing

	// for the actual resizing
	// windows
	// possibly also the HDWP
	// gtk
	shouldVAlignTop	bool
	// mac os x
	// neighbor control alignment rect/baseline info
}

func (s *sysData) beginResize() *sysSizeData {
	// windows: get baseX/baseY for window and compute margin and spacing
	// gtk, mac: return zero
}

func (s *sysData) endResize(d *sysSizeData) {
	// redraw
}

func (s *sysData) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// windows, gtk: nothing
	// mac
	for _, a := range allocations {
		// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
		// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
		a.y = (winheight - a.y) - a.height
	}
}

// windows
func (s *sysData) doResize(c *allocation, d *sysSizeData) {
	if s.ctype == c_label {
		// add additional offset of 4 dialog units
	}
	// resize
}
func (s *sysData) getAuxResizeInfo(d *sysSizeData) {
	// do nothing
}

// gtk+
func (s *sysData) doResize(c *allocation, d *sysSizeData) {
	if s.ctype == c_label && !s.alternate && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.shouldVAlignTop {
			// TODO should it be center-aligned to the first line or not
			gtk_misc_set_align(s.widget, 0, 0)
		} else {
			gtk_misc_set_align(s.widget, 0, 0.5)
		}
	}
	// resize
}
func (s *sysData) getAuxResizeInfo(d *sysSizeData) {
	d.shouldVAlignTop = (s.ctype == c_listbox) || (s.ctype == c_area)
}

// cocoa
func (s *sysData) doResize(c *allocation, d *sysSizeData) {
	if s.ctype == c_label && !s.alternate && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		// get this control's alignment rect and baseline
		// align
	}
	// resize
}
func (s *sysData) getAuxResizeInfo(d *sysSizeData) {
	// get this control's alignment rect and baseline
}
