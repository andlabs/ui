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

	// stuff required by GtkTreeModel
	nColumns		C.gint
	old			C.gint
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	widget := C.gtk_tree_view_new()
	t := &table{
		scroller:			newScroller(widget, true),
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
		C.tableAppendColumn(t.treeview, C.gint(i), cname)
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
	// TODO RACE CONDITION HERE
	// not sure about this one...
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	new := C.gint(d.Len())
	C.tableUpdate(t.model, t.old, new)
}

//export goTableModel_get_n_columns
func goTableModel_get_n_columns(model *C.GtkTreeModel) C.gint {
	tm := (*C.goTableModel)(unsafe.Pointer(model))
	t := (*table)(tm.gotable)
	return t.nColumns
}

//export goTableModel_do_get_value
func goTableModel_do_get_value(data unsafe.Pointer, row C.gint, col C.gint) *C.gchar {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	s := fmt.Sprintf("%v", datum)
	return togstr(s)
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
