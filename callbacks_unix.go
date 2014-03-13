// +build !windows,!darwin,!plan9

// 16 february 2014
package ui

import (
	"unsafe"
)

/*
cgo doesn't support calling Go functions by default; we have to mark them for export. Not a problem, except arguments to GTK+ callbacks depend on the callback itself. Since we're generating callback functions as simple closures of one type, this file will wrap the generated callbacks in the appropriate callback type. We pass the actual generated pointer to the extra data parameter of the callback.

while we're at it the callback for our idle function will be handled here too
*/

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// extern gboolean our_callback(gpointer);
// extern gboolean our_window_callback(GtkWidget *, GdkEvent *, gpointer);
// extern void our_clicked_callback(GtkButton *, gpointer);
// extern gboolean our_idle_callback(gpointer);
import "C"

//export our_callback
func our_callback(what C.gpointer) C.gboolean {
	f := *(*func() bool)(unsafe.Pointer(what))
	return togbool(f())
}

//export our_window_callback
func our_window_callback(widget *C.GtkWidget, event *C.GdkEvent, what C.gpointer) C.gboolean {
	return our_callback(what)
}

//export our_clicked_callback
func our_clicked_callback(button *C.GtkButton, what C.gpointer) {
	our_callback(what)
}

//export our_idle_callback
func our_idle_callback(what C.gpointer) C.gboolean {
	// there are two issues we solve here:
	// 1) we need to make sure the uitask request gets garbage collected when we're done so as to not waste memory, but only when we're done so as to not have craziness happen
	// 2) we need to make sure one idle function runs and finishes running before we start the next; otherwise we could wind up with weird things like the ret channel being closed early
	// so this function calls the uitask function and sends a message back to the dispatcher that it finished running; the dispatcher is still holding onto the uitask function so it won't be collected
	idleop := (*gtkIdleOp)(unsafe.Pointer(what))
	idleop.what()
	idleop.done <- struct{}{}
	return C.FALSE		// remove this idle function; we're finished
}

var callbacks = map[string]C.GCallback{
	"idle":			C.GCallback(C.our_idle_callback),
	"delete-event":		C.GCallback(C.our_window_callback),
	"configure-event":	C.GCallback(C.our_window_callback),
	"clicked":			C.GCallback(C.our_clicked_callback),
}
