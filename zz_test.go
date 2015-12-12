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
		c := NewCheckbox("Click Me")
		c.OnToggled(func(c *Checkbox) {
			w.SetMargined(c.Checked())
		})
		w.SetChild(c)
		w.Show()
	})
	if err != nil {
		t.Fatal(err)
	}
}
