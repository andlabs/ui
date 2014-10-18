// 16 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type group struct {
	*controlSingleObject

	child			Control
	container		*container

	margined		bool

	chainresize	func(x int, y int, width int, height int, d *sizing)
}

func newGroup(text string, control Control) Group {
	g := new(group)
	g.container = newContainer()
	g.controlSingleObject = newControlSingleObject(C.newGroup(g.container.id))
	g.child = control
	g.child.setParent(g.container.parent())
	g.SetText(text)
	g.chainresize = g.fresize
	g.fresize = g.xresize
	return g
}

func (g *group) Text() string {
	return C.GoString(C.groupText(g.id))
}

func (g *group) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.groupSetText(g.id, ctext)
}

func (g *group) Margined() bool {
	return g.margined
}

func (g *group) SetMargined(margined bool) {
	g.margined = margined
}

func (g *group) xresize(x int, y int, width int, height int, d *sizing) {
	// first, chain up to change the GtkFrame and its child container
	g.chainresize(x, y, width, height, d)

	// now that the container has the correct size, we can resize the child
	a := g.container.allocation(g.margined)
	g.child.resize(int(a.x), int(a.y), int(a.width), int(a.height), d)
}
