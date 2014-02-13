// 11 february 2014
package main

import (
	"fmt"
)

func main() {
	w := NewWindow("Main Window", 320, 240)
	w.Closing = make(chan struct{})
	b := NewButton("Click Me")
	c := NewCheckbox("Check Me")
	s := NewStack(Vertical)
	s.Controls = []Control{b, c}
	err := w.Open(s)
	if err != nil {
		panic(err)
	}

mainloop:
	for {
		select {
		case <-w.Closing:
			break mainloop
		case <-b.Clicked:
			err := w.SetTitle(fmt.Sprintf("Check State: %v", c.Checked()))
			if err != nil {
				panic(err)
			}
		}
	}
	w.Hide()
}

