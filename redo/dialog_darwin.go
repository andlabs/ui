// 19 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

func openFile() string {
	fname := C.openFile()
	if fname == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(fname))
	return C.GoString(fname)
}
