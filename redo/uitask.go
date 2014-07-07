// 6 july 2014

package ui

// TODO Go, Start, Stop

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
		issue(req)
	}
	stall <- struct{}{}		// leave event handler
}
