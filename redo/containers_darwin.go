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
	t.containers = append(t.containers, c)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	tabview := C.tabAppend(t.id, cname)
	c.child = control
	c.child.setParent(&controlParent{tabview})
}

func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	// only prepared the tabbed control; its children will be reallocated when that one is resized
	return t.widgetbase.allocate(x, y, width, height, d)
}

//export tabResized
func tabResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
	t := (*tab)(unsafe.Pointer(data))
	for _, c := range t.containers {
		// the tab area's coordinate system is localized, so the origin is (0, 0)
		c.resize(0, 0, int(width), int(height))
	}
}
