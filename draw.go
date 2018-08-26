// 13 december 2015

package ui

// #include "pkgui.h"
import "C"

// DrawPath represents a geometric path in a drawing context.
// This is the basic unit of drawing: all drawing operations consist of
// forming a path, then stroking, filling, or clipping to that path.
// A path is an OS resource; you must explicitly free it when finished.
// Paths consist of multiple figures. Once you have added all the
// figures to a path, you must "end" the path to make it ready to draw
// with.
// TODO rewrite all that
// 
// Or more visually, the lifecycle of a Path is
// 	p := DrawNewPath()
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
// 
// A DrawPath also defines its fill mode. (This should ideally be a fill
// parameter, but some implementations prevent it.)
// TODO talk about fill modes
type DrawPath struct {
	p	*C.uiDrawPath
}

// TODO
// 
// TODO disclaimer
type DrawFillMode uint
const (
	DrawFillModeWinding DrawFillMode = iota
	DrawFillModeAlternate
)

// DrawNewPath creates a new DrawPath with the given fill mode.
func DrawNewPath(fillMode DrawFillMode) *DrawPath {
	var fm C.uiDrawFillMode

	switch fillMode {
	case DrawFillModeWinding:
		fm = C.uiDrawFillModeWinding
	case DrawFillModeAlternate:
		fm = C.uiDrawFillModeAlternate
	default:
		panic("invalid fill mode passed to ui.NewPath()")
	}
	return &DrawPath{
		p:	C.uiDrawNewPath(fm),
	}
}

// Free destroys a DrawPath. After calling Free the DrawPath cannot
// be used.
func (p *DrawPath) Free() {
	C.uiDrawFreePath(p.p)
}

// NewFigure starts a new figure in the DrawPath. The current point
// is set to the given point.
func (p *DrawPath) NewFigure(x float64, y float64) {
	C.uiDrawPathNewFigure(p.p, C.double(x), C.double(y))
}

// NewFigureWithArc starts a new figure in the DrawPath and adds
// an arc as the first element of the figure. Unlike ArcTo,
// NewFigureWithArc does not draw an initial line segment.
// Otherwise, see ArcTo.
func (p *DrawPath) NewFigureWithArc(xCenter float64, yCenter float64, radius float64, startAngle float64, sweep float64, isNegative bool) {
	C.uiDrawPathNewFigureWithArc(p.p,
		C.double(xCenter), C.double(yCenter),
		C.double(radius),
		C.double(startAngle), C.double(sweep),
		frombool(isNegative))
}

// LineTo adds a line to the current figure of the DrawPath starting
// from the current point and ending at the given point. The current
// point is set to the ending point.
func (p *DrawPath) LineTo(x float64, y float64) {
	C.uiDrawPathLineTo(p.p, C.double(x), C.double(y))
}

// ArcTo adds a circular arc to the current figure of the DrawPath.
// You pass it the center of the arc, its radius in radians, the starting
// angle (couterclockwise) in radians, and the number of radians the
// arc should sweep (counterclockwise). A line segment is drawn from
// the current point to the start of the arc. The current point is set to
// the end of the arc.
func (p *DrawPath) ArcTo(xCenter float64, yCenter float64, radius float64, startAngle float64, sweep float64, isNegative bool) {
	C.uiDrawPathArcTo(p.p,
		C.double(xCenter), C.double(yCenter),
		C.double(radius),
		C.double(startAngle), C.double(sweep),
		frombool(isNegative))
}

// BezierTo adds a cubic Bezier curve to the current figure of the
// DrawPath. Its start point is the current point. c1x and c1y are the
// first control point. c2x and c2y are the second control point. endX
// and endY are the end point. The current point is set to the end
// point.
func (p *DrawPath) BezierTo(c1x float64, c1y float64, c2x float64, c2y float64, endX float64, endY float64) {
	C.uiDrawPathBezierTo(p.p,
		C.double(c1x), C.double(c1y),
		C.double(c2x), C.double(c2y),
		C.double(endX), C.double(endY))
}

// CloseFigure draws a line segment from the current point of the
// current figure of the DrawPath back to its initial point. After calling
// this, the current figure is over and you must either start a new
// figure or end the DrawPath. If this is not called and you start a
// new figure or end the DrawPath, then the current figure will not
// have this closing line segment added to it (but the figure will still
// be over).
func (p *DrawPath) CloseFigure() {
	C.uiDrawPathCloseFigure(p.p)
}

// AddRectangle creates a new figure in the DrawPath that consists
// entirely of a rectangle whose top-left corner is at the given point
// and whose size is the given size. The rectangle is a closed figure;
// you must either start a new figure or end the Path after calling
// this method.
func (p *DrawPath) AddRectangle(x float64, y float64, width float64, height float64) {
	C.uiDrawPathAddRectangle(p.p, C.double(x), C.double(y), C.double(width), C.double(height))
}

// End ends the current DrawPath. You cannot add figures to a
// DrawPath that has been ended. You cannot draw with a
// DrawPath that has not been ended.
func (p *DrawPath) End() {
	C.uiDrawPathEnd(p.p)
}

// DrawContext represents a drawing surface that you can draw to.
// At present the only DrawContexts are surfaces associated with
// Areas and are provided by package ui; see AreaDrawParams.
type DrawContext struct {
	c	*C.uiDrawContext
}

// DrawBrushType defines the various types of brushes.
// 
// TODO disclaimer
type DrawBrushType int
const (
	DrawBrushTypeSolid DrawBrushType = iota
	DrawBrushTypeLinearGradient
	DrawBrushTypeRadialGradient
	DrawBrushTypeImage		// presently unimplemented
)

// TODO
// 
// TODO disclaimer
// TODO rename these to put LineCap at the beginning? or just Cap?
type DrawLineCap int
const (
	DrawLineCapFlat DrawLineCap = iota
	DrawLineCapRound
	DrawLineCapSquare
)

// TODO
// 
// TODO disclaimer
type DrawLineJoin int
const (
	DrawLineJoinMiter DrawLineJoin = iota
	DrawLineJoinRound
	DrawLineJoinBevel
)

// TODO document
const DrawDefaultMiterLimit = 10.0

// TODO
type DrawBrush struct {
	Type		DrawBrushType

	// If Type is Solid.
	// TODO
	R		float64
	G		float64
	B		float64
	A		float64

	// If Type is LinearGradient or RadialGradient.
	// TODO
	X0			float64	// start point for both
	Y0			float64
	X1			float64	// linear: end point; radial: circle center
	Y1			float64
	OuterRadius	float64	// for radial gradients only
	Stops		[]DrawGradientStop
}

// TODO
type DrawGradientStop struct {
	Pos	float64		// between 0 and 1 inclusive
	R	float64
	G	float64
	B	float64
	A	float64
}

func (b *DrawBrush) toLibui() *C.uiDrawBrush {
	cb := C.pkguiAllocBrush()
	cb.Type = C.uiDrawBrushType(b.Type)
	switch b.Type {
	case DrawBrushTypeSolid:
		cb.R = C.double(b.R)
		cb.G = C.double(b.G)
		cb.B = C.double(b.B)
		cb.A = C.double(b.A)
	case DrawBrushTypeLinearGradient, DrawBrushTypeRadialGradient:
		cb.X0 = C.double(b.X0)
		cb.Y0 = C.double(b.Y0)
		cb.X1 = C.double(b.X1)
		cb.Y1 = C.double(b.Y1)
		cb.OuterRadius = C.double(b.OuterRadius)
		cb.NumStops = C.size_t(len(b.Stops))
		cb.Stops = C.pkguiAllocGradientStops(cb.NumStops)
		for i, s := range b.Stops {
			C.pkguiSetGradientStop(cb.Stops, C.size_t(i),
				C.double(s.Pos),
				C.double(s.R),
				C.double(s.G),
				C.double(s.B),
				C.double(s.A))
		}
	case DrawBrushTypeImage:
		panic("unimplemented")
	default:
		panic("invalid brush type in Brush.toLibui()")
	}
	return cb
}

func freeBrush(cb *C.uiDrawBrush) {
	if cb.Type == C.uiDrawBrushTypeLinearGradient || cb.Type == C.uiDrawBrushTypeRadialGradient {
		C.pkguiFreeGradientStops(cb.Stops)
	}
	C.pkguiFreeBrush(cb)
}

// TODO
type DrawStrokeParams struct {
	Cap			DrawLineCap
	Join			DrawLineJoin
	Thickness		float64
	MiterLimit		float64
	Dashes		[]float64
	DashPhase	float64
}

func (sp *DrawStrokeParams) toLibui() *C.uiDrawStrokeParams {
	csp := C.pkguiAllocStrokeParams()
	csp.Cap = C.uiDrawLineCap(sp.Cap)
	csp.Join = C.uiDrawLineJoin(sp.Join)
	csp.Thickness = C.double(sp.Thickness)
	csp.MiterLimit = C.double(sp.MiterLimit)
	csp.Dashes = nil
	csp.NumDashes = C.size_t(len(sp.Dashes))
	if csp.NumDashes != 0 {
		csp.Dashes = C.pkguiAllocDashes(csp.NumDashes)
		for i, d := range sp.Dashes {
			C.pkguiSetDash(csp.Dashes, C.size_t(i), C.double(d))
		}
	}
	csp.DashPhase = C.double(sp.DashPhase)
	return csp
}

func freeStrokeParams(csp *C.uiDrawStrokeParams) {
	if csp.Dashes != nil {
		C.pkguiFreeDashes(csp.Dashes)
	}
	C.pkguiFreeStrokeParams(csp)
}

// TODO
func (c *DrawContext) Stroke(p *DrawPath, b *DrawBrush, sp *DrawStrokeParams) {
	cb := b.toLibui()
	defer freeBrush(cb)
	csp := sp.toLibui()
	defer freeStrokeParams(csp)
	C.uiDrawStroke(c.c, p.p, cb, csp)
}

// TODO
func (c *DrawContext) Fill(p *DrawPath, b *DrawBrush) {
	cb := b.toLibui()
	defer freeBrush(cb)
	C.uiDrawFill(c.c, p.p, cb)
}

// TODO
// TODO should the methods of these return self for chaining?
type DrawMatrix struct {
	M11		float64
	M12		float64
	M21		float64
	M22		float64
	M31		float64
	M32		float64
}

// TODO identity matrix
func DrawNewMatrix() *DrawMatrix {
	m := new(DrawMatrix)
	m.SetIdentity()
	return m
}

// TODO
func (m *DrawMatrix) SetIdentity() {
	m.M11 = 1
	m.M12 = 0
	m.M21 = 0
	m.M22 = 1
	m.M31 = 0
	m.M32 = 0
}

func (m *DrawMatrix) toLibui() *C.uiDrawMatrix {
	cm := C.pkguiAllocMatrix()
	cm.M11 = C.double(m.M11)
	cm.M12 = C.double(m.M12)
	cm.M21 = C.double(m.M21)
	cm.M22 = C.double(m.M22)
	cm.M31 = C.double(m.M31)
	cm.M32 = C.double(m.M32)
	return cm
}

func (m *DrawMatrix) fromLibui(cm *C.uiDrawMatrix) {
	m.M11 = float64(cm.M11)
	m.M12 = float64(cm.M12)
	m.M21 = float64(cm.M21)
	m.M22 = float64(cm.M22)
	m.M31 = float64(cm.M31)
	m.M32 = float64(cm.M32)
	C.pkguiFreeMatrix(cm)
}

// TODO
func (m *DrawMatrix) Translate(x float64, y float64) {
	cm := m.toLibui()
	C.uiDrawMatrixTranslate(cm, C.double(x), C.double(y))
	m.fromLibui(cm)
}

// TODO
func (m *DrawMatrix) Scale(xCenter float64, yCenter float64, x float64, y float64) {
	cm := m.toLibui()
	C.uiDrawMatrixScale(cm,
		C.double(xCenter), C.double(yCenter),
		C.double(x), C.double(y))
	m.fromLibui(cm)
}

// TODO
func (m *DrawMatrix) Rotate(x float64, y float64, amount float64) {
	cm := m.toLibui()
	C.uiDrawMatrixRotate(cm, C.double(x), C.double(y), C.double(amount))
	m.fromLibui(cm)
}

// TODO
func (m *DrawMatrix) Skew(x float64, y float64, xamount float64, yamount float64) {
	cm := m.toLibui()
	C.uiDrawMatrixSkew(cm,
		C.double(x), C.double(y),
		C.double(xamount), C.double(yamount))
	m.fromLibui(cm)
}

// TODO
func (m *DrawMatrix) Multiply(m2 *DrawMatrix) {
	cm := m.toLibui()
	cm2 := m2.toLibui()
	C.uiDrawMatrixMultiply(cm, cm2)
	C.pkguiFreeMatrix(cm2)
	m.fromLibui(cm)
}

// TODO
func (m *DrawMatrix) Invertible() bool {
	cm := m.toLibui()
	res := C.uiDrawMatrixInvertible(cm)
	C.pkguiFreeMatrix(cm)
	return tobool(res)
}

// TODO
// 
// If m is not invertible, false is returned and m is left unchanged.
func (m *DrawMatrix) Invert() bool {
	cm := m.toLibui()
	res := C.uiDrawMatrixInvert(cm)
	m.fromLibui(cm)
	return tobool(res)
}

// TODO unimplemented
func (m *DrawMatrix) TransformPoint(x float64, y float64) (xout float64, yout float64) {
	panic("TODO")
}

// TODO unimplemented
func (m *DrawMatrix) TransformSize(x float64, y float64) (xout float64, yout float64) {
	panic("TODO")
}

// TODO
func (c *DrawContext) Transform(m *DrawMatrix) {
	cm := m.toLibui()
	C.uiDrawTransform(c.c, cm)
	C.pkguiFreeMatrix(cm)
}

// TODO
func (c *DrawContext) Clip(p *DrawPath) {
	C.uiDrawClip(c.c, p.p)
}

// TODO
func (c *DrawContext) Save() {
	C.uiDrawSave(c.c)
}

// TODO
func (c *DrawContext) Restore() {
	C.uiDrawRestore(c.c)
}
