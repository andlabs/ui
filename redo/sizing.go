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

// this ensures that all *windows across all platforms contain the necessary functions
// if this fails to compile, we have a problem
var windowSizeEnsure interface {
	beginResize() *sizing
	endResize(*sizing)
	translateAllocationCoords([]*allocation, int, int)
} = &window{}

type controlSizing interface {
	allocate(x int, y int, width int, height int, d *sizing) []*allocation
	preferredSize(*sizing) (int, int)
	commitResize(*allocation, *sizing)
	getAuxResizeInfo(*sizing)
}

// on Windows, this is only embedded by window, as all other containers cannot have their own children
// on GTK+ and Mac OS X, one is embedded by window and all containers; the containers call container.continueResize()
type container struct {
	child			Control
	spaced		bool
	beginResize	func() (d *sizing)
}

func (c *container) resize(width, height int) {
	if c.child == nil {		// no children; nothing to do
		return
	}
	d := c.beginResize()
	c.continueResize(width, height, d)
}

func (c *container) continueResize(width, height int, d *sizing) {
	if c.child == nil {		// no children; nothing to do
		return
	}
	allocations := c.child.allocate(0, 0, width, height, d)
	c.translateAllocationCoords(allocations, width, height)
	// move in reverse so as to approximate right->left order so neighbors make sense
	for i := len(allocations) - 1; i >= 0; i-- {
		allocations[i].this.commitResize(allocations[i], d)
	}
	c.endResize(d)
}
