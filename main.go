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
	e := NewLineEdit("Enter text here too")
	l := NewLabel("This is a label")
	s0 := NewStack(Vertical, b, c, cb1, cb2, e, l)
	lb := NewListbox(true, "Select One", "Or More", "To Continue")
	lb2 := NewListbox(false, "Select", "Only", "One", "Please")
	s1 := NewStack(Vertical, lb2, lb)
	s := NewStack(Horizontal, s1, s0)
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
			err = w.SetTitle(fmt.Sprintf("%v | %s | %s | %s",
				c.Checked(),
				cb1.Selection(),
				cb2.Selection(),
				e.Text()))
			if err != nil {
				panic(err)
			}
		}
	}
	w.Hide()
}

