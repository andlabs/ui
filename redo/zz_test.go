// 8 july 2014

package ui

// This file is called zz_test.go to keep it separate from the other files in this package (and because go test won't accept just test.go)

import (
	"testing"
)

// because Cocoa hates being run off the main thread, even if it's run exclusively off the main thread
func init() {
	go func() {
		w := GetNewWindow(Do, "Hello", 320, 240)
		done := make(chan struct{})
		Wait(Do, w.OnClosing(func(c Doer) bool {
			Wait(c, Stop())
			done <- struct{}{}
			return true
		}))
		Wait(Do, w.Show())
		<-done
	}()
	err := Go()
	if err != nil {
		panic(err)
	}
}

func TestDummy(t *testing.T) {
	// do nothing
}
