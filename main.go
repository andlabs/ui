// 11 december 2015

package ui

import (
	"errors"
)

// #include "interop.h"
import "C"

// Main initializes package ui, runs f to set up the program,
// and executes the GUI main loop. f should set up the program's
// initial state: open the main window, create controls, and set up
// events. It should then return, at which point Main will
// process events until Quit is called, at which point Main will return
// nil. If package ui fails to initialize, Main returns an appropriate
// error.
func Main(f func()) error {
	errchan := make(chan error)
	go start(errchan, f)
	return <-errchan
}

func start(errchan chan error, f func()) {
	estr := C.interopInit()
	if estr != "" {
		errchan <- errors.New(C.GoString(estr))
		C.interopFreeStr(estr)
		return
	}
	QueueMain(f)
	C.interopRun()
	errchan <- nil
}

// Quit queues an exit from the GUI thread. It does not exit the
// program. Quit must be called from the GUI thread.
func Quit() {
	C.interopQuit()
}

// QueueMain queues f to be executed on the GUI thread when
// next possible. It returns immediately; that is, it does not wait
// for the function to actually be executed. QueueMain is the only
// function that can be called from other goroutines, and its
// primary purpose is to allow communication between other
// goroutines and the GUI thread. Calling QueueMain after Quit
// has been called results in undefined behavior.
func QueueMain(f func()) {
	n := interoperAdd(f)
	C.interopQueueMain(n)
}

//export interopDoQueued
func interopDoQueued(n C.uintptr_t) {
	ff := interoperTake(n)
	f := ff.(func())
	f()
}
