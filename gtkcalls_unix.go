// +build !windows,!darwin,!plan9
// TODO is there a way to simplify the above? :/

// 16 february 2014
package ui

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
	return fromgtkwidget(C.gtk_window_new(0))
}

// the garbage collector has been found to eat my callback functions; this will stop it
var callbackstore = make([]*func() bool, 0, 50)

func g_signal_connect(obj *gtkWidget, sig string, callback func() bool) {
	callbackstore = append(callbackstore, &callback)
	ccallback := callbacks[sig]
	csig := C.CString(sig)
	defer C.free(unsafe.Pointer(csig))
	C.gSignalConnect(togtkwidget(obj), csig, ccallback, unsafe.Pointer(&callback))
}

// TODO ensure this works if called on an individual control
func gtk_widget_show(widget *gtkWidget) {
	C.gtk_widget_show_all(togtkwidget(widget))
}

func gtk_widget_hide(widget *gtkWidget) {
	C.gtk_widget_hide(togtkwidget(widget))
}

func gtk_window_set_title(window *gtkWidget, title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.gtk_window_set_title(togtkwindow(window), togchar(ctitle))
}

func gtk_window_get_title(window *gtkWidget) string {
	return C.GoString(fromgchar(C.gtk_window_get_title(togtkwindow(window))))
}

func gtk_window_resize(window *gtkWidget, width int, height int) {
	C.gtk_window_resize(togtkwindow(window), C.gint(width), C.gint(height))
}

func gtk_window_get_size(window *gtkWidget) (int, int) {
	var width, height C.gint

	C.gtk_window_get_size(togtkwindow(window), &width, &height)
	return int(width), int(height)
}

func gtk_fixed_new() *gtkWidget {
	return fromgtkwidget(C.gtk_fixed_new())
}

func gtk_container_add(container *gtkWidget, widget *gtkWidget) {
	C.gtk_container_add(togtkcontainer(container), togtkwidget(widget))
}

func gtk_fixed_move(container *gtkWidget, widget *gtkWidget, x int, y int) {
	C.gtk_fixed_move(togtkfixed(container), togtkwidget(widget),
		C.gint(x), C.gint(y))
}

func gtk_widget_set_size_request(widget *gtkWidget, width int, height int) {
	C.gtk_widget_set_size_request(togtkwidget(widget), C.gint(width), C.gint(height))
}

func gtk_button_new() *gtkWidget {
	return fromgtkwidget(C.gtk_button_new())
}

func gtk_button_set_label(button *gtkWidget, label string) {
	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))
	C.gtk_button_set_label(togtkbutton(button), togchar(clabel))
}

func gtk_button_get_label(button *gtkWidget) string {
	return C.GoString(fromgchar(C.gtk_button_get_label(togtkbutton(button))))
}

func gtk_check_button_new() *gtkWidget {
	return fromgtkwidget(C.gtk_check_button_new())
}

func gtk_toggle_button_get_active(widget *gtkWidget) bool {
	return fromgbool(C.gtk_toggle_button_get_active(togtktogglebutton(widget)))
}

func gtk_combo_box_text_new() *gtkWidget {
	return fromgtkwidget(C.gtk_combo_box_text_new())
}

func gtk_combo_box_text_new_with_entry() *gtkWidget {
	return fromgtkwidget(C.gtk_combo_box_text_new_with_entry())
}

func gtk_combo_box_text_get_active_text(widget *gtkWidget) string {
	return C.GoString(fromgchar(C.gtk_combo_box_text_get_active_text(togtkcombobox(widget))))
}

func gtk_combo_box_text_append_text(widget *gtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_combo_box_text_append_text(togtkcombobox(widget), togchar(ctext))
}

func gtk_combo_box_text_insert_text(widget *gtkWidget, index int, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_combo_box_text_insert_text(togtkcombobox(widget), C.gint(index), togchar(ctext))
}

func gtk_combo_box_get_active(widget *gtkWidget) int {
	cb := (*C.GtkComboBox)(unsafe.Pointer(widget))
	return int(C.gtk_combo_box_get_active(cb))
}

func gtk_combo_box_text_remove(widget *gtkWidget, index int) {
	C.gtk_combo_box_text_remove(togtkcombobox(widget), C.gint(index))
}

func gtk_entry_new() *gtkWidget {
	return fromgtkwidget(C.gtk_entry_new())
}

func gtk_entry_set_text(widget *gtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_entry_set_text(togtkentry(widget), togchar(ctext))
}

func gtk_entry_get_text(widget *gtkWidget) string {
	return C.GoString(fromgchar(C.gtk_entry_get_text(togtkentry(widget))))
}

var _emptystring = [1]C.gchar{0}
var emptystring = &_emptystring[0]

func gtk_label_new() *gtkWidget {
	return fromgtkwidget(C.gtk_label_new(emptystring))
	// TODO left-justify?
}

func gtk_label_set_text(widget *gtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_label_set_text(togtklabel(widget), togchar(ctext))
}

func gtk_label_get_text(widget *gtkWidget) string {
	return C.GoString(fromgchar(C.gtk_label_get_text(togtklabel(widget))))
}
