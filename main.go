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
	cb1 := NewCombobox(true, "You can edit me!", "Yes you can!", "Yes you will!")
	cb2 := NewCombobox(false, "You can't edit me!", "No you can't!", "No you won't!")
	s := NewStack(Vertical, b, c, cb1, cb2)
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
			cs1, err := cb1.Selection()
			if err != nil {
				panic(err)
			}
			cs2, err := cb2.Selection()
			if err != nil {
				panic(err)
			}
			err = w.SetTitle(fmt.Sprintf("%v | %s | %s", c.Checked(), cs1, cs2))
			if err != nil {
				panic(err)
			}
		}
	}
	w.Hide()
}

