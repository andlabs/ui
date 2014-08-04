// 29 july 2014

package ui

import (
	"fmt"
	"reflect"
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type table struct {
	*tablebase

	_id			C.id
	scroller		*scroller
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	id := C.newTable()
	t := &table{
		_id:			id,
		scroller:		newScroller(id),
		tablebase:		b,
	}
	C.tableMakeDataSource(t._id, unsafe.Pointer(t))
	for i := 0; i < ty.NumField(); i++ {
		cname := C.CString(ty.Field(i).Name)
		C.tableAppendColumn(t._id, cname)
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
	C.tableUpdate(t._id)
}

//export goTableDataSource_getValue
func goTableDataSource_getValue(data unsafe.Pointer, row C.intptr_t, col C.intptr_t) *C.char {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	s := fmt.Sprintf("%v", datum)
	return C.CString(s)
}

//export goTableDataSource_getRowCount
func goTableDataSource_getRowCount(data unsafe.Pointer) C.intptr_t {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	return C.intptr_t(d.Len())
}

func (t *table) id() C.id {
	return t._id
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
	basegetAuxResizeInfo(t, d)
}
