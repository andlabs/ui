// +build !windows,!darwin,!plan9

// 17 february 2014
package main

import (
	"unsafe"
)

// TODOs:
// - document the magic stuff in the listbox code

// #cgo pkg-config: gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// /* because cgo seems to choke on ... */
// void gtkTreeModelGet(GtkTreeModel *model, GtkTreeIter *iter, gchar **gs) { gtk_tree_model_get(model, iter, 0, gs, -1); }
// GtkListStore *gtkListStoreNew(void) { return gtk_list_store_new(1, G_TYPE_STRING); }
// void gtkListStoreSet(GtkListStore *ls, GtkTreeIter *iter, char *gs) { gtk_list_store_set(ls, iter, 0, (gchar *) gs, -1); }
// GtkTreeViewColumn *gtkTreeViewColumnNewWithAttributes(GtkCellRenderer *renderer) { return gtk_tree_view_column_new_with_attributes("", renderer, "text", 0, NULL); }
import "C"

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
