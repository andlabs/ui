// 12 december 2015

package ui

import (
	"unsafe"
)

// #include <stdlib.h>
import "C"

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
