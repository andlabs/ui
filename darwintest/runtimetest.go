// 27 february 2014
package main

import (
	"fmt"
	"time"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
// #include "objc_darwin.h"
// Class NilClass = Nil; /* for newtypes.go */
// id Nilid = nil;
import "C"

var (
	_NSApplication = objc_getClass("NSApplication")
	_NSNotificationCenter = objc_getClass("NSNotificationCenter")

	_sharedApplication = sel_getUid("sharedApplication")
	_defaultCenter = sel_getUid("defaultCenter")
	_run = sel_getUid("run")
)

var NSApp C.id
var defNC C.id
var delegate C.id
var notesel C.SEL

func init() {
	// need an NSApplication first - see https://github.com/TooTallNate/NodObjC/issues/21
	NSApp = C.objc_msgSend_noargs(_NSApplication, _sharedApplication)

	defNC = C.objc_msgSend_noargs(_NSNotificationCenter, _defaultCenter)

	selW := sel_getUid("windowShouldClose:")
	selB := sel_getUid("buttonClicked:")
	selN := sel_getUid("gotNotification:")
	mk("hello", selW, selB, selN)
	delegate = objc_alloc(objc_getClass("hello"))

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
	// TODO copy in the rest?
)

var (
	_NSAutoreleasePool = objc_getClass("NSAutoreleasePool")

	_performSelectorOnMainThread =
		sel_getUid("performSelectorOnMainThread:withObject:waitUntilDone:")
)

func notify(source string) {
	// we need to make an NSAutoreleasePool, otherwise we get leak warnings on stderr
	pool := objc_new(_NSAutoreleasePool)
	src := toNSString(source)
	C.objc_msgSend_sel_id_bool(
		delegate,
		_performSelectorOnMainThread,
		notesel,
		src,
		C.BOOL(C.YES))			// wait so we can properly drain the autorelease pool; on other platforms we wind up waiting anyway (since the main thread can only handle one thing at a time) so
	objc_release(pool)
}

var (
	_NSWindow = objc_getClass("NSWindow")
	_NSButton = objc_getClass("NSButton")

	_initWithContentRect = sel_getUid("initWithContentRect:styleMask:backing:defer:")
	_setDelegate = sel_getUid("setDelegate:")
	_makeKeyAndOrderFront = sel_getUid("makeKeyAndOrderFront:")
	_contentView = sel_getUid("contentView")
	_initWithFrame = sel_getUid("initWithFrame:")
	_setTarget = sel_getUid("setTarget:")
	_setAction = sel_getUid("setAction:")
	_setBezelStyle = sel_getUid("setBezelStyle:")
	_addSubview = sel_getUid("addSubview:")
)

func main() {
	style := uintptr(NSTitledWindowMask | NSClosableWindowMask)
	backing := uintptr(NSBackingStoreBuffered)
	deferx := C.BOOL(C.YES)
	window := objc_alloc(_NSWindow)
	window = objc_msgSend_rect_uint_uint_bool(window,
		_initWithContentRect,
		100, 100, 320, 240,
		style, backing, deferx)
	C.objc_msgSend_id(window, _makeKeyAndOrderFront, window)
	C.objc_msgSend_id(window, _setDelegate, delegate)
	windowView := C.objc_msgSend_noargs(window, _contentView)

	button := objc_alloc(_NSButton)
	button = objc_msgSend_rect(button,
		_initWithFrame,
		20, 20, 200, 200)
	C.objc_msgSend_id(button, _setTarget, delegate)
	C.objc_msgSend_sel(button,
		_setAction,
		sel_getUid("buttonClicked:"))
	objc_msgSend_uint(button, _setBezelStyle, NSRoundedBezelStyle)
	C.objc_msgSend_id(windowView, _addSubview, button)

	go func() {
		for {
			<-time.After(5 * time.Second)
			fmt.Println("five seconds passed; sending notification...")
			notify("timer")
		}
	}()

	C.objc_msgSend_noargs(NSApp, _run)
}
