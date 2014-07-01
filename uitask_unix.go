// +build !windows,!darwin,!plan9

// 16 february 2014

package ui

import (
	"fmt"
	"unsafe"
)

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
// /* unfortunately, there's no way to differentiate between the main thread and other threads; in fact, doing what we do on other platforms is discouraged by the GTK+ devs!
// but I can't avoid this any other way... so we have structures defined on the C side to skirt the garbage collector */
// struct uitaskParams {
// 	void *window;		/* createWindow */
// 	void *control;		/* createWindow */
// 	gboolean show;	/* createWindow */
// };
// static struct uitaskParams *mkParams(void)
// {
// 	/* g_malloc0() will abort on not enough memory */
// 	return (struct uitaskParams *) g_malloc0(sizeof (struct uitaskParams));
// }
// static void freeParams(struct uitaskParams *p)
// {
// 	g_free(p);
// }
// extern gboolean our_createWindow_callback(gpointer);
import "C"

//export our_createWindow_callback
func our_createWindow_callback(what C.gpointer) C.gboolean {
	uc := (*C.struct_uitaskParams)(unsafe.Pointer(what))
	w := (*Window)(unsafe.Pointer(uc.window))
	c := *(*Control)(unsafe.Pointer(uc.control))
	s := fromgbool(uc.show)
	w.create(c, s)
	C.freeParams(uc)
	return C.FALSE // remove this idle function; we're finished
}

func (_uitask) createWindow(w *Window, c Control, s bool) {
	uc := C.mkParams()
	uc.window = unsafe.Pointer(w)
	uc.control = unsafe.Pointer(&c)
	uc.show = togbool(s)
	gdk_threads_add_idle(C.our_createWindow_callback, unsafe.Pointer(uc))
}

func gdk_threads_add_idle(f unsafe.Pointer, what unsafe.Pointer) {
	C.gdk_threads_add_idle(C.GCallback(f), C.gpointer(what))
}

func uiinit() error {
	err := gtk_init()
	if err != nil {
		return fmt.Errorf("gtk_init() failed: %v", err)
	}

	return nil
}

func ui() {
	// thanks to tristan and Daniel_S in irc.gimp.net/#gtk
	// see our_idle_callback in callbacks_unix.go for details
	go func() {
		// TODO do differently
		<-Stop
		// using gdk_threads_add_idle() will make sure it gets run when any events currently being handled finish running
		f := func() {
			C.gtk_main_quit()
		}
		done := make(chan struct{})
		gdk_threads_add_idle_op(&gtkIdleOp{
			what: f,
			done: done,
		})
		<-done
		close(done)
	}()

	C.gtk_main()
}
