// +build !windows,!darwin,!plan9

// 14 march 2014

package ui

import (
	"unsafe"
	"image"
)

// #cgo pkg-config: gtk+-3.0
// /* GTK+ 3.8 deprecates gtk_scrolled_window_add_with_viewport(); we need 3.4 miniimum though
// setting MIN_REQUIRED ensures nothing older; setting MAX_ALLOWED disallows newer functions - thanks to desrt in irc.gimp.net/#gtk+
// TODO add this to the other files too */
// #define GDK_VERSION_MIN_REQUIRED GDK_VERSION_3_4
// #define GDK_VERSION_MAX_ALLOWED GDK_VERSION_3_4
// #include <gtk/gtk.h>
// extern gboolean our_area_draw_callback(GtkWidget *, cairo_t *, gpointer);
// /* HACK - see https://code.google.com/p/go/issues/detail?id=7548 */
// struct _cairo {};
import "C"

func gtkAreaNew() *gtkWidget {
	drawingarea := C.gtk_drawing_area_new()
	C.gtk_widget_set_size_request(drawingarea, 320, 240)
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
