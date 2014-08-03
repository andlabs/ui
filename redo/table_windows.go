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
	*controlbase
	*tablebase
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	t := &table{
		controlbase:	newControl(C.xWC_LISTVIEW,
			C.LVS_REPORT | C.LVS_OWNERDATA | C.LVS_NOSORTHEADER | C.LVS_SHOWSELALWAYS | C.WS_HSCROLL | C.WS_VSCROLL,
			C.WS_EX_CLIENTEDGE),		// WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
		tablebase:		b,
	}
	C.setTableSubclass(t.hwnd, unsafe.Pointer(t))
	// LVS_EX_FULLROWSELECT gives us selection across the whole row, not just the leftmost column; this makes the list view work like on other platforms
	// LVS_EX_SUBITEMIMAGES gives us images in subitems, which will be important when both images and checkboxes are added
	C.tableAddExtendedStyles(t.hwnd, C.LVS_EX_FULLROWSELECT | C.LVS_EX_SUBITEMIMAGES)
	for i := 0; i < ty.NumField(); i++ {
		C.tableAppendColumn(t.hwnd, C.int(i), toUTF16(ty.Field(i).Name))
	}
	return t
}

func (t *table) Unlock() {
	t.unlock()
	// TODO RACE CONDITION HERE
	// I think there's a way to set the item count without causing a refetch of data that works around this...
	t.RLock()
	defer t.RUnlock()
	C.tableUpdate(t.hwnd, C.int(reflect.Indirect(reflect.ValueOf(t.data)).Len()))
}

//export tableGetCellText
func tableGetCellText(data unsafe.Pointer, row C.int, col C.int, str *C.LPWSTR) {
	t := (*table)(data)
	t.RLock()
	defer t.RUnlock()
	d := reflect.Indirect(reflect.ValueOf(t.data))
	datum := d.Index(int(row)).Field(int(col))
	s := fmt.Sprintf("%v", datum)
	*str = toUTF16(s)
}

func (t *table) setParent(p *controlParent) {
	basesetParent(t.controlbase, p)
}

func (t *table) containerShow() {
	basecontainerShow(t.controlbase)
}

func (t *table) containerHide() {
	basecontainerHide(t.controlbase)
}

func (t *table) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

const (
	// from C++ Template 05 in http://msdn.microsoft.com/en-us/library/windows/desktop/bb226818%28v=vs.85%29.aspx as this is the best I can do for now... (TODO see if I can reliably get item width/height from text size)
	tableWidth = 183
	tableHeight = 50
)

func (t *table) preferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(tableWidth, d), fromdlgunitsY(tableHeight, d)
}

func (t *table) commitResize(a *allocation, d *sizing) {
	basecommitResize(t.controlbase, a, d)
}

func (t *table) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(d)
}
