// +build !windows,!darwin

// 23 february 2014

package ui

import (
	"unsafe"
"fmt"
)

// #include "gtk_unix.h"
// extern void layoutResizing(GtkWidget *, GdkRectangle *, gpointer);
import "C"

type sizing struct {
	sizingbase

	// for size calculations
	// gtk+ needs nothing

	// for the actual resizing
	shouldVAlignTop	bool
}

const (
	gtkXMargin = 12
	gtkYMargin = 12
	gtkXPadding = 12
	gtkYPadding = 6
)

func (s *sizer) beginResize() (d *sizing) {
	d = new(sizing)
	if spaced {
		d.xmargin = gtkXMargin
		d.ymargin = gtkYMargin
		d.xpadding = gtkXPadding
		d.ypadding = gtkYPadding
	}
	return d
}

func (s *sizer) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// no need for coordinate conversion with gtk+
}

// layout maintains the widget hierarchy by containing all of a sizer's children in a single layout widget
type layout struct {
	*sizer
	layoutwidget		*C.GtkWidget
	layoutcontainer	*C.GtkContainer
	layout			*C.GtkLayout
}

func newLayout(child Control) *layout {
	widget := C.gtk_layout_new(nil, nil)
	l := &layout{
		sizer:			new(sizer),
		layoutwidget:		widget,
		layoutcontainer:	(*C.GtkContainer)(unsafe.Pointer(widget)),
		layout:			(*C.GtkLayout)(unsafe.Pointer(widget)),
	}
	l.child = child
	l.child.setParent(&controlParent{l.layoutcontainer})
	// we connect to the layout's size-allocate, not to the window's configure-event
	// this allows us to handle client-side decoration-based configurations (such as GTK+ on Wayland) properly
	// also see commitResize() in sizing_unix.go for additional notes
	// thanks to many people in irc.gimp.net/#gtk+ for help (including tristan for suggesting g_signal_connect_after())
	g_signal_connect_after(
		C.gpointer(unsafe.Pointer(l.layout)),
		"size-allocate",
		C.GCallback(C.layoutResizing),
		C.gpointer(unsafe.Pointer(l)))
	return l
}

//export layoutResizing
func layoutResizing(wid *C.GtkWidget, r *C.GdkRectangle, data C.gpointer) {
	l := (*layout)(unsafe.Pointer(data))
	// the layout's coordinate system is localized, so the origin is (0, 0)
	l.resize(0, 0, int(r.width), int(r.height))
fmt.Printf("new size %d x %d\n", r.width, r.height)
}
