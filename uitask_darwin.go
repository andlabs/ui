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
// extern void appDelegate_applicationDidFinishLaunching(id, SEL, id);
// extern void appDelegate_uitask(id, SEL, id);
import "C"

// temporary for now
func msgBox(string, string){}
func msgBoxError(string, string){}

var uitask chan func()

var mtret chan interface{}

var (
	_NSAutoreleasePool = objc_getClass("NSAutoreleasePool")
	_NSValue = objc_getClass("NSValue")

	_uitask = sel_getUid("uitask:")
	_valueWithPointer = sel_getUid("valueWithPointer:")
	_performSelectorOnMainThread =
		sel_getUid("performSelectorOnMainThread:withObject:waitUntilDone:")
	_pointerValue = sel_getUid("pointerValue")
)

func ui(initDone chan error) {
	runtime.LockOSThread()

	uitask = make(chan func())
	mtret = make(chan interface{})
	go mainThread()
	v := <-mtret
	if err, ok := v.(error); ok {
		initDone <- fmt.Errorf("error initializing Cocoa: %v", err)
		return
	}
	appDelegate := v.(C.id)

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
}

const (
	_goAppDelegate = "goAppDelegate"
)

var (
	_NSApplication = objc_getClass("NSApplication")

	_sharedApplication = sel_getUid("sharedApplication")
	_applicationDidFinishLaunching = sel_getUid("applicationDidFinishLaunching:")
	_run = sel_getUid("run")
)

func mainThread() {
	runtime.LockOSThread()

	_NSApp := C.objc_msgSend_noargs(_NSApplication, _sharedApplication)
	appdelegateclass, err := makeDelegateClass(_goAppDelegate)
	if err != nil {
		mtret <- fmt.Errorf("error creating NSApplication delegate: %v", err)
		return
	}
	err = addDelegateMethod(appdelegateclass, _applicationDidFinishLaunching,
		C.appDelegate_applicationDidFinishLaunching)
	if err != nil {
		mtret <- fmt.Errorf("error adding NSApplication delegate applicationDidFinishLaunching: method (to start UI loop): %v", err)
		return
	}
	err = addDelegateMethod(appdelegateclass, _uitask, C.appDelegate_uitask)
	if err != nil {
		mtret <- fmt.Errorf("error adding NSApplication delegate uitask: method (to do UI tasks): %v", err)
		return
	}
	appDelegate := objc_new(objc_getClass(_goAppDelegate))
	objc_setDelegate(_NSApp, appDelegate)
	// and that's it, really
	C.objc_msgSend_noargs(_NSApp, _run)
}

//export appDelegate_applicationDidFinishLaunching
func appDelegate_applicationDidFinishLaunching(self C.id, sel C.SEL, arg C.id) {
	mtret <- self
}

//export appDelegate_uitask
func appDelegate_uitask(self C.id, sel C.SEL, arg C.id) {
	p := C.objc_msgSend_noargs(arg, _pointerValue)
	f := (*func ())(unsafe.Pointer(p))
	(*f)()
}
