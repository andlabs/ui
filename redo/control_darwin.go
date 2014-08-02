// 30 july 2014

package ui

// #include "objc_darwin.h"
import "C"

type controlbase struct {
	*controldefs
	id	C.id
}

type controlParent struct {
	id	C.id
}

func newControl(id C.id) *controlbase {
	c := new(controlbase)
	c.id = id
	c.controldefs = new(controldefs)
	c.fsetParent = func(p *controlParent) {
		// redrawing the new window handled by C.parent()
		C.parent(c.id, p.id)
	}
	c.fcontainerShow = func() {
		C.controlSetHidden(c.id, C.NO)
	}
	c.fcontainerHide = func() {
		C.controlSetHidden(c.id, C.YES)
	}
	c.fallocate = baseallocate(c)
	c.fpreferredSize = func(d *sizing) (int, int) {
		s := C.controlPrefSize(c.id)
		return int(s.width), int(s.height)
	}
	c.fcommitResize = func(a *allocation, d *sizing) {
		C.moveControl(c.id, C.intptr_t(a.x), C.intptr_t(a.y), C.intptr_t(a.width), C.intptr_t(a.height))
	}
	c.fgetAuxResizeInfo = func(d *sizing) {
		d.neighborAlign = C.alignmentInfo(c.id, C.frame(c.id))
	}
	return c
}

type scrolledcontrol struct {
	*controlbase
	scroller			*controlbase
}

func newScrolledControl(id C.id) *scrolledcontrol {
	scroller := C.newScrollView(id)
	s := &scrolledcontrol{
		controlbase:		newControl(id),
		scroller:			newControl(scroller),
	}
	s.fsetParent = s.scroller.fsetParent
	s.fcommitResize = s.scroller.fcommitResize
	return s
}
