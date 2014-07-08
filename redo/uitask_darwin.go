// 8 july 2014

package ui

import (
	"unsafe"
)

// #cgo CFLAGS: -DTODO
// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
import "C"

func uiinit() error {
	// TODO check error
	C.uiinit()
	return nil
}

func uimsgloop() {
	C.uimsgloop()
}

func issue(req *Request) {
	C.issue(unsafe.Pointer(req))
}

//export doissue
func doissue(r unsafe.Pointer) {
	req := (*Request)(unsafe.Pointer(r))
	perform(req)
}
