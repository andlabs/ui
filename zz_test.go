// 11 december 2015

package ui

import "testing"

func TestIt(t *testing.T) {
	err := Main(func() {
		w := NewWindow("Hello", 320, 240, false)
		w.OnClosing(func(w *Window) bool {
			Quit()
			return true
		})
		s := NewDateTimePicker()
		w.SetChild(s)
		w.SetMargined(true)
		w.Show()
	})
	if err != nil {
		t.Fatal(err)
	}
}
