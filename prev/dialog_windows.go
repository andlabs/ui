// 18 august 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

func (w *window) openFile(f func(filename string)) {
	C.openFile(w.hwnd, unsafe.Pointer(&f))
}

//export finishOpenFile
func finishOpenFile(name *C.WCHAR, fp unsafe.Pointer) {
	f := (*func(string))(fp)
	if name == nil {
		(*f)("")
		return
	}
	defer C.free(unsafe.Pointer(name))
	(*f)(wstrToString(name))
}
