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

func initCocoa() (err error) {
	makeAppDelegate()
	if C.initCocoa(appDelegate) != C.YES {
		return fmt.Errorf("error setting NSApplication activation policy (basically identifies our program as a separate program; needed for several things, such as Dock icon, application menu, window resizing, etc.) (unknown reason)")
	}
	return nil
}

//export appDelegate_uitask
func appDelegate_uitask(p unsafe.Pointer) {
	f := (*func ())(unsafe.Pointer(p))
	(*f)()
}
