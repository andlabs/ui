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

func (w *window) doresize(width, height int) {
	if w.child == nil {		// no children; nothing to do
		return
	}
	d := w.beginResize()
	allocations := w.child.allocate(0, 0, width, height, d)
	w.translateAllocationCoords(allocations, width, height)
	// move in reverse so as to approximate right->left order so neighbors make sense
	for i := len(allocations) - 1; i >= 0; i-- {
		allocations[i].this.commitResize(allocations[i], d)
	}
	w.endResize(d)
}
