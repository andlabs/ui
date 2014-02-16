// +build !windows,!darwin,!plan9

// 16 february 2014
package main

import (
	"unsafe"
)

/*
cgo doesn't support calling Go functions by default; we have to mark them for export. Not a problem, except arguments to GTK+ callbacks depend on the callback itself. Since we're generating callback functions as simple closures of one type, this file will wrap the generated callbacks in the appropriate callback type. We pass the actual generated pointer to the extra data parameter of the callback.

while we're at it the callback for our idle function will be handled here too for cleanliness purposes
*/

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// extern gboolean our_callback(gpointer);
// extern gboolean our_delete_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern void our_clicked_callback(GtkButton *, gpointer);
import "C"

//export our_callback
func our_callback(what C.gpointer) C.gboolean {
	f := *(*func() bool)(unsafe.Pointer(what))
	return togbool(f())
}

//export our_delete_event_callback
func our_delete_event_callback(widget *C.GtkWidget, event *C.GdkEvent, what C.gpointer) C.gboolean {
	return our_callback(what)
}

//export our_clicked_callback
func our_clicked_callback(button *C.GtkButton, what C.gpointer) {
	our_callback(what)
}

var callbacks = map[string]C.GCallback{
	"idle":			C.GCallback(C.our_callback),
	"delete-event":		C.GCallback(C.our_delete_event_callback),
	"clicked":			C.GCallback(C.our_clicked_callback),
}
