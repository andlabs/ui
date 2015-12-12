// 30 july 2014

package ui

// #include "objc_darwin.h"
import "C"

type controlParent struct {
	id C.id
}

type controlSingleObject struct {
	*controlbase
	id	C.id
}

func newControlSingleObject(id C.id) *controlSingleObject {
	c := new(controlSingleObject)
	c.controlbase = &controlbase{
		fsetParent:		c.xsetParent,
		fpreferredSize:		c.xpreferredSize,
		fresize:			c.xresize,
	}
	c.id = id
	return c
}

func (c *controlSingleObject) xsetParent(p *controlParent) {
	// redrawing the new window handled by C.parent()
	C.parent(c.id, p.id)
}

func (c *controlSingleObject) xpreferredSize(d *sizing) (int, int) {
	s := C.controlPreferredSize(c.id)
	return int(s.width), int(s.height)
}

func (c *controlSingleObject) xresize(x int, y int, width int, height int, d *sizing) {
	C.moveControl(c.id, C.intptr_t(x), C.intptr_t(y), C.intptr_t(width), C.intptr_t(height))
}

type scroller struct {
	*controlSingleObject
	scroller	*controlSingleObject
}

func newScroller(child C.id, bordered bool) *scroller {
	sid := C.newScrollView(child, toBOOL(bordered))
	s := &scroller{
		controlSingleObject:		newControlSingleObject(child),
		scroller:				newControlSingleObject(sid),
	}
	s.fsetParent = s.scroller.fsetParent
	s.fresize = s .scroller.fresize
	return s
}
