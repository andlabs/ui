// 6 july 2014

package ui

import (
	"runtime"
)

// Go initializes package ui.
// TODO write this bit
func Go() error {
	runtime.LockOSThread()
	if err := uiinit(); err != nil {
		return err
	}
	go uitask()
	uimsgloop()
	return nil
}

// TODO Stop

// This is the ui main loop.
// It is spawned by Go as a goroutine.
func uitask() {
	for {
		select {
		case req := <-Do:
			// TODO foreign event
			issue(req)
		case <-stall:		// wait for event to finish
			<-stall		// see below for information
		}
	}
}

// At each event, this is pulsed twice: once when the event begins, and once when the event ends.
// Do is not processed in between.
var stall = make(chan struct{})

// This is the common code for running an event.
// It runs on the main thread without a message pump; it provides its own.
// TODO
// - define event
// - figure out how to return values from event handlers
func doevent(e event) {
	stall <- struct{}{}		// enter event handler
	c := make(Doer)
	go func() {
		e.do(c)
		close(c)
	}()
	for req := range c {
		// note: this is perform, not issue!
		// doevent runs on the main thread without a message pump!
		perform(req)
	}
	// leave the event handler; leave it only after returning from an event handler so we must issue it like a normal Request
	issue(&Request{
		op:		func() {
			stall <- struct{}{}
		},
		// unfortunately, closing a nil channel causes a panic
		// therefore, we have to make a dummy channel
		// TODO add conditional checks to the request handler instead?
		resp:		make(chan interface{}),
	})
}

// Common code for performing a Request.
// This should run on the main thread.
// Implementations of issue() should call this.
func perform(req *Request) {
	req.op()
	close(req.resp)
}
