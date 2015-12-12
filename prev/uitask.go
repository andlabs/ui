// 6 july 2014

package ui

import (
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

// Go initializes and runs package ui.
// It returns a non-nil error if initialization fails.
// Otherwise, it will run the event loop and not return until Stop is called.
// Due to platform-specific issues, it must be called from the main OS thread; in general, do not call Go() from anywhere except main() (including any goroutines).
func Go() error {
	runtime.LockOSThread()
	if err := uiinit(); err != nil {
		return err
	}
	go uiissueloop()
	uimsgloop()
	return nil
}

// To ensure that Do() and Stop() only do things after Go() has been called, this channel accepts the requests to issue. The issuing is done by uiissueloop() below.
// Notice that this is a pointer ot a function. See Do() below for details.
var issuer = make(chan *func())

// Do performs f on the main loop, as if it were an event handler.
// It waits for f to execute before returning.
// Do cannot be called within event handlers or within Do itself.
func Do(f func()) {
	done := make(chan struct{})
	defer close(done)
	// THIS MUST BE A POINTER.
	// Previously, the pointer was constructed within issue().
	// This meant that if the Do() was stalled, the garbage collector came in and reused the pointer value too soon!
	call := func() {
		f()
		done <- struct{}{}
	}
	issuer <- &call
	<-done
}

// Stop informs package ui that it should stop.
// Stop then returns immediately.
// Some time after this request is received, Go() will return without performing any final cleanup.
// Stop will not have an effect until any event handlers return.
func Stop() {
	// can't send this directly across issuer
	go func() {
		Do(uistop)
	}()
}

func uiissueloop() {
	for f := range issuer {
		issue(f)
	}
}

type event struct {
	// All events internally return bool; those that don't will be wrapped around to return a dummy value.
	do   func() bool
	lock sync.Mutex
}

func newEvent() *event {
	return &event{
		do: func() bool {
			return false
		},
	}
}

func (e *event) set(f func()) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if f == nil {
		f = func() {}
	}
	e.do = func() bool {
		f()
		return false
	}
}

func (e *event) setbool(f func() bool) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if f == nil {
		f = func() bool {
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

	return e.do()
}

// Common code for performing a requested action (ui.Do() or ui.Stop()).
// This should run on the main thread.
// Implementations of issue() should call this.
func perform(fp unsafe.Pointer) {
	f := (*func())(fp)
	(*f)()
}

// ForeignEvent wraps a channel in such a way that it can be used safely with package ui.
type ForeignEvent struct {
	c reflect.Value
	e *event
	d interface{}
}

// NewForeignEvent creates a new ForeignEvent with the specified channel.
// It panics if the argument is not a receivable channel.
// The returned ForeignEvent assumes ownership of the channel.
// Each time a value is received on the channel, the returned function is invoked on the main thread.
func NewForeignEvent(channel interface{}, handler func(data interface{})) *ForeignEvent {
	c := reflect.ValueOf(channel)
	t := c.Type()
	if t.Kind() != reflect.Chan || (t.ChanDir()&reflect.RecvDir) == 0 {
		panic("non-channel or non-receivable channel passed to NewForeignEvent()")
	}
	fe := &ForeignEvent{
		c: c,
		e: newEvent(),
	}
	fe.e.set(func() {
		handler(fe.d)
	})
	go fe.do()
	return fe
}

func (fe *ForeignEvent) do() {
	for {
		v, ok := fe.c.Recv()
		if !ok {
			break
		}
		fe.d = v.Interface()
		Do(func() {
			fe.e.fire()
		})
	}
}

// Stop ceases all future invocations of the handler passed to NewForeignEvent() on fe; the values read from the channel are merely discarded.
func (fe *ForeignEvent) Stop() {
	fe.e.set(nil)
}
