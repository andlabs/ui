// +build !windows,!darwin

// 16 august 2014

package ui

import (
	"fmt"
	"image"
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

// this is what GtkFileChooserWidget uses
// technically it uses max(width from that, height from that) if the call below fails and 16x16 otherwise, but we won't worry about that here (yet?)
const scaleTo = C.GTK_ICON_SIZE_MENU

func toIconSizedGdkPixbuf(img *image.RGBA) *C.GdkPixbuf {
	var width, height C.gint

	surface := C.cairo_image_surface_create(C.CAIRO_FORMAT_ARGB32,
		C.int(img.Rect.Dx()),
		C.int(img.Rect.Dy()))
	if status := C.cairo_surface_status(surface); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("cairo_create_image_surface() failed in toIconSizedGdkPixbuf(): %s\n",
			C.GoString(C.cairo_status_to_string(status))))
	}
	C.cairo_surface_flush(surface)
	toARGB(img, uintptr(unsafe.Pointer(C.cairo_image_surface_get_data(surface))),
		int(C.cairo_image_surface_get_stride(surface)), false) // not NRGBA
	C.cairo_surface_mark_dirty(surface)
	basepixbuf := C.gdk_pixbuf_get_from_surface(surface, 0, 0, C.gint(img.Rect.Dx()), C.gint(img.Rect.Dy()))
	if basepixbuf == nil {
		panic(fmt.Errorf("gdk_pixbuf_get_from_surface() failed in toIconSizedGdkPixbuf() (no reason available)"))
	}

	if C.gtk_icon_size_lookup(scaleTo, &width, &height) == C.FALSE {
		panic(fmt.Errorf("gtk_icon_size_lookup() failed in toIconSizedGdkPixbuf() (no reason available)"))
	}
	if int(width) == img.Rect.Dx() && int(height) == img.Rect.Dy() {
		// just return the base pixbuf; we're good
		C.cairo_surface_destroy(surface)
		return basepixbuf
	}
	// else scale
	pixbuf := C.gdk_pixbuf_scale_simple(basepixbuf, C.int(width), C.int(height), C.GDK_INTERP_NEAREST)
	if pixbuf == nil {
		panic(fmt.Errorf("gdk_pixbuf_scale_simple() failed in toIconSizedGdkPixbuf() (no reason available)"))
	}

	C.g_object_unref(C.gpointer(unsafe.Pointer(basepixbuf)))
	C.cairo_surface_destroy(surface)
	return pixbuf
}
