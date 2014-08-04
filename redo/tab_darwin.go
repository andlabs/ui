// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type tab struct {
	_id		C.id
	tabs		[]*container
}

func newTab() Tab {
	t := new(tab)
	t._id = C.newTab(unsafe.Pointer(t))
	return t
}

func (t *tab) Append(name string, control Control) {
	c := newContainer(control)
	t.tabs = append(t.tabs, c)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.tabAppend(t._id, cname, c.view)
}

//export tabResized
func tabResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
//	t := (*tab)(unsafe.Pointer(data))
//	for _, c := range t.tabs {
		// the tab area's coordinate system is localized, so the origin is (0, 0)
//		c.resize(0, 0, int(width), int(height))
//	}
}

func (t *tab) id() C.id {
	return t._id
}

func (t *tab) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *tab) preferredSize(d *sizing) (width, height int) {
	s := C.tabPrefSize(t._id)
	return int(s.width), int(s.height)
}

// no need to override Control.commitResize() as only prepared the tabbed control; its children will be reallocated when that one is resized
func (t *tab) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *tab) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
