// 11 february 2014
package main

import (
	"fmt"
	"flag"
	. "github.com/andlabs/ui"
)

var prefsizetest = flag.Bool("prefsize", false, "")
func listboxPreferredSizeTest() (*Window, error) {
	lb := NewListbox("xxxxx", "y", "zzz")
	g := NewGrid(1, lb)
	w := NewWindow("Listbox Preferred Size Test", 300, 300)
	return w, w.Open(g)
}

var gridtest = flag.Bool("grid", false, "")
func gridWindow() (*Window, error) {
	w := NewWindow("Grid Test", 400, 400)
	b00 := NewButton("0,0")
	b01 := NewButton("0,1")
	b02 := NewButton("0,2")
	l11 := NewMultiSelListbox("1,1")
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

var macCrashTest = flag.Bool("maccrash", false, "attempt crash on Mac OS X on deleting too far (debug lack of panic on 32-bit)")

func invalidTest(c *Combobox, l *Listbox, s *Stack, g *Grid) {
	x := func(what string ) {
		if j := recover(); j == nil {
			MsgBoxError("test", "%s: no panic", what)
			panic("invalid test fail")
		} else {
			println("got", j.(error).Error())
		}
	}

	func() {
		defer x("Combobox.InsertBefore < 0"); c.InsertBefore("xxx", -5); panic(nil)
	}()
	func() {
		defer x("Combobox.InsertBefore > len"); c.InsertBefore("xxx", c.Len() + 5); panic(nil)
	}()
	func() {
		defer x("Combobox.Delete < 0"); c.Delete(-5); panic(nil)
	}()
	func() {
		defer x("Combobox.Delete > len"); c.Delete(c.Len() + 5); panic(nil)
	}()
	func() {
		defer x("Listbox.InsertBefore < 0"); l.InsertBefore("xxx", -5); panic(nil)
	}()
	func() {
		defer x("Listbox.InsertBefore > len"); l.InsertBefore("xxx", l.Len() + 5); panic(nil)
	}()
	func() {
		defer x("Listbox.Delete < 0"); l.Delete(-5); panic(nil)
	}()
	func() {
		defer x("Listbox.Delete > len"); l.Delete(c.Len() + 5); panic(nil)
	}()
	if s != nil {
		func() {
			defer x("Stack.SetStretchy < 0"); s.SetStretchy(-5); panic(nil)
		}()
		func() {
			defer x("Stack.SetStretchy > len"); s.SetStretchy(5555); panic(nil)
		}()
	}
	if g != nil {
		func() {
			defer x("Grid.SetFilling x < 0"); g.SetFilling(-5, 0); panic(nil)
		}()
		func() {
			defer x("Grid.SetFilling x > len"); g.SetFilling(5555, 0); panic(nil)
		}()
		func() {
			defer x("Grid.SetFilling y < 0"); g.SetFilling(0, -5); panic(nil)
		}()
		func() {
			defer x("Grid.SetFilling y > len"); g.SetFilling(0, 5555); panic(nil)
		}()
		func() {
			defer x("Grid.SetStretchy x < 0"); g.SetStretchy(-5, 0); panic(nil)
		}()
		func() {
			defer x("Grid.SetStretchy x > len"); g.SetStretchy(5555, 0); panic(nil)
		}()
		func() {
			defer x("Grid.SetStretchy y < 0"); g.SetStretchy(0, -5); panic(nil)
		}()
		func() {
			defer x("Grid.SetStretchy y > len"); g.SetStretchy(0, 5555); panic(nil)
		}()
	}
	MsgBox("test", "all working as intended")
}

var invalidBefore = flag.Bool("invalid", false, "run invalid test before opening window")

func myMain() {
	w := NewWindow("Main Window", 320, 240)
	w.Closing = Event()
	b := NewButton("Click Me")
	b2 := NewButton("Or Me")
	s2 := NewHorizontalStack(b, b2)
	c := NewCheckbox("Check Me")
	cb1 := NewEditableCombobox("You can edit me!", "Yes you can!", "Yes you will!")
	cb2 := NewCombobox("You can't edit me!", "No you can't!", "No you won't!")
	e := NewLineEdit("Enter text here too")
	l := NewLabel("This is a label")
	b3 := NewButton("List Info")
	s3 := NewHorizontalStack(l, b3)
	s3.SetStretchy(0)
//	s3.SetStretchy(1)
	pbar := NewProgressBar()
	prog := 0
	incButton := NewButton("Inc")
	decButton := NewButton("Dec")
	invalidButton := NewButton("Run Invalid Test")
	sincdec := NewHorizontalStack(incButton, decButton, invalidButton)
	password := NewPasswordEdit()
	s0 := NewVerticalStack(s2, c, cb1, cb2, e, s3, pbar, sincdec, Space(), password)
	s0.SetStretchy(8)
	lb1 := NewMultiSelListbox("Select One", "Or More", "To Continue")
	lb2 := NewListbox("Select", "Only", "One", "Please")
	i := 0
	doAdjustments := func() {
		cb1.Append("append")
		cb2.InsertBefore(fmt.Sprintf("before %d", i), 1)
		lb1.InsertBefore(fmt.Sprintf("%d", i), 2)
		lb2.Append("Please")
		i++
	}
	doAdjustments()
	cb1.Append("append multi 1", "append multi 2")
	lb2.Append("append multi 1", "append multi 2")
	s1 := NewVerticalStack(lb2, lb1)
	s1.SetStretchy(0)
	s1.SetStretchy(1)
	s := NewHorizontalStack(s1, s0)
	s.SetStretchy(0)
	s.SetStretchy(1)
	if *invalidBefore {
		invalidTest(cb1, lb1, s, NewGrid(1, Space()))
	}
	err := w.Open(s)
	if err != nil {
		panic(err)
	}
	if *gridtest {
		_, err := gridWindow()
		if err != nil {
			panic(err)
		}
	}
	if *prefsizetest {
		_, err = listboxPreferredSizeTest()
		if err != nil {
			panic(err)
		}
	}

mainloop:
	for {
		select {
		case <-w.Closing:
			break mainloop
		case <-b.Clicked:
			w.SetTitle(fmt.Sprintf("%v | %s | %s | %s | %s",
				c.Checked(),
				cb1.Selection(),
				cb2.Selection(),
				e.Text(),
				password.Text()))
			doAdjustments()
		case <-b2.Clicked:
			if cb1.Len() > 1 {
				cb1.Delete(1)
			}
			if cb2.Len() > 2 {
				cb2.Delete(2)
			}
			if lb1.Len() > 3 || *macCrashTest {
				lb1.Delete(3)
			}
			if lb2.Len() > 4 {
				lb2.Delete(4)
			}
		case <-b3.Clicked:
			MsgBox("List Info",
				"cb1: %d %q (len %d)\ncb2: %d %q (len %d)\nlb1: %d %q (len %d)\nlb2: %d %q (len %d)",
				cb1.SelectedIndex(), cb1.Selection(), cb1.Len(),
				cb2.SelectedIndex(), cb2.Selection(), cb2.Len(),
				lb1.SelectedIndices(), lb1.Selection(), lb1.Len(),
				lb2.SelectedIndices(), lb2.Selection(), lb2.Len())
		case <-incButton.Clicked:
			prog++
			if prog > 100 {
				prog = 100
			}
			pbar.SetProgress(prog)
			cb1.Append("append multi 1", "append multi 2")
			lb2.Append("append multi 1", "append multi 2")
		case <-decButton.Clicked:
			prog--
			if prog < 0 {
				prog = 0
			}
			pbar.SetProgress(prog)
		case <-invalidButton.Clicked:
			invalidTest(cb1, lb1, nil, nil)
		}
	}
	w.Hide()
}

func main() {
	flag.Parse()
	err := Go(myMain)
	if err != nil {
		panic(err)
	}
}
