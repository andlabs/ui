// 29 march 2014

package ui

import (
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
//// #include <HIToolbox/Events.h>
// #include "objc_darwin.h"
// extern void areaView_drawRect(id, struct xrect);
import "C"

const (
	__goArea = "goArea"
)

var (
	_goArea C.id

	_drawRect = sel_getUid("drawRect:")
)

func mkAreaClass() error {
	areaclass, err := makeAreaClass(__goArea)
	if err != nil {
		return fmt.Errorf("error creating Area backend class: %v", err)
	}
//	// TODO rename this function (it overrides anyway)
	// addAreaViewDrawMethod() is in bleh_darwin.m
	ok := C.addAreaViewDrawMethod(areaclass)
	if ok != C.BOOL(C.YES) {
		return fmt.Errorf("error overriding Area drawRect: method; reason unknown")
	}
	_goArea = objc_getClass(__goArea)
	return nil
}

//export areaView_drawRect
func areaView_drawRect(self C.id, rect C.struct_xrect) {
	// TODO
fmt.Println(rect)
}

// TODO combine these with the listbox functions?

func newAreaScrollView(area C.id) C.id {
	scrollview := objc_alloc(_NSScrollView)
	scrollview = objc_msgSend_rect(scrollview, _initWithFrame,
		0, 0, 100, 100)
	C.objc_msgSend_bool(scrollview, _setHasHorizontalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setHasVerticalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setAutohidesScrollers, C.BOOL(C.YES))
	C.objc_msgSend_id(scrollview, _setDocumentView, area)
	return scrollview
}

func areaInScrollView(scrollview C.id) C.id {
	return C.objc_msgSend_noargs(scrollview, _documentView)
}

func makeArea(parentWindow C.id, alternate bool) C.id {
	area := objc_alloc(_goArea)
println(area)
	area = objc_msgSend_rect(area, _initWithFrame,
		0, 0, 100, 100)
println("out")
	// TODO others?
	area = newAreaScrollView(area)
	addControl(parentWindow, area)
	return area
}

// TODO combine the below with the delegate stuff

var (
	_NSView = objc_getClass("NSView")
	_NSView_Class = C.object_getClass(_NSView)
)

func makeAreaClass(name string) (C.Class, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	c := C.objc_allocateClassPair(_NSView_Class, cname, 0)
	if c == C.NilClass {
		return C.NilClass, fmt.Errorf("unable to create Objective-C class %s for Area; reason unknown", name)
	}
	C.objc_registerClassPair(c)
	return c, nil
}

var (
	// delegate_rect in bleh_darwin.m
)
