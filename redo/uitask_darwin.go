// 8 july 2014

package ui

import (
	"unsafe"
)

// #cgo CFLAGS: -mmacosx-version-min=10.7 -DMACOSX_DEPLOYMENT_TARGET=10.7
// #cgo LDFLAGS: -mmacosx-version-min=10.7 -lobjc -framework Foundation -framework AppKit
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

func uistop() {
	C.uistop()
}

func issue(f func()) {
	C.issue(unsafe.Pointer(&f))
}

//export doissue
func doissue(fp unsafe.Pointer) {
	perform(fp)
}
