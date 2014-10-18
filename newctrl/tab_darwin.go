// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type tab struct {
	*controlSingleObject
	tabs			[]*container
	children		[]Control
	chainresize	func(x int, y int, width int, height int
}

func newTab() Tab {
	t := &tab{
		controlSingleObject:		newControlSingleObject(C.newTab()),
	}
	t.fpreferredsize = t.xpreferredsize
	t.chainresize = t.fresize
	t.fresize = t.xresize
	return t
}

func (t *tab) Append(name string, control Control) {
	c := newContainer()
	t.tabs = append(t.tabs, c)
	control.setParent(c.parent())
	t.children = append(t.children, control)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.tabAppend(t.id, cname, c.id)
}

func (t *tab) xpreferredSize(d *sizing) (width, height int) {
	s := C.tabPreferredSize(t.id)
	return int(s.width), int(s.height)
}

func (t *tab) xresize(x int, y int, width int, height int, d *sizing) {
	// first, chain up to change the GtkFrame and its child container
	t.chainresize(x, y, width, height, d)

	// now that the containers have the correct size, we can resize the children
	for i, _ := range t.tabs {
		a := t.tabs[i].allocation(false/*TODO*/)
		t.children[i].resize(int(a.x), int(a.y), int(a.width), int(a.height), d)
	}
}
