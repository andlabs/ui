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
}

func newTab() Tab {
	t := &tab{
		controlSingleObject:		newControlSingleObject(C.newTab()),
	}
	t.fpreferredSize = t.xpreferredSize
	return t
}

func (t *tab) Append(name string, control Control) {
	c := newContainer(control.resize)
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

// no need to handle resize; the children containers handle that for us
