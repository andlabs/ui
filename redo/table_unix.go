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
	*widgetbase
	*tablebase

	treewidget	*C.GtkWidget
	treeview		*C.GtkTreeView

	scrollc		*C.GtkContainer
	scrollwindow	*C.GtkScrolledWindow

	model		*C.goTableModel
	modelgtk		*C.GtkTreeModel

	// stuff required by GtkTreeModel
	nColumns		C.gint
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	widget := C.gtk_tree_view_new()
	scroller := C.gtk_scrolled_window_new(nil, nil)
	t := &table{
		// TODO kludge
		widgetbase:	newWidget(scroller),
		tablebase:		b,
		treewidget:	widget,
		treeview:		(*C.GtkTreeView)(unsafe.Pointer(widget)),
		scrollc:		(*C.GtkContainer)(unsafe.Pointer(scroller)),
		scrollwindow:	(*C.GtkScrolledWindow)(unsafe.Pointer(scroller)),
	}
	// give the scrolled window a border (thanks to jlindgren in irc.gimp.net/#gtk+)
	C.gtk_scrolled_window_set_shadow_type(t.scrollwindow, C.GTK_SHADOW_IN)
	C.gtk_container_add(t.scrollc, t.treewidget)
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

func (t *table) preferredSize(d *sizing) (width int, height int) {
	var r C.GtkRequisition

	C.gtk_widget_get_preferred_size(t.treewidget, nil, &r)
	return int(r.width), int(r.height)
}


func (t *table) Unlock() {
	t.unlock()
	// TODO RACE CONDITION HERE
	// not sure about this one...
	t.RLock()
	defer t.RUnlock()
	// TODO
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
