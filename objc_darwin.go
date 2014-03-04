// 28 february 2014

//
package ui

import (
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation
// #include <stdlib.h>
// #include "objc_darwin.h"
import "C"

func objc_getClass(class string) C.id {
	cclass := C.CString(class)
	defer C.free(unsafe.Pointer(cclass))

	return C.objc_getClass(cclass)
}

func sel_getUid(sel string) C.SEL {
	csel := C.CString(sel)
	defer C.free(unsafe.Pointer(csel))

	return C.sel_getUid(csel)
}

// Common Objective-C types and selectors.
var (
	_NSObject = objc_getClass("NSObject")
	_NSString = objc_getClass("NSString")

	_alloc                = sel_getUid("alloc")
	_new                  = sel_getUid("new")
	_release              = sel_getUid("release")
	_stringWithUTF8String = sel_getUid("stringWithUTF8String:")
	_UTF8String           = sel_getUid("UTF8String")
	_setDelegate          = sel_getUid("setDelegate:")
)

// some helper functions

func objc_alloc(class C.id) C.id {
	return C.objc_msgSend_noargs(class, _alloc)
}

func objc_new(class C.id) C.id {
	return C.objc_msgSend_noargs(class, _new)
}

func objc_release(obj C.id) {
	C.objc_msgSend_noargs(obj, _release)
}

func toNSString(str string) C.id {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	return C.objc_msgSend_str(_NSString,
		_stringWithUTF8String,
		cstr)
}

func fromNSString(str C.id) string {
	cstr := C.objc_msgSend_noargs(str, _UTF8String)
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

func objc_setDelegate(obj C.id, delegate C.id) {
	C.objc_msgSend_id(obj, _setDelegate, delegate)
}

/*
These are wrapper functions for the functions in bleh_darwin.m to wrap around stdint.h type casting.
*/

func objc_msgSend_rect(obj C.id, sel C.SEL, x int, y int, w int, h int) C.id {
	return C._objc_msgSend_rect(obj, sel,
		C.int64_t(x), C.int64_t(y), C.int64_t(w), C.int64_t(h))
}

func objc_msgSend_uint(obj C.id, sel C.SEL, a uintptr) C.id {
	return C._objc_msgSend_uint(obj, sel, C.uintptr_t(a))
}

func objc_msgSend_rect_bool(obj C.id, sel C.SEL, x int, y int, w int, h int, b C.BOOL) C.id {
	return C._objc_msgSend_rect_bool(obj, sel,
		C.int64_t(x), C.int64_t(y), C.int64_t(w), C.int64_t(h),
		b)
}

func objc_msgSend_rect_uint_uint_bool(obj C.id, sel C.SEL, x int, y int, w int, h int, b uintptr, c uintptr, d C.BOOL) C.id {
	return C._objc_msgSend_rect_uint_uint_bool(obj, sel,
		C.int64_t(x), C.int64_t(y), C.int64_t(w), C.int64_t(h),
		C.uintptr_t(b), C.uintptr_t(c), d)
}
