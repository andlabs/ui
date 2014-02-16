// +build !windows,!darwin,!plan9
// TODO is there a way to simplify the above? :/

// 16 february 2014
package main

import (
	"fmt"
	"unsafe"
	"reflect"
)

// #cgo pkg-config: gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// /* because cgo is flaky with macros */
// static inline void gSignalConnect(GtkWidget *widget, char *signal, GCallback callback, void *data) { g_signal_connect(widget, signal, callback, data); }
// /* so we can call uistep */
// extern gboolean our_thread_callback(gpointer);
import "C"

type (
	gtkWidget C.GtkWidget

	// these are needed for signals
	gdkEvent C.GdkEvent
	gpointer C.gpointer
)

func fromgbool(b C.gboolean) bool {
	return b != C.FALSE
}

func togbool(b bool) C.gboolean {
	if b {
		return C.TRUE
	}
	return C.FALSE
}

//export our_thread_callback
func our_thread_callback(C.gpointer) C.gboolean {
	uistep()
	return C.TRUE
}

func gtk_init() bool {
	// TODO allow GTK+ standard command-line argument processing
	b := fromgbool(C.gtk_init_check((*C.int)(nil), (***C.char)(nil)))
	if !b {
		return false
	}
	// thanks to tristan in irc.gimp.net/#gtk
	C.gdk_threads_add_idle(C.GSourceFunc(C.our_thread_callback), C.gpointer(unsafe.Pointer(nil)))
	return true
}

func gtk_main() {
	C.gtk_main()
}

func gtk_window_new() *gtkWidget {
	// 0 == GTK_WINDOW_TOPLEVEL (the only other type, _POPUP, should not be used)
	return (*gtkWidget)(unsafe.Pointer(C.gtk_window_new(0)))
}

// because *gtkWidget and *C.GtkWidget are not compatible
func gtkwidget(g *gtkWidget) (*C.GtkWidget) {
	return (*C.GtkWidget)(unsafe.Pointer(g))
}

// TODO do we need the argument?
// TODO fine-tune the function type
func g_signal_connect(obj *gtkWidget, sig string, callback interface{}) {
	v := reflect.ValueOf(callback)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("UI library internal error: callback %v given to g_signal_connect not a function", v))
	}
	ccallback := C.GCallback(unsafe.Pointer(v.Pointer()))
	csig := C.CString(sig)
	defer C.free(unsafe.Pointer(csig))
	C.gSignalConnect(gtkwidget(obj), csig, ccallback, unsafe.Pointer(nil))
}

// TODO ensure this works if called on an individual control
func gtk_widget_show(widget *gtkWidget) {
	C.gtk_widget_show_all(gtkwidget(widget))
}

func gtk_widget_hide(widget *gtkWidget) {
	C.gtk_widget_hide(gtkwidget(widget))
}

func gtk_window_set_title(window *gtkWidget, title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(window)),
		(*C.gchar)(unsafe.Pointer(ctitle)))
}

func gtk_window_resize(window *gtkWidget, width int, height int) {
	C.gtk_window_resize((*C.GtkWindow)(unsafe.Pointer(window)), C.gint(width), C.gint(height))
}

func gtk_fixed_new() *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(C.gtk_fixed_new()))
}

func gtk_container_add(container *gtkWidget, widget *gtkWidget) {
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(container)), gtkwidget(widget))
}

func gtk_fixed_put(container *gtkWidget, widget *gtkWidget, x int, y int) {
	C.gtk_fixed_put((*C.GtkFixed)(unsafe.Pointer(container)), gtkwidget(widget),
		C.gint(x), C.gint(y))
}

func gtk_widget_set_size_request(widget *gtkWidget, width int, height int) {
	C.gtk_widget_set_size_request(gtkwidget(widget), C.gint(width), C.gint(height))
}
