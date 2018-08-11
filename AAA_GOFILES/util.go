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

// We want Go itself to complain when we're out of memory.
// The allocators in cgo *should* do this, but there isn't a
// C.CMalloc(). There *is* a C.CBytes(), however, for transferring
// binary blobs from Go to C. If we pass this an arbitrary slice
// of the desired length, we get our C.CMalloc(). Using a slice
// that's always initialized to zero gives us the ZeroMemory()
// for free.
var uimallocBytes = make([]byte, 1024)		// 1024 bytes first

//export uimalloc
func uimalloc(n C.size_t) unsafe.Pointer {
	if n > C.size_t(len(uimallocBytes)) {
		// TODO round n up to a multiple of a power of 2?
		// for instance 0x1234 bytes -> 0x1800 bytes
		uimallocBytes = make([]byte, n)
	}
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
