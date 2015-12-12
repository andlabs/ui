// +build !windows,!darwin

// 15 august 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

// TODO on which sides do margins get applied?

type group struct {
	*controlSingleWidget
	gcontainer *C.GtkContainer
	frame      *C.GtkFrame

	child			Control
	container		*container
}

func newGroup(text string, control Control) Group {
	ctext := togstr(text)
	defer freegstr(ctext)
	widget := C.gtk_frame_new(ctext)
	g := &group{
		controlSingleWidget:	newControlSingleWidget(widget),
		gcontainer: (*C.GtkContainer)(unsafe.Pointer(widget)),
		frame:      (*C.GtkFrame)(unsafe.Pointer(widget)),
		child:	control,
	}

	// with GTK+, groupboxes by default have frames and slightly x-offset regular text
	// they should have no frame and fully left-justified, bold text
	var yalign C.gfloat

	// preserve default y-alignment
	C.gtk_frame_get_label_align(g.frame, nil, &yalign)
	C.gtk_frame_set_label_align(g.frame, 0, yalign)
	C.gtk_frame_set_shadow_type(g.frame, C.GTK_SHADOW_NONE)
	label := (*C.GtkLabel)(unsafe.Pointer(C.gtk_frame_get_label_widget(g.frame)))
	// this is the boldness level used by GtkPrintUnixDialog
	// (it technically uses "bold" but see pango's pango-enum-types.c for the name conversion; GType is weird)
	bold := C.pango_attr_weight_new(C.PANGO_WEIGHT_BOLD)
	boldlist := C.pango_attr_list_new()
	C.pango_attr_list_insert(boldlist, bold)
	C.gtk_label_set_attributes(label, boldlist)
	C.pango_attr_list_unref(boldlist) // thanks baedert in irc.gimp.net/#gtk+

	g.container = newContainer()
	g.child.setParent(g.container.parent())
	g.container.resize = g.child.resize
	C.gtk_container_add(g.gcontainer, g.container.widget)

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

func (g *group) Margined() bool {
	return g.container.margined
}

func (g *group) SetMargined(margined bool) {
	g.container.margined = margined
}

// no need to override resize; the child container handles that for us
