// 28 february 2014
package main

/*
These are wrapper functions for the functions in bleh_darwin.m to wrap around stdint.h type casting.

This will eventually be expanded to include the other Objective-C runtime support functions.
*/

// #cgo LDFLAGS: -lobjc -framework Foundation
// #include "objc_darwin.h"
import "C"

func objc_msgSend_rect(obj C.id, sel C.SEL, x int, y int, w int, h int) C.id {
	return C._objc_msgSend_rect(obj, sel,
		C.int64_t(x), C.int64_t(y), C.int64_t(w), C.int64_t(h))
}

func objc_msgSend_uint(obj C.id, sel C.SEL, a uintptr) C.id {
	return C._objc_msgSend_uint(obj, sel, C.uintptr_t(a))
}

func objc_msgSend_rect_uint_uint_bool(obj C.id, sel C.SEL, x int, y int, w int, h int, b uintptr, c uintptr, d C.BOOL) C.id {
	return C._objc_msgSend_rect_uint_uint_bool(obj, sel,
		C.int64_t(x), C.int64_t(y), C.int64_t(w), C.int64_t(h),
		C.uintptr_t(b), C.uintptr_t(c), d)
}
