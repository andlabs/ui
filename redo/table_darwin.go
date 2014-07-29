// 29 july 2014

package ui

import (
//	"fmt"
	"reflect"
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type table struct {
	*widgetbase
	*tablebase

	// TODO kludge
	table		C.id
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	id := C.newTable()
	t := &table{
		widgetbase:	newWidget(C.newScrollView(id)),
		tablebase:		b,
		table:		id,
	}
	// TODO model
	for i := 0; i < ty.NumField(); i++ {
		cname := C.CString(ty.Field(i).Name)
		C.tableAppendColumn(t.table, cname)
		C.free(unsafe.Pointer(cname))		// free now (not deferred) to conserve memory
	}
	return t
}

func (t *table) Unlock() {
	t.unlock()
	// TODO RACE CONDITION HERE
	// not sure about this one...
	t.RLock()
	defer t.RUnlock()
	C.tableUpdate(t.table)
}
