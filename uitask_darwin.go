// 28 february 2014

//
package ui

import (
	"fmt"
	"runtime"
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
import "C"

var uitask chan func()

var (
	_NSAutoreleasePool = objc_getClass("NSAutoreleasePool")
	_NSValue           = objc_getClass("NSValue")

	_valueWithPointer            = sel_getUid("valueWithPointer:")
	_performSelectorOnMainThread = sel_getUid("performSelectorOnMainThread:withObject:waitUntilDone:")
	_pointerValue = sel_getUid("pointerValue")
	_run          = sel_getUid("run")
)

func ui(main func()) error {
	runtime.LockOSThread()

	uitask = make(chan func())

	NSApp, err := initCocoa()
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
				C.BOOL(C.YES)) // wait so we can properly drain the autorelease pool; on other platforms we wind up waiting anyway (since the main thread can only handle one thing at a time) so
			objc_release(pool)
		}
	}()

	go main()

	C.objc_msgSend_noargs(NSApp, _run)
	return nil
}

// TODO move to init_darwin.go?

var (
	_NSApplication = objc_getClass("NSApplication")

	_sharedApplication   = sel_getUid("sharedApplication")
	_setActivationPolicy = sel_getUid("setActivationPolicy:")
)

func initCocoa() (NSApp C.id, err error) {
	NSApp = C.objc_msgSend_noargs(_NSApplication, _sharedApplication)
	r := C.objc_msgSend_int(NSApp, _setActivationPolicy,
		0) // NSApplicationActivationPolicyRegular
	if C.BOOL(uintptr(unsafe.Pointer(r))) != C.BOOL(C.YES) {
		err = fmt.Errorf("error setting NSApplication activation policy (basically identifies our program as a separate program; needed for several things, such as Dock icon, application menu, window resizing, etc.) (unknown reason)")
		return
	}
	// TODO we need to call [NSApp activateIgnoringOtherApps:YES] here when we are ready to have the program become active (for now I won't)
	err = mkAppDelegate()
	return
}

//export appDelegate_uitask
func appDelegate_uitask(self C.id, sel C.SEL, arg C.id) {
	p := C.objc_msgSend_noargs(arg, _pointerValue)
	f := (*func())(unsafe.Pointer(p))
	(*f)()
}
