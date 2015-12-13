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
		t := NewTab()
		w.SetChild(t)
		w.SetMargined(true)
		t.Append("First Page", NewButton("Click Me"))
		t.Append("Second Page", NewButton("Click Me Too"))
		t.SetMargined(0, true)
		w.Show()
	})
	if err != nil {
		t.Fatal(err)
	}
}
