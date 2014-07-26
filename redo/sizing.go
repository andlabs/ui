// 25 june 2014

package ui

type allocation struct {
	x		int
	y		int
	width	int
	height	int
	this		Control
	neighbor	Control
}

type sizingbase struct {
	xmargin		int
	ymargin		int
	xpadding		int
	ypadding		int
}

type controlSizing interface {
	allocate(x int, y int, width int, height int, d *sizing) []*allocation
	preferredSize(*sizing) (int, int)
	commitResize(*allocation, *sizing)
	getAuxResizeInfo(*sizing)
}

// on Windows, this is only embedded by window, as all other containers cannot have their own children; beginResize() points to an instance method literal (TODO get correct term) from window
// on GTK+ and Mac OS X, one is embedded by window and all containers; beginResize() points to a global function (TODO NOT GOOD; ideally the sizing data should be passed across size-allocate requests)
type container struct {
	child			Control
	spaced		bool
	beginResize	func() (d *sizing)	// for the initial call
	d			*sizing			// for recursive calls
}

func (c *container) resize(width, height int) {
	if c.child == nil {		// no children; nothing to do
		return
	}
	if c.d == nil {			// initial call
		if c.beginResize == nil {
			// should be a recursive call, but is not
			// TODO get rid of this
			return
		}
		c.d = c.beginResize()
	}
	d := c.d
	allocations := c.child.allocate(0, 0, width, height, d)
	c.translateAllocationCoords(allocations, width, height)
	// move in reverse so as to approximate right->left order so neighbors make sense
	for i := len(allocations) - 1; i >= 0; i-- {
		allocations[i].this.commitResize(allocations[i], d)
	}
	// always set c.d to nil so it can be garbage-collected
	// the c.endResize() above won't matter since the c.d there is evaluated then, not when c.endResize() is called
	c.d = nil
}
