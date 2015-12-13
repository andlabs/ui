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
		b := NewVerticalBox()
		w.SetChild(b)
		w.SetMargined(true)
		button := NewButton("Click Me")
		button.OnClicked(func(*Button) {
			b.SetPadded(!b.Padded())
		})
		b.Append(button, true)
		b.Append(NewButton("Button 2"), false)
		b.Append(NewButton("Button 3"), true)
		w.Show()
	})
	if err != nil {
		t.Fatal(err)
	}
}
