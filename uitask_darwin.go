// 28 february 2014

package ui

import (
	"fmt"
	"runtime"
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
// #include "delegateuitask_darwin.h"
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

	err := initCocoa()
	if err != nil {
		return err
	}

	// Cocoa must run on the first thread created by the program, so we run our dispatcher on another thread instead
	go func() {
		for f := range uitask {
			C.douitask(appDelegate, unsafe.Pointer(&f))
		}
	}()

	go func() {
		main()
		uitask <- func() {
			C.breakMainLoop()
		}
	}()

	C.cocoaMainLoop()
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

func initCocoa() (err error) {
	C.initBleh()		// initialize bleh_darwin.m functions
	err = mkAppDelegate()
	if err != nil {
		return
	}
	if C.initCocoa(appDelegate) != C.YES {
		err = fmt.Errorf("error setting NSApplication activation policy (basically identifies our program as a separate program; needed for several things, such as Dock icon, application menu, window resizing, etc.) (unknown reason)")
		return
	}
	err = mkAreaClass()
	return
}

//export appDelegate_uitask
func appDelegate_uitask(p unsafe.Pointer) {
	f := (*func ())(unsafe.Pointer(p))
	(*f)()
}
