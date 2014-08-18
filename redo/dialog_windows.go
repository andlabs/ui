// 18 august 2014

package ui

import (
	"syscall"
	"unsafe"
	"reflect"
)

// #include "winapi_windows.h"
import "C"

// TODO move to common_windows.go
func wstrToString(wstr *C.WCHAR) string {
	n := C.wcslen((*C.wchar_t)(unsafe.Pointer(wstr)))
	xbuf := &reflect.SliceHeader{
		Data:	uintptr(unsafe.Pointer(wstr)),
		Len:		int(n + 1),
		Cap:		int(n + 1),
	}
	buf := (*[]uint16)(unsafe.Pointer(xbuf))
	return syscall.UTF16ToString(*buf)
}

func openFile() string {
	name := C.openFile()
	if name == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(name))
	return wstrToString(name)
}
