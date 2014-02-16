// +build !windows,!darwin,!plan9
// TODO is there a way to simplify the above? :/

// 16 february 2014
package main

import (
	"unsafe"
)

// #cgo pkg-config: gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// /* because cgo is flaky with macros */
// void gSignalConnect(GtkWidget *widget, char *signal, GCallback callback, void *data) { g_signal_connect(widget, signal, callback, data); }
import "C"

type (
	gtkWidget C.GtkWidget
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

func gtk_init() bool {
	// TODO allow GTK+ standard command-line argument processing
	return fromgbool(C.gtk_init_check((*C.int)(nil), (***C.char)(nil)))
}

func gdk_threads_add_idle(what func() bool) {
	C.gdk_threads_add_idle(callbacks["idle"], C.gpointer(unsafe.Pointer(&what)))
}

func gtk_main() {
	C.gtk_main()
}

func gtk_window_new() *gtkWidget {
	// 0 == GTK_WINDOW_TOPLEVEL (the only other type, _POPUP, should not be used)
	return (*gtkWidget)(unsafe.Pointer(C.gtk_window_new(0)))
}

// shorthand
func gtkwidget(what *gtkWidget) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(what))
}

func g_signal_connect(obj *gtkWidget, sig string, callback func() bool) {
	ccallback := callbacks[sig]
	csig := C.CString(sig)
	defer C.free(unsafe.Pointer(csig))
	C.gSignalConnect(gtkwidget(obj), csig, ccallback, unsafe.Pointer(&callback))
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

func gtk_window_get_size(window *gtkWidget) (int, int) {
	var width, height C.gint

	C.gtk_window_get_size((*C.GtkWindow)(unsafe.Pointer(window)), &width, &height)
	return int(width), int(height)
}

func gtk_fixed_new() *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(C.gtk_fixed_new()))
}

func gtk_container_add(container *gtkWidget, widget *gtkWidget) {
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(container)), gtkwidget(widget))
}

func gtk_fixed_move(container *gtkWidget, widget *gtkWidget, x int, y int) {
	C.gtk_fixed_move((*C.GtkFixed)(unsafe.Pointer(container)), gtkwidget(widget),
		C.gint(x), C.gint(y))
}

func gtk_widget_set_size_request(widget *gtkWidget, width int, height int) {
	C.gtk_widget_set_size_request(gtkwidget(widget), C.gint(width), C.gint(height))
}

func gtk_button_new() *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(C.gtk_button_new()))
}

func gtk_button_set_label(button *gtkWidget, label string) {
	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))
	C.gtk_button_set_label((*C.GtkButton)(unsafe.Pointer(button)),
		(*C.gchar)(unsafe.Pointer(clabel)))
}

func gtk_check_button_new() *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(C.gtk_check_button_new()))
}

func gtk_toggle_button_get_active(widget *gtkWidget) bool {
	return fromgbool(C.gtk_toggle_button_get_active((*C.GtkToggleButton)(unsafe.Pointer(widget))))
}
