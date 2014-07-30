// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type tab struct {
	*controlbase

	containers	[]*container
}

func newTab() Tab {
	t := new(tab)
	id := C.newTab(unsafe.Pointer(t))
	t.controlbase = newControl(id)
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

// no need to override Control.allocate() as only prepared the tabbed control; its children will be reallocated when that one is resized

//export tabResized
func tabResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
	t := (*tab)(unsafe.Pointer(data))
	for _, c := range t.containers {
		// the tab area's coordinate system is localized, so the origin is (0, 0)
		c.resize(0, 0, int(width), int(height))
	}
}
