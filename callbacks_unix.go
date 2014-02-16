// +build !windows,!darwin,!plan9

// 16 february 2014
package main

import (
	"unsafe"
)

/*
cgo doesn't support calling Go functions by default; we have to mark them for export. Not a problem, except arguments to GTK+ callbacks depend on the callback itself. Since we're generating callback functions as simple closures of one type, this file will wrap the generated callbacks in the appropriate callback type. We pass the actual generated pointer to the extra data parameter of the callback.
*/

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// extern gboolean our_delete_event_callback(GtkWidget *, GdkEvent *, gpointer);
import "C"

//export our_delete_event_callback
func our_delete_event_callback(widget *C.GtkWidget, event *C.GdkEvent, what C.gpointer) C.gboolean {
	f := *(*func() bool)(unsafe.Pointer(what))
	return togbool(f())
}

var callbacks = map[string]C.GCallback{
	"delete-event":		C.GCallback(C.our_delete_event_callback),
}
