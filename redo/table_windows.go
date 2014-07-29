// 28 july 2014

package ui

import (
	"unsafe"
	"reflect"
)

// #include "winapi_windows.h"
import "C"

type table struct {
	*widgetbase
	*tablebase
}

func finishNewTable(b *tablebase, ty reflect.Type) Table {
	t := &table{
		widgetbase:	newWidget(C.xWC_LISTVIEW,
			C.LVS_REPORT | C.LVS_OWNERDATA | C.LVS_NOSORTHEADER | C.LVS_SHOWSELALWAYS | C.WS_HSCROLL | C.WS_VSCROLL,
			C.WS_EX_CLIENTEDGE),		// WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
		tablebase:		b,
	}
	C.setTableSubclass(t.hwnd, unsafe.Pointer(&t))
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
