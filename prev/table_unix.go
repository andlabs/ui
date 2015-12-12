// +build !windows,!darwin

// 29 july 2014

package ui

import (
	"fmt"
	"reflect"
	"unsafe"
	"image"
)

// #include "gtk_unix.h"
// extern void goTableModel_toggled(GtkCellRendererToggle *, gchar *, gpointer);
// extern void tableSelectionChanged(GtkTreeSelection *, gpointer);
import "C"

type table struct {
	*tablebase

	*scroller
	treeview *C.GtkTreeView

	model     *C.goTableModel
	modelgtk  *C.GtkTreeModel
	selection *C.GtkTreeSelection

	selected *event

	// stuff required by GtkTreeModel
	nColumns C.gint
	old      C.gint
	types    []C.GType
	crtocol  map[*C.GtkCellRendererToggle]int
}

var (
	attribText   = togstr("text")
	attribPixbuf = togstr("pixbuf")
	attribActive = togstr("active")
)

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	widget := C.gtk_tree_view_new()
	t := &table{
		scroller:  newScroller(widget, true, true, false), // natively scrollable; has a border; no overlay
		tablebase: b,
		treeview:  (*C.GtkTreeView)(unsafe.Pointer(widget)),
		crtocol:   make(map[*C.GtkCellRendererToggle]int),
		selected:  newEvent(),
	}
	model := C.newTableModel(unsafe.Pointer(t))
	t.model = model
	t.modelgtk = (*C.GtkTreeModel)(unsafe.Pointer(model))
	t.selection = C.gtk_tree_view_get_selection(t.treeview)
	g_signal_connect(
		C.gpointer(unsafe.Pointer(t.selection)),
		"changed",
		C.GCallback(C.tableSelectionChanged),
		C.gpointer(unsafe.Pointer(t)))
	C.gtk_tree_view_set_model(t.treeview, t.modelgtk)
	for i := 0; i < ty.NumField(); i++ {
		colname := ty.Field(i).Tag.Get("uicolumn")
		if colname == "" {
			colname = ty.Field(i).Name
		}
		cname := togstr(colname)
		switch {
		case ty.Field(i).Type == reflect.TypeOf((*image.RGBA)(nil)):
			// can't use GDK_TYPE_PIXBUF here because it's a macro that expands to a function and cgo hates that
			t.types = append(t.types, C.gdk_pixbuf_get_type())
			C.tableAppendColumn(t.treeview, C.gint(i), cname,
				C.gtk_cell_renderer_pixbuf_new(), attribPixbuf)
		case ty.Field(i).Type.Kind() == reflect.Bool:
			t.types = append(t.types, C.G_TYPE_BOOLEAN)
			cr := C.gtk_cell_renderer_toggle_new()
			crt := (*C.GtkCellRendererToggle)(unsafe.Pointer(cr))
			t.crtocol[crt] = i
			g_signal_connect(C.gpointer(unsafe.Pointer(cr)),
				"toggled",
				C.GCallback(C.goTableModel_toggled),
				C.gpointer(unsafe.Pointer(t)))
			C.tableAppendColumn(t.treeview, C.gint(i), cname,
				cr, attribActive)
		default:
			t.types = append(t.types, C.G_TYPE_STRING)
			C.tableAppendColumn(t.treeview, C.gint(i), cname,
				C.gtk_cell_renderer_text_new(), attribText)
		}
		freegstr(cname) // free now (not deferred) to conserve memory
	}
	// and for some GtkTreeModel boilerplate
	t.nColumns = C.gint(ty.NumField())
	return t
}

func (t *table) Lock() {
	t.tablebase.Lock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	t.old = C.gint(d.Len())
}

func (t *table) Unlock() {
	t.unlock()
	// there's a possibility that user actions can happen at this point, before the view is updated
	// alas, this is something we have to deal with, because Unlock() can be called from any thread
	go func() {
		Do(func() {
			t.RLock()
			defer t.RUnlock()
			d := reflect.Indirect(reflect.ValueOf(t.data))
			new := C.gint(d.Len())
			C.tableUpdate(t.model, t.old, new)
		})
	}()
}

func (t *table) Selected() int {
	var iter C.GtkTreeIter

	t.RLock()
	defer t.RUnlock()
	if C.gtk_tree_selection_get_selected(t.selection, nil, &iter) == C.FALSE {
		return -1
	}
	path := C.gtk_tree_model_get_path(t.modelgtk, &iter)
	if path == nil {
		panic(fmt.Errorf("invalid iter in Table.Selected()"))
	}
	defer C.gtk_tree_path_free(path)
	return int(*C.gtk_tree_path_get_indices(path))
}

func (t *table) Select(index int) {
	t.RLock()
	defer t.RUnlock()
	C.gtk_tree_selection_unselect_all(t.selection)
	if index == -1 {
		return
	}
	path := C.gtk_tree_path_new()
	defer C.gtk_tree_path_free(path)
	C.gtk_tree_path_append_index(path, C.gint(index))
	C.gtk_tree_selection_select_path(t.selection, path)
}

func (t *table) OnSelected(f func()) {
	t.selected.set(f)
}

//export goTableModel_get_n_columns
func goTableModel_get_n_columns(model *C.GtkTreeModel) C.gint {
	tm := (*C.goTableModel)(unsafe.Pointer(model))
	t := (*table)(tm.gotable)
	return t.nColumns
}

//export goTableModel_get_column_type
func goTableModel_get_column_type(model *C.GtkTreeModel, column C.gint) C.GType {
	tm := (*C.goTableModel)(unsafe.Pointer(model))
	t := (*table)(tm.gotable)
	return t.types[column]
}

//export goTableModel_do_get_value
func goTableModel_do_get_value(data unsafe.Pointer, row C.gint, col C.gint, value *C.GValue) {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	switch {
	case datum.Type() == reflect.TypeOf((*image.RGBA)(nil)):
		d := datum.Interface().(*image.RGBA)
		pixbuf := toIconSizedGdkPixbuf(d)
		C.g_value_init(value, C.gdk_pixbuf_get_type())
		object := C.gpointer(unsafe.Pointer(pixbuf))
		// use g_value_take_object() so the GtkTreeView becomes the pixbuf's owner
		C.g_value_take_object(value, object)
	case datum.Kind() == reflect.Bool:
		d := datum.Interface().(bool)
		C.g_value_init(value, C.G_TYPE_BOOLEAN)
		C.g_value_set_boolean(value, togbool(d))
	default:
		s := fmt.Sprintf("%v", datum)
		str := togstr(s)
		defer freegstr(str)
		C.g_value_init(value, C.G_TYPE_STRING)
		C.g_value_set_string(value, str)
	}
}

//export goTableModel_getRowCount
func goTableModel_getRowCount(data unsafe.Pointer) C.gint {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	return C.gint(d.Len())
}

//export goTableModel_toggled
func goTableModel_toggled(cr *C.GtkCellRendererToggle, pathstr *C.gchar, data C.gpointer) {
	t := (*table)(unsafe.Pointer(data))
	t.Lock()
	defer t.Unlock()
	path := C.gtk_tree_path_new_from_string(pathstr)
	if len := C.gtk_tree_path_get_depth(path); len != 1 {
		panic(fmt.Errorf("invalid path of depth %d given to goTableModel_toggled()", len))
	}
	// dereference return value to get our sole member
	row := *C.gtk_tree_path_get_indices(path)
	col := t.crtocol[cr]
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	datum.SetBool(!datum.Bool())
}

//export tableSelectionChanged
func tableSelectionChanged(sel *C.GtkTreeSelection, data C.gpointer) {
	t := (*table)(unsafe.Pointer(data))
	t.selected.fire()
}
