// +build !windows,!darwin,!plan9
// TODO is there a way to simplify the above? :/

// 16 february 2014

package ui

import (
	"unsafe"
)

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
import "C"

type (
	gtkWidget C.GtkWidget
)

func gtk_init() bool {
	// TODO allow GTK+ standard command-line argument processing
	return fromgbool(C.gtk_init_check((*C.int)(nil), (***C.char)(nil)))
}

func gtk_main() {
	C.gtk_main()
}

func gtk_main_quit() {
	C.gtk_main_quit()
}

func gtk_window_new() *gtkWidget {
	// 0 == GTK_WINDOW_TOPLEVEL (the only other type, _POPUP, should not be used)
	return fromgtkwidget(C.gtk_window_new(0))
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

// this should allow us to resize the window arbitrarily
// thanks to Company in irc.gimp.net/#gtk+
func gtkNewWindowLayout() *gtkWidget {
	layout := C.gtk_layout_new(nil, nil)
	scrollarea := C.gtk_scrolled_window_new((*C.GtkAdjustment)(nil), (*C.GtkAdjustment)(nil))
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(scrollarea)), layout)
	// never show scrollbars; we're just doing this to allow arbitrary resizes
	C.gtk_scrolled_window_set_policy((*C.GtkScrolledWindow)(unsafe.Pointer(scrollarea)),
		C.GTK_POLICY_NEVER, C.GTK_POLICY_NEVER)
	return fromgtkwidget(scrollarea)
}

func gtk_container_add(container *gtkWidget, widget *gtkWidget) {
	C.gtk_container_add(togtkcontainer(container), togtkwidget(widget))
}

func gtkAddWidgetToLayout(container *gtkWidget, widget *gtkWidget) {
	layout := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(container)))
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(layout)), togtkwidget(widget))
}

func gtkMoveWidgetInLayout(container *gtkWidget, widget *gtkWidget, x int, y int) {
	layout := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(container)))
	C.gtk_layout_move((*C.GtkLayout)(unsafe.Pointer(layout)), togtkwidget(widget),
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

func gtkComboBoxLen(widget *gtkWidget) int {
	cb := (*C.GtkComboBox)(unsafe.Pointer(widget))
	model := C.gtk_combo_box_get_model(cb)
	// this is the same as with a Listbox so
	return gtkTreeModelListLen(model)
}

func gtk_entry_new() *gtkWidget {
	return fromgtkwidget(C.gtk_entry_new())
}

func gtkPasswordEntryNew() *gtkWidget {
	e := gtk_entry_new()
	C.gtk_entry_set_visibility(togtkentry(e), C.FALSE)
	return e
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

func gtk_widget_get_preferred_size(w *gtkWidget) (minWidth int, minHeight int, natWidth int, natHeight int) {
	var minimum, natural C.GtkRequisition

	C.gtk_widget_get_preferred_size(togtkwidget(w), &minimum, &natural)
	return int(minimum.width), int(minimum.height),
		int(natural.width), int(natural.height)
}

func gtk_progress_bar_new() *gtkWidget {
	return fromgtkwidget(C.gtk_progress_bar_new())
}

func gtk_progress_bar_set_fraction(w *gtkWidget, percent int) {
	p := C.gdouble(percent) / 100
	C.gtk_progress_bar_set_fraction(togtkprogressbar(w), p)
}

func gtk_progress_bar_pulse(w *gtkWidget) {
	C.gtk_progress_bar_pulse(togtkprogressbar(w))
}
