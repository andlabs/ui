// 11 february 2014
package ui

import (
	"fmt"
	"testing"
)

func gridWindow() (*Window, error) {
	w := NewWindow("Grid Test", 400, 400)
	b00 := NewButton("0,0")
	b01 := NewButton("0,1")
	b02 := NewButton("0,2")
	l11 := NewListbox(true, "1,1")
	b12 := NewButton("1,2")
	l20 := NewLabel("2,0")
	c21 := NewCheckbox("2,1")
	l22 := NewLabel("2,2")
	g := NewGrid(3,
		b00, b01, b02,
		Space(), l11, b12,
		l20, c21, l22)
	g.SetFilling(1, 2)
	g.SetStretchy(1, 1)
	return w, w.Open(g)
}

func TestMain(t *testing.T) {
	w := NewWindow("Main Window", 320, 240)
	w.Closing = Event()
	b := NewButton("Click Me")
	b2 := NewButton("Or Me")
	s2 := NewStack(Horizontal, b, b2)
	c := NewCheckbox("Check Me")
	cb1 := NewEditableCombobox("You can edit me!", "Yes you can!", "Yes you will!")
	cb2 := NewCombobox("You can't edit me!", "No you can't!", "No you won't!")
	e := NewLineEdit("Enter text here too")
	l := NewLabel("This is a label")
	b3 := NewButton("List Info")
	s3 := NewStack(Horizontal, l, b3)
	s3.SetStretchy(0)
//	s3.SetStretchy(1)
	pbar := NewProgressBar()
	prog := 0
	incButton := NewButton("Inc")
	decButton := NewButton("Dec")
	sincdec := NewStack(Horizontal, incButton, decButton)
	password := NewPasswordEdit()
	s0 := NewStack(Vertical, s2, c, cb1, cb2, e, s3, pbar, sincdec, Space(), password)
	s0.SetStretchy(8)
	lb1 := NewListbox(true, "Select One", "Or More", "To Continue")
	lb2 := NewListbox(false, "Select", "Only", "One", "Please")
	i := 0
	doAdjustments := func() {
		cb1.Append("append")
		cb2.InsertBefore(fmt.Sprintf("before %d", i), 1)
		lb1.InsertBefore(fmt.Sprintf("%d", i), 2)
		lb2.Append("Please")
		i++
	}
	doAdjustments()
	s1 := NewStack(Vertical, lb2, lb1)
	s1.SetStretchy(0)
	s1.SetStretchy(1)
	s := NewStack(Horizontal, s1, s0)
	s.SetStretchy(0)
	s.SetStretchy(1)
	err := w.Open(s)
	if err != nil {
		panic(err)
	}
	gw, err := gridWindow()
	if err != nil {
		panic(err)
	}

mainloop:
	for {
		select {
		case <-w.Closing:
			break mainloop
		case <-b.Clicked:
			err = w.SetTitle(fmt.Sprintf("%v | %s | %s | %s | %s",
				c.Checked(),
				cb1.Selection(),
				cb2.Selection(),
				e.Text(),
				password.Text()))
			if err != nil {
				panic(err)
			}
			doAdjustments()
		case <-b2.Clicked:
			cb1.Delete(1)
			cb2.Delete(2)
			lb1.Delete(3)
			lb2.Delete(4)
		case <-b3.Clicked:
			MsgBox("List Info",
				"cb1: %d %q\ncb2: %d %q\nlb1: %d %q\nlb2: %d %q",
				cb1.SelectedIndex(), cb1.Selection(),
				cb2.SelectedIndex(), cb2.Selection(),
				lb1.SelectedIndices(), lb1.Selection(),
				lb2.SelectedIndices(), lb2.Selection())
		case <-incButton.Clicked:
			prog++
			if prog > 100 {
				prog = 100
			}
			pbar.SetProgress(prog)
		case <-decButton.Clicked:
			prog--
			if prog < 0 {
				prog = 0
			}
			pbar.SetProgress(prog)
		}
	}
	gw.Hide()
	w.Hide()
}
