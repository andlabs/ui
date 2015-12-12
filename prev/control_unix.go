// +build !windows,!darwin

// 30 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type controlParent struct {
	c *C.GtkContainer
}

type controlSingleWidget struct {
	*controlbase
	widget	*C.GtkWidget
}

func newControlSingleWidget(widget *C.GtkWidget) *controlSingleWidget {
	c := new(controlSingleWidget)
	c.controlbase = &controlbase{
		fsetParent:		c.xsetParent,
		fpreferredSize:		c.xpreferredSize,
		fresize:			c.xresize,
	}
	c.widget = widget
	return c
}

func (c *controlSingleWidget) xsetParent(p *controlParent) {
	C.gtk_container_add(p.c, c.widget)
	// make sure the new widget is shown if not explicitly hidden
	// TODO why did I have this again?
	C.gtk_widget_show_all(c.widget)
}

func (c *controlSingleWidget) xpreferredSize(d *sizing) (int, int) {
	// GTK+ 3 makes this easy: controls can tell us what their preferred size is!
	// ...actually, it tells us two things: the "minimum size" and the "natural size".
	// The "minimum size" is the smallest size we /can/ display /anything/. The "natural size" is the smallest size we would /prefer/ to display.
	// The difference? Minimum size takes into account things like truncation with ellipses: the minimum size of a label can allot just the ellipses!
	// So we use the natural size instead.
	// There is a warning about height-for-width controls, but in my tests this isn't an issue.
	var r C.GtkRequisition

	C.gtk_widget_get_preferred_size(c.widget, nil, &r)
	return int(r.width), int(r.height)
}

func (c *controlSingleWidget) xresize(x int, y int, width int, height int, d *sizing) {
	// as we resize on size-allocate, we have to also use size-allocate on our children
	// this is fine anyway; in fact, this allows us to move without knowing what the container is!
	// this is what GtkBox does anyway
	// thanks to tristan in irc.gimp.net/#gtk+

	var r C.GtkAllocation

	r.x = C.int(x)
	r.y = C.int(y)
	r.width = C.int(width)
	r.height = C.int(height)
	C.gtk_widget_size_allocate(c.widget, &r)
}

type scroller struct {
	*controlSingleWidget

	scroller	*controlSingleWidget
	scrollwidget    *C.GtkWidget
	scrollcontainer *C.GtkContainer
	scrollwindow    *C.GtkScrolledWindow

	overlay	*controlSingleWidget
	overlaywidget    *C.GtkWidget
	overlaycontainer *C.GtkContainer
	overlayoverlay      *C.GtkOverlay
}

func newScroller(widget *C.GtkWidget, native bool, bordered bool, overlay bool) *scroller {
	s := new(scroller)
	s.controlSingleWidget = newControlSingleWidget(widget)
	s.scrollwidget = C.gtk_scrolled_window_new(nil, nil)
	s.scrollcontainer = (*C.GtkContainer)(unsafe.Pointer(s.scrollwidget))
	s.scrollwindow = (*C.GtkScrolledWindow)(unsafe.Pointer(s.scrollwidget))

	// any actual changing operations need to be done to the GtkScrolledWindow
	// that is, everything /except/ preferredSize() are done to the GtkScrolledWindow
	s.scroller = newControlSingleWidget(s.scrollwidget)
	s.fsetParent = s.scroller.fsetParent
	s.fresize = s.scroller.fresize

	// in GTK+ 3.4 we still technically need to use the separate gtk_scrolled_window_add_with_viewpoint()/gtk_container_add() spiel for adding the widget to the scrolled window
	if native {
		C.gtk_container_add(s.scrollcontainer, s.widget)
	} else {
		C.gtk_scrolled_window_add_with_viewport(s.scrollwindow, s.widget)
	}

	// give the scrolled window a border (thanks to jlindgren in irc.gimp.net/#gtk+)
	if bordered {
		C.gtk_scrolled_window_set_shadow_type(s.scrollwindow, C.GTK_SHADOW_IN)
	}

	if overlay {
		// ok things get REALLY fun now
		// we now have to do all of the above again
		s.overlaywidget = C.gtk_overlay_new()
		s.overlaycontainer = (*C.GtkContainer)(unsafe.Pointer(s.overlaywidget))
		s.overlayoverlay = (*C.GtkOverlay)(unsafe.Pointer(s.overlaywidget))
		s.overlay = newControlSingleWidget(s.overlaywidget)
		s.fsetParent = s.overlay.fsetParent
		s.fresize = s.overlay.fresize
		C.gtk_container_add(s.overlaycontainer, s.scrollwidget)
	}

	return s
}
