// +build !windows,!darwin,!plan9

// 17 february 2014

package ui

import (
	"unsafe"
)

/*
GTK+ 3.10 introduces a dedicated GtkListView type for simple listboxes like our Listbox. Unfortunately, since I want to target at least GTK+ 3.4, I need to do things the old, long, and hard way: manually with a GtkTreeView and GtkListStore model.

You are not expected to understand this.

if you must though:
GtkTreeViews are model/view. We use a GtkListStore as a model.
GtkTreeViews also separate selections into another type, but the GtkTreeView creates the selection object for us.
GtkTreeViews can scroll, but do not draw scrollbars or borders; we need to use a GtkScrolledWindow to hold the GtkTreeView to do so. We return the GtkScrolledWindow and get its control out when we want to access the GtkTreeView.
Like with Windows, there's a difference between signle-selection and multi-selection GtkTreeViews when it comes to getting the list of selections that we can exploit. The GtkTreeSelection class hands us an iterator and the model (for some reason). We pull a GtkTreePath out of the iterator, which we can then use to get the indices or text data.

For more information, read
	https://developer.gnome.org/gtk3/3.4/TreeWidget.html
	http://ubuntuforums.org/showthread.php?t=1208655
	http://scentric.net/tutorial/sec-treemodel-remove-row.html
	http://gtk.10911.n7.nabble.com/Scrollbars-in-a-GtkTreeView-td58076.html
	http://stackoverflow.com/questions/11407447/gtk-treeview-get-current-row-index-in-python (I think; I don't remember if I wound up using this one as a reference or not; I know after that I found the ubuntuforums link above)
and the GTK+ reference documentation.
*/

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
// /* because cgo seems to choke on ... */
// void gtkTreeModelGet(GtkTreeModel *model, GtkTreeIter *iter, gchar **gs)
// {
// 	/* 0 is the column #; we only have one column here */
// 	gtk_tree_model_get(model, iter, 0, gs, -1);
// }
// GtkListStore *gtkListStoreNew(void)
// {
// 	/* 1 column that stores strings */
// 	return gtk_list_store_new(1, G_TYPE_STRING);
// }
// void gtkListStoreSet(GtkListStore *ls, GtkTreeIter *iter, char *gs)
// {
// 	/* same parameters as in gtkTreeModelGet() */
// 	gtk_list_store_set(ls, iter, 0, (gchar *) gs, -1);
// }
// GtkTreeViewColumn *gtkTreeViewColumnNewWithAttributes(GtkCellRenderer *renderer)
// {
// 	/* "" is the column header; "text" associates the text of the column with column 0 */
// 	return gtk_tree_view_column_new_with_attributes("", renderer, "text", 0, NULL);
// }
import "C"

func fromgtktreemodel(x *C.GtkTreeModel) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtktreemodel(what *C.GtkWidget) *C.GtkTreeModel {
	return (*C.GtkTreeModel)(unsafe.Pointer(what))
}

func fromgtktreeview(x *C.GtkTreeView) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtktreeview(what *C.GtkWidget) *C.GtkTreeView {
	return (*C.GtkTreeView)(unsafe.Pointer(what))
}

func gListboxNew(multisel bool) *C.GtkWidget {
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
	// thanks to jlindgren in irc.gimp.net/#gtk+
	C.gtk_scrolled_window_set_shadow_type((*C.GtkScrolledWindow)(unsafe.Pointer(scrollarea)), C.GTK_SHADOW_IN)
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(scrollarea)), widget)
	return scrollarea
}

func gListboxNewSingle() *C.GtkWidget {
	return gListboxNew(false)
}

func gListboxNewMulti() *C.GtkWidget {
	return gListboxNew(true)
}

func getTreeViewFrom(widget *C.GtkWidget) *C.GtkTreeView {
	wid := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(widget)))
	return (*C.GtkTreeView)(unsafe.Pointer(wid))
}

func gListboxText(widget *C.GtkWidget) string {
	var model *C.GtkTreeModel
	var iter C.GtkTreeIter
	var gs *C.gchar

	tv := getTreeViewFrom(widget)
	sel := C.gtk_tree_view_get_selection(tv)
	if !fromgbool(C.gtk_tree_selection_get_selected(sel, &model, &iter)) {
		return ""
	}
	C.gtkTreeModelGet(model, &iter, &gs)
	return C.GoString(fromgchar(gs))
}

func gListboxAppend(widget *C.GtkWidget, what string) {
	var iter C.GtkTreeIter

	tv := getTreeViewFrom(widget)
	ls := (*C.GtkListStore)(unsafe.Pointer(C.gtk_tree_view_get_model(tv)))
	C.gtk_list_store_append(ls, &iter)
	cwhat := C.CString(what)
	defer C.free(unsafe.Pointer(cwhat))
	C.gtkListStoreSet(ls, &iter, cwhat)
}

func gListboxInsert(widget *C.GtkWidget, index int, what string) {
	var iter C.GtkTreeIter

	tv := getTreeViewFrom(widget)
	ls := (*C.GtkListStore)(unsafe.Pointer(C.gtk_tree_view_get_model(tv)))
	C.gtk_list_store_insert(ls, &iter, C.gint(index))
	cwhat := C.CString(what)
	defer C.free(unsafe.Pointer(cwhat))
	C.gtkListStoreSet(ls, &iter, cwhat)
}

func gListboxSelectedMulti(widget *C.GtkWidget) (indices []int) {
	var model *C.GtkTreeModel

	tv := getTreeViewFrom(widget)
	sel := C.gtk_tree_view_get_selection(tv)
	rows := C.gtk_tree_selection_get_selected_rows(sel, &model)
	defer C.g_list_free_full(rows, C.GDestroyNotify(unsafe.Pointer(C.gtk_tree_path_free)))
	// TODO needed?
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

func gListboxSelMultiTexts(widget *C.GtkWidget) (texts []string) {
	var model *C.GtkTreeModel
	var iter C.GtkTreeIter
	var gs *C.gchar

	tv := getTreeViewFrom(widget)
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
		texts[i] = C.GoString(fromgchar(gs))
		rows = rows.next
	}
	return texts
}

func gListboxDelete(widget *C.GtkWidget, index int) {
	var iter C.GtkTreeIter

	tv := getTreeViewFrom(widget)
	ls := (*C.GtkListStore)(unsafe.Pointer(C.gtk_tree_view_get_model(tv)))
	if !fromgbool(C.gtk_tree_model_iter_nth_child((*C.GtkTreeModel)(unsafe.Pointer(ls)), &iter, (*C.GtkTreeIter)(nil), C.gint(index))) {		// no such index
		// TODO
		return
	}
	C.gtk_list_store_remove(ls, &iter)
}

// this is a separate function because Combobox uses it too
func gtkTreeModelListLen(model *C.GtkTreeModel) int {
	// "As a special case, if iter is NULL, then the number of toplevel nodes is returned."
	return int(C.gtk_tree_model_iter_n_children(model, (*C.GtkTreeIter)(nil)))
}

func gListboxLen(widget *C.GtkWidget) int {
	tv := getTreeViewFrom(widget)
	model := C.gtk_tree_view_get_model(tv)
	return gtkTreeModelListLen(model)
}
