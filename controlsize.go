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

// for verification; see sysdata.go
type sysDataSizeFuncs interface {
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
	for _, c := range s.allocations {
		c.this.commitResize(c, d)
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

// vertical stack: no concept of neighbor, but not too hard to add a vertical neighbor
// horizontal stack:
	var current *allocation
	// ...
	as := s.controls[i].allocate(...)
	if current != nil {
		current.neighbor = as[0].self
	}
	current = as[0]
	// append all of as
// grid:
	// same as above, except current is set to nil on each new row
	// adding a vertical neighbor would require storing an extra list
