// 12 december 2015

package ui

import (
	"unsafe"
)

// #include <stdlib.h>
// // TODO remove when switching to Go 1.7
// #include <string.h>
import "C"

// TODO move this to C.CBytes() when switching to Go 1.7

//export uimalloc
func uimalloc(n C.size_t) unsafe.Pointer {
	p := C.malloc(n)
	if p == nil {
		panic("out of memory in uimalloc()")
	}
	C.memset(p, 0, n)
	return p
}

func freestr(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func tobool(b C.int) bool {
	return b != 0
}

func frombool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
