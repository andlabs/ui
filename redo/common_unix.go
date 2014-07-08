// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

func fromgstr(s *C.gchar) string {
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

func togstr(s string) *C.gchar {
	return (*C.gchar)(unsafe.Pointer(C.CString(s)))
}

func freegstr(s *C.gchar) {
	C.free(unsafe.Pointer(s))
}
