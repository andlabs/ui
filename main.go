// 11 february 2014
package main

import (
	"fmt"
)

func main() {
	w := NewWindow("Main Window", 320, 240)
	w.Closing = make(chan struct{})
	b := NewButton("Click Me")
	b2 := NewButton("Or Me")
	s2 := NewStack(Horizontal, b, b2)
	c := NewCheckbox("Check Me")
	cb1 := NewCombobox(true, "You can edit me!", "Yes you can!", "Yes you will!")
	cb2 := NewCombobox(false, "You can't edit me!", "No you can't!", "No you won't!")
	e := NewLineEdit("Enter text here too")
	l := NewLabel("This is a label")
	s0 := NewStack(Vertical, s2, c, cb1, cb2, e, l)
	lb := NewListbox(true, "Select One", "Or More", "To Continue")
	lb2 := NewListbox(false, "Select", "Only", "One", "Please")
	i := 0
	doAdjustments := func() {
		cb1.Append("append")
		cb2.InsertBefore(fmt.Sprintf("before %d", i), 1)
		lb.InsertBefore(fmt.Sprintf("%d", i), 2)
		lb2.Append("Please")
		i++
	}
	doAdjustments()
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
			doAdjustments()
		case <-b2.Clicked:
			cb1.Delete(1)
			cb2.Delete(2)
			lb.Delete(3)
			lb2.Delete(4)
		}
	}
	w.Hide()
}

