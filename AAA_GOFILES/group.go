// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// Group is a Control that holds another Control and wraps it around
// a labelled box (though some systems make this box invisible).
// You can use this to group related controls together.
type Group struct {
	c	*C.uiControl
	g	*C.uiGroup

	child		Control
}

// NewGroup creates a new Group.
func NewGroup(title string) *Group {
	g := new(Group)

	ctitle := C.CString(title)
	g.g = C.uiNewGroup(ctitle)
	g.c = (*C.uiControl)(unsafe.Pointer(g.g))
	freestr(ctitle)

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
	C.uiControlDestroy(g.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Group. This is only used by package ui itself and should
// not be called by programs.
func (g *Group) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(g.c))
}

// Handle returns the OS-level handle associated with this Group.
// On Windows this is an HWND of a standard Windows API BUTTON
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkFrame.
// On OS X this is a pointer to a NSBox.
func (g *Group) Handle() uintptr {
	return uintptr(C.uiControlHandle(g.c))
}

// Show shows the Group.
func (g *Group) Show() {
	C.uiControlShow(g.c)
}

// Hide hides the Group.
func (g *Group) Hide() {
	C.uiControlHide(g.c)
}

// Enable enables the Group.
func (g *Group) Enable() {
	C.uiControlEnable(g.c)
}

// Disable disables the Group.
func (g *Group) Disable() {
	C.uiControlDisable(g.c)
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
