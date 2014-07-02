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

type cSysSizeData struct {
	xmargin		int
	ymargin		int
	xpadding		int
	ypadding		int
}

// for verification; see sysdata.go
type sysDataSizingFunctions interface {
	beginResize() *sysSizeData
	endResize(*sysSizeData)
	translateAllocationCoords([]*allocation, int, int)
	preferredSize(*sysSizeData) (int, int)
	commitResize(*allocation, *sysSizeData)
	getAuxResizeInfo(*sysSizeData)
}

func (s *sysData) resizeWindow(width, height int) {
	d := s.beginResize()
	allocations := s.allocate(0, 0, width, height, d)
	s.translateAllocationCoords(allocations, width, height)
	// move in reverse so as to approximate right->left order so neighbors make sense
	for i := len(allocations) - 1; i >= 0; i-- {
		allocations[i].this.commitResize(allocations[i], d)
	}
	s.endResize(d)
}

// non-layout controls: allocate() should just return a one-element slice; preferredSize(), commitResize(), and getAuxResizeInfo() should defer to their sysData equivalents
type controlSizing interface {
	allocate(x int, y int, width int, height int, d *sysSizeData) []*allocation
	preferredSize(d *sysSizeData) (width, height int)
	commitResize(c *allocation, d *sysSizeData)
	getAuxResizeInfo(d *sysSizeData)
}
