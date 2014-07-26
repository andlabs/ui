// 8 july 2014

package ui

// This file is called zz_test.go to keep it separate from the other files in this package (and because go test won't accept just test.go)

import (
	"fmt"
	"flag"
	"testing"
)

var closeOnClick = flag.Bool("close", false, "close on click")
var spaced = flag.Bool("spaced", false, "enable spacing")

// because Cocoa hates being run off the main thread, even if it's run exclusively off the main thread
func init() {
	flag.Parse()
	go func() {
		done := make(chan struct{})
		Do(func() {
			t := NewTab()
			w := NewWindow("Hello", 320, 240, t)
			// TODO use a method here
			w.(*window).spaced = *spaced
			w.OnClosing(func() bool {
				if *closeOnClick {
					panic("window closed normally in close on click mode (should not happen)")
				}
				println("window close event received")
				Stop()
				done <- struct{}{}
				return true
			})
			b := NewButton("There")
			if *closeOnClick {
				b.SetText("Click to Close")
			}
			// GTK+ TODO: this is causing a resize event to happen afterward?!
			b.OnClicked(func() {
				println("in OnClicked()")
				if *closeOnClick {
					w.Close()
					Stop()
					done <- struct{}{}
				}
			})
			t.Append("Button", b)
			c := NewCheckbox("You Should Now See Me Instead")
			c.OnClicked(func() {
				w.SetTitle(fmt.Sprint(c.Checked()))
			})
			t.Append("Checkbox", c)
			e := NewTextField()
			t.Append("Text Field", e)
			e = NewPasswordField()
			t.Append("Password Field", e)
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
