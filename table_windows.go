// 28 july 2014

package ui

import (
	"fmt"
	"reflect"
	"unsafe"
	"sync"
)

// #include "winapi_windows.h"
import "C"

type table struct {
	*tablebase
	*controlSingleHWND
	noautosize bool
	colcount   C.int
	selected   *event
	chainresize		func(x int, y int, width int, height int, d *sizing)
	freeTexts		map[unsafe.Pointer]bool
	freeLock		sync.Mutex
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	hwnd := C.newControl(C.xtableWindowClass,
		C.WS_HSCROLL|C.WS_VSCROLL|C.WS_TABSTOP,
		C.WS_EX_CLIENTEDGE)		// WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
	t := &table{
		controlSingleHWND:		newControlSingleHWND(hwnd),
		tablebase: b,
		selected:  newEvent(),
		free:		make(map[unsafe.Pointer]bool),
	}
	t.fpreferredSize = t.xpreferredSize
	t.chainresize = t.fresize
	t.fresize = t.xresize
	C.setTableSubclass(t.hwnd, unsafe.Pointer(t))
	for i := 0; i < ty.NumField(); i++ {
		coltype := C.WPARAM(C.tableColumnText)
		switch ty.Field(i).Type {
		case ty.Field(i).Type == reflect.TypeOf((*image.RGBA)(nil)):
			coltype = C.tableColumnImage
		case ty.Field(i).Type.Kind() == reflect.Bool:
			coltype = C.tableColumnCheckbox
		}
		colname := toUTF16(ty.Field(i).Name)
		C.SendMessageW(t.hwnd, C.tableAddColumn, coltype, C.LPARAM(uintptr(unsafe.Pointer(colname))))
	}
	t.colcount = C.int(ty.NumField())
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
			C.SendMessageW(t.hwnd, C.tableSetRowCount, 0, C.LPARAM(C.intptr_t(reflect.Indirect(reflect.ValueOf(t.data)).Len())))
		})
	}()
}

func (t *table) Selected() int {
	t.RLock()
	defer t.RUnlock()
//TODO	return int(C.tableSelectedItem(t.hwnd))
	return -1
}

func (t *table) Select(index int) {
	t.RLock()
	defer t.RUnlock()
//TODO	C.tableSelectItem(t.hwnd, C.intptr_t(index))
}

func (t *table) OnSelected(f func()) {
	t.selected.set(f)
}

//export tableGetCell
func tableGetCell(data unsafe.Pointer, item *C.LVITEMW) C.LRESULT {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(item.iItem)).Field(int(item.iSubItem))
	isText := true
	switch {
	case datum.Type() == reflect.TypeOf((*image.RGBA)(nil)):
		i := datum.Interface().(*image.RGBA)
		hbitmap := C.toBitmap(unsafe.Pointer(i), C.intptr_t(i.Dx()), C.intptr_t(i.Dy()))
		bitmap := unsafe.Pointer(hbitmap)
		t.freeLock.Lock()
		t.free[bitmap] = true		// bitmap freed with C.freeBitmap()
		t.freeLock.Unlock()
		return C.LRESULT(uintptr(bmp))
	case datum.Kind() == reflect.Bool:
		if datum.Bool() == true {
			return C.TRUE
		}
		return C.FALSE
	default:
		s := fmt.Sprintf("%v", datum)
		text := unsafe.Pointer(toUTF16(s))
		t.freeLock.Lock()
		t.free[text] = false		// text freed with C.free()
		t.freeLock.Unlock()
		return C.LRESULT(uintptr(text))
	}
}

//export tableFreeData
func tableFreeData(gotable unsafe.Pointer, data unsafe.Pointer) {
	t := (*table)(gotable)
	t.freeLock.Lock()
	defer t.freeLock.Unlock()
	b, ok := t.free[data]
	if !ok {
		panic(fmt.Errorf("undefined data %p in tableFreeData()", data))
	}
	if b == false {
		C.free(data)
	} else {
		C.freeBitmap(data)
	}
	delete(t.free, data)
}

// the column autoresize policy is simple:
// on every table.commitResize() call, if the columns have not been resized by the user, autoresize
func (t *table) autoresize() {
	t.RLock()
	defer t.RUnlock()
	if !t.noautosize {
//TODO		C.tableAutosizeColumns(t.hwnd, t.colcount)
	}
}

//export tableStopColumnAutosize
func tableStopColumnAutosize(data unsafe.Pointer) {
	t := (*table)(data)
	t.noautosize = true
}

//export tableColumnCount
func tableColumnCount(data unsafe.Pointer) C.int {
	t := (*table)(data)
	return t.colcount
}

//export tableToggled
func tableToggled(data unsafe.Pointer, row C.int, col C.int) {
	t := (*table)(data)
	t.Lock()
	defer t.Unlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	if datum.Kind() == reflect.Bool {
		datum.SetBool(!datum.Bool())
		return
	}
	panic(fmt.Errorf("tableSetHot() on non-checkbox at (%d, %d)", row, col))
}

//export tableSelectionChanged
func tableSelectionChanged(data unsafe.Pointer) {
	t := (*table)(data)
	t.selected.fire()
}

const (
	// from C++ Template 05 in http://msdn.microsoft.com/en-us/library/windows/desktop/bb226818%28v=vs.85%29.aspx as this is the best I can do for now
	// there IS a message LVM_APPROXIMATEVIEWRECT that can do calculations, but it doesn't seem to work right when asked to base its calculations on the current width/height on Windows and wine...
	tableWidth  = 183
	tableHeight = 50
)

func (t *table) xpreferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(tableWidth, d), fromdlgunitsY(tableHeight, d)
}

func (t *table) xresize(x int, y int, width int, height int, d *sizing) {
	t.chainresize(x, y, width, height, d)
	t.RLock()
	defer t.RUnlock()
	t.autoresize()
}
