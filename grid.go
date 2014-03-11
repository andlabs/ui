// 25 february 2014
package ui

import (
	"fmt"
	"sync"
)

// A Grid arranges Controls in a two-dimensional grid.
// The height of each row and the width of each column is the maximum preferred height and width (respectively) of all the controls in that row or column (respectively).
// Controls are aligned to the top left corner of each cell.
// All Controls in a Grid maintain their preferred sizes by default; if a Control is marked as being "filling", it will be sized to fill its cell.
// Even if a Control is marked as filling, its preferred size is used to calculate cell sizes.
// One Control can be marked as "stretchy": when the Window containing the Grid is resized, the cell containing that Control resizes to take any remaining space; its row and column are adjusted accordingly (so other filling controls in the same row and column will fill to the new height and width, respectively).
// A stretchy Control implicitly fills its cell.
// All cooridnates in a Grid are given in (row,column) form with (0,0) being the top-left cell.
// Unlike other UI toolkit Grids, this Grid does not (yet? TODO) allow Controls to span multiple rows or columns.
// TODO differnet row/column control alignment
type Grid struct {
	lock					sync.Mutex
	created				bool
	controls				[][]Control
	filling				[][]bool
	stretchyrow, stretchycol	int
	widths, heights			[][]int		// caches to avoid reallocating each time
	rowheights, colwidths	[]int
}

// NewGrid creates a new Grid with the given Controls.
// NewGrid needs to know the number of Controls in a row (alternatively, the number of columns); it will determine the number in a column from the number of Controls given.
// NewGrid panics if not given a full grid of Controls.
// Example:
// 	grid := NewGrid(3,
// 		control00, control01, control02,
// 		control10, control11, control12,
// 		control20, control21, control22)
func NewGrid(nPerRow int, controls ...Control) *Grid {
	if len(controls) % nPerRow != 0 {
		panic(fmt.Errorf("incomplete grid given to NewGrid() (not enough controls to evenly divide %d controls into rows of %d controls each)", len(controls), nPerRow))
	}
	nRows := len(controls) / nPerRow
	cc := make([][]Control, nRows)
	cf := make([][]bool, nRows)
	cw := make([][]int, nRows)
	ch := make([][]int, nRows)
	i := 0
	for row := 0; row < nRows; row++ {
		cc[row] = make([]Control, nPerRow)
		cf[row] = make([]bool, nPerRow)
		cw[row] = make([]int, nPerRow)
		ch[row] = make([]int, nPerRow)
		for x := 0; x < nPerRow; x++ {
			cc[row][x] = controls[i]
			i++
		}
	}
	return &Grid{
		controls:		cc,
		filling:		cf,
		stretchyrow:	-1,
		stretchycol:	-1,
		widths:		cw,
		heights:		ch,
		rowheights:	make([]int, nRows),
		colwidths:		make([]int, nPerRow),
	}
}

// SetFilling marks the given Control of the Grid as filling its cell instead of staying at its preferred size.
// This function cannot be called after the Window that contains the Grid has been created.
// It panics if the given coordinate is invalid.
func (g *Grid) SetFilling(row int, column int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.created {
		panic(fmt.Errorf("Grid.SetFilling() called after window create"))
	}
	if row < 0 || column < 0 || row > len(g.filling) || column > len(g.filling[row]) {
		panic(fmt.Errorf("coordinate (%d,%d) out of range passed to Grid.SetFilling()", row, column))
	}
	g.filling[row][column] = true
}

// SetStretchy marks the given Control of the Grid as stretchy.
// Stretchy implies filling.
// Only one control can be stretchy per Grid; calling SetStretchy multiple times merely changes which control is stretchy.
// This function cannot be called after the Window that contains the Grid has been created.
// It panics if the given coordinate is invalid.
func (g *Grid) SetStretchy(row int, column int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.created {
		panic(fmt.Errorf("Grid.SetFilling() called after window create"))
	}
	if row < 0 || column < 0 || row > len(g.filling) || column > len(g.filling[row]) {
		panic(fmt.Errorf("coordinate (%d,%d) out of range passed to Grid.SetStretchy()", row, column))
	}
	g.stretchyrow = row
	g.stretchycol = column
	g.filling[row][column] = true
}

func (g *Grid) make(window *sysData) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	for row, xcol := range g.controls {
		for col, c := range xcol {
			err := c.make(window)
			if err != nil {
				return fmt.Errorf("error adding control (%d,%d) to Grid: %v", row, col, err)
			}
		}
	}
	g.created = true
	return nil
}

func (g *Grid) setRect(x int, y int, width int, height int, winheight int) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	max := func(a int, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// 1) clear data structures
	for i := range g.rowheights {
		g.rowheights[i] = 0
	}
	for i := range g.colwidths {
		g.colwidths[i] = 0
	}
	// 2) get preferred sizes; compute row/column sizes
	for row, xcol := range g.controls {
		for col, c := range xcol {
			w, h := c.preferredSize()
			g.widths[row][col] = w
			g.heights[row][col] = h
			g.rowheights[row] = max(g.rowheights[row], h)
			g.colwidths[col] = max(g.colwidths[col], w)
		}
	}
	// 3) handle the stretchy control
	if g.stretchyrow != -1 && g.stretchycol != -1 {
		for i, w := range g.colwidths {
			if i != g.stretchycol {
				width -= w
			}
		}
		for i, h := range g.rowheights {
			if i != g.stretchyrow {
				height -= h
			}
		}
		g.colwidths[g.stretchycol] = width
		g.rowheights[g.stretchyrow] = height
	}
	// TODO add a sanity check for g.stretchyrow xor g.stretchycol == -1?
	// 4) draw
	startx := x
	for row, xcol := range g.controls {
		for col, c := range xcol {
			w := g.widths[row][col]
			h := g.heights[row][col]
			if g.filling[row][col] {
				w = g.colwidths[col]
				h = g.rowheights[row]
			}
			err := c.setRect(x, y, w, h, winheight)
			if err != nil {
				return fmt.Errorf("error setting size of control (%d,%d) in Grid.setRect(): %v", row, col, err)
			}
			x += g.colwidths[col]
		}
		x = startx
		y += g.rowheights[row]
	}
	return nil
}

// filling and stretchy are ignored for preferred size calculation
func (g *Grid) preferredSize() (width int, height int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	max := func(a int, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// 1) clear data structures
	for i := range g.rowheights {
		g.rowheights[i] = 0
	}
	for i := range g.colwidths {
		g.colwidths[i] = 0
	}
	// 2) get preferred sizes; compute row/column sizes
	for row, xcol := range g.controls {
		for col, c := range xcol {
			w, h := c.preferredSize()
			g.widths[row][col] = w
			g.heights[row][col] = h
			g.rowheights[row] = max(g.rowheights[row], h)
			g.colwidths[col] = max(g.colwidths[col], w)
		}
	}
	// 3) now compute
	for _, w := range g.colwidths {
		width += w
	}
	for _, h := range g.rowheights {
		height += h
	}
	return width, height
}
