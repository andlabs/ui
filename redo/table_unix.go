// 29 july 2014

package ui

import (
//	"fmt"
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
	// TODO model
	for i := 0; i < ty.NumField(); i++ {
		cname := togstr(ty.Field(i).Name)
		C.tableAppendColumn(t.treeview, cname)
		freegstr(cname)		// free now (not deferred) to conserve memory
	}
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
