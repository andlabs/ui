// 28 february 2014

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
	_NSValue = objc_getClass("NSValue")

	_valueWithPointer = sel_getUid("valueWithPointer:")
	_performSelectorOnMainThread =
		sel_getUid("performSelectorOnMainThread:withObject:waitUntilDone:")
	_stop = sel_getUid("stop:")
	_postEventAtStart = sel_getUid("postEvent:atStart:")
	_pointerValue = sel_getUid("pointerValue")
	_run = sel_getUid("run")
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
			pool := C.objc_msgSend_noargs(_NSAutoreleasePool, _new)
			fp := C.objc_msgSend_ptr(_NSValue, _valueWithPointer,
				unsafe.Pointer(&f))
			C.objc_msgSend_sel_id_bool(
				appDelegate,
				_performSelectorOnMainThread,
				_uitask,
				fp,
				C.BOOL(C.YES))			// wait so we can properly drain the autorelease pool; on other platforms we wind up waiting anyway (since the main thread can only handle one thing at a time) so
			C.objc_msgSend_noargs(pool, _release)
		}
	}()

	go func() {
		main()
		uitask <- func() {
			// -[NSApplication stop:] stops the event loop; it won't do a clean termination, but we're not too concerned with that (at least not on the other platforms either so)
			// we can't call -[NSApplication terminate:] because that will just quit the program, ensuring we never leave ui.Go()
			C.objc_msgSend_id(NSApp, _stop, NSApp)
			// simply calling -[NSApplication stop:] is not good enough, as the stop flag is only checked when an event comes in
			// we have to create a "proper" event; a blank event will just throw an exception
			C.objc_msgSend_id_bool(NSApp,
				_postEventAtStart,
				C.makeDummyEvent(),
				C.BOOL(C.NO))			// not at start, just in case there are other events pending (TODO is this correct?)
		}
	}()

	C.objc_msgSend_noargs(NSApp, _run)
	return nil
}

// TODO move to init_darwin.go?

var (
	_NSApplication = objc_getClass("NSApplication")

	_sharedApplication = sel_getUid("sharedApplication")
	_setActivationPolicy = sel_getUid("setActivationPolicy:")
	_activateIgnoringOtherApps = sel_getUid("activateIgnoringOtherApps:")
	// _setDelegate in sysdata_darwin.go
)

func initCocoa() (NSApp C.id, err error) {
	C.initBleh()		// initialize bleh_darwin.m functions
	NSApp = C.objc_msgSend_noargs(_NSApplication, _sharedApplication)
	r := C.objc_msgSend_int(NSApp, _setActivationPolicy,
		0)			// NSApplicationActivationPolicyRegular
	if C.BOOL(uintptr(unsafe.Pointer(r))) != C.BOOL(C.YES) {
		err = fmt.Errorf("error setting NSApplication activation policy (basically identifies our program as a separate program; needed for several things, such as Dock icon, application menu, window resizing, etc.) (unknown reason)")
		return
	}
	C.objc_msgSend_bool(NSApp, _activateIgnoringOtherApps, C.BOOL(C.YES))		// TODO actually do C.NO here? Russ Cox does YES in his devdraw; the docs say the Finder does NO
	err = mkAppDelegate()
	if err != nil {
		return
	}
	C.objc_msgSend_id(NSApp, _setDelegate, appDelegate)
	err = mkAreaClass()
	return
}

//export appDelegate_uitask
func appDelegate_uitask(self C.id, sel C.SEL, arg C.id) {
	p := C.objc_msgSend_noargs(arg, _pointerValue)
	f := (*func ())(unsafe.Pointer(p))
	(*f)()
}
