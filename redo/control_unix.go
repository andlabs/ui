// +build !windows,!darwin

// 30 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type controlbase struct {
	*controldefs
	widget		*C.GtkWidget
}

type controlParent struct {
	c	*C.GtkContainer
}

func newControl(widget *C.GtkWidget) *controlbase {
	c := new(controlbase)
	c.widget = widget
	c.controldefs = new(controldefs)
	c.fsetParent = func(p *controlParent) {
		C.gtk_container_add(p.c, c.widget)
		// make sure the new widget is shown if not explicitly hidden
		c.containerShow()
	}
	c.fcontainerShow = func() {
		C.gtk_widget_show_all(c.widget)
	}
	c.fcontainerHide = func() {
		C.gtk_widget_hide(c.widget)
	}
	c.fallocate = baseallocate(c)
	c.fpreferredSize = func(d *sizing) (int, int) {
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

		C.gtk_widget_get_preferred_size(c.widget, nil, &r)
		return int(r.width), int(r.height)
	}
	c.fcommitResize = func(a *allocation, d *sizing) {
		// as we resize on size-allocate, we have to also use size-allocate on our children
		// this is fine anyway; in fact, this allows us to move without knowing what the container is!
		// this is what GtkBox does anyway
		// thanks to tristan in irc.gimp.net/#gtk+

		var r C.GtkAllocation

		r.x = C.int(a.x)
		r.y = C.int(a.y)
		r.width = C.int(a.width)
		r.height = C.int(a.height)
		C.gtk_widget_size_allocate(c.widget, &r)
	}
	c.fgetAuxResizeInfo = func(d *sizing) {
		// controls set this to true if a Label to its left should be vertically aligned to the control's top
		d.shouldVAlignTop = false
	}
	return c
}

type scrolledcontrol struct {
	*controlbase
	scroller			*controlbase
	scrollcontainer		*C.GtkContainer
	scrollwindow		*C.GtkScrolledWindow
}

func newScrolledControl(widget *C.GtkWidget, native bool) *scrolledcontrol {
	scroller := C.gtk_scrolled_window_new(nil, nil)
	s := &scrolledcontrol{
		controlbase:		newControl(widget),
		scroller:			newControl(scroller),
		scrollcontainer:	(*C.GtkContainer)(unsafe.Pointer(scroller)),
		scrollwindow:		(*C.GtkScrolledWindow)(unsafe.Pointer(scroller)),
	}
	// give the scrolled window a border (thanks to jlindgren in irc.gimp.net/#gtk+)
	C.gtk_scrolled_window_set_shadow_type(s.scrollwindow, C.GTK_SHADOW_IN)
	C.gtk_container_add(s.scrollcontainer, s.widget)
	s.fsetParent = s.scroller.fsetParent
	s.fcommitResize = s.scroller.fcommitResize
	return s
}
