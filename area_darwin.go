// 29 march 2014

package ui

import (
	"fmt"
	"unsafe"
	"image"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
//// #include <HIToolbox/Events.h>
// #include "objc_darwin.h"
// extern void areaView_drawRect(id, struct xrect);
// extern BOOL areaView_isFlipped(id, SEL);
import "C"

const (
	__goArea = "goArea"
)

var (
	_goArea C.id

	_drawRect = sel_getUid("drawRect:")
	_isFlipped = sel_getUid("isFlipped")
)

func mkAreaClass() error {
	areaclass, err := makeAreaClass(__goArea)
	if err != nil {
		return fmt.Errorf("error creating Area backend class: %v", err)
	}
	// addAreaViewDrawMethod() is in bleh_darwin.m
	ok := C.addAreaViewDrawMethod(areaclass)
	if ok != C.BOOL(C.YES) {
		return fmt.Errorf("error overriding Area drawRect: method; reason unknown")
	}
	// TODO rename this function (it overrides anyway)
	err = addDelegateMethod(areaclass, _isFlipped,
		C.areaView_isFlipped, area_boolret)
	if err != nil {
		return fmt.Errorf("error overriding Area isFlipped method: %v", err)
	}
	_goArea = objc_getClass(__goArea)
	return nil
}

var (
	_drawAtPoint = sel_getUid("drawAtPoint:")
)

//export areaView_drawRect
func areaView_drawRect(self C.id, rect C.struct_xrect) {
	s := getSysData(self)
	// TODO clear clip rect
	// rectangles in Cocoa are origin/size, not point0/point1; if we don't watch for this, weird things will happen when scrolling
	// TODO change names EVERYWHERE ELSE to match
	cliprect := image.Rect(int(rect.x), int(rect.y), int(rect.x + rect.width), int(rect.y + rect.height))
	max := C.objc_msgSend_stret_rect_noargs(self, _frame)
	cliprect = image.Rect(0, 0, int(max.width), int(max.height)).Intersect(cliprect)
	if cliprect.Empty() {			// no intersection; nothing to paint
		return
	}
	i := s.handler.Paint(cliprect)
	C.drawImage(
		unsafe.Pointer(&i.Pix[0]), C.int64_t(i.Rect.Dx()), C.int64_t(i.Rect.Dy()), C.int64_t(i.Stride),
		C.int64_t(cliprect.Min.X), C.int64_t(cliprect.Min.Y))
}

//export areaView_isFlipped
func areaView_isFlipped(self C.id, sel C.SEL) C.BOOL {
	return C.BOOL(C.YES)
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
	area = objc_msgSend_rect(area, _initWithFrame,
		0, 0, 100, 100)
	// TODO others?
	area = newAreaScrollView(area)
	addControl(parentWindow, area)
	return area
}

// TODO combine the below with the delegate stuff

var (
	_NSView = objc_getClass("NSView")
	_NSView_Class = C.Class(unsafe.Pointer(_NSView))
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
	area_boolret = []C.char{'c', '@', ':', 0}			// BOOL (*)(id, SEL)
)
