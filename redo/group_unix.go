// +build !windows,!darwin

// 15 august 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type group struct {
	_widget		*C.GtkWidget
	gcontainer	*C.GtkContainer
	frame		*C.GtkFrame

	*container
}

func newGroup(text string, control Control) Group {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_frame_new(ctext)
	g := &group{
		_widget:		widget,
		gcontainer:	(*C.GtkContainer)(unsafe.Pointer(widget)),
		frame:		(*C.GtkFrame)(unsafe.Pointer(widget)),
	}

	// with GTK+, groupboxes by default have frames and slightly x-offset regular text
	// they should have no frame and fully left-justified, bold text
	var yalign C.gfloat

	// preserve default y-alignment
	C.gtk_frame_get_label_align(g.frame, nil, &yalign)
	C.gtk_frame_set_label_align(g.frame, 0, yalign)
	C.gtk_frame_set_shadow_type(g.frame, C.GTK_SHADOW_NONE)
	label := (*C.GtkLabel)(unsafe.Pointer(C.gtk_frame_get_label_widget(g.frame)))
	// TODO confirm this boldness level against Glade
	bold := C.pango_attr_weight_new(C.PANGO_WEIGHT_BOLD)
	boldlist := C.pango_attr_list_new()
	C.pango_attr_list_insert(boldlist, bold)
	C.gtk_label_set_attributes(label, boldlist)
	// TODO free either bold or boldlist?

	g.container = newContainer(control)
	g.container.setParent(&controlParent{g.gcontainer})

	return g
}

func (g *group) Text() string {
	return fromgstr(C.gtk_frame_get_label(g.frame))
}

func (g *group) SetText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_frame_set_label(g.frame, ctext)
}

func (g *group) widget() *C.GtkWidget {
	return g._widget
}

func (g *group) setParent(p *controlParent) {
	basesetParent(g, p)
}

func (g *group) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(g, x, y, width, height, d)
}

func (g *group) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(g, d)
}

func (g *group) commitResize(a *allocation, d *sizing) {
	basecommitResize(g, a, d)
}

func (g *group) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(g, d)
}
