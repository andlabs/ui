// 28 february 2014
package ui

import (
	"fmt"
	"runtime"
	"unsafe"
)

/*
We will create an Objective-C class goAppDelegate. It contains two methods:
	- (void)applicationDidFinishLoading:(NSNotification *)unused
		will signal to ui() that we are now in the Cocoa event loop; we make our goAppDelegate instance the NSApplication delegate
	- (void)uitask:(NSValue *)functionPointer
		the function that actually performs our UI task functions; it is called with NSObject's performSelectorOnMainThread system
*/

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
// extern void appDelegate_uitask(id, SEL, id);
import "C"

// temporary for now
func msgBox(string, string){}
func msgBoxError(string, string){}

var uitask chan func()

var (
	_NSAutoreleasePool = objc_getClass("NSAutoreleasePool")
	_NSValue = objc_getClass("NSValue")

	_uitask = sel_getUid("uitask:")
	_valueWithPointer = sel_getUid("valueWithPointer:")
	_performSelectorOnMainThread =
		sel_getUid("performSelectorOnMainThread:withObject:waitUntilDone:")
	_pointerValue = sel_getUid("pointerValue")
	_run = sel_getUid("run")
)

func ui(main func()) error {
	runtime.LockOSThread()

	uitask = make(chan func())

	NSApp, appDelegate, err := initCocoa()
	if err != nil {
		return err
	}

	// Cocoa must run on the first thread created by the program, so we run our dispatcher on another thread instead
	go func() {
		for f := range uitask {
			// we need to make an NSAutoreleasePool, otherwise we get leak warnings on stderr
			pool := objc_new(_NSAutoreleasePool)
			fp := C.objc_msgSend_ptr(_NSValue, _valueWithPointer,
				unsafe.Pointer(&f))
			C.objc_msgSend_sel_id_bool(
				appDelegate,
				_performSelectorOnMainThread,
				_uitask,
				fp,
				C.BOOL(C.YES))			// wait so we can properly drain the autorelease pool; on other platforms we wind up waiting anyway (since the main thread can only handle one thing at a time) so
			objc_release(pool)
		}
	}()

	go main()

	C.objc_msgSend_noargs(NSApp, _run)
	return nil
}

// TODO move to init_darwin.go?

const (
	_goAppDelegate = "goAppDelegate"
)

var (
	_NSApplication = objc_getClass("NSApplication")

	_sharedApplication = sel_getUid("sharedApplication")
)

func initCocoa() (NSApp C.id, appDelegate C.id, err error) {
	var appdelegateclass C.Class

	NSApp = C.objc_msgSend_noargs(_NSApplication, _sharedApplication)
	err = addDelegateMethod(appdelegateclass, _uitask, C.appDelegate_uitask)
	if err != nil {
		err = fmt.Errorf("error adding NSApplication delegate uitask: method (to do UI tasks): %v", err)
		return
	}
	// TODO using objc_new() causes a segfault; find out why
	// TODO make alloc followed by init (I thought NSObject provided its own init?)
	appDelegate = objc_alloc(objc_getClass(_goAppDelegate))

	return
}

//export appDelegate_uitask
func appDelegate_uitask(self C.id, sel C.SEL, arg C.id) {
	p := C.objc_msgSend_noargs(arg, _pointerValue)
	f := (*func ())(unsafe.Pointer(p))
	(*f)()
}
