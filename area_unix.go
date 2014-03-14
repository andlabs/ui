// +build ignore
// (ignored pending Go issues 7409 and 7548)
// x+build !windows,!darwin,!plan9

// 14 march 2014

package ui

import (
	"image"
)

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// extern gboolean our_draw_callback(GtkWidget *, cairo_t *, gpointer);
import "C"

func gtkAreaNew() *gtkWidget {
	drawingarea := C.gtk_drawing_area_new()
	scrollarea := C.gtk_scrolled_window_new((*C.GtkAdjustment)(nil), (*C.GtkAdjustment)(nil))
	// need a viewport because GtkDrawingArea isn't natively scrollable
	C.gtk_scrolled_window_add_with_viewport(scrollarea, drawingarea)
	return fromgtkwidget(scrollarea)
}

//export our_draw_callback
func our_draw_callback(widget *C.GtkWidget, cr *C.cairo_t, data C.gpointer) C.gboolean {
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
	C.g_object_unref((C.gpointer)(unsafe.Pointer(pixbuf)))		// free pixbuf
	return C.FALSE		// TODO what does this return value mean? docs don't say
}

var draw_callback = C.GCallback(C.our_draw_callback)
