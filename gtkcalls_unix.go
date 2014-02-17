// +build !windows,!darwin,!plan9
// TODO is there a way to simplify the above? :/

// 16 february 2014
package main

import (
	"unsafe"
)

// TODOs:
// - document the magic stuff in the listbox code

// #cgo pkg-config: gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// /* because cgo is flaky with macros */
// void gSignalConnect(GtkWidget *widget, char *signal, GCallback callback, void *data) { g_signal_connect(widget, signal, callback, data); }
// /* because cgo seems to choke on ... */
// void gtkTreeModelGet(GtkTreeModel *model, GtkTreeIter *iter, gchar **gs) { gtk_tree_model_get(model, iter, 0, gs, -1); }
// GtkListStore *gtkListStoreNew(void) { return gtk_list_store_new(1, G_TYPE_STRING); }
// void gtkListStoreSet(GtkListStore *ls, GtkTreeIter *iter, char *gs) { gtk_list_store_set(ls, iter, 0, (gchar *) gs, -1); }
// GtkTreeViewColumn *gtkTreeViewColumnNewWithAttributes(GtkCellRenderer *renderer) { return gtk_tree_view_column_new_with_attributes("", renderer, "text", 0, NULL); }
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

func g_signal_connect(obj *gtkWidget, sig string, callback func() bool) {
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

// TODO split all this out into its own file?

func gListboxNew(multisel bool) *gtkWidget {
	store := C.gtkListStoreNew()
	widget := C.gtk_tree_view_new_with_model((*C.GtkTreeModel)(unsafe.Pointer(store)))
	tv := (*C.GtkTreeView)(unsafe.Pointer(widget))
	column := C.gtkTreeViewColumnNewWithAttributes(C.gtk_cell_renderer_text_new())
	// TODO set AUTOSIZE?
	C.gtk_tree_view_append_column(tv, column)
	C.gtk_tree_view_set_headers_visible(tv, C.FALSE)
	sel := C.GTK_SELECTION_SINGLE
	if multisel {
		sel = C.GTK_SELECTION_MULTIPLE
	}
	C.gtk_tree_selection_set_mode(C.gtk_tree_view_get_selection(tv), C.GtkSelectionMode(sel))
	scrollarea := C.gtk_scrolled_window_new((*C.GtkAdjustment)(nil), (*C.GtkAdjustment)(nil))
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(scrollarea)), widget)
	return (*gtkWidget)(unsafe.Pointer(scrollarea))
}

func gListboxNewSingle() *gtkWidget {
	return gListboxNew(false)
}

func gListboxNewMulti() *gtkWidget {
	return gListboxNew(true)
}

func getTreeViewFrom(widget *gtkWidget) *C.GtkWidget {
	return C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(widget)))
}

func gListboxText(widget *gtkWidget) string {
	var model *C.GtkTreeModel
	var iter C.GtkTreeIter
	var gs *C.gchar

	tv := (*C.GtkTreeView)(unsafe.Pointer(getTreeViewFrom(widget)))
	sel := C.gtk_tree_view_get_selection(tv)
	if !fromgbool(C.gtk_tree_selection_get_selected(sel, &model, &iter)) {
		return ""
	}
	C.gtkTreeModelGet(model, &iter, &gs)
	return C.GoString((*C.char)(unsafe.Pointer(gs)))
}

func gListboxAppend(widget *gtkWidget, what string) {
	var iter C.GtkTreeIter

	tv := (*C.GtkTreeView)(unsafe.Pointer(getTreeViewFrom(widget)))
	ls := (*C.GtkListStore)(unsafe.Pointer(C.gtk_tree_view_get_model(tv)))
	C.gtk_list_store_append(ls, &iter)
	cwhat := C.CString(what)
	defer C.free(unsafe.Pointer(cwhat))
	C.gtkListStoreSet(ls, &iter, cwhat)
}

func gListboxInsert(widget *gtkWidget, index int, what string) {
	var iter C.GtkTreeIter

	tv := (*C.GtkTreeView)(unsafe.Pointer(getTreeViewFrom(widget)))
	ls := (*C.GtkListStore)(unsafe.Pointer(C.gtk_tree_view_get_model(tv)))
	C.gtk_list_store_insert(ls, &iter, C.gint(index))
	cwhat := C.CString(what)
	defer C.free(unsafe.Pointer(cwhat))
	C.gtkListStoreSet(ls, &iter, cwhat)
}

func gListboxSelected(widget *gtkWidget) int {
	var model *C.GtkTreeModel
	var iter C.GtkTreeIter

	tv := (*C.GtkTreeView)(unsafe.Pointer(getTreeViewFrom(widget)))
	sel := C.gtk_tree_view_get_selection(tv)
	if !fromgbool(C.gtk_tree_selection_get_selected(sel, &model, &iter)) {
		return -1
	}
	path := C.gtk_tree_model_get_path(model, &iter)
	idx := C.gtk_tree_path_get_indices(path)
	return int(*idx)
}

func gListboxSelectedMulti(widget *gtkWidget) (indices []int) {
	var model *C.GtkTreeModel

	tv := (*C.GtkTreeView)(unsafe.Pointer(getTreeViewFrom(widget)))
	sel := C.gtk_tree_view_get_selection(tv)
	rows := C.gtk_tree_selection_get_selected_rows(sel, &model)
	defer C.g_list_free_full(rows, C.GDestroyNotify(unsafe.Pointer(C.gtk_tree_path_free)))
	len := C.g_list_length(rows)
	if len == 0 {
		return nil
	}
	indices = make([]int, len)
	for i := C.guint(0); i < len; i++ {
		path := (*C.GtkTreePath)(unsafe.Pointer(rows.data))
		idx := C.gtk_tree_path_get_indices(path)
		indices[i] = int(*idx)
		rows = rows.next
	}
	return indices
}

func gListboxSelMultiTexts(widget *gtkWidget) (texts []string) {
	var model *C.GtkTreeModel
	var iter C.GtkTreeIter
	var gs *C.gchar

	tv := (*C.GtkTreeView)(unsafe.Pointer(getTreeViewFrom(widget)))
	sel := C.gtk_tree_view_get_selection(tv)
	rows := C.gtk_tree_selection_get_selected_rows(sel, &model)
	defer C.g_list_free_full(rows, C.GDestroyNotify(unsafe.Pointer(C.gtk_tree_path_free)))
	len := C.g_list_length(rows)
	if len == 0 {
		return nil
	}
	texts = make([]string, len)
	for i := C.guint(0); i < len; i++ {
		path := (*C.GtkTreePath)(unsafe.Pointer(rows.data))
		if !fromgbool(C.gtk_tree_model_get_iter(model, &iter, path)) {
			// TODO
			return
		}
		C.gtkTreeModelGet(model, &iter, &gs)
		texts[i] = C.GoString((*C.char)(unsafe.Pointer(gs)))
		rows = rows.next
	}
	return texts
}

func gListboxDelete(widget *gtkWidget, index int) {
	var iter C.GtkTreeIter

	tv := (*C.GtkTreeView)(unsafe.Pointer(getTreeViewFrom(widget)))
	ls := (*C.GtkListStore)(unsafe.Pointer(C.gtk_tree_view_get_model(tv)))
	if !fromgbool(C.gtk_tree_model_iter_nth_child((*C.GtkTreeModel)(unsafe.Pointer(ls)), &iter, (*C.GtkTreeIter)(nil), C.gint(index))) {		// no such index
		// TODO
		return
	}
	C.gtk_list_store_remove(ls, &iter)
}
