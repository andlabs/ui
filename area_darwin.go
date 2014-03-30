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
	_goArea = "goArea"
)

var (
	_drawRect = sel_getUid("drawRect:")
)

func mkAreaClass() error {
	areaclass, err := makeAreaClass(_goArea)
	if err != nil {
		return fmt.Errorf("error creating Area backend class: %v", err)
	}
//	// TODO rename this function (it overrides anyway)
	// addAreaViewDrawMethod() is in bleh_darwin.m
	ok := C.addAreaViewDrawMethod(areaclass)
	if ok != C.BOOL(C.YES) {
		return fmt.Errorf("error overriding Area drawRect: method; reason unknown")
	}
	return nil
}

//export areaView_drawRect
func areaView_drawRect(self C.id, rect C.struct_xrect) {
	// TODO
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
