// 8 july 2014

package ui

// This file is called zz_test.go to keep it separate from the other files in this package (and because go test won't accept just test.go)

import (
	"testing"
)

func TestPackage(t *testing.T) {
	go func() {
		w := GetNewWindow(Do, "Hello", 320, 240)
		done := make(chan struct{})
//		Wait(Do, w.OnClosing(func(Doer) bool {
//			done <- struct{}{}
//			return true
//		}))
		Wait(Do, w.Show())
		<-done
	}()
	err := Go()
	if err != nil {
		t.Error(err)
	}
}
