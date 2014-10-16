// 25 june 2014

package ui

type sizingbase struct {
	xpadding      int
	ypadding      int
}

// The container type, which is defined per-platform, is an internal Control that is only used to house other Controls from the underlying UI toolkit's point of view

/* TODO
func (c *container) resize(x, y, width, height int) {
	if c.child == nil { // no children; nothing to do
		return
	}
	d := c.beginResize()
	allocations := c.child.allocate(x+d.xmargin, y+d.ymargintop,
		width-(2*d.xmargin), height-d.ymargintop-d.ymarginbottom, d)
	c.translateAllocationCoords(allocations, width, height)
	// move in reverse so as to approximate right->left order so neighbors make sense
	for i := len(allocations) - 1; i >= 0; i-- {
		allocations[i].this.commitResize(allocations[i], d)
	}
}
*/
