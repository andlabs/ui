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
		done := make(chan struct{})
		Do(func() {
			w := NewWindow("Hello", 320, 240)
			b := NewButton("There")
			w.SetControl(b)
			if *closeOnClick {
				b.SetText("Click to Close")
			}
			w.OnClosing(func() bool {
				if *closeOnClick {
					panic("window closed normally in close on click mode (should not happen)")
				}
				println("window close event received")
				Stop()
				done <- struct{}{}
				return true
			})
			// GTK+ TODO: this is causing a resize event to happen afterward?!
			b.OnClicked(func() {
				println("in OnClicked()")
				if *closeOnClick {
					w.Close()
					Stop()
					done <- struct{}{}
				}
			})
			w.Show()
		})
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
