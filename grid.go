// 25 february 2014
package ui

import (
	"fmt"
	"sync"
)

// A Grid arranges Controls in a two-dimensional grid.
// All Controls in a Grid maintain their preferred sizes.
// The height of each row and the width of each column is the maximum preferred height and width (respectively) of all the controls in that row or column (respectively).
// Unlike other UI toolkit Grids, this Grid does not (yet?) allow Controls to span multiple rows or columns.
type Grid struct {
	lock					sync.Mutex
	created				bool
	controls				[][]Control
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
	cw := make([][]int, nRows)
	ch := make([][]int, nRows)
	i := 0
	for row := 0; row < nRows; row++ {
		cc[row] = make([]Control, nPerRow)
		cw[row] = make([]int, nPerRow)
		ch[row] = make([]int, nPerRow)
		for x := 0; x < nPerRow; x++ {
			cc[row][x] = controls[i]
			i++
		}
	}
	return &Grid{
		controls:		cc,
		widths:		cw,
		heights:		ch,
		rowheights:	make([]int, nRows),
		colwidths:		make([]int, nPerRow),
	}
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

func (g *Grid) setRect(x int, y int, width int, height int) error {
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
			w, h, err := c.preferredSize()
			if err != nil {
				return fmt.Errorf("error getting preferred size of control (%d,%d) in Grid.setRect(): %v", row, col, err)
			}
			g.widths[row][col] = w
			g.heights[row][col] = h
			g.rowheights[row] = max(g.rowheights[row], h)
			g.colwidths[col] = max(g.colwidths[col], w)
		}
	}
	// 3) draw
	startx := x
	for row, xcol := range g.controls {
		for col, c := range xcol {
			err := c.setRect(x, y, g.widths[row][col], g.heights[row][col])
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

func (g *Grid) preferredSize() (width int, height int, err error) {
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
			w, h, err := c.preferredSize()
			if err != nil {
				return 0, 0, fmt.Errorf("error getting preferred size of control (%d,%d) in Grid.setRect(): %v", row, col, err)
			}
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
	return width, height, nil
}
