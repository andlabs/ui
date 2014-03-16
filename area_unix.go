// +build !windows,!darwin,!plan9

// 14 march 2014

package ui

import (
	"unsafe"
	"image"
)

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
// extern gboolean our_area_draw_callback(GtkWidget *, cairo_t *, gpointer);
// extern gboolean our_area_button_press_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_area_button_release_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_area_motion_notify_event_callback(GtkWidget *, GdkEvent *, gpointer);
// /* HACK - see https://code.google.com/p/go/issues/detail?id=7548 */
// struct _cairo {};
import "C"

func gtkAreaNew() *gtkWidget {
	drawingarea := C.gtk_drawing_area_new()
	C.gtk_widget_set_size_request(drawingarea, 320, 240)
	// we need to explicitly subscribe to mouse events with GtkDrawingArea
	C.gtk_widget_add_events(drawingarea,
		C.GDK_BUTTON_PRESS_MASK | C.GDK_BUTTON_RELEASE_MASK | C.GDK_POINTER_MOTION_MASK | C.GDK_BUTTON_MOTION_MASK)
	scrollarea := C.gtk_scrolled_window_new((*C.GtkAdjustment)(nil), (*C.GtkAdjustment)(nil))
	// need a viewport because GtkDrawingArea isn't natively scrollable
	C.gtk_scrolled_window_add_with_viewport((*C.GtkScrolledWindow)(unsafe.Pointer(scrollarea)), drawingarea)
	return fromgtkwidget(scrollarea)
}

func gtkAreaGetControl(scrollarea *gtkWidget) *gtkWidget {
	viewport := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(scrollarea)))
	control := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(viewport)))
	return fromgtkwidget(control)
}

//export our_area_draw_callback
func our_area_draw_callback(widget *C.GtkWidget, cr *C.cairo_t, data C.gpointer) C.gboolean {
	var x, y, w, h C.double

	s := (*sysData)(unsafe.Pointer(data))
	// thanks to desrt in irc.gimp.net/#gtk+
	C.cairo_clip_extents(cr, &x, &y, &w, &h)
	cliprect := image.Rect(int(x), int(y), int(w), int(h))
	imgret := make(chan *image.NRGBA)
	defer close(imgret)
	s.paint <- PaintRequest{
		Rect:		cliprect,
		Out:		imgret,
	}
	i := <-imgret
	// pixel order is [R G B A] (see Example 1 on https://developer.gnome.org/gdk-pixbuf/2.26/gdk-pixbuf-The-GdkPixbuf-Structure.html) so we don't have to convert anything
	// gdk-pixbuf is not alpha-premultiplied (thanks to desrt in irc.gimp.net/#gtk+)
	pixbuf := C.gdk_pixbuf_new_from_data(
		(*C.guchar)(unsafe.Pointer(&i.Pix[0])),
		C.GDK_COLORSPACE_RGB,
		C.TRUE,			// has alpha channel
		8,				// bits per sample
		C.int(i.Rect.Dx()),
		C.int(i.Rect.Dy()),
		C.int(i.Stride),
		nil, nil)			// do not free data
	C.gdk_cairo_set_source_pixbuf(cr,
		pixbuf,
		C.gdouble(cliprect.Min.X),
		C.gdouble(cliprect.Min.Y))
	// that just set the brush that cairo uses: we have to actually draw now
	// (via https://developer.gnome.org/gtkmm-tutorial/stable/sec-draw-images.html.en)
	C.cairo_rectangle(cr, x, y, w, h)		// breaking the nrom here since we have the double data already
	C.cairo_fill(cr)
	C.g_object_unref((C.gpointer)(unsafe.Pointer(pixbuf)))		// free pixbuf
	return C.FALSE		// signals handled without stopping the event chain (thanks to desrt again)
}

var area_draw_callback = C.GCallback(C.our_area_draw_callback)

// shared code for finishing up and sending a mouse event
func finishMouseEvent(data C.gpointer, me MouseEvent, mb uint, x C.gdouble, y C.gdouble, state C.guint, gdkwindow *C.GdkWindow) {
	s := (*sysData)(unsafe.Pointer(data))
	// GDK doesn't initialize the modifier flags fully; we have to explicitly tell it to (thanks to Daniel_S and daniels (two different people) in irc.gimp.net/#gtk+)
	C.gdk_keymap_add_virtual_modifiers(
		C.gdk_keymap_get_for_display(C.gdk_window_get_display(gdkwindow)),
		(*C.GdkModifierType)(unsafe.Pointer(&state)))
	if (state & C.GDK_CONTROL_MASK) != 0 {
		me.Modifiers |= Ctrl
	}
	if (state & C.GDK_META_MASK) != 0 {
		me.Modifiers |= Alt
	}
	if (state & C.GDK_SHIFT_MASK) != 0 {
		me.Modifiers |= Shift
	}
	// the mb != # checks exclude the Up/Down button from Held
	if mb != 1 && (state & C.GDK_BUTTON1_MASK) != 0 {
		me.Held = append(me.Held, 1)
	}
	if mb != 2 && (state & C.GDK_BUTTON2_MASK) != 0 {
		me.Held = append(me.Held, 2)
	}
	if mb != 3 && (state & C.GDK_BUTTON3_MASK) != 0 {
		me.Held = append(me.Held, 3)
	}
	// TODO keep?
	if mb != 4 && (state & C.GDK_BUTTON4_MASK) != 0 {
		me.Held = append(me.Held, 4)
	}
	if mb != 5 && (state & C.GDK_BUTTON5_MASK) != 0 {
		me.Held = append(me.Held, 5)
	}
	me.Pos = image.Pt(int(x), int(y))
	// see cSysData.signal() in sysdata.go
	go func() {
		select {
		case s.mouse <- me:
		default:
		}
	}()
}

//export our_area_button_press_event_callback
func our_area_button_press_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	e := (*C.GdkEventButton)(unsafe.Pointer(event))
	me := MouseEvent{
		// GDK button ID == our button ID
		Down:	uint(e.button),
	}
	switch e._type {
	case C.GDK_BUTTON_PRESS:
		me.Count = 1
	case C.GDK_2BUTTON_PRESS:
		me.Count = 2
	default:		// ignore triple-clicks and beyond; we don't handle those
		return C.FALSE		// TODO really false?
	}
	finishMouseEvent(data, me, me.Down, e.x, e.y, e.state, e.window)
	return C.FALSE			// TODO really false?
}

var area_button_press_event_callback = C.GCallback(C.our_area_button_press_event_callback)

//export our_area_button_release_event_callback
func our_area_button_release_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	e := (*C.GdkEventButton)(unsafe.Pointer(event))
	me := MouseEvent{
		// GDK button ID == our button ID
		Up:		uint(e.button),
	}
	finishMouseEvent(data, me, me.Up, e.x, e.y, e.state, e.window)
	return C.FALSE			// TODO really false?
}

var area_button_release_event_callback = C.GCallback(C.our_area_button_release_event_callback)

//export our_area_motion_notify_event_callback
func our_area_motion_notify_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	e := (*C.GdkEventMotion)(unsafe.Pointer(event))
	me := MouseEvent{}
	finishMouseEvent(data, me, 0, e.x, e.y, e.state, e.window)
	return C.FALSE			// TODO really false?
}

var area_motion_notify_event_callback = C.GCallback(C.our_area_motion_notify_event_callback)
