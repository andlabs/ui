// 8 july 2014

package ui

// This file is called zz_test.go to keep it separate from the other files in this package (and because go test won't accept just test.go)

import (
	"flag"
	"testing"
)

var closeOnClick = flag.Bool("close", false, "close on click")

// because Cocoa hates being run off the main thread, even if it's run exclusively off the main thread
func init() {
	flag.Parse()
	go func() {
		w := GetNewWindow(Do, "Hello", 320, 240)
		b := GetNewButton(Do, "There")
		Wait(Do, w.SetControl(b))
		if *closeOnClick {
			Wait(Do, b.SetText("Click to Close"))
		}
		done := make(chan struct{})
		Wait(Do, w.OnClosing(func(c Doer) bool {
			if *closeOnClick {
				panic("window closed normally in close on click mode (should not happen)")
			}
			println("window close event received")
			Stop()
			done <- struct{}{}
			return true
		}))
		Wait(Do, b.OnClicked(func(c Doer) {
			println("in OnClicked()")
			if *closeOnClick {
				Wait(c, w.Close())
				Stop()
				done <- struct{}{}
			}
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
