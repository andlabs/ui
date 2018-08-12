// 12 december 2015

package ui

import (
	"unsafe"
)

// #include <stdlib.h>
// #include "util.h"
import "C"

// We want Go itself to complain when we're out of memory.
// The allocators in cgo *should* do this, but there isn't a
// C.CMalloc(). There *is* a C.CBytes(), however, for transferring
// binary blobs from Go to C. If we pass this an arbitrary slice
// of the desired length, we get our C.CMalloc(). Using a slice
// that's always initialized to zero gives us the memset(0)
// (or ZeroMemory()) for free.
var allocBytes = make([]byte, 1024)		// 1024 bytes first

//export pkguiAlloc
func pkguiAlloc(n C.size_t) unsafe.Pointer {
	if n > C.size_t(len(allocBytes)) {
		// TODO round n up to a multiple of a power of 2?
		// for instance 0x1234 bytes -> 0x1800 bytes
		allocBytes = make([]byte, n)
	}
	return C.CBytes(allocBytes[:n])
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
