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

type imagelist struct {
	list []*C.GdkPixbuf
}

func newImageList() ImageList {
	return new(imagelist)
}

// this is what GtkFileChooserWidget uses
// technically it uses max(width from that, height from that) if the call below fails and 16x16 otherwise, but we won't worry about that here (yet?)
const scaleTo = C.GTK_ICON_SIZE_MENU

func (i *imagelist) Append(img *image.RGBA) {
	var width, height C.gint

	surface := C.cairo_image_surface_create(C.CAIRO_FORMAT_ARGB32,
		C.int(img.Rect.Dx()),
		C.int(img.Rect.Dy()))
	if status := C.cairo_surface_status(surface); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("cairo_create_image_surface() failed in ImageList.Append(): %s\n",
			C.GoString(C.cairo_status_to_string(status))))
	}
	C.cairo_surface_flush(surface)
	toARGB(img, uintptr(unsafe.Pointer(C.cairo_image_surface_get_data(surface))),
		int(C.cairo_image_surface_get_stride(surface)), false) // not NRGBA
	C.cairo_surface_mark_dirty(surface)
	basepixbuf := C.gdk_pixbuf_get_from_surface(surface, 0, 0, C.gint(img.Rect.Dx()), C.gint(img.Rect.Dy()))
	if basepixbuf == nil {
		panic(fmt.Errorf("gdk_pixbuf_get_from_surface() failed in ImageList.Append() (no reason available)"))
	}

	if C.gtk_icon_size_lookup(scaleTo, &width, &height) == C.FALSE {
		panic(fmt.Errorf("gtk_icon_size_lookup() failed in ImageList.Append() (no reason available)"))
	}
	if int(width) == img.Rect.Dx() && int(height) == img.Rect.Dy() {
		// just add the base pixbuf; we're good
		i.list = append(i.list, basepixbuf)
		C.cairo_surface_destroy(surface)
		return
	}
	// else scale
	pixbuf := C.gdk_pixbuf_scale_simple(basepixbuf, C.int(width), C.int(height), C.GDK_INTERP_NEAREST)
	if pixbuf == nil {
		panic(fmt.Errorf("gdk_pixbuf_scale_simple() failed in ImageList.Append() (no reason available)"))
	}

	i.list = append(i.list, pixbuf)
	C.g_object_unref(C.gpointer(unsafe.Pointer(basepixbuf)))
	C.cairo_surface_destroy(surface)
}

func (i *imagelist) Len() ImageIndex {
	return ImageIndex(len(i.list))
}

type imageListApply interface {
	apply(*[]*C.GdkPixbuf)
}

func (i *imagelist) apply(out *[]*C.GdkPixbuf) {
	*out = make([]*C.GdkPixbuf, len(i.list))
	copy(*out, i.list)
}
