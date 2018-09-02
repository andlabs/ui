// 12 august 2018

// +build OMIT

package main

import (
	"math/rand"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var (
	histogram *ui.Area
	datapoints [10]*ui.Spinbox
	colorButton *ui.ColorButton

	currentPoint = -1
)

// some metrics
const (
	xoffLeft = 20		// histogram margins
	yoffTop = 20
	xoffRight = 20
	yoffBottom = 20
	pointRadius = 5
)

// helper to quickly set a brush color
func mkSolidBrush(color uint32, alpha float64) *ui.DrawBrush {
	brush := new(ui.DrawBrush)
	brush.Type = ui.DrawBrushTypeSolid
	component := uint8((color >> 16) & 0xFF)
	brush.R = float64(component) / 255
	component = uint8((color >> 8) & 0xFF)
	brush.G = float64(component) / 255
	component = uint8(color & 0xFF)
	brush.B = float64(component) / 255
	brush.A = alpha
	return brush
}

// and some colors
// names and values from https://msdn.microsoft.com/en-us/library/windows/desktop/dd370907%28v=vs.85%29.aspx
const (
	colorWhite = 0xFFFFFF
	colorBlack = 0x000000
	colorDodgerBlue = 0x1E90FF
)

func pointLocations(width, height float64) (xs, ys [10]float64) {
	xincr := width / 9		// 10 - 1 to make the last point be at the end
	yincr := height / 100
	for i := 0; i < 10; i++ {
		// get the value of the point
		n := datapoints[i].Value()
		// because y=0 is the top but n=0 is the bottom, we need to flip
		n = 100 - n
		xs[i] = xincr * float64(i)
		ys[i] = yincr * float64(n)
	}
	return xs, ys
}

func constructGraph(width, height float64, extend bool) *ui.DrawPath {
	xs, ys := pointLocations(width, height)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)

	path.NewFigure(xs[0], ys[0])
	for i := 1; i < 10; i++ {
		path.LineTo(xs[i], ys[i])
	}

	if extend {
		path.LineTo(width, height)
		path.LineTo(0, height)
		path.CloseFigure()
	}

	path.End()
	return path
}

func graphSize(clientWidth, clientHeight float64) (graphWidth, graphHeight float64) {
	return clientWidth - xoffLeft - xoffRight,
		clientHeight - yoffTop - yoffBottom
}

type areaHandler struct{}

func (areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	// fill the area with white
	brush := mkSolidBrush(colorWhite, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(0, 0, p.AreaWidth, p.AreaHeight)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()

	graphWidth, graphHeight := graphSize(p.AreaWidth, p.AreaHeight)

	sp := &ui.DrawStrokeParams{
		Cap:			ui.DrawLineCapFlat,
		Join:			ui.DrawLineJoinMiter,
		Thickness:	2,
		MiterLimit:	ui.DrawDefaultMiterLimit,
	}

	// draw the axes
	brush = mkSolidBrush(colorBlack, 1.0)
	path = ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(xoffLeft, yoffTop)
	path.LineTo(xoffLeft, yoffTop + graphHeight)
	path.LineTo(xoffLeft + graphWidth, yoffTop + graphHeight)
	path.End()
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now transform the coordinate space so (0, 0) is the top-left corner of the graph
	m := ui.DrawNewMatrix()
	m.Translate(xoffLeft, yoffTop)
	p.Context.Transform(m)

	// now get the color for the graph itself and set up the brush
	graphR, graphG, graphB, graphA := colorButton.Color()
	brush.Type = ui.DrawBrushTypeSolid
	brush.R = graphR
	brush.G = graphG
	brush.B = graphB
	// we set brush.A below to different values for the fill and stroke

	// now create the fill for the graph below the graph line
	path = constructGraph(graphWidth, graphHeight, true)
	brush.A = graphA / 2
	p.Context.Fill(path, brush)
	path.Free()

	// now draw the histogram line
	path = constructGraph(graphWidth, graphHeight, false)
	brush.A = graphA
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now draw the point being hovered over
	if currentPoint != -1 {
		xs, ys := pointLocations(graphWidth, graphHeight)
		path = ui.DrawNewPath(ui.DrawFillModeWinding)
		path.NewFigureWithArc(
			xs[currentPoint], ys[currentPoint],
			pointRadius,
			0, 6.23,		// TODO pi
			false)
		path.End()
		// use the same brush as for the histogram lines
		p.Context.Fill(path, brush)
		path.Free()
	}
}

func inPoint(x, y float64, xtest, ytest float64) bool {
	// TODO switch to using a matrix
	x -= xoffLeft
	y -= yoffTop
	return (x >= xtest - pointRadius) &&
		(x <= xtest + pointRadius) &&
		(y >= ytest - pointRadius) &&
		(y <= ytest + pointRadius)
}

func (areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	graphWidth, graphHeight := graphSize(me.AreaWidth, me.AreaHeight)
	xs, ys := pointLocations(graphWidth, graphHeight)

	currentPoint = -1
	for i := 0; i < 10; i++ {
		if inPoint(me.X, me.Y, xs[i], ys[i]) {
			currentPoint = i
			break
		}
	}

	// TODO only redraw the relevant area
	histogram.QueueRedrawAll()
}

func (areaHandler) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (areaHandler) DragBroken(a *ui.Area) {
	// do nothing
}

func (areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}

func setupUI() {
	mainwin := ui.NewWindow("libui Histogram Example", 640, 480, true)
	mainwin.SetMargined(true)
	mainwin.OnClosing(func(*ui.Window) bool {
		mainwin.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	mainwin.SetChild(hbox)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	histogram = ui.NewArea(areaHandler{})

	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		datapoints[i] = ui.NewSpinbox(0, 100)
		datapoints[i].SetValue(rand.Intn(101))
		datapoints[i].OnChanged(func(*ui.Spinbox) {
			histogram.QueueRedrawAll()
		})
		vbox.Append(datapoints[i], false)
	}

	colorButton = ui.NewColorButton()
	// TODO inline these
	brush := mkSolidBrush(colorDodgerBlue, 1.0)
	colorButton.SetColor(brush.R,
		brush.G,
		brush.B,
		brush.A)
	colorButton.OnChanged(func(*ui.ColorButton) {
		histogram.QueueRedrawAll()
	})
	vbox.Append(colorButton, false)

	hbox.Append(histogram, true)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
