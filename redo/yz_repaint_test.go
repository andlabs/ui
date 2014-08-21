// 21 august 2014

package ui

import (
	"image"
)

type repainter struct {
	img		*image.RGBA
	area		Area
	x		TextField
	y		TextField
	width	TextField
	height	TextField
	repaint	Button
	all		Button
	stack	Stack
}

func newRepainter(times int) *repainter {
	r := new(repainter)
	r.img = tileImage(times)
	r.area = NewArea(r.img.Rect.Dx(), r.img.Rect.Dy(), r)
	r.x = NewTextField()
	r.y = NewTextField()
	r.width = NewTextField()
	r.height = NewTextField()
	r.repaint = NewButton("Rect")
	r.all = NewButton("All")
	r.stack = NewHorizontalStack(r.x, r.y, r.width, r.height, r.repaint, r.all)
	r.stack.SetStretchy(0)
	r.stack.SetStretchy(1)
	r.stack.SetStretchy(2)
	r.stack.SetStretchy(3)
	r.stack = NewVerticalStack(r.area, r.stack)
	r.stack.SetStretchy(0)
	return r
}

func  (r *repainter) Paint(rect image.Rectangle) *image.RGBA {
	return r.img.SubImage(rect).(*image.RGBA)
}

func (r *repainter) Mouse(me MouseEvent) {}
func (r *repainter) Key(ke KeyEvent) bool { return false }
