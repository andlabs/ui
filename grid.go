// 31 august 2014

package ui

import (
	"fmt"
)

// Grid is a Control that arranges other Controls in a grid.
// Grid is a very powerful container: it can position and size each Control in several ways and can (and must) have Controls added to it at any time.
// [TODO it can also have Controls spanning multiple rows and columns.]
type Grid interface {
	Control

	// Add adds a Control to the Grid.
	// If this is the first Control in the Grid, it is merely added; nextTo should be nil.
	// Otherwise, it is placed relative to nextTo.
	// If nextTo is nil, it is placed next to the previously added Control,
	// The effect of adding the same Control multiple times is undefined, as is the effect of adding a Control next to one not present in the Grid.
	Add(control Control, nextTo Control, side Side, hexpand bool, halign Align, vexpand bool, valign Align)
}

// Align represents the alignment of a Control in its cell of a Grid.
type Align uint
const (
	LeftTop Align = iota
	Center
	RightBottom
	Fill
)

// Side represents a side of a Control to add other Controls to a Grid to.
type Side uint
const (
	// this arrangement is important
	// it makes finding the opposite side as easy as ^ 1
	West Side = iota
	East
	North
	South
	nSides
)

type grid struct {
	controls		map[Control]*gridCell
	prev			Control
	parent		*controlParent

	// for allocate() and preferredSize()
	xoff, yoff		int
	xmax, ymax	int
	grid			[][]gridCellAllocation
}

type gridCell struct {
	hexpand		bool
	halign		Align
	vexpand		bool
	valign		Align
	neighbors		[nSides]Control

	// for allocate() and preferredSize()
	gridx		int
	gridy		int
	width		int
	height		int
	visited		bool
}

// NewGrid creates a new Grid with no Controls.
func NewGrid() Grid {
	return &grid{
		controls:		map[Control]*gridCell{},
	}
}

func (g *grid) Add(control Control, nextTo Control, side Side, hexpand bool, halign Align, vexpand bool, valign Align) {
	cell := &gridCell{
		hexpand:		hexpand,
		halign:		halign,
		vexpand:		vexpand,
		valign:		valign,
	}
	// if this is the first control, just add it in directly
	if len(g.controls) != 0 {
		if nextTo == nil {
			nextTo = g.prev
		}
		next := g.controls[nextTo]
		// squeeze any control previously on the same side out of the way
		temp := next.neighbors[side]
		next.neighbors[side] = control
		cell.neighbors[side] = temp
		cell.neighbors[side ^ 1] = nextTo		// doubly-link
	}
	g.controls[control] = cell
	g.prev = control
	if g.parent != nil {
		control.setParent(g.parent)
	}
}

func (g *grid) setParent(p *controlParent) {
	g.parent = p
	for c, _ := range g.controls {
		c.setParent(g.parent)
	}
}

func (g *grid) trasverse(c Control, x int, y int) {
	cell := g.controls[c]
	if cell.visited {
		return
	}
	cell.visited = true
	cell.gridx = x
	cell.gridy = y
	if x < g.xoff {
		g.xoff = x
	}
	if y < g.yoff {
		g.yoff = y
	}
	if cell.neighbors[West] != nil {
		g.trasverse(cell.neighbors[West], x - 1, y)
	}
	if cell.neighbors[North] != nil {
		g.trasverse(cell.neighbors[North], x, y - 1)
	}
	if cell.neighbors[East] != nil {
		g.trasverse(cell.neighbors[East], x + 1, y)
	}
	if cell.neighbors[South] != nil {
		g.trasverse(cell.neighbors[South], x, y + 1)
	}
}

type gridCellAllocation struct {
	width	int
	height	int
	c		Control
}

func (g *grid) buildGrid() {
	// thanks to http://programmers.stackexchange.com/a/254968/147812
	// before we do anything, reset the visited bits
	for _, cell := range g.controls {
		cell.visited = false
	}
	// we first mark the previous control as the origin...
	g.xoff = 0
	g.yoff = 0
	g.trasverse(g.prev, 0, 0)		// start at the last control added
	// now we need to make all offsets zero-based
	g.xoff = -g.xoff
	g.yoff = -g.yoff
	g.xmax = 0
	g.ymax = 0
	for _, cell := range g.controls {
		cell.gridx += g.xoff
		cell.gridy += g.yoff
		if cell.gridx > g.xmax {
			g.xmax = cell.gridx
		}
		if cell.gridy > g.ymax {
			g.ymax = cell.gridy
		}
	}
	// g.xmax and g.ymax are the last valid index; make them one over to make everything work
	g.xmax++
	g.ymax++
	// and finally build the matrix
	g.grid = make([][]gridCellAllocation, g.ymax)
	for y := 0; y < g.ymax; y++ {
		g.grid[y] = make([]gridCellAllocation, g.xmax)
		// the field c is assigned below for efficiency
	}
}

func (g *grid) allocate(x int, y int, width int, height int, d *sizing) (allocations []*allocation) {
	if len(g.controls) == 0 {
		// nothing to do
		return nil
	}

	// 1) compute the resultant grid
	g.buildGrid()
	width -= d.xpadding * g.xmax
	height -= d.ypadding * g.ymax

	// 2) for every control that doesn't expand, set the width of each cell of its column/height of each cell of its row to the largest such
	nhexpand := make([]bool, g.xmax)
	nvexpand := make([]bool, g.ymax)
	for c, cell := range g.controls {
		width, height := c.preferredSize(d)
		cell.width = width
		cell.height = height
		if !cell.hexpand {
			g.grid[cell.gridy][cell.gridx].width = width
		} else {
			nhexpand[cell.gridx] = true
		}
		if !cell.vexpand {
			g.grid[cell.gridy][cell.gridx].height = height
		} else {
			nvexpand[cell.gridy] = true
		}
		g.grid[cell.gridy][cell.gridx].c = c
	}
	// cells on the same row have the same height
	for y := 0; y < g.ymax; y++ {
		max := 0
		for x := 0; x < g.xmax; x++ {
			if max < g.grid[y][x].height {
				max = g.grid[y][x].height
			}
		}
		for x := 0; x < g.xmax; x++ {
			g.grid[y][x].height = max
		}
	}
	// cells on the same column have the same width
	for x := 0; x < g.xmax; x++ {
		max := 0
		for y := 0; y < g.ymax; y++ {
			if max < g.grid[y][x].width {
				max = g.grid[y][x].width
			}
		}
		for y := 0; y < g.ymax; y++ {
			g.grid[y][x].width = max
		}
	}

	// 3) distribute the remaining space equally to expanding cells, adjusting widths and heights as needed
	nh := 0
	for x, b := range nhexpand {
		if b {
			nh++
		} else {		// column width known; subtract it
			width -= g.grid[0][x].width
		}
	}
	if nh > 0 {
		h := width / nh
		for y, b := range nhexpand {
			if b {
				for x := 0; x < g.xmax; x++ {
					if g.grid[y][x].width < h {
						g.grid[y][x].width = h
					}
				}
			}
		}
	}
	nv := 0
	for y, b := range nvexpand {
		if b {
			nv++
		} else {		// column height known; subtract it
			height -= g.grid[y][0].height
		}
	}
	if nv > 0 {
		v := height / nv
		for x, b := range nvexpand {
			if b {
				for y := 0; y < g.ymax; y++ {
					if g.grid[y][x].height < v {
						g.grid[y][x].height = v
					}
				}
			}
		}
	}

	// all right, now we have the size of each cell

	// 4) handle alignment
	for _, cell := range g.controls {
		if cell.hexpand {
			switch cell.halign {
			case LeftTop:
				// do nothing; this is the default
			case Center:
			case RightBottom:
				// TODO
			case Fill:
				cell.width = g.grid[cell.gridy][cell.gridx].width
			default:
				panic(fmt.Errorf("invalid halign %d in Grid.allocate()", cell.halign))
			}
		}
		if cell.vexpand {
			switch cell.valign {
			case LeftTop:
				// do nothing; this is the default
			case Center:
			case RightBottom:
				// TODO
			case Fill:
				cell.height = g.grid[cell.gridy][cell.gridx].height
			default:
				panic(fmt.Errorf("invalid valign %d in Grid.allocate()", cell.valign))
			}
		}
	}

	// 5) draw
	var current *allocation

	startx := x
	for row, xcol := range g.grid {
		current = nil
		for col, ca := range xcol {
			cell := g.controls[ca.c]
			as := ca.c.allocate(x, y, cell.width, cell.height, d)
			if current != nil {			// connect first left to first right
				current.neighbor = ca.c
			}
			if len(as) != 0 {
				current = as[0]			// next left is first subwidget
			} else {
				current = nil			// spaces don't have allocation data
			}
			allocations = append(allocations, as...)
			x += g.grid[0][col].width + d.xpadding
		}
		x = startx
		y += g.grid[row][0].height + d.ypadding
	}

	return allocations
}

func (g *grid) preferredSize(d *sizing) (width, height int) {
	if len(g.controls) == 0 {
		// nothing to do
		return 0, 0
	}

	// 1) compute the resultant grid
	g.buildGrid()

	// 2) for every control (including those that don't expand), set the width of each cell of its column/height of each cell of its row to the largest such
	for c, cell := range g.controls {
		width, height := c.preferredSize(d)
		g.grid[cell.gridy][cell.gridx].width = width
		g.grid[cell.gridy][cell.gridx].height = height
	}
	// cells on the same row have the same height
	maxy := 0
	for y := 0; y < g.ymax; y++ {
		max := 0
		for x := 0; x < g.xmax; x++ {
			if max < g.grid[y][x].height {
				max = g.grid[y][x].height
			}
		}
		for x := 0; x < g.xmax; x++ {
			g.grid[y][x].height = max
		}
		maxy += max
	}
	// cells on the same column have the same width
	maxx := 0
	for x := 0; x < g.xmax; x++ {
		max := 0
		for y := 0; y < g.ymax; y++ {
			if max < g.grid[y][x].width {
				max = g.grid[y][x].width
			}
		}
		for y := 0; y < g.ymax; y++ {
			g.grid[y][x].width = max
		}
		maxx += max
	}

	// and that's it really; just discount the padding
	return maxx + (g.xmax - 1) * d.xpadding,
		maxy + (g.ymax - 1) * d.ypadding
}

func (g *grid) commitResize(a *allocation, d *sizing) {
	// do nothing; needed to satisfy Control
}

func (g *grid) getAuxResizeInfo(d *sizing) {
	// do nothing; needed to satisfy Control
}
