// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type tab struct {
	*widgetbase

	containers	[]*container
}

func newTab() Tab {
	t := new(tab)
	id := C.newTab(unsafe.Pointer(t))
	t.widgetbase = newWidget(id)
	return t
}

func (t *tab) Append(name string, control Control) {
	// TODO isolate and standardize
	c := new(container)
	// don't set beginResize; this container's resize() will be a recursive call
	t.containers = append(t.containers, c)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	tabview := C.tabAppend(t.id, cname)
	c.child = control
	c.child.setParent(tabview)
}

func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	// set up the recursive calls
	for _, c := range t.containers {
		c.d = d
	}
	// and prepare the tabbed control itself
	return t.widgetbase.allocate(x, y, width, height, d)
}

//export tabResized
func tabResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
	t := (*tab)(unsafe.Pointer(data))
	for _, c := range t.containers {
		c.resize(int(width), int(height))
	}
}
