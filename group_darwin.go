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
}

func newGroup(text string, control Control) Group {
	g := new(group)
	g.container = newContainer()
	g.controlSingleObject = newControlSingleObject(C.newGroup(g.container.id))
	g.child = control
	g.child.setParent(g.container.parent())
	g.container.resize = g.child.resize
	g.SetText(text)
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
	return g.container.margined
}

func (g *group) SetMargined(margined bool) {
	g.container.margined = margined
}

// no need to override resize; the child container handles that for us
