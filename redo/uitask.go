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
	go uitask(Do)
	uimsgloop()
	return nil
}

// Stop issues a Request for package ui to stop and returns immediately.
// Some time after this request is received, Go() will return without performing any final cleanup.
// Stop is internally issued to Do, so it will not be registered until any event handlers and dialog boxes return.
func Stop() {
	go func() {
		c := make(chan interface{})
		Do <- &Request{
			op:		func() {
				uistop()
				c <- struct{}{}
			},
			resp:		c,
		}
		<-c
	}()
}

// This is the ui main loop.
// It is spawned by Go as a goroutine.
// It can also be called recursively using the recur/unrecur chain.
func uitask(doer Doer) {
	for {
		select {
		case req := <-doer:
			// TODO foreign event
			issue(req)
		case rec := <-recur:		// want to perform event
			c := make(Doer)
			rec <- c
			uitask(c)
		case <-unrecur:		// finished with event
			close(doer)
			return
		}
	}
}

// Send a channel over recur to have uitask() above enter a recursive loop in which the Doer sent back becomes the active request handler.
// Pulse unrecur when finished.
var recur = make(chan chan Doer)
var unrecur = make(chan struct{})

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
	cc := make(chan Doer)
	recur <- cc
	c := <-cc
	result := false
	finished := make(chan struct{})
	go func() {
		e.lock.Lock()
		defer e.lock.Unlock()

		result = e.do(c)
		finished <- struct{}{}
	}()
	<-finished
	close(finished)
	// leave the event handler; leave it only after returning from the OS-side event handler so we must issue it like a normal Request
	issue(&Request{
		op:		func() {
			unrecur <- struct{}{}
		},
		// unfortunately, closing a nil channel causes a panic
		// therefore, we have to make a dummy channel
		// TODO add conditional checks to the request handler instead?
		resp:		make(chan interface{}),
	})
	return result
}

// Common code for performing a Request.
// This should run on the main thread.
// Implementations of issue() should call this.
func perform(req *Request) {
	req.op()
	close(req.resp)
}
