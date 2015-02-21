// 28 july 2014

package ui

// TODO
// - why are we not getting keyboard input (focus?) on a mouse click?
// - are we getting keyboard input (focus?) on tab?
// - random freezes on Windows 7 when resizing headers or clicking new rows; likely another package ui infrastructure issue though...
// - investigate japanese display in the table headers on wine (it has worked in the past in ANSI mode via Shift-JIS with what I assume is the same font, so huh?)

import (
	"fmt"
	"reflect"
	"unsafe"
	"sync"
	"image"
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
	free			map[C.uintptr_t]bool
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
		free:		make(map[C.uintptr_t]bool),
	}
	t.fpreferredSize = t.xpreferredSize
	t.chainresize = t.fresize
	t.fresize = t.xresize
	C.setTableSubclass(t.hwnd, unsafe.Pointer(t))
	// TODO listview didn't need this; someone mentioned (TODO) it uses the small caption font???
	C.controlSetControlFont(t.hwnd)
	for i := 0; i < ty.NumField(); i++ {
		coltype := C.WPARAM(C.tableColumnText)
		switch {
		case ty.Field(i).Type == reflect.TypeOf((*image.RGBA)(nil)):
			coltype = C.tableColumnImage
		case ty.Field(i).Type.Kind() == reflect.Bool:
			coltype = C.tableColumnCheckbox
		}
		colname := ty.Field(i).Tag.Get("uicolumn")
		if colname == "" {
			colname = ty.Field(i).Name
		}
		ccolname := toUTF16(colname)
		C.SendMessageW(t.hwnd, C.tableAddColumn, coltype, C.LPARAM(uintptr(unsafe.Pointer(ccolname))))
		// TODO free ccolname
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
			C.gotableSetRowCount(t.hwnd, C.intptr_t(reflect.Indirect(reflect.ValueOf(t.data)).Len()))
		})
	}()
}

func (t *table) Selected() int {
	t.RLock()
	defer t.RUnlock()
	return int(C.tableSelectedItem(t.hwnd))
}

func (t *table) Select(index int) {
	t.RLock()
	defer t.RUnlock()
	C.tableSelectItem(t.hwnd, C.intptr_t(index))
}

func (t *table) OnSelected(f func()) {
	t.selected.set(f)
}

//export tableGetCell
func tableGetCell(data unsafe.Pointer, tnm *C.tableNM) C.LRESULT {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(tnm.row)).Field(int(tnm.column))
	switch {
	case datum.Type() == reflect.TypeOf((*image.RGBA)(nil)):
		i := datum.Interface().(*image.RGBA)
		hbitmap := C.toBitmap(unsafe.Pointer(i), C.intptr_t(i.Rect.Dx()), C.intptr_t(i.Rect.Dy()))
		bitmap := C.uintptr_t(uintptr(unsafe.Pointer(hbitmap)))
		t.freeLock.Lock()
		t.free[bitmap] = true		// bitmap freed with C.freeBitmap()
		t.freeLock.Unlock()
		return C.LRESULT(bitmap)
	case datum.Kind() == reflect.Bool:
		if datum.Bool() == true {
			return C.TRUE
		}
		return C.FALSE
	default:
		s := fmt.Sprintf("%v", datum)
		text := C.uintptr_t(uintptr(unsafe.Pointer(toUTF16(s))))
		t.freeLock.Lock()
		t.free[text] = false		// text freed with C.free()
		t.freeLock.Unlock()
		return C.LRESULT(text)
	}
}

//export tableFreeCellData
func tableFreeCellData(gotable unsafe.Pointer, data C.uintptr_t) {
	t := (*table)(gotable)
	t.freeLock.Lock()
	defer t.freeLock.Unlock()
	b, ok := t.free[data]
	if !ok {
		panic(fmt.Errorf("undefined data %p in tableFreeData()", data))
	}
	if b == false {
		C.free(unsafe.Pointer(uintptr(data)))
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
func tableToggled(data unsafe.Pointer, row C.intptr_t, col C.intptr_t) {
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
