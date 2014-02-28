// 27 february 2014
package main

import (
	"fmt"
	"unsafe"
	"time"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
// #include <objc/message.h>
// #include <objc/objc.h>
// #include <objc/runtime.h>
// /* TODO <objc/NSObjCRuntime.h not found?!?! */
// /* TODO this HAS to be unsafe */
// typedef unsigned long NSUInteger;
// /* avoid depending on Objective-C */
// #include <CoreGraphics/CGGeometry.h>
// /* cgo doesn't handle ... */
// id objc_msgSend_noargs(id obj, SEL sel) { return objc_msgSend(obj, sel); }
// id objc_msgSend_strarg(id obj, SEL sel, char *a) { return objc_msgSend(obj, sel, a); }
// id objc_msgSend_NSRect_uint_uint_bool(id obj, SEL sel, CGRect a, NSUInteger b, NSUInteger c, BOOL d) { return objc_msgSend(obj, sel, a, b, c, d); }
// id objc_msgSend_id(id obj, SEL sel, id a) { return objc_msgSend(obj, sel, a); }
// id objc_msgSend_NSRect(id obj, SEL sel, CGRect a) { return objc_msgSend(obj, sel, a); }
// id objc_msgSend_sel(id obj, SEL sel, SEL a) { return objc_msgSend(obj, sel, a); }
// id objc_msgSend_uint(id obj, SEL sel, NSUInteger a) { return objc_msgSend(obj, sel, a); }
// id objc_msgSend_id_sel_id_id(id obj, SEL sel, id a, SEL b, id c, id d) { return objc_msgSend(obj, sel, a, b, c, d); }
// id objc_msgSend_id_id_id(id obj, SEL sel, id a, id b, id c) { return objc_msgSend(obj, sel, a, b, c); }
// id objc_msgSend_id_id(id obj, SEL sel, id a, id b) { return objc_msgSend(obj, sel, a, b); }
// id objc_msgSend_sel_id_bool(id obj, SEL sel, SEL a, id b, BOOL c) { return objc_msgSend(obj, sel, a, b, c); }
// Class NilClass = Nil; /* for newtypes.go */
// id Nilid = nil;
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

var NSApp C.id
var defNC C.id
var delegate C.id
var notesel C.SEL

func init() {
	// need an NSApplication first - see https://github.com/TooTallNate/NodObjC/issues/21
	NSApplication := objc_getClass("NSApplication")
	sharedApplication := sel_getUid("sharedApplication")
	NSApp = C.objc_msgSend_noargs(NSApplication, sharedApplication)

	defNC = C.objc_msgSend_noargs(
		objc_getClass("NSNotificationCenter"),
		sel_getUid("defaultCenter"))

	selW := sel_getUid("windowShouldClose:")
	selB := sel_getUid("buttonClicked:")
	selN := sel_getUid("gotNotification:")
	mk("hello", selW, selB, selN)
	delegate = C.objc_msgSend_noargs(
		objc_getClass("hello"),
		alloc)

	notesel = selN
}

const (
	NSBorderlessWindowMask = 0
	NSTitledWindowMask = 1 << 0
	NSClosableWindowMask = 1 << 1
	NSMiniaturizableWindowMask = 1 << 2
	NSResizableWindowMask = 1 << 3
	NSTexturedBackgroundWindowMask = 1 << 8
)

const (
//	NSBackingStoreRetained = 0			// "You should not use this mode."
//	NSBackingStoreNonretained = 1		// "You should not use this mode."
	NSBackingStoreBuffered = 2
)

const (
	NSRoundedBezelStyle = 1
)

var alloc = sel_getUid("alloc")

func notify(source string) {
	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))

	src := C.objc_msgSend_strarg(
		objc_getClass("NSString"),
		sel_getUid("stringWithUTF8String:"),
		csource)
	C.objc_msgSend_sel_id_bool(
		delegate,
		sel_getUid("performSelectorOnMainThread:withObject:waitUntilDone:"),
		notesel,
		src,
		C.BOOL(C.FALSE))		// don't wait; we're using a channel for this (in the final ui code)
}

func main() {
	NSWindow := objc_getClass("NSWindow")
	NSWindowinit :=
		sel_getUid("initWithContentRect:styleMask:backing:defer:")
	setDelegate := sel_getUid("setDelegate:")
	makeKeyAndOrderFront := sel_getUid("makeKeyAndOrderFront:")

	rect := C.CGRect{
		origin:	C.CGPoint{100, 100},
		size:		C.CGSize{320, 240},
	}
	style := C.NSUInteger(NSTitledWindowMask | NSClosableWindowMask)
	backing := C.NSUInteger(NSBackingStoreBuffered)
	deferx := C.BOOL(C.YES)
	window := C.objc_msgSend_noargs(NSWindow, alloc)
	window = C.objc_msgSend_NSRect_uint_uint_bool(window, NSWindowinit, rect, style, backing, deferx)
	C.objc_msgSend_id(window, makeKeyAndOrderFront, window)
	C.objc_msgSend_id(window, setDelegate,
		delegate)
	windowView := C.objc_msgSend_noargs(window,
		sel_getUid("contentView"))

	NSButton := objc_getClass("NSButton")
	button := C.objc_msgSend_noargs(NSButton, alloc)
	button = C.objc_msgSend_NSRect(button,
		sel_getUid("initWithFrame:"),
		C.CGRect{
			origin:	C.CGPoint{20, 20},
			size:		C.CGSize{200, 200},
		})
	C.objc_msgSend_id(button,
		sel_getUid("setTarget:"),
		delegate)
	C.objc_msgSend_sel(button,
		sel_getUid("setAction:"),
		sel_getUid("buttonClicked:"))
	C.objc_msgSend_uint(button,
		sel_getUid("setBezelStyle:"),
		C.NSUInteger(NSRoundedBezelStyle))
	C.objc_msgSend_id(windowView,
		sel_getUid("addSubview:"),
		button)

	go func() {
		for {
			<-time.After(5 * time.Second)
			fmt.Println("five seconds passed; sending notification...")
			notify("timer")
		}
	}()

	C.objc_msgSend_noargs(NSApp,
		sel_getUid("run"))
}

func helloworld() {
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
