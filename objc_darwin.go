// 28 february 2014

package ui

import (
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation
// #include <stdlib.h>
// #include "objc_darwin.h"
// /* cgo doesn't like Nil */
// Class NilClass = Nil;
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

	_alloc = sel_getUid("alloc")
	_new = sel_getUid("new")
	_release = sel_getUid("release")
	_stringWithUTF8String = sel_getUid("stringWithUTF8String:")
	_UTF8String = sel_getUid("UTF8String")
)

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

// These consolidate the NSScrollView code (used by listbox_darwin.go and area_darwin.go) into a single place.

var (
	_NSScrollView = objc_getClass("NSScrollView")

	_setHasHorizontalScroller = sel_getUid("setHasHorizontalScroller:")
	_setHasVerticalScroller = sel_getUid("setHasVerticalScroller:")
	_setAutohidesScrollers = sel_getUid("setAutohidesScrollers:")
	_setDocumentView = sel_getUid("setDocumentView:")
	_documentView = sel_getUid("documentView")
)

func newScrollView(content C.id) C.id {
	scrollview := C.objc_msgSend_noargs(_NSScrollView, _alloc)
	scrollview = initWithDummyFrame(scrollview)
	C.objc_msgSend_bool(scrollview, _setHasHorizontalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setHasVerticalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setAutohidesScrollers, C.BOOL(C.YES))
	C.objc_msgSend_id(scrollview, _setDocumentView, content)
	return scrollview
}

func getScrollViewContent(scrollview C.id) C.id {
	return C.objc_msgSend_noargs(scrollview, _documentView)
}

// These create new classes.

// selector contains the information for a new selector.
type selector struct {
	name	string
	imp		uintptr	// not unsafe.Pointer because https://code.google.com/p/go/issues/detail?id=7665
	itype		itype
	desc		string	// for error reporting
}

// sel_[returntype] or sel_[returntype]_[arguments] (after the required self/sel arguments)
type itype uint
const (
	sel_void_id itype = iota
	sel_bool_id
	sel_bool
	sel_void_rect
	sel_terminatereply_id
	nitypes
)

var itypes = [nitypes][]C.char{
	sel_void_id:			[]C.char{'v', '@', ':', '@', 0},
	sel_bool_id:			[]C.char{'c', '@', ':', '@', 0},
	sel_bool:				[]C.char{'c', '@', ':', 0},
	sel_void_rect:			nil,			// see init() below
	sel_terminatereply_id:	nil,
}

func init() {
	// see encodedNSRect in bleh_darwin.m
	x := make([]C.char, 0, 256)	// more than enough
	x = append(x, 'v', '@', ':')
	y := C.GoString(C.encodedNSRect)
	for _, b := range y {
		x = append(x, C.char(b))
	}
	x = append(x, 0)
	itypes[sel_void_rect] = x

	x = make([]C.char, 0, 256)	// more than enough
	y = C.GoString(C.encodedTerminateReply)
	for _, b := range y {
		x = append(x, C.char(b))
	}
	x = append(x, '@', ':', '@', 0)
	itypes[sel_terminatereply_id] = x
}

func makeClass(name string, super C.id, sels []selector, desc string) (id C.id, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	// an id that describes a class is itself a Class
	// thanks to Psy| in irc.freenode.net/##objc
	c := C.objc_allocateClassPair(C.Class(unsafe.Pointer(super)), cname, 0)
	if c == C.NilClass {
		err = fmt.Errorf("unable to create Objective-C class %s for %s; reason unknown", name, desc)
		return
	}
	C.objc_registerClassPair(c)
	for _, v := range sels {
		ok := C.class_addMethod(c, sel_getUid(v.name),
			C.IMP(unsafe.Pointer(v.imp)), &itypes[v.itype][0])
		if ok == C.BOOL(C.NO) {
			err = fmt.Errorf("unable to add selector %s to class %s (needed for %s; reason unknown)", v.name, name, v.desc)
			return
		}
	}
	return objc_getClass(name), nil
}
