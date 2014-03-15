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
// #include <stdlib.h>
// extern gboolean our_window_delete_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_window_configure_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern void our_button_clicked_callback(GtkButton *, gpointer);
// extern gboolean our_idle_callback(gpointer);
// /* because cgo is flaky with macros; static inline because we have //exports */
// static inline void gSignalConnect(GtkWidget *widget, char *signal, GCallback callback, void *data) { g_signal_connect(widget, signal, callback, data); }
import "C"

//export our_window_delete_event_callback
func our_window_delete_event_callback(widget *C.GtkWidget, event *C.GdkEvent, what C.gpointer) C.gboolean {
	// called when the user tries to close the window
	s := (*sysData)(unsafe.Pointer(what))
	s.signal()
	return C.TRUE		// do not close the window
}

var window_delete_event_callback = C.GCallback(C.our_window_delete_event_callback)

//export our_window_configure_event_callback
func our_window_configure_event_callback(widget *C.GtkWidget, event *C.GdkEvent, what C.gpointer) C.gboolean {
	// called when the window is resized
	s := (*sysData)(unsafe.Pointer(what))
	if s.container != nil && s.resize != nil {		// wait for init
		width, height := gtk_window_get_size(s.widget)
		// top-left is (0,0) so no need for winheight
		err := s.resize(0, 0, width, height, 0)
		if err != nil {
			panic("child resize failed: " + err.Error())
		}
	}
	// returning false indicates that we continue processing events related to configure-event; if we choose not to, then after some controls have been added, the layout fails completely and everything stays in the starting position/size
	// TODO make sure this is the case
	return C.FALSE
}

var window_configure_event_callback = C.GCallback(C.our_window_configure_event_callback)

//export our_button_clicked_callback
func our_button_clicked_callback(button *C.GtkButton, what C.gpointer) {
	// called when the user clicks a button
	s := (*sysData)(unsafe.Pointer(what))
	s.signal()
}

var button_clicked_callback = C.GCallback(C.our_button_clicked_callback)

// this is the type of the signals fields in classData; here to avoid needing to import C
type callbackMap map[string]C.GCallback

// this is what actually connects a signal
func g_signal_connect(obj *gtkWidget, sig string, callback C.GCallback, sysData *sysData) {
	csig := C.CString(sig)
	defer C.free(unsafe.Pointer(csig))
	C.gSignalConnect(togtkwidget(obj), csig, callback, unsafe.Pointer(sysData))
}

// there are two issues we solve here:
// 1) we need to make sure the uitask request gets garbage collected when we're done so as to not waste memory, but only when we're done so as to not have craziness happen
// 2) we need to make sure one idle function runs and finishes running before we start the next; otherwise we could wind up with weird things like the ret channel being closed early
// so our_idle_callback() calls the uitask function in what and sends a message back to the dispatcher over done that it finished running; the dispatcher is still holding onto the uitask function so it won't be collected
type gtkIdleOp struct {
	what		func()
	done		chan struct{}
}

//export our_idle_callback
func our_idle_callback(what C.gpointer) C.gboolean {
	idleop := (*gtkIdleOp)(unsafe.Pointer(what))
	idleop.what()
	idleop.done <- struct{}{}
	return C.FALSE		// remove this idle function; we're finished
}

func gdk_threads_add_idle(idleop *gtkIdleOp) {
	C.gdk_threads_add_idle(C.GCallback(C.our_idle_callback),
		C.gpointer(unsafe.Pointer(idleop)))
}
