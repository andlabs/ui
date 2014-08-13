// 8 july 2014

package ui

// This file is called zz_test.go to keep it separate from the other files in this package (and because go test won't accept just test.go)

import (
	"fmt"
	"flag"
	"reflect"
	"testing"
	"image"
	"image/color"
	"image/draw"
)

var closeOnClick = flag.Bool("close", false, "close on click")
var smallWindow = flag.Bool("small", false, "open a small window (test Mac OS X initial control sizing)")
var defocuses = flag.Bool("defocuses", false, "if the Area in the small window (see -small) should defocus")

type dtype struct {
	Name	string
	Address	string
}
var ddata = []dtype{
	{ "alpha", "beta" },
	{ "gamma", "delta" },
	{ "epsilon", "zeta" },
	{ "eta", "theta" },
	{ "iota", "kappa" },
}

type testwin struct {
	t		Tab
	w		Window
	nt		Tab
	a		Area
	spw		Stack
	sph		Stack
	s		Stack
	l		Label
	table		Table
	b		Button
	c		Checkbox
	e		TextField
	e2		TextField

	wsmall	Window
}

type areaHandler struct{defocuses bool}
func (a *areaHandler) Paint(r image.Rectangle) *image.RGBA {
	i := image.NewRGBA(r)
	draw.Draw(i, r, &image.Uniform{color.RGBA{128,0,128,255}}, image.ZP, draw.Src)
	return i
}
func (a *areaHandler) Mouse(me MouseEvent) { fmt.Printf("%#v\n", me) }
func (a *areaHandler) Key(ke KeyEvent) { fmt.Printf("%#v %q\n", ke, ke.Key) }
func (a *areaHandler) Defocuses() bool { return a.defocuses }

func (tw *testwin) make(done chan struct{}) {
	tw.t = NewTab()
	tw.w = NewWindow("Hello", 320, 240, tw.t)
	tw.w.OnClosing(func() bool {
		if *closeOnClick {
			panic("window closed normally in close on click mode (should not happen)")
		}
		println("window close event received")
		Stop()
		done <- struct{}{}
		return true
	})
	tw.t.Append("Blank Tab", NewTab())
	tw.nt = NewTab()
	tw.nt.Append("Tab 1", Space())
	tw.nt.Append("Tab 2", Space())
	tw.t.Append("Tab", tw.nt)
	tw.t.Append("Space", Space())
	tw.a = NewArea(200, 200, &areaHandler{})
	tw.t.Append("Area", tw.a)
	tw.spw = NewHorizontalStack(
		NewButton("hello"),
		NewCheckbox("hello"),
		NewTextField(),
		NewPasswordField(),
		NewTable(reflect.TypeOf(struct{A,B,C int}{})),
		NewStandaloneLabel("hello"))
	tw.t.Append("Pref Width", tw.spw)
	tw.sph = NewVerticalStack(
		NewButton("hello"),
		NewCheckbox("hello"),
		NewTextField(),
		NewPasswordField(),
		NewTable(reflect.TypeOf(struct{A,B,C int}{})),
		NewStandaloneLabel("hello ÉÀÔ"))
	tw.t.Append("Pref Height", tw.sph)
	stack1 := NewHorizontalStack(NewLabel("Test"), NewTextField())
	stack1.SetStretchy(1)
	stack2 := NewHorizontalStack(NewLabel("ÉÀÔ"), NewTextField())
	stack2.SetStretchy(1)
	stack3 := NewHorizontalStack(NewLabel("Test 2"),
		NewTable(reflect.TypeOf(struct{A,B,C int}{})))
	stack3.SetStretchy(1)
	tw.s = NewVerticalStack(stack1, stack2, stack3)
	tw.s.SetStretchy(2)
	tw.t.Append("Stack", tw.s)
	tw.l = NewStandaloneLabel("hello")
	tw.t.Append("Label", tw.l)
	tw.table = NewTable(reflect.TypeOf(ddata[0]))
	tw.table.Lock()
	dq := tw.table.Data().(*[]dtype)
	*dq = ddata
	tw.table.Unlock()
	tw.t.Append("Table", tw.table)
	tw.b = NewButton("There")
	if *closeOnClick {
		tw.b.SetText("Click to Close")
	}
	// GTK+ TODO: this is causing a resize event to happen afterward?!
	tw.b.OnClicked(func() {
		println("in OnClicked()")
		if *closeOnClick {
			tw.w.Close()
			Stop()
			done <- struct{}{}
		}
	})
	tw.t.Append("Button", tw.b)
	tw.c = NewCheckbox("You Should Now See Me Instead")
	tw.c.OnToggled(func() {
		tw.w.SetTitle(fmt.Sprint(tw.c.Checked()))
	})
	tw.t.Append("Checkbox", tw.c)
	tw.e = NewTextField()
	tw.t.Append("Text Field", tw.e)
	tw.e2 = NewPasswordField()
	tw.t.Append("Password Field", tw.e2)
	tw.w.Show()
	if *smallWindow {
		tw.wsmall = NewWindow("Small", 80, 80,
			NewVerticalStack(
				NewButton("Small"),
				NewButton("Small"),
				NewArea(200, 200, &areaHandler{*defocuses})))
		tw.wsmall.Show()
	}
}

// this must be on the heap thanks to moving stacks
// soon even this won't be enough...
var tw *testwin

// because Cocoa hates being run off the main thread, even if it's run exclusively off the main thread
func init() {
	flag.BoolVar(&spaced, "spaced", false, "enable spacing")
	flag.Parse()
	go func() {
		tw = new(testwin)
		done := make(chan struct{})
		Do(func() { tw.make(done) })
		<-done
	}()
	err := Go()
	if err != nil {
		panic(err)
	}
}

func TestDummy(t *testing.T) {
	// do nothing
}
