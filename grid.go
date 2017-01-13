// 14 june 2016

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

type Align int
const (
	AlignFill Align = iota
	AlignStart
	AlignCenter
	AlignEnd
)

type At int
const (
	AtLeading At = iota
	AtTop
	AtTrailing
	AtBottom
)

// Grid is a container Control that arranges controls in rows and columns, with
// stretchy ("expanding") rows, stretchy ("expanding") columns, cells that span
// rows and columns, and cells whose content is aligned in either direction
// rather than just filling.
type Grid struct {
	c	*C.uiControl
	g	*C.uiGrid

	children	[]Control
}

// NewGrid creates a new Grid.
func NewGrid() *Grid {
	g := new(Grid)

	g.g = C.uiNewGrid()
	g.c = (*C.uiControl)(unsafe.Pointer(g.g))

	return g
}

// Destroy destroys the Grid. If the Grid has children,
// Destroy calls Destroy on those Controls as well.
func (g *Grid) Destroy() {
	for len(g.children) != 0 {
		c := g.children[0]
		g.Delete(0)
		c.Destroy()
	}
	C.uiControlDestroy(g.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Grid. This is only used by package ui itself and should
// not be called by programs.
func (g *Grid) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(g.c))
}

// Handle returns the OS-level handle associated with this Grid.
func (g *Grid) Handle() uintptr {
	return uintptr(C.uiControlHandle(g.c))
}

// Show shows the Grid.
func (g *Grid) Show() {
	C.uiControlShow(g.c)
}

// Hide hides the Grid.
func (g *Grid) Hide() {
	C.uiControlHide(g.c)
}

// Enable enables the Grid.
func (g *Grid) Enable() {
	C.uiControlEnable(g.c)
}

// Disable disables the Grid.
func (g *Grid) Disable() {
	C.uiControlDisable(g.c)
}

// Append adds the given control to the end of the Grid.
func (g *Grid) Append(child Control, left, top, xspan, yspan int, hexpand bool, uialign Align, vexpand bool, valign Align) {
	c := (*C.uiControl)(nil)
	if child != nil {
		c = touiControl(child.LibuiControl())
	}
	C.uiGridAppend(g.g, c, C.int(left), C.int(top), C.int(xspan), C.int(yspan), frombool(hexpand), C.uiAlign(uialign), frombool(vexpand), C.uiAlign(valign))
	g.children = append(g.children, child)
}

// Delete deletes the nth control of the Grid.
func (g *Grid) Delete(n int) {
	g.children = append(g.children[:n], g.children[n + 1:]...)
	//C.uiGridDelete(g.g, C.int(n))
}

// TODO: InsertAt

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
