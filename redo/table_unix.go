// +build !windows,!darwin

// 29 july 2014

package ui

import (
	"fmt"
	"reflect"
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type table struct {
	*tablebase

	_widget		*C.GtkWidget
	treeview		*C.GtkTreeView
	scroller		*scroller

	model		*C.goTableModel
	modelgtk		*C.GtkTreeModel

	pixbufs		[]*C.GdkPixbuf

	// stuff required by GtkTreeModel
	nColumns		C.gint
	old			C.gint
	types		[]C.GType
}

var (
	attribText = togstr("text")
	attribPixbuf = togstr("pixbuf")
)

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	widget := C.gtk_tree_view_new()
	t := &table{
		scroller:			newScroller(widget, true, true),			// natively scrollable; has a border
		tablebase:			b,
		_widget:			widget,
		treeview:			(*C.GtkTreeView)(unsafe.Pointer(widget)),
	}
	model := C.newTableModel(unsafe.Pointer(t))
	t.model = model
	t.modelgtk = (*C.GtkTreeModel)(unsafe.Pointer(model))
	C.gtk_tree_view_set_model(t.treeview, t.modelgtk)
	for i := 0; i < ty.NumField(); i++ {
		cname := togstr(ty.Field(i).Name)
		switch ty.Field(i).Type {
		case reflect.TypeOf(ImageIndex(0)):
			// can't use GDK_TYPE_PIXBUF here because it's a macro that expands to a function and cgo hates that
			t.types = append(t.types, C.gdk_pixbuf_get_type())
			C.tableAppendColumn(t.treeview, C.gint(i), cname,
				C.gtk_cell_renderer_pixbuf_new(), attribPixbuf)
		default:
			t.types = append(t.types, C.G_TYPE_STRING)
			C.tableAppendColumn(t.treeview, C.gint(i), cname,
				C.gtk_cell_renderer_text_new(), attribText)
		}
		freegstr(cname)		// free now (not deferred) to conserve memory
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

func (t *table) LoadImageList(i ImageList) {
	i.apply(&t.pixbufs)
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
	switch d := datum.Interface().(type) {
	case ImageIndex:
		C.g_value_init(value, C.gdk_pixbuf_get_type())
		C.g_value_set_object(value, C.gpointer(unsafe.Pointer(t.pixbufs[d])))
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

func (t *table) widget() *C.GtkWidget {
	return t._widget
}

func (t *table) setParent(p *controlParent) {
	t.scroller.setParent(p)
}

func (t *table) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *table) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(t, d)
}

func (t *table) commitResize(c *allocation, d *sizing) {
	t.scroller.commitResize(c, d)
}

func (t *table) getAuxResizeInfo(d *sizing) {
	// a Label to the left of a Table should be vertically aligned to the top
	d.shouldVAlignTop = true
}
