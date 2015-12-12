// 30 august 2014

package main

import (
	"image"
	"reflect"
	"github.com/andlabs/ui"
)

var w ui.Window

type areaHandler struct {
	img		*image.RGBA
}

func (a *areaHandler) Paint(rect image.Rectangle) *image.RGBA {
	return a.img.SubImage(rect).(*image.RGBA)
}

func (a *areaHandler) Mouse(me ui.MouseEvent) {}
func (a *areaHandler) Key(ke ui.KeyEvent) bool { return false }

func initGUI() {
	b := ui.NewButton("Button")
	c := ui.NewCheckbox("Checkbox")
	tf := ui.NewTextField()
	tf.SetText("Text Field")
	pf := ui.NewPasswordField()
	pf.SetText("Password Field")
	l := ui.NewLabel("Label")

	t := ui.NewTab()
	t.Append("Tab 1", ui.Space())
	t.Append("Tab 2", ui.Space())
	t.Append("Tab 3", ui.Space())

	g := ui.NewGroup("Group", ui.Space())

	icons := readIcons()
	table := ui.NewTable(reflect.TypeOf(icons[0]))
	table.Lock()
	d := table.Data().(*[]icon)
	*d = icons
	table.Unlock()

	area := ui.NewArea(200, 200, &areaHandler{tileImage(20)})

	stack := ui.NewVerticalStack(
		b,
		c,
		tf,
		pf,
		l,
		t,
		g,
		table,
		area)
	stack.SetStretchy(5)
	stack.SetStretchy(6)
	stack.SetStretchy(7)
	stack.SetStretchy(8)

	w = ui.NewWindow("Window", 400, 500, stack)
	w.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	w.Show()
}

func main() {
	go ui.Do(initGUI)
	err := ui.Go()
	if err != nil {
		panic(err)
	}
}
