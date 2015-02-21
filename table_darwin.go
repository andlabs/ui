// 29 july 2014

package ui

import (
	"fmt"
	"reflect"
	"unsafe"
	"image"
)

// #include "objc_darwin.h"
import "C"

type table struct {
	*tablebase

	*scroller

	selected *event
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	id := C.newTable()
	t := &table{
		scroller:  newScroller(id, true), // border on Table
		tablebase: b,
		selected:  newEvent(),
	}
	t.fpreferredSize = t.xpreferredSize
	// also sets the delegate
	C.tableMakeDataSource(t.id, unsafe.Pointer(t))
	for i := 0; i < ty.NumField(); i++ {
		colname := ty.Field(i).Tag.Get("uicolumn")
		if colname == "" {
			colname = ty.Field(i).Name
		}
		cname := C.CString(colname)
		coltype := C.colTypeText
		editable := false
		switch {
		case ty.Field(i).Type == reflect.TypeOf((*image.RGBA)(nil)):
			coltype = C.colTypeImage
		case ty.Field(i).Type.Kind() == reflect.Bool:
			coltype = C.colTypeCheckbox
			editable = true
		}
		C.tableAppendColumn(t.id, C.intptr_t(i), cname, C.int(coltype), toBOOL(editable))
		C.free(unsafe.Pointer(cname)) // free now (not deferred) to conserve memory
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
			C.tableUpdate(t.id)
		})
	}()
}

func (t *table) Selected() int {
	t.RLock()
	defer t.RUnlock()
	return int(C.tableSelected(t.id))
}

func (t *table) Select(index int) {
	t.RLock()
	defer t.RUnlock()
	C.tableSelect(t.id, C.intptr_t(index))
}

func (t *table) OnSelected(f func()) {
	t.selected.set(f)
}

//export goTableDataSource_getValue
func goTableDataSource_getValue(data unsafe.Pointer, row C.intptr_t, col C.intptr_t, outtype *C.int) unsafe.Pointer {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	switch {
	case datum.Type() == reflect.TypeOf((*image.RGBA)(nil)):
		*outtype = C.colTypeImage
		d := datum.Interface().(*image.RGBA)
		img := C.toTableImage(unsafe.Pointer(pixelData(d)), C.intptr_t(d.Rect.Dx()), C.intptr_t(d.Rect.Dy()), C.intptr_t(d.Stride))
		return unsafe.Pointer(img)
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

//export tableSelectionChanged
func tableSelectionChanged(data unsafe.Pointer) {
	t := (*table)(data)
	t.selected.fire()
}

func (t *table) xpreferredSize(d *sizing) (width, height int) {
	s := C.tablePreferredSize(t.id)
	return int(s.width), int(s.height)
}
