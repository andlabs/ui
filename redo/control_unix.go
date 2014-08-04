// +build !windows,!darwin

// 30 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

// all Controls that call base methods must be this
type controlPrivate interface {
	widget() *C.GtkWidget
	Control
}

type controlParent struct {
	c	*C.GtkContainer
}

func basesetParent(c controlPrivate, p *controlParent) {
	C.gtk_container_add(p.c, c.widget())
	// make sure the new widget is shown if not explicitly hidden
	c.containerShow()
}

func basecontainerShow(c controlPrivate) {
	C.gtk_widget_show_all(c.widget())
}

func basecontainerHide(c controlPrivate) {
	C.gtk_widget_hide(c.widget())
}

func basepreferredSize(c controlPrivate, d *sizing) (int, int) {
	// GTK+ 3 makes this easy: controls can tell us what their preferred size is!
	// ...actually, it tells us two things: the "minimum size" and the "natural size".
	// The "minimum size" is the smallest size we /can/ display /anything/. The "natural size" is the smallest size we would /prefer/ to display.
	// The difference? Minimum size takes into account things like truncation with ellipses: the minimum size of a label can allot just the ellipses!
	// So we use the natural size instead.
	// There is a warning about height-for-width controls, but in my tests this isn't an issue.
	// For Areas, we manually save the Area size and use that, just to be safe.

//TODO
/*
	if s.ctype == c_area {
		return s.areawidth, s.areaheight
	}
*/

	var r C.GtkRequisition

	C.gtk_widget_get_preferred_size(c.widget(), nil, &r)
	return int(r.width), int(r.height)
}

func basecommitResize(c controlPrivate, a *allocation, d *sizing) {
	dobasecommitResize(c.widget(), a, d)
}

func dobasecommitResize(w *C.GtkWidget, c *allocation, d *sizing) {
	// as we resize on size-allocate, we have to also use size-allocate on our children
	// this is fine anyway; in fact, this allows us to move without knowing what the container is!
	// this is what GtkBox does anyway
	// thanks to tristan in irc.gimp.net/#gtk+

	var r C.GtkAllocation

	r.x = C.int(c.x)
	r.y = C.int(c.y)
	r.width = C.int(c.width)
	r.height = C.int(c.height)
	C.gtk_widget_size_allocate(w, &r)
}

func basegetAuxResizeInfo(c Control, d *sizing) {
	// controls set this to true if a Label to its left should be vertically aligned to the control's top
	d.shouldVAlignTop = false
}

type scroller struct {
	scrollwidget		*C.GtkWidget
	scrollcontainer		*C.GtkContainer
	scrollwindow		*C.GtkScrolledWindow
}

func newScroller(widget *C.GtkWidget, native bool) *scroller {
	scrollwidget := C.gtk_scrolled_window_new(nil, nil)
	s := &scroller{
		scrollwidget:		scrollwidget,
		scrollcontainer:	(*C.GtkContainer)(unsafe.Pointer(scrollwidget)),
		scrollwindow:		(*C.GtkScrolledWindow)(unsafe.Pointer(scrollwidget)),
	}
	// give the scrolled window a border (thanks to jlindgren in irc.gimp.net/#gtk+)
	C.gtk_scrolled_window_set_shadow_type(s.scrollwindow, C.GTK_SHADOW_IN)
	// TODO use native here
	C.gtk_container_add(s.scrollcontainer, widget)
	return s
}

func (s *scroller) setParent(p *controlParent) {
	C.gtk_container_add(p.c, s.scrollwidget)
	// TODO for when hiding/showing is implemented
	C.gtk_widget_show_all(s.scrollwidget)
}

func (s *scroller) commitResize(c *allocation, d *sizing) {
	dobasecommitResize(s.scrollwidget, c, d)
}
