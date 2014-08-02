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

// A sizer hosts a Control and resizes that Control based on changes in size to the parent Window.
// sizer is used by Window, Tab, and [TODO implement] Group to contain and control their respective controls.
// Window is the beginning of the resize chain; resizes happen on the system side.
// Tab and Group are Controls and thus implement controlSizing; they should call their internal sizers's resize() method in their own commitResize().
type sizer struct {
	child		Control
}

// set to true to apply spacing to all windows
var spaced bool = false

func (c *sizer) resize(x, y, width, height int) {
	if c.child == nil {		// no children; nothing to do
		return
	}
	d := c.beginResize()
	allocations := c.child.allocate(x + d.xmargin, y + d.ymargin, width - (2 * d.xmargin), height - (2 * d.ymargin), d)
	c.translateAllocationCoords(allocations, width, height)
	// move in reverse so as to approximate right->left order so neighbors make sense
	for i := len(allocations) - 1; i >= 0; i-- {
		allocations[i].this.commitResize(allocations[i], d)
	}
}
