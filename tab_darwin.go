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
	return &tab{
		_id:		C.newTab(),
	}
}

func (t *tab) Append(name string, control Control) {
	c := newContainer(control)
	t.tabs = append(t.tabs, c)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.tabAppend(t._id, cname, c.id)
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
	s := C.tabPreferredSize(t._id)
	return int(s.width), int(s.height)
}

// no need to override Control.commitResize() as only prepared the tabbed control; its children will be resized when that one is resized (and NSTabView itself will call setFrame: for us)
func (t *tab) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *tab) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
