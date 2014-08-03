// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type tab struct {
	*controlbase

	tabs			[]*sizer
}

func newTab() Tab {
	t := new(tab)
	id := C.newTab(unsafe.Pointer(t))
	t.controlbase = newControl(id)
	t.fpreferredSize = t.tabpreferredSize
	return t
}

func (t *tab) Append(name string, control Control) {
	s := new(sizer)
	t.tabs = append(t.tabs, s)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	tabview := C.tabAppend(t.id, cname)
	s.child = control
	s.child.setParent(&controlParent{tabview})
}

func (t *tab) tabpreferredSize(d *sizing) (width, height int) {
	s := C.tabPrefSize(t.id)
	return int(s.width), int(s.height)
}

// no need to override Control.commitResize() as only prepared the tabbed control; its children will be reallocated when that one is resized

//export tabResized
func tabResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
	t := (*tab)(unsafe.Pointer(data))
	for _, s := range t.tabs {
		// the tab area's coordinate system is localized, so the origin is (0, 0)
		s.resize(0, 0, int(width), int(height))
	}
}
