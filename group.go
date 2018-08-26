// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Group is a Control that holds another Control and wraps it around
// a labelled box (though some systems make this box invisible).
// You can use this to group related controls together.
type Group struct {
	ControlBase
	g	*C.uiGroup
	child		Control
}

// NewGroup creates a new Group.
func NewGroup(title string) *Group {
	g := new(Group)

	ctitle := C.CString(title)
	g.g = C.uiNewGroup(ctitle)
	freestr(ctitle)

	g.ControlBase = NewControlBase(g, uintptr(unsafe.Pointer(g.g)))
	return g
}

// Destroy destroys the Group. If the Group has a child,
// Destroy calls Destroy on that as well.
func (g *Group) Destroy() {
	if g.child != nil {
		c := g.child
		g.SetChild(nil)
		c.Destroy()
	}
	g.ControlBase.Destroy()
}

// Title returns the Group's title.
func (g *Group) Title() string {
	ctitle := C.uiGroupTitle(g.g)
	title := C.GoString(ctitle)
	C.uiFreeText(ctitle)
	return title
}

// SetTitle sets the Group's title to title.
func (g *Group) SetTitle(title string) {
	ctitle := C.CString(title)
	C.uiGroupSetTitle(g.g, ctitle)
	freestr(ctitle)
}

// SetChild sets the Group's child to child. If child is nil, the Group
// will not have a child.
func (g *Group) SetChild(child Control) {
	g.child = child
	c := (*C.uiControl)(nil)
	if g.child != nil {
		c = touiControl(g.child.LibuiControl())
	}
	C.uiGroupSetChild(g.g, c)
}

// Margined returns whether the Group has margins around its child.
func (g *Group) Margined() bool {
	return tobool(C.uiGroupMargined(g.g))
}

// SetMargined controls whether the Group has margins around its
// child. The size of the margins are determined by the OS and its
// best practices.
func (g *Group) SetMargined(margined bool) {
	C.uiGroupSetMargined(g.g, frombool(margined))
}
