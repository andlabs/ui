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
		s := NewGroup("Group")
		w.SetChild(s)
		w.SetMargined(true)
		b := NewButton("Click Me")
		b.OnClicked(func(*Button) {
			s.SetMargined(!s.Margined())
		})
		s.SetChild(b)
		w.Show()
	})
	if err != nil {
		t.Fatal(err)
	}
}
