// 11 december 2015

package ui

import (
	"runtime"
	"errors"
	"sync"
	"unsafe"
)

// #include "ui.h"
// extern void doQueued(void *);
// extern int doOnShouldQuit(void *);
// /* I forgot how dumb cgo is... ./main.go:73: cannot use _Cgo_ptr(_Cfpvar_fp_doQueued) (type unsafe.Pointer) as type *[0]byte in argument to _Cfunc_uiQueueMain */
// /* I'm pretty sure this worked before... */
// static inline void realQueueMain(void *x)
// {
// 	uiQueueMain(doQueued, x);
// }
// static inline void realOnShouldQuit(void)
// {
// 	uiOnShouldQuit(doOnShouldQuit, NULL);
// }
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
	runtime.LockOSThread()

	ensureMainThread()

	// TODO HEAP SAFETY
	opts := C.uiInitOptions{}
	estr := C.uiInit(&opts)
	if estr != nil {
		errchan <- errors.New(C.GoString(estr))
		C.uiFreeInitError(estr)
		return
	}
	// set up OnShouldQuit()
	C.realOnShouldQuit()
	QueueMain(f)
	C.uiMain()
	errchan <- nil
}

// Quit queues a return from Main. It does not exit the program.
// It also does not immediately cause Main to return; Main will
// return when it next can. Quit must be called from the GUI thread.
func Quit() {
	C.uiQuit()
}

// These prevent the passing of Go functions into C land.
// TODO make an actual sparse list instead of this monotonic map thingy
var (
	qmmap = make(map[uintptr]func())
	qmcurrent = uintptr(0)
	qmlock sync.Mutex
)

// QueueMain queues f to be executed on the GUI thread when
// next possible. It returns immediately; that is, it does not wait
// for the function to actually be executed. QueueMain is the only
// function that can be called from other goroutines, and its
// primary purpose is to allow communication between other
// goroutines and the GUI thread. Calling QueueMain after Quit
// has been called results in undefined behavior.
// 
// If you start a goroutine in f, it also cannot call package ui
// functions. So for instance, the following will result in
// undefined behavior:
// 
// 	ui.QueueMain(func() {
// 		go ui.MsgBox(...)
// 	})
func QueueMain(f func()) {
	qmlock.Lock()
	defer qmlock.Unlock()

	n := uintptr(0)
	for {
		n = qmcurrent
		qmcurrent++
		if qmmap[n] == nil {
			break
		}
	}
	qmmap[n] = f
	C.realQueueMain(unsafe.Pointer(n))
}

//export doQueued
func doQueued(nn unsafe.Pointer) {
	qmlock.Lock()

	n := uintptr(nn)
	f := qmmap[n]
	delete(qmmap, n)

	// allow uiQueueMain() to be called by a queued function
	// TODO explicitly allow this in libui too
	qmlock.Unlock()

	f()
}

// no need to lock this; this API is only safe on the main thread
var shouldQuitFunc func() bool

// OnShouldQuit schedules f to be exeucted when the OS wants
// the program to quit or when a Quit menu item has been clicked.
// Only one function may be registered at a time. If the function
// returns true, Quit will be called. If the function returns false, or
// if OnShouldQuit is never called. Quit will not be called and the
// OS will be told that the program needs to continue running.
func OnShouldQuit(f func() bool) {
	shouldQuitFunc = f
}

//export doOnShouldQuit
func doOnShouldQuit(unused unsafe.Pointer) C.int {
	if shouldQuitFunc == nil {
		return 0
	}
	return frombool(shouldQuitFunc())
}
