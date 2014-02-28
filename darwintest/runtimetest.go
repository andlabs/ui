// 27 february 2014
package main

import (
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation
// #include <stdlib.h>
// #include <objc/message.h>
// #include <objc/objc.h>
// #include <objc/runtime.h>
// /* cgo doesn't handle ... */
// id objc_msgSend_noargs(id obj, SEL sel) { return objc_msgSend(obj, sel); }
// id objc_msgSend_strarg(id obj, SEL sel, char *a) { return objc_msgSend(obj, sel, a); }
import "C"

func main() {
	_NSString := C.CString("NSString")
	defer C.free(unsafe.Pointer(_NSString))
	_stringWithUTF8String := C.CString("stringWithUTF8String:")
	defer C.free(unsafe.Pointer(_stringWithUTF8String))
	_UTF8String := C.CString("UTF8String")
	defer C.free(unsafe.Pointer(_UTF8String))
	_hello := C.CString("hello, world\n")
	defer C.free(unsafe.Pointer(_hello))

	NSString := C.objc_getClass(_NSString)
	stringWithUTF8String :=
		C.sel_getUid(_stringWithUTF8String)
	str := C.objc_msgSend_strarg(NSString,
		stringWithUTF8String,
		_hello)
	UTF8String :=
		C.sel_getUid(_UTF8String)
	res := C.objc_msgSend_noargs(str,
			UTF8String)
	cres := (*C.char)(unsafe.Pointer(res))
	fmt.Printf("%s", C.GoString(cres))
}
