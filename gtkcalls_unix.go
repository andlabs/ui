// +build !windows,!darwin,!plan9
// TODO is there a way to simplify the above? :/

// 16 february 2014

package ui

import (
	"fmt"
	"unsafe"
)

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
import "C"

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

func gtk_window_new() *C.GtkWidget {
	return C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
}

// TODO ensure this works if called on an individual control
func gtk_widget_show(widget *C.GtkWidget) {
	C.gtk_widget_show_all(widget)
}

func gtk_widget_hide(widget *C.GtkWidget) {
	C.gtk_widget_hide(widget)
}

func gtk_window_set_title(window *C.GtkWidget, title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.gtk_window_set_title(togtkwindow(window), togchar(ctitle))
}

func gtk_window_get_title(window *C.GtkWidget) string {
	return C.GoString(fromgchar(C.gtk_window_get_title(togtkwindow(window))))
}

func gtk_window_resize(window *C.GtkWidget, width int, height int) {
	C.gtk_window_resize(togtkwindow(window), C.gint(width), C.gint(height))
}

func gtk_window_get_size(window *C.GtkWidget) (int, int) {
	var width, height C.gint

	C.gtk_window_get_size(togtkwindow(window), &width, &height)
	return int(width), int(height)
}

// on some themes, such as oxygen-gtk, GtkLayout draws a solid-color background, not the window background (as GtkFixed and GtkDrawingArea do)
// this CSS fixes it
// thanks to drahnr and ptomato in http://stackoverflow.com/questions/22940588/how-do-i-really-make-a-gtk-3-gtklayout-transparent-draw-theme-background
// TODO report to oxygen-gtk devs
var gtkLayoutCSS = []byte(`GtkLayout {
	background-color: transparent;
}
`)

func makeTransparent(layout *C.GtkWidget) {
	var err *C.GError = nil		// redundant in Go, but let's explicitly assign it anyway

	provider := C.gtk_css_provider_new()
	added := C.gtk_css_provider_load_from_data(provider,
		(*C.gchar)(unsafe.Pointer(&gtkLayoutCSS[0])), C.gssize(len(gtkLayoutCSS)), &err)
	if added == C.FALSE {
		message := C.GoString(fromgchar(err.message))
		panic(fmt.Errorf("error loading transparent background CSS for GtkLayout: %s", message))
	}
	C.gtk_style_context_add_provider(C.gtk_widget_get_style_context(layout),
		(*C.GtkStyleProvider)(unsafe.Pointer(provider)),
		C.GTK_STYLE_PROVIDER_PRIORITY_USER)
}

// this should allow us to resize the window arbitrarily
// thanks to Company in irc.gimp.net/#gtk+
func gtkNewWindowLayout() *C.GtkWidget {
	layout := C.gtk_layout_new(nil, nil)
	makeTransparent(layout)
	scrollarea := C.gtk_scrolled_window_new((*C.GtkAdjustment)(nil), (*C.GtkAdjustment)(nil))
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(scrollarea)), layout)
	// never show scrollbars; we're just doing this to allow arbitrary resizes
	C.gtk_scrolled_window_set_policy((*C.GtkScrolledWindow)(unsafe.Pointer(scrollarea)),
		C.GTK_POLICY_NEVER, C.GTK_POLICY_NEVER)
	return scrollarea
}

func gtk_container_add(container *C.GtkWidget, widget *C.GtkWidget) {
	C.gtk_container_add(togtkcontainer(container), widget)
}

func gtkAddWidgetToLayout(container *C.GtkWidget, widget *C.GtkWidget) {
	layout := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(container)))
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(layout)), widget)
}

func gtkMoveWidgetInLayout(container *C.GtkWidget, widget *C.GtkWidget, x int, y int) {
	layout := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(container)))
	C.gtk_layout_move((*C.GtkLayout)(unsafe.Pointer(layout)), widget,
		C.gint(x), C.gint(y))
}

func gtk_widget_set_size_request(widget *C.GtkWidget, width int, height int) {
	C.gtk_widget_set_size_request(widget, C.gint(width), C.gint(height))
}

func gtk_button_new() *C.GtkWidget {
	return C.gtk_button_new()
}

func gtk_button_set_label(button *C.GtkWidget, label string) {
	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))
	C.gtk_button_set_label(togtkbutton(button), togchar(clabel))
}

func gtk_button_get_label(button *C.GtkWidget) string {
	return C.GoString(fromgchar(C.gtk_button_get_label(togtkbutton(button))))
}

func gtk_check_button_new() *C.GtkWidget {
	return C.gtk_check_button_new()
}

func gtk_toggle_button_get_active(widget *C.GtkWidget) bool {
	return fromgbool(C.gtk_toggle_button_get_active(togtktogglebutton(widget)))
}

func gtk_combo_box_text_new() *C.GtkWidget {
	return C.gtk_combo_box_text_new()
}

func gtk_combo_box_text_new_with_entry() *C.GtkWidget {
	return C.gtk_combo_box_text_new_with_entry()
}

func gtk_combo_box_text_get_active_text(widget *C.GtkWidget) string {
	return C.GoString(fromgchar(C.gtk_combo_box_text_get_active_text(togtkcombobox(widget))))
}

func gtk_combo_box_text_append_text(widget *C.GtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_combo_box_text_append_text(togtkcombobox(widget), togchar(ctext))
}

func gtk_combo_box_text_insert_text(widget *C.GtkWidget, index int, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_combo_box_text_insert_text(togtkcombobox(widget), C.gint(index), togchar(ctext))
}

func gtk_combo_box_get_active(widget *C.GtkWidget) int {
	cb := (*C.GtkComboBox)(unsafe.Pointer(widget))
	return int(C.gtk_combo_box_get_active(cb))
}

func gtk_combo_box_text_remove(widget *C.GtkWidget, index int) {
	C.gtk_combo_box_text_remove(togtkcombobox(widget), C.gint(index))
}

func gtkComboBoxLen(widget *C.GtkWidget) int {
	cb := (*C.GtkComboBox)(unsafe.Pointer(widget))
	model := C.gtk_combo_box_get_model(cb)
	// this is the same as with a Listbox so
	return gtkTreeModelListLen(model)
}

func gtk_entry_new() *C.GtkWidget {
	return C.gtk_entry_new()
}

func gtkPasswordEntryNew() *C.GtkWidget {
	e := gtk_entry_new()
	C.gtk_entry_set_visibility(togtkentry(e), C.FALSE)
	return e
}

func gtk_entry_set_text(widget *C.GtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_entry_set_text(togtkentry(widget), togchar(ctext))
}

func gtk_entry_get_text(widget *C.GtkWidget) string {
	return C.GoString(fromgchar(C.gtk_entry_get_text(togtkentry(widget))))
}

var _emptystring = [1]C.gchar{0}
var emptystring = &_emptystring[0]

func gtk_label_new() *C.GtkWidget {
	return C.gtk_label_new(emptystring)
	// TODO left-justify?
}

func gtk_label_set_text(widget *C.GtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_label_set_text(togtklabel(widget), togchar(ctext))
}

func gtk_label_get_text(widget *C.GtkWidget) string {
	return C.GoString(fromgchar(C.gtk_label_get_text(togtklabel(widget))))
}

func gtk_widget_get_preferred_size(widget *C.GtkWidget) (minWidth int, minHeight int, natWidth int, natHeight int) {
	var minimum, natural C.GtkRequisition

	C.gtk_widget_get_preferred_size(widget, &minimum, &natural)
	return int(minimum.width), int(minimum.height),
		int(natural.width), int(natural.height)
}

func gtk_progress_bar_new() *C.GtkWidget {
	return C.gtk_progress_bar_new()
}

func gtk_progress_bar_set_fraction(w *C.GtkWidget, percent int) {
	p := C.gdouble(percent) / 100
	C.gtk_progress_bar_set_fraction(togtkprogressbar(w), p)
}

func gtk_progress_bar_pulse(w *C.GtkWidget) {
	C.gtk_progress_bar_pulse(togtkprogressbar(w))
}
