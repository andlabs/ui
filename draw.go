// 13 december 2015

package ui

// #include "ui.h"
import "C"

// Path represents a geometric path in a drawing context.
// This is the basic unit of drawing: all drawing operations consist of
// forming a path, then stroking, filling, or clipping to that path.
// A path is an OS resource; you must explicitly free it when finished.
// Paths consist of multiple figures. Once you have added all the
// figures to a path, you must "end" the path to make it ready to draw
// with.
// TODO rewrite all that
// 
// Or more visually, the lifecycle of a Path is
// 	p := NewPath()
// 	for every figure {
// 		p.NewFigure(...) // or NewFigureWithArc
// 		p.LineTo(...)    // any number of these in any order
// 		p.ArcTo(...)
// 		p.BezierTo(...)
// 		if figure should be closed {
// 			p.CloseFigure()
// 		}
// 	}
// 	p.End()
// 	// ...
// 	dp.Context.Stroke(p, ...) // any number of these in any order
// 	dp.Context.Fill(p, ...)
// 	dp.Context.Clip(p)
// 	// ...
// 	p.Free() // when done with the path
type Path struct {
	p	*C.uiDrawPath
}

// NewPath creates a new Path.
func NewPath() *Path {
	return &Path{
		p:	C.uiDrawNewPath(),
	}
}

// Free destroys a Path. After calling Free the Path cannot be used.
func (p *Path) Free() {
	C.uiDrawFreePath(p.p)
}

// NewFigure starts a new figure in the Path. The current point
// is set to the given point.
func (p *Path) NewFigure(x float64, y float64) {
	C.uiDrawPathNewFigure(p.p, C.double(x), C.double(y))
}

// NewFigureWithArc starts a new figure in the Path and adds an arc
// as the first element of the figure. Unlike ArcTo, NewFigureWithArc
// does not draw an initial line segment. Otherwise, see ArcTo.
func (p *Path) NewFigureWithArc(xCenter float64, yCenter float64, radius float64, startAngle float64, sweep float64, isNegative bool) {
	C.uiDrawPathNewFigureWithArc(p.p,
		C.double(xCenter), C.double(yCenter),
		C.double(radius),
		C.double(startAngle), C.double(sweep),
		frombool(isNegative))
}

// LineTo adds a line to the current figure of the Path starting from
// the current point and ending at the given point. The current point
// is set to the ending point.
func (p *Path) LineTo(x float64, y float64) {
	C.uiDrawPathLineTo(p.p, C.double(x), C.double(y))
}

// ArcTo adds a circular arc to the current figure of the Path.
// You pass it the center of the arc, its radius in radians, the starting
// angle (couterclockwise) in radians, and the number of radians the
// arc should sweep (counterclockwise). A line segment is drawn from
// the current point to the start of the arc. The current point is set to
// the end of the arc.
func (p *Path) ArcTo(xCenter float64, yCenter float64, radius float64, startAngle float64, sweep float64, isNegative bool) {
	C.uiDrawPathArcTo(p.p,
		C.double(xCenter), C.double(yCenter),
		C.double(radius),
		C.double(startAngle), C.double(sweep),
		frombool(isNegative))
}

// BezierTo adds a cubic Bezier curve to the current figure of the Path.
// Its start point is the current point. c1x and c1y are the first control
// point. c2x and c2y are the second control point. endX and endY
// are the end point. The current point is set to the end point.
func (p *Path) BezierTo(c1x float64, c1y float64, c2x float64, c2y float64, endX float64, endY float64) {
	C.uiDrawPathBezierTo(p.p,
		C.double(c1x), C.double(c1y),
		C.double(c2x), C.double(c2y),
		C.double(endX), C.double(endY))
}

// CloseFigure draws a line segment from the current point of the
// current figure of the Path back to its initial point. After calling this,
// the current figure is over and you must either start a new figure
// or end the Path. If this is not called and you start a new figure or
// end the Path, then the current figure will not have this closing line
// segment added to it (but the figure will still be over).
func (p *Path) CloseFigure() {
	C.uiDrawPathCloseFigure(p.p)
}

// AddRectangle creates a new figure in the Path that consists entirely
// of a rectangle whose top-left corner is at the given point and whose
// size is the given size. The rectangle is a closed figure; you must
// either start a new figure or end the Path after calling this method.
func (p *Path) AddRectangle(x float64, y float64, width float64, height float64) {
	C.uiDrawPathAddRectangle(p.p, C.double(x), C.double(y), C.double(width), C.double(height))
}

// End ends the current Path. You cannot add figures to a Path that has
// been ended. You cannot draw with a Path that has not been ended.
func (p *Path) End() {
	C.uiDrawPathEnd(p.p)
}

// DrawContext represents a drawing surface that you can draw to.
// At present the only DrawContexts are surfaces associated with
// Areas and are provided by package ui; see AreaDrawParams.
type DrawContext struct {
	c	*C.uiDrawContext
}
