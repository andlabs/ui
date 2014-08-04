// 30 july 2014

package ui

// #include "objc_darwin.h"
import "C"

// all Controls that call base methods must be this
type controlPrivate interface {
	id() C.id
	Control
}

type controlParent struct {
	id	C.id
}

func basesetParent(c controlPrivate, p *controlParent) {
	// redrawing the new window handled by C.parent()
	C.parent(c.id(), p.id)
}

func basecontainerShow(c controlPrivate) {
	C.controlSetHidden(c.id(), C.NO)
}

func basecontainerHide(c controlPrivate) {
	C.controlSetHidden(c.id(), C.YES)
}

func basepreferredSize(c controlPrivate, d *sizing) (int, int) {
	s := C.controlPrefSize(c.id())
	return int(s.width), int(s.height)
}

func basecommitResize(c controlPrivate, a *allocation, d *sizing) {
	dobasecommitResize(c.id(), a, d)
}

func dobasecommitResize(id C.id, c *allocation, d *sizing) {
	C.moveControl(id, C.intptr_t(c.x), C.intptr_t(c.y), C.intptr_t(c.width), C.intptr_t(c.height))
}

func basegetAuxResizeInfo(c controlPrivate, d *sizing) {
	id := c.id()
	d.neighborAlign = C.alignmentInfo(id, C.frame(id))
}

type scroller struct {
	id	C.id
}

func newScroller(child C.id) *scroller {
	id := C.newScrollView(child)
	s := &scroller{
		id:	id,
	}
	return s
}

func (s *scroller) setParent(p *controlParent) {
	C.parent(s.id, p.id)
}

func (s *scroller) commitResize(c *allocation, d *sizing) {
	dobasecommitResize(s.id, c, d)
}
