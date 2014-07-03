// 28 february 2014

package ui

import (
	"unsafe"
)

// #include <stdlib.h>
// #include "objc_darwin.h"
import "C"

func toNSString(str string) C.id {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	return C.toNSString(cstr)
}

func fromNSString(str C.id) string {
	return C.GoString(C.fromNSString(str))
}

func toBOOL(what bool) C.BOOL {
	if what {
		return C.YES
	}
	return C.NO
}

// These consolidate the NSScrollView code (used by listbox_darwin.go and area_darwin.go) into a single place.

func makeScrollView(content C.id) C.id {
	return C.makeScrollView(content)
}

func getScrollViewContent(scrollview C.id) C.id {
	return C.scrollViewContent(scrollview)
}
