// 6 july 2014

package ui

import (
	"runtime"
	"sync"
)

// Go initializes package ui.
// TODO write this bit
func Go() error {
	runtime.LockOSThread()
	if err := uiinit(); err != nil {
		return err
	}
	uimsgloop()
	return nil
}

// Do performs f on the main loop, as if it were an event handler.
// It waits for f to execute before returning.
// Do cannot be called within event handlers or within Do itself.
func Do(f func()) {
	done := make(chan struct{})
	defer close(done)
	issue(func() {
		f()
		done <- struct{}{}
	})
	<-done
}

// Stop informs package ui that it should stop.
// Stop then returns immediately.
// Some time after this request is received, Go() will return without performing any final cleanup.
// Stop will not have an effect until any event handlers or dialog boxes presently active return.
// (TODO make sure this is the case for dialog boxes)
func Stop() {
	issue(uistop)
}

type event struct {
	// All events internally return bool; those that don't will be wrapped around to return a dummy value.
	do		func(c Doer) bool
	lock		sync.Mutex
}

// do should never be nil; TODO should we make setters panic instead?

func newEvent() *event {
	return &event{
		do:	func(c Doer) bool {
			return false
		},
	}
}

func (e *event) set(f func(Doer)) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if f == nil {
		f = func(c Doer) {}
	}
	e.do = func(c Doer) bool {
		f(c)
		return false
	}
}

func (e *event) setbool(f func(Doer) bool) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if f == nil {
		f = func(c Doer) bool {
			return false
		}
	}
	e.do = f
}

// This is the common code for running an event.
// It runs on the main thread without a message pump; it provides its own.
func (e *event) fire() bool {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.do(c)
}

// Common code for performing a requested action (ui.Do() or ui.Stop()).
// This should run on the main thread.
// Implementations of issue() should call this.
func perform(fp unsafe.Pointer) {
	f := (*func())(fp)
	(*f)()
}
