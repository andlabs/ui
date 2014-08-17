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

	images		[]C.id
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	id := C.newTable()
	t := &table{
		_id:			id,
		scroller:		newScroller(id, true),		// border on Table
		tablebase:		b,
	}
	C.tableMakeDataSource(t._id, unsafe.Pointer(t))
	for i := 0; i < ty.NumField(); i++ {
		cname := C.CString(ty.Field(i).Name)
		coltype := C.colTypeText
		editable := false
		switch {
		case ty.Field(i).Type == reflect.TypeOf(ImageIndex(0)):
			coltype = C.colTypeImage
		case ty.Field(i).Type.Kind() == reflect.Bool:
			coltype = C.colTypeCheckbox
			editable = true
		}
		C.tableAppendColumn(t._id, C.intptr_t(i), cname, C.int(coltype), toBOOL(editable))
		C.free(unsafe.Pointer(cname))		// free now (not deferred) to conserve memory
	}
	return t
}

func (t *table) Unlock() {
	t.unlock()
	// there's a possibility that user actions can happen at this point, before the view is updated
	// alas, this is something we have to deal with, because Unlock() can be called from any thread
	go func() {
		Do(func() {
			t.RLock()
			defer t.RUnlock()
			C.tableUpdate(t._id)
		})
	}()
}

func (t *table) LoadImageList(i ImageList) {
	i.apply(&t.images)
}

//export goTableDataSource_getValue
func goTableDataSource_getValue(data unsafe.Pointer, row C.intptr_t, col C.intptr_t, outtype *C.int) unsafe.Pointer {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	switch {
	case datum.Type() == reflect.TypeOf(ImageIndex(0)):
		*outtype = C.colTypeImage
		d := datum.Interface().(ImageIndex)
		return unsafe.Pointer(t.images[d])
	case datum.Kind() == reflect.Bool:
		*outtype = C.colTypeCheckbox
		if datum.Bool() == true {
			// return a non-nil pointer
			// outtype isn't Go-side so it'll work
			return unsafe.Pointer(outtype)
		}
		return nil
	default:
		s := fmt.Sprintf("%v", datum)
		return unsafe.Pointer(C.CString(s))
	}
}

//export goTableDataSource_getRowCount
func goTableDataSource_getRowCount(data unsafe.Pointer) C.intptr_t {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	return C.intptr_t(d.Len())
}

//export goTableDataSource_toggled
func goTableDataSource_toggled(data unsafe.Pointer, row C.intptr_t, col C.intptr_t, checked C.BOOL) {
	t := (*table)(data)
	t.Lock()
	defer t.Unlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	datum.SetBool(fromBOOL(checked))
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
	s := C.tablePreferredSize(t._id)
	return int(s.width), int(s.height)
}

func (t *table) commitResize(c *allocation, d *sizing) {
	t.scroller.commitResize(c, d)
}

func (t *table) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
