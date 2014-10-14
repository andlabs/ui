// 14 october 2014

package ui

import "flag"
import "testing"

var twindow *window

func maketw(done chan struct{}) {
	button := newButton("Greet")
	twindow = newWindow("Hello", 200, 100, button)
	twindow.OnClosing(func() bool {
		Stop()
		return true
	})
	twindow.Show()
}

// because Cocoa hates being run off the main thread, even if it's run exclusively off the main thread
func init() {
	flag.Parse()
	go func() {
		done := make(chan struct{})
		Do(func() { maketw(done) })
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
