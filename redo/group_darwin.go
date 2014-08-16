// 16 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type group struct {
	_id		C.id

	*container
}

func newGroup(text string, control Control) Group {
	g := new(group)
	g.container = newContainer(control)
	g._id = C.newGroup(g.container.id)
	g.SetText(text)
	return g
}

func (g *group) Text() string {
	return C.GoString(C.groupText(g._id))
}

func (g *group) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.groupSetText(g._id, ctext)
}

func (g *group) id() C.id {
	return g._id
}

func (g *group) setParent(p *controlParent) {
	basesetParent(g, p)
}

func (g *group) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(g, x, y, width, height, d)
}

func (g *group) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(g, d)
}

func (g *group) commitResize(a *allocation, d *sizing) {
	basecommitResize(g, a, d)
}

func (g *group) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(g, d)
}
