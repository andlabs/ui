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
	widget := c.widget()		// avoid multiple interface lookups
	C.gtk_container_add(p.c, widget)
	// make sure the new widget is shown if not explicitly hidden
	C.gtk_widget_show_all(widget)
}

func basepreferredSize(c controlPrivate, d *sizing) (int, int) {
	// GTK+ 3 makes this easy: controls can tell us what their preferred size is!
	// ...actually, it tells us two things: the "minimum size" and the "natural size".
	// The "minimum size" is the smallest size we /can/ display /anything/. The "natural size" is the smallest size we would /prefer/ to display.
	// The difference? Minimum size takes into account things like truncation with ellipses: the minimum size of a label can allot just the ellipses!
	// So we use the natural size instead.
	// There is a warning about height-for-width controls, but in my tests this isn't an issue.
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

	overlaywidget		*C.GtkWidget
	overlaycontainer	*C.GtkContainer
	overlay			*C.GtkOverlay

	addShowWhich		*C.GtkWidget
}

func newScroller(widget *C.GtkWidget, native bool, bordered bool, overlay bool) *scroller {
	var o *C.GtkWidget

	scrollwidget := C.gtk_scrolled_window_new(nil, nil)
	if overlay {
		o = C.gtk_overlay_new()
	}
	s := &scroller{
		scrollwidget:		scrollwidget,
		scrollcontainer:	(*C.GtkContainer)(unsafe.Pointer(scrollwidget)),
		scrollwindow:		(*C.GtkScrolledWindow)(unsafe.Pointer(scrollwidget)),
		overlaywidget:		o,
		overlaycontainer:	(*C.GtkContainer)(unsafe.Pointer(o)),
		overlay:			(*C.GtkOverlay)(unsafe.Pointer(o)),
	}
	// give the scrolled window a border (thanks to jlindgren in irc.gimp.net/#gtk+)
	if bordered {
		C.gtk_scrolled_window_set_shadow_type(s.scrollwindow, C.GTK_SHADOW_IN)
	}
	if native {
		C.gtk_container_add(s.scrollcontainer, widget)
	} else {
		C.gtk_scrolled_window_add_with_viewport(s.scrollwindow, widget)
	}
	s.addShowWhich = s.scrollwidget
	if overlay {
		C.gtk_container_add(s.overlaycontainer, s.scrollwidget)
		s.addShowWhich = s.overlaywidget
	}
	return s
}

func (s *scroller) setParent(p *controlParent) {
	C.gtk_container_add(p.c, s.addShowWhich)
	// see basesetParent() above for why we call gtk_widget_show_all()
	C.gtk_widget_show_all(s.addShowWhich)
}

func (s *scroller) commitResize(c *allocation, d *sizing) {
	dobasecommitResize(s.addShowWhich, c, d)
}
