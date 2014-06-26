// 26 june 2014
package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	. "github.com/andlabs/ui"
)

// spacing test

type solidColor struct {
	c	color.Color
}
func (s solidColor) Paint(r image.Rectangle) *image.RGBA {
	i := image.NewRGBA(r)
	draw.Draw(i, r, &image.Uniform{s.c}, image.ZP, draw.Src)
	return i
}
func (s solidColor) Mouse(m MouseEvent) bool { return false }
func (s solidColor) Key(e KeyEvent) bool { return false }

var spacetest = flag.String("spacetest", "", "test space idempotency; arg is x or y; overrides -area")
func spaceTest() {
	w := 100
	h := 50
	ng := 1
	gsx, gsy := 1, 0
	f := NewVerticalStack
	if *spacetest == "x" {
		w = 50
		h = 100
		ng = 2
		gsx, gsy = 0, 1
		f = NewHorizontalStack
	}
	ah := solidColor{color.NRGBA{0,0,255,255}}
	a1 := NewArea(w, h, ah)
	a2 := NewArea(w, h, ah)
	a3 := NewArea(w, h, ah)
	a4 := NewArea(w, h, ah)
	win := NewWindow("Stack", 250, 250)
	win.SetSpaced(true)
	win.Open(f(a1, a2))
	win = NewWindow("Grid", 250, 250)
	win.SetSpaced(true)
	g := NewGrid(ng, a3, a4)
	g.SetFilling(0, 0)
	g.SetStretchy(gsx, gsy)
	win.Open(g)
	select {}
}
