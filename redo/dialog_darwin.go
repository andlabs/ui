// 19 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

func (w *window) openFile(f func(filename string)) {
	C.openFile(w.id, unsafe.Pointer(&f))
}

//export finishOpenFile
func finishOpenFile(fname *C.char, data unsafe.Pointer) {
	f := (*func(string))(data)
	if fname == nil {
		(*f)("")
		return
	}
	defer C.free(unsafe.Pointer(fname))
	(*f)(C.GoString(fname))
}
