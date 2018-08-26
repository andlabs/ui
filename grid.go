// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Grid is a Control that arranges other Controls in a grid.
// Grid is a very powerful container: it can position and size each
// Control in several ways and can (and must) have Controls added
// to it in any direction. It can also have Controls spanning multiple
// rows and columns.
//
// Each Control in a Grid has associated "expansion" and
// "alignment" values in both the X and Y direction.
// Expansion determines whether all cells in the same row/column
// are given whatever space is left over after figuring out how big
// the rest of the Grid should be. Alignment determines the position
// of a Control relative to its cell after computing the above. The
// special alignment Fill can be used to grow a Control to fit its cell.
// Note that expansion and alignment are independent variables.
// For more information on expansion and alignment, read
// https://developer.gnome.org/gtk3/unstable/ch28s02.html.
type Grid struct {
	ControlBase
	g	*C.uiGrid
	children	[]Control
}

// Align represents the alignment of a Control in its cell of a Grid.
type Align int
const (
	AlignFill Align = iota
	AlignStart
	AlignCenter
	AlignEnd
)

// At represents a side of a Control to add other Controls to a Grid to.
type At int
const (
	Leading At = iota
	Top
	Trailing
	Bottom
)

// NewGrid creates a new Grid.
func NewGrid() *Grid {
	g := new(Grid)

	g.g = C.uiNewGrid()

	g.ControlBase = NewControlBase(g, uintptr(unsafe.Pointer(g.g)))
	return g
}

// TODO Destroy

// Append adds the given control to the Grid, at the given coordinate.
func (g *Grid) Append(child Control, left, top int, xspan, yspan int, hexpand bool, halign Align, vexpand bool, valign Align) {
	C.uiGridAppend(g.g, touiControl(child.LibuiControl()),
		C.int(left), C.int(top),
		C.int(xspan), C.int(yspan),
		frombool(hexpand), C.uiAlign(halign),
		frombool(vexpand), C.uiAlign(valign))
	g.children = append(g.children, child)
}

// InsertAt adds the given control to the Grid relative to an existing
// control.
func (g *Grid) InsertAt(child Control, existing Control, at At, xspan, yspan int, hexpand bool, halign Align, vexpand bool, valign Align) {
	C.uiGridInsertAt(g.g, touiControl(child.LibuiControl()),
		touiControl(existing.LibuiControl()), C.uiAt(at),
		C.int(xspan), C.int(yspan),
		frombool(hexpand), C.uiAlign(halign),
		frombool(vexpand), C.uiAlign(valign))
	g.children = append(g.children, child)
}

// Padded returns whether there is space between each control
// of the Grid.
func (g *Grid) Padded() bool {
	return tobool(C.uiGridPadded(g.g))
}

// SetPadded controls whether there is space between each control
// of the Grid. The size of the padding is determined by the OS and
// its best practices.
func (g *Grid) SetPadded(padded bool) {
	C.uiGridSetPadded(g.g, frombool(padded))
}
