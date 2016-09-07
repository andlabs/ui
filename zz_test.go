// 11 december 2015

package ui

import (
	"flag"
	"testing"
)

var (
	nomenus     = flag.Bool("nomenus", false, "No menus")
	startspaced = flag.Bool("startspaced", false, "Start with spacing")
	swaphv      = flag.Bool("swaphv", false, "Swap horizontal and vertical boxes")
)

var mainbox *Box
var mainTab *Tab

func xmain() {
	if !*nomenus {
		// TODO
	}

	w := newWindow("Main Window", 320, 240, true)
	w.OnClosing(func(*Window) bool {
		Quit()
		return true
	})

	OnShouldQuit(func() bool {
		// TODO
		return true
	})

	mainbox = newHorizontalBox()
	w.SetChild(mainbox)

	outerTab := newTab()
	mainbox.Append(outerTab, true)

	mainTab = newTab()
	outerTab.Append("Original", mainTab)

	makePage1(w)
	mainTab.Append("Page 1", page1)

	mainTab.Append("Page 2", makePage2())

	// TODO

	if *startspaced {
		setSpaced(true)
	}

	w.Show()
}

func TestIt(t *testing.T) {
	err := Main(xmain)
	if err != nil {
		t.Fatal(err)
	}
}

var (
	spwindows []*Window
	sptabs    []*Tab
	spgroups  []*Group
	spboxes   []*Box
)

func newWindow(title string, width int, height int, hasMenubar bool) *Window {
	w := NewWindow(title, width, height, hasMenubar)
	spwindows = append(spwindows, w)
	return w
}

func newTab() *Tab {
	t := NewTab()
	sptabs = append(sptabs, t)
	return t
}

func newGroup(title string) *Group {
	g := NewGroup(title)
	spgroups = append(spgroups, g)
	return g
}

func newHorizontalBox() *Box {
	var b *Box

	if *swaphv {
		b = NewVerticalBox()
	} else {
		b = NewHorizontalBox()
	}
	spboxes = append(spboxes, b)
	return b
}

func newVerticalBox() *Box {
	var b *Box

	if *swaphv {
		b = NewHorizontalBox()
	} else {
		b = NewVerticalBox()
	}
	spboxes = append(spboxes, b)
	return b
}

func setSpaced(spaced bool) {
	for _, w := range spwindows {
		w.SetMargined(spaced)
	}
	for _, t := range sptabs {
		for i := 0; i < t.NumPages(); i++ {
			t.SetMargined(i, spaced)
		}
	}
	for _, g := range spgroups {
		g.SetMargined(spaced)
	}
	for _, b := range spboxes {
		b.SetPadded(spaced)
	}
}
