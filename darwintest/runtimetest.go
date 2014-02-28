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

func main() {
	_hello := C.CString("hello, world\n")
	defer C.free(unsafe.Pointer(_hello))

	NSString := objc_getClass("NSString")
	stringWithUTF8String :=
		sel_getUid("stringWithUTF8String:")
	str := C.objc_msgSend_strarg(NSString,
		stringWithUTF8String,
		_hello)
	UTF8String := sel_getUid("UTF8String")
	res := C.objc_msgSend_noargs(str,
			UTF8String)
	cres := (*C.char)(unsafe.Pointer(res))
	fmt.Printf("%s", C.GoString(cres))
}
