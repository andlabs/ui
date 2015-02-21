// 8 july 2014

package ui

// This file is called zz_test.go to keep it separate from the other files in this package (and because go test won't accept just test.go)

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"strings"
	"testing"
	"time"
)

var closeOnClick = flag.Bool("close", false, "close on click")
var smallWindow = flag.Bool("small", false, "open a small window (test Mac OS X initial control sizing)")
var spaced = flag.Bool("spaced", false, "enable spacing")

func newHorizontalStack(c ...Control) Stack {
	s := NewHorizontalStack(c...)
	s.SetPadded(*spaced)
	return s
}

func newVerticalStack(c ...Control) Stack {
	s := NewVerticalStack(c...)
	s.SetPadded(*spaced)
	return s
}

func newSimpleGrid(n int, c ...Control) SimpleGrid {
	g := NewSimpleGrid(n, c...)
	g.SetPadded(*spaced)
	return g
}

type colnametest struct {
	IncorrectColumnName	string	`uicolumn:"Correct Column Name"`
	IncorrectColumnName2	string	`uicolumn:"正解なコラムネーム"`		// thanks GlitterBerri in irc.badnik.net/#zelda
	AlreadyCorrect			string
}

type dtype struct {
	Name    string
	Address string
}

var ddata = []dtype{
	{"alpha", "beta"},
	{"gamma", "delta"},
	{"epsilon", "zeta"},
	{"eta", "theta"},
	{"iota", "kappa"},
}

type testwin struct {
	t          Tab
	w          Window
	roenter	TextField
	roro		TextField
	repainter  *repainter
	fe         *ForeignEvent
	festack    Stack
	festack2		Stack
	festart    Button
	felabel    Label
	festop     Button
	vedit      TextField
	openbtn    Button
	fnlabel    Label
	icons      []icon
	icontbl    Table
	group2     Group
	group      Group
	simpleGrid SimpleGrid
	nt         Tab
	a          Area
	spw        Stack
	sph        Stack
	s          Stack
	l          Label
	table      Table
	b          Button
	c          Checkbox
	e          TextField
	e2         TextField

	wsmall Window
}

type areaHandler struct {
	handled bool
}

func (a *areaHandler) Paint(r image.Rectangle) *image.RGBA {
	i := image.NewRGBA(r)
	draw.Draw(i, r, &image.Uniform{color.RGBA{128, 0, 128, 255}}, image.ZP, draw.Src)
	return i
}
func (a *areaHandler) Mouse(me MouseEvent)  { fmt.Printf("%#v\n", me) }
func (a *areaHandler) Key(ke KeyEvent) bool { fmt.Printf("%#v %q\n", ke, ke.Key); return a.handled }

func (tw *testwin) openFile(fn string) {
	if fn == "" {
		fn = "<no file selected>"
	}
	tw.fnlabel.SetText(fn)
}

func (tw *testwin) addfe() {
	tw.festart = NewButton("Start")
	tw.festart.OnClicked(func() {
		if tw.fe != nil {
			tw.fe.Stop()
		}
		ticker := time.NewTicker(1 * time.Second)
		tw.fe = NewForeignEvent(ticker.C, func(d interface{}) {
			t := d.(time.Time)
			tw.felabel.SetText(t.String())
		})
	})
	tw.felabel = NewLabel("<stopped>")
	tw.festop = NewButton("Stop")
	tw.festop.OnClicked(func() {
		if tw.fe != nil {
			tw.fe.Stop()
			tw.felabel.SetText("<stopped>")
			tw.fe = nil
		}
	})
	tw.vedit = NewTextField()
	tw.vedit.OnChanged(func() {
		if strings.Contains(tw.vedit.Text(), "bad") {
			tw.vedit.Invalid("bad entered")
		} else {
			tw.vedit.Invalid("")
		}
	})
	tw.openbtn = NewButton("Open")
	tw.openbtn.OnClicked(func() {
		OpenFile(tw.w, tw.openFile)
	})
	tw.fnlabel = NewLabel("<no file selected>")
	tw.festack = newVerticalStack(tw.festart,
		tw.felabel,
		tw.festop,
		NewCheckbox("This is a checkbox test"),
		Space(),
		tw.vedit,
		Space(),
		NewCheckbox("This is a checkbox test"),
		tw.openbtn, tw.fnlabel)
	tw.festack.SetStretchy(4)
	tw.festack.SetStretchy(6)
	sb := NewSpinbox(0, 100)
	sp := NewProgressBar()
	sb.OnChanged(func() {
		sp.SetPercent(sb.Value())
	})
	tw.festack2 = newVerticalStack(sb, sp, Space(), Space(), NewTextbox())
	tw.festack2.SetStretchy(3)
	tw.festack2.SetStretchy(4)
	tw.festack = newHorizontalStack(tw.festack, tw.festack2)
	tw.festack.SetStretchy(0)
	tw.festack.SetStretchy(1)
	tw.t.Append("Foreign Events", tw.festack)
}

func (tw *testwin) make(done chan struct{}) {
	tw.t = NewTab()
	tw.w = NewWindow("Hello", 320, 240, tw.t)
	tw.w.SetMargined(*spaced)
	tw.w.OnClosing(func() bool {
		if *closeOnClick {
			panic("window closed normally in close on click mode (should not happen)")
		}
		println("window close event received")
		Stop()
		done <- struct{}{}
		return true
	})
	tw.roenter = NewTextField()
	tw.roro = NewTextField()
	tw.roro.SetReadOnly(true)
	tw.roenter.OnChanged(func() {
		tw.roro.SetText(tw.roenter.Text())
	})
	s := newVerticalStack(tw.roenter, tw.roro, NewTable(reflect.TypeOf(colnametest{})))
	s.SetStretchy(2)
	tw.t.Append("Read-Only", s)
	tw.icons = readIcons() // repainter uses these
	tw.repainter = newRepainter(15)
	tw.t.Append("Repaint", tw.repainter.grid)
	tw.addfe()
	tw.icontbl = NewTable(reflect.TypeOf(icon{}))
	tw.icontbl.Lock()
	idq := tw.icontbl.Data().(*[]icon)
	*idq = tw.icons
	tw.icontbl.Unlock()
	tw.icontbl.OnSelected(func() {
		s := fmt.Sprintf("%d ", tw.icontbl.Selected())
		tw.icontbl.RLock()
		defer tw.icontbl.RUnlock()
		idq := tw.icontbl.Data().(*[]icon)
		for _, v := range *idq {
			s += strings.ToUpper(fmt.Sprintf("%v", v.Bool)[0:1]) + " "
		}
		tw.w.SetTitle(s)
	})
	tw.t.Append("Image List Table", tw.icontbl)
	tw.group2 = NewGroup("Group", NewButton("Button in Group"))
	tw.t.Append("Empty Group", NewGroup("Group", Space()))
	tw.t.Append("Filled Group", tw.group2)
	tw.group2.SetMargined(*spaced)
	tw.group = NewGroup("Group", newVerticalStack(NewCheckbox("Checkbox in Group")))
	tw.group.SetMargined(*spaced)
	tw.t.Append("Group", tw.group)
	tw.simpleGrid = newSimpleGrid(3,
		NewLabel("0,0"), NewTextField(), NewLabel("0,2"),
		NewButton("1,0"), NewButton("1,1"), NewButton("1,2"),
		NewLabel("2,0"), NewTextField(), NewLabel("2,2"))
	tw.simpleGrid.SetFilling(2, 1)
	tw.simpleGrid.SetFilling(1, 2)
	tw.simpleGrid.SetStretchy(1, 1)
	tw.t.Append("Simple Grid", tw.simpleGrid)
	tw.t.Append("Blank Tab", NewTab())
	tw.nt = NewTab()
	tw.nt.Append("Tab 1", Space())
	tw.nt.Append("Tab 2", Space())
	tw.t.Append("Tab", tw.nt)
	tw.t.Append("Space", Space())
	tw.a = NewArea(200, 200, &areaHandler{false})
	tw.t.Append("Area", tw.a)
	tw.spw = newHorizontalStack(
		NewButton("hello"),
		NewCheckbox("hello"),
		NewTextField(),
		NewPasswordField(),
		NewTable(reflect.TypeOf(struct{ A, B, C int }{})),
		NewLabel("hello"))
	tw.t.Append("Pref Width", tw.spw)
	tw.sph = newVerticalStack(
		NewButton("hello"),
		NewCheckbox("hello"),
		NewTextField(),
		NewPasswordField(),
		NewTable(reflect.TypeOf(struct{ A, B, C int }{})),
		NewLabel("hello ÉÀÔ"))
	tw.t.Append("Pref Height", tw.sph)
	stack1 := newHorizontalStack(NewLabel("Test"), NewTextField())
	stack1.SetStretchy(1)
	stack2 := newHorizontalStack(NewLabel("ÉÀÔ"), NewTextField())
	stack2.SetStretchy(1)
	stack3 := newHorizontalStack(NewLabel("Test 2"),
		NewTable(reflect.TypeOf(struct{ A, B, C int }{})))
	stack3.SetStretchy(1)
	tw.s = newVerticalStack(stack1, stack2, stack3)
	tw.s.SetStretchy(2)
	tw.t.Append("Stack", tw.s)
	tw.l = NewLabel("hello")
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
			newVerticalStack(
				NewButton("Small"),
				NewButton("Small 2"),
				NewArea(200, 200, &areaHandler{true})))
		tw.wsmall.Show()
	}
}

// this must be on the heap thanks to moving stacks
// soon even this won't be enough...
var tw *testwin

// because Cocoa hates being run off the main thread, even if it's run exclusively off the main thread
func init() {
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
