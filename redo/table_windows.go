// 28 july 2014

package ui

import (
	"fmt"
	"unsafe"
	"reflect"
)

// #include "winapi_windows.h"
import "C"

type table struct {
	*tablebase
	_hwnd		C.HWND
	noautosize	bool
	colcount		C.int
	hotrow		C.int
	hotcol		C.int
	pushedrow	C.int
	pushedcol	C.int
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	t := &table{
		_hwnd:		C.newControl(C.xWC_LISTVIEW,
			C.LVS_REPORT | C.LVS_OWNERDATA | C.LVS_NOSORTHEADER | C.LVS_SHOWSELALWAYS | C.WS_HSCROLL | C.WS_VSCROLL | C.WS_TABSTOP,
			C.WS_EX_CLIENTEDGE),		// WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
		tablebase:		b,
		hotrow:		-1,
		hotcol:		-1,
		pushedrow:	-1,
		pushedcol:	-1,
	}
	C.setTableSubclass(t._hwnd, unsafe.Pointer(t))
	// LVS_EX_FULLROWSELECT gives us selection across the whole row, not just the leftmost column; this makes the list view work like on other platforms
	// LVS_EX_SUBITEMIMAGES gives us images in subitems, which will be important when both images and checkboxes are added
	C.tableAddExtendedStyles(t._hwnd, C.LVS_EX_FULLROWSELECT | C.LVS_EX_SUBITEMIMAGES)
	C.tableSetCheckboxImageList(t._hwnd)
	for i := 0; i < ty.NumField(); i++ {
		C.tableAppendColumn(t._hwnd, C.int(i), toUTF16(ty.Field(i).Name))
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
			C.tableUpdate(t._hwnd, C.int(reflect.Indirect(reflect.ValueOf(t.data)).Len()))
		})
	}()
}

func (t *table) LoadImageList(il ImageList) {
	il.apply(t._hwnd, C.LVM_SETIMAGELIST, C.LVSIL_SMALL)
}

//export tableGetCell
func tableGetCell(data unsafe.Pointer, item *C.LVITEMW) {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(item.iItem)).Field(int(item.iSubItem))
	// TODO figure out why changing item.mask causes crashes or why "it just works"
	switch {
	case datum.Type() == reflect.TypeOf(ImageIndex(0)):
		item.iImage = C.int(datum.Interface().(ImageIndex))
	case datum.Kind() == reflect.Bool:
		item.stateMask = C.LVIS_STATEIMAGEMASK
		// state image index is 1-based
		curstate := ((item.state & C.LVIS_STATEIMAGEMASK) >> 12)
		if curstate > 0 {
			curstate--
		}
		if datum.Bool() == true {
			curstate |= C.checkboxStateChecked
		} else {
			curstate &^= C.checkboxStateChecked
		}
		if item.iItem == t.hotrow && item.iSubItem == t.hotcol {
			curstate |= C.checkboxStateHot
		} else {
			curstate &^= C.checkboxStateHot
		}
		if item.iItem == t.pushedrow && item.iSubItem == t.pushedcol {
			curstate |= C.checkboxStatePushed
		} else {
			curstate &^= C.checkboxStatePushed
		}
		item.state = (curstate + 1) << 12
	default:
		s := fmt.Sprintf("%v", datum)
		item.pszText = toUTF16(s)
	}
}

// the column autoresize policy is simple:
// on every table.commitResize() call, if the columns have not been resized by the user, autoresize
func (t *table) autoresize() {
	t.RLock()
	defer t.RUnlock()
	if !t.noautosize {
		C.tableAutosizeColumns(t._hwnd, t.colcount)
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

//export tableSetHot
func tableSetHot(data unsafe.Pointer, row C.int, col C.int) {
	t := (*table)(data)
	redraw := (row != t.hotrow || col != t.hotcol)
	t.hotrow = row
	t.hotcol = col
	if redraw {
		C.tableUpdate(t._hwnd, C.int(reflect.Indirect(reflect.ValueOf(t.data)).Len()))
	}
}

//export tablePushed
func tablePushed(data unsafe.Pointer, row C.int, col C.int) {
	t := (*table)(data)
	t.pushedrow = row
	t.pushedcol = col
	C.tableUpdate(t._hwnd, C.int(reflect.Indirect(reflect.ValueOf(t.data)).Len()))
}

//export tableToggled
func tableToggled(data unsafe.Pointer, row C.int, col C.int) {
	t := (*table)(data)
	t.Lock()
	defer func() {
		// reset for next time
		t.pushedrow = -1
		t.pushedcol = -1
		// and THEN unlock so the reset takes effect
		t.Unlock()
	}()
	if row == -1 || col == -1 {		// discard extras sent by handle() in table_windows.c
		return
	}
	if row != t.pushedrow || col != t.pushedcol {				// mouse moved out
		return
	}
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	if datum.Kind() == reflect.Bool {
		datum.SetBool(!datum.Bool())
		return
	}
	panic(fmt.Errorf("tableSetHot() on non-checkbox at (%d, %d)", row, col))
}

func (t *table) hwnd() C.HWND {
	return t._hwnd
}

func (t *table) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *table) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

const (
	// from C++ Template 05 in http://msdn.microsoft.com/en-us/library/windows/desktop/bb226818%28v=vs.85%29.aspx as this is the best I can do for now
	// there IS a message LVM_APPROXIMATEVIEWRECT that can do calculations, but it doesn't seem to work right when asked to base its calculations on the current width/height on Windows and wine...
	tableWidth = 183
	tableHeight = 50
)

func (t *table) preferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(tableWidth, d), fromdlgunitsY(tableHeight, d)
}

func (t *table) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
	t.RLock()
	defer t.RUnlock()
	t.autoresize()
}

func (t *table) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}
