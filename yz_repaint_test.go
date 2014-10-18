// 21 august 2014

package ui

import (
	"image"
	"image/color"
	"image/draw"
	"strconv"
)

type repainter struct {
	img     *image.RGBA
	area    Area
	x       TextField
	y       TextField
	width   TextField
	height  TextField
	repaint Button
	all     Button
	grid    Grid

	xv int
	yv int
	wv int
	hv int
}

func newRepainter(times int) *repainter {
	r := new(repainter)
	r.img = tileImage(times)
	r.area = NewArea(r.img.Rect.Dx(), r.img.Rect.Dy(), r)
	r.area.OnTextFieldDismissed(r.tfdone)
	r.x = NewTextField()
	r.x.OnChanged(r.setx)
	r.y = NewTextField()
	r.y.OnChanged(r.sety)
	r.width = NewTextField()
	r.width.OnChanged(r.setwidth)
	r.height = NewTextField()
	r.height.OnChanged(r.setheight)
	r.repaint = NewButton("Rect")
	r.repaint.OnClicked(r.dorect)
	r.all = NewButton("All")
	r.all.OnClicked(r.doall)
	grid := NewGrid()
	grid.Add(r.area, nil, South, true, Fill, true, Fill, 3, 6)
	grid.Add(r.x, nil, East, true, Fill, false, LeftTop, 1, 1)
	grid.Add(r.y, nil, South, true, Fill, false, LeftTop, 1, 1)
	grid.Add(r.width, nil, South, true, Fill, false, LeftTop, 1, 1)
	grid.Add(r.height, nil, South, true, Fill, false, LeftTop, 1, 1)
	grid.Add(r.repaint, nil, South, true, Fill, true, Fill, 1, 1)
	grid.Add(r.all, nil, South, true, Center, false, LeftTop, 1, 1)
	r.grid = grid
	r.grid.SetPadded(*spaced)
	return r
}

func (r *repainter) Paint(rect image.Rectangle) *image.RGBA {
	return r.img.SubImage(rect).(*image.RGBA)
}

func (r *repainter) Mouse(me MouseEvent) {
	if me.Up == 1 {
		r.area.OpenTextFieldAt(me.Pos.X, me.Pos.Y)
	}
}

func (r *repainter) tfdone() {
	println(r.area.TextFieldText())
}

func (r *repainter) Key(ke KeyEvent) bool { return false }

func (r *repainter) setx() {
	i, err := strconv.Atoi(r.x.Text())
	if err != nil {
		r.x.Invalid(err.Error())
		return
	}
	r.x.Invalid("")
	r.xv = i
}

func (r *repainter) sety() {
	i, err := strconv.Atoi(r.y.Text())
	if err != nil {
		r.y.Invalid(err.Error())
		return
	}
	r.y.Invalid("")
	r.yv = i
}

func (r *repainter) setwidth() {
	i, err := strconv.Atoi(r.width.Text())
	if err != nil {
		r.width.Invalid(err.Error())
		return
	}
	r.width.Invalid("")
	r.wv = i
}

func (r *repainter) setheight() {
	i, err := strconv.Atoi(r.height.Text())
	if err != nil {
		r.height.Invalid(err.Error())
		return
	}
	r.height.Invalid("")
	r.hv = i
}

func (r *repainter) alter(rect image.Rectangle, c color.Color) {
	draw.Draw(r.img, rect, &image.Uniform{c}, image.ZP, draw.Over)
}

func (r *repainter) dorect() {
	rect := image.Rect(r.xv, r.yv, r.xv+r.wv, r.yv+r.hv)
	r.alter(rect, color.RGBA{255, 0, 255, 128})
	r.area.Repaint(rect)
}

func (r *repainter) doall() {
	r.alter(r.img.Rect, color.RGBA{255, 255, 0, 128})
	r.area.RepaintAll()
}
