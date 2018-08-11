// 13 december 2015

package ui

// #include <stdlib.h>
// #include "ui.h"
// // TODO figure this one out
// extern void *uimalloc(size_t);
// static uiDrawBrush *newBrush(void)
// {
// 	uiDrawBrush *b;
// 
// 	b = (uiDrawBrush *) uimalloc(sizeof (uiDrawBrush));
// 	return b;
// }
// static uiDrawBrushGradientStop *newStops(size_t n)
// {
// 	uiDrawBrushGradientStop *stops;
// 
// 	stops = (uiDrawBrushGradientStop *) malloc(n * sizeof (uiDrawBrushGradientStop));
// 	// TODO
// 	return stops;
// }
// static void setStop(uiDrawBrushGradientStop *stops, size_t i, double pos, double r, double g, double b, double a)
// {
// 	stops[i].Pos = pos;
// 	stops[i].R = r;
// 	stops[i].G = g;
// 	stops[i].B = b;
// 	stops[i].A = a;
// }
// static void freeBrush(uiDrawBrush *b)
// {
// 	if (b->Type == uiDrawBrushTypeLinearGradient || b->Type == uiDrawBrushTypeRadialGradient)
// 		free(b->Stops);
// 	free(b);
// }
// static uiDrawStrokeParams *newStrokeParams(void)
// {
// 	uiDrawStrokeParams *b;
// 
// 	b = (uiDrawStrokeParams *) malloc(sizeof (uiDrawStrokeParams));
// 	// TODO
// 	return b;
// }
// static double *newDashes(size_t n)
// {
// 	double *dashes;
// 
// 	dashes = (double *) malloc(n * sizeof (double));
// 	// TODO
// 	return dashes;
// }
// static void setDash(double *dashes, size_t i, double dash)
// {
// 	dashes[i] = dash;
// }
// static void freeStrokeParams(uiDrawStrokeParams *sp)
// {
// 	if (sp->Dashes != NULL)
// 		free(sp->Dashes);
// 	free(sp);
// }
// static uiDrawMatrix *newMatrix(void)
// {
// 	uiDrawMatrix *m;
// 
// 	m = (uiDrawMatrix *) malloc(sizeof (uiDrawMatrix));
// 	// TODO
// 	return m;
// }
// static void freeMatrix(uiDrawMatrix *m)
// {
// 	free(m);
// }
// static uiDrawTextFontDescriptor *newFontDescriptor(void)
// {
// 	uiDrawTextFontDescriptor *desc;
// 
// 	desc = (uiDrawTextFontDescriptor *) malloc(sizeof (uiDrawTextFontDescriptor));
// 	// TODO
// 	return desc;
// }
// static uiDrawTextFont *newFont(uiDrawTextFontDescriptor *desc)
// {
// 	uiDrawTextFont *font;
// 
// 	font = uiDrawLoadClosestFont(desc);
// 	free((char *) (desc->Family));
// 	free(desc);
// 	return font;
// }
// static uiDrawTextLayout *newTextLayout(char *text, uiDrawTextFont *defaultFont, double width)
// {
// 	uiDrawTextLayout *layout;
// 
// 	layout = uiDrawNewTextLayout(text, defaultFont, width);
// 	free(text);
// 	return layout;
// }
// static uiDrawTextFontMetrics *newFontMetrics(void)
// {
// 	uiDrawTextFontMetrics *m;
// 
// 	m = (uiDrawTextFontMetrics *) malloc(sizeof (uiDrawTextFontMetrics));
// 	// TODO
// 	return m;
// }
// static void freeFontMetrics(uiDrawTextFontMetrics *m)
// {
// 	free(m);
// }
// static double *newDouble(void)
// {
// 	double *d;
// 
// 	d = (double *) malloc(sizeof (double));
// 	// TODO
// 	return d;
// }
// static void freeDoubles(double *a, double *b)
// {
// 	free(a);
// 	free(b);
// }
import "C"

// BUG(andlabs): Ideally, all the drawing APIs should be in another package ui/draw (they all have the "uiDraw" prefix in C to achieve a similar goal of avoiding confusing programmers via namespace pollution); managing the linkage of the libui shared library itself across multiple packages is likely going to be a pain, though. (Custom controls implemented using libui won't have this issue, as they *should* only need libui present when linking the shared object, not when linking the Go wrapper. I'm not sure; I'd have to find out first.)

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
// 
// A Path also defines its fill mode. (This should ideally be a fill
// parameter, but some implementations prevent it.)
// TODO talk about fill modes
type Path struct {
	p	*C.uiDrawPath
}

// TODO
// 
// TODO disclaimer
type FillMode uint
const (
	Winding FillMode = iota
	Alternate
)

// NewPath creates a new Path with the given fill mode.
func NewPath(fillMode FillMode) *Path {
	var fm C.uiDrawFillMode

	switch fillMode {
	case Winding:
		fm = C.uiDrawFillModeWinding
	case Alternate:
		fm = C.uiDrawFillModeAlternate
	default:
		panic("invalid fill mode passed to ui.NewPath()")
	}
	return &Path{
		p:	C.uiDrawNewPath(fm),
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

// BrushType defines the various types of brushes.
// 
// TODO disclaimer
type BrushType int
const (
	Solid BrushType = iota
	LinearGradient
	RadialGradient
	Image		// presently unimplemented
)

// TODO
// 
// TODO disclaimer
// TODO rename these to put LineCap at the beginning? or just Cap?
type LineCap int
const (
	FlatCap LineCap = iota
	RoundCap
	SquareCap
)

// TODO
// 
// TODO disclaimer
type LineJoin int
const (
	MiterJoin LineJoin = iota
	RoundJoin
	BevelJoin
)

// TODO document
const DefaultMiterLimit = 10.0

// TODO
type Brush struct {
	Type		BrushType

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
	Stops		[]GradientStop
}

// TODO
type GradientStop struct {
	Pos	float64		// between 0 and 1 inclusive
	R	float64
	G	float64
	B	float64
	A	float64
}

func (b *Brush) toC() *C.uiDrawBrush {
	cb := C.newBrush()
	cb.Type = C.uiDrawBrushType(b.Type)
	switch b.Type {
	case Solid:
		cb.R = C.double(b.R)
		cb.G = C.double(b.G)
		cb.B = C.double(b.B)
		cb.A = C.double(b.A)
	case LinearGradient, RadialGradient:
		cb.X0 = C.double(b.X0)
		cb.Y0 = C.double(b.Y0)
		cb.X1 = C.double(b.X1)
		cb.Y1 = C.double(b.Y1)
		cb.OuterRadius = C.double(b.OuterRadius)
		cb.NumStops = C.size_t(len(b.Stops))
		cb.Stops = C.newStops(cb.NumStops)
		for i, s := range b.Stops {
			C.setStop(cb.Stops, C.size_t(i),
				C.double(s.Pos),
				C.double(s.R),
				C.double(s.G),
				C.double(s.B),
				C.double(s.A))
		}
	case Image:
		panic("unimplemented")
	default:
		panic("invalid brush type in Brush.toC()")
	}
	return cb
}

// TODO
type StrokeParams struct {
	Cap			LineCap
	Join			LineJoin
	Thickness		float64
	MiterLimit		float64
	Dashes		[]float64
	DashPhase	float64
}

func (sp *StrokeParams) toC() *C.uiDrawStrokeParams {
	csp := C.newStrokeParams()
	csp.Cap = C.uiDrawLineCap(sp.Cap)
	csp.Join = C.uiDrawLineJoin(sp.Join)
	csp.Thickness = C.double(sp.Thickness)
	csp.MiterLimit = C.double(sp.MiterLimit)
	csp.Dashes = nil
	csp.NumDashes = C.size_t(len(sp.Dashes))
	if csp.NumDashes != 0 {
		csp.Dashes = C.newDashes(csp.NumDashes)
		for i, d := range sp.Dashes {
			C.setDash(csp.Dashes, C.size_t(i), C.double(d))
		}
	}
	csp.DashPhase = C.double(sp.DashPhase)
	return csp
}

// TODO
func (c *DrawContext) Stroke(p *Path, b *Brush, sp *StrokeParams) {
	cb := b.toC()
	csp := sp.toC()
	C.uiDrawStroke(c.c, p.p, cb, csp)
	C.freeBrush(cb)
	C.freeStrokeParams(csp)
}

// TODO
func (c *DrawContext) Fill(p *Path, b *Brush) {
	cb := b.toC()
	C.uiDrawFill(c.c, p.p, cb)
	C.freeBrush(cb)
}

// TODO
// TODO should the methods of these return self for chaining?
type Matrix struct {
	M11		float64
	M12		float64
	M21		float64
	M22		float64
	M31		float64
	M32		float64
}

// TODO identity matrix
func NewMatrix() *Matrix {
	m := new(Matrix)
	m.SetIdentity()
	return m
}

// TODO
func (m *Matrix) SetIdentity() {
	m.M11 = 1
	m.M12 = 0
	m.M21 = 0
	m.M22 = 1
	m.M31 = 0
	m.M32 = 0
}

func (m *Matrix) toC() *C.uiDrawMatrix {
	cm := C.newMatrix()
	cm.M11 = C.double(m.M11)
	cm.M12 = C.double(m.M12)
	cm.M21 = C.double(m.M21)
	cm.M22 = C.double(m.M22)
	cm.M31 = C.double(m.M31)
	cm.M32 = C.double(m.M32)
	return cm
}

func (m *Matrix) fromC(cm *C.uiDrawMatrix) {
	m.M11 = float64(cm.M11)
	m.M12 = float64(cm.M12)
	m.M21 = float64(cm.M21)
	m.M22 = float64(cm.M22)
	m.M31 = float64(cm.M31)
	m.M32 = float64(cm.M32)
	C.freeMatrix(cm)
}

// TODO
func (m *Matrix) Translate(x float64, y float64) {
	cm := m.toC()
	C.uiDrawMatrixTranslate(cm, C.double(x), C.double(y))
	m.fromC(cm)
}

// TODO
func (m *Matrix) Scale(xCenter float64, yCenter float64, x float64, y float64) {
	cm := m.toC()
	C.uiDrawMatrixScale(cm,
		C.double(xCenter), C.double(yCenter),
		C.double(x), C.double(y))
	m.fromC(cm)
}

// TODO
func (m *Matrix) Rotate(x float64, y float64, amount float64) {
	cm := m.toC()
	C.uiDrawMatrixRotate(cm, C.double(x), C.double(y), C.double(amount))
	m.fromC(cm)
}

// TODO
func (m *Matrix) Skew(x float64, y float64, xamount float64, yamount float64) {
	cm := m.toC()
	C.uiDrawMatrixSkew(cm,
		C.double(x), C.double(y),
		C.double(xamount), C.double(yamount))
	m.fromC(cm)
}

// TODO
func (m *Matrix) Multiply(m2 *Matrix) {
	cm := m.toC()
	cm2 := m2.toC()
	C.uiDrawMatrixMultiply(cm, cm2)
	C.freeMatrix(cm2)
	m.fromC(cm)
}

// TODO
func (m *Matrix) Invertible() bool {
	cm := m.toC()
	res := C.uiDrawMatrixInvertible(cm)
	C.freeMatrix(cm)
	return tobool(res)
}

// TODO
// 
// If m is not invertible, false is returned and m is left unchanged.
func (m *Matrix) Invert() bool {
	cm := m.toC()
	res := C.uiDrawMatrixInvert(cm)
	m.fromC(cm)
	return tobool(res)
}

// TODO unimplemented
func (m *Matrix) TransformPoint(x float64, y float64) (xout float64, yout float64) {
	panic("TODO")
}

// TODO unimplemented
func (m *Matrix) TransformSize(x float64, y float64) (xout float64, yout float64) {
	panic("TODO")
}

// TODO
func (c *DrawContext) Transform(m *Matrix) {
	cm := m.toC()
	C.uiDrawTransform(c.c, cm)
	C.freeMatrix(cm)
}

// TODO
func (c *DrawContext) Clip(p *Path) {
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

// FontFamilies represents an enumerator over the font families
// available for use by package ui. A FontFamilies object behaves
// similarly to a []string, except that since family names are loaded
// on demand (depending on the operating system), it is not an
// actual []string. You call ListFontFamilies to obtain a FontFamilies
// object, which should reflect the available fonts at the time of the
// call (TODO verify). Use NumFamilies to get the number of families,
// and Family to get the name of a given family by index. When
// finished, call Free.
// 
// There is no guarantee that the list of families is sorted. You will
// need to do sorting yourself if you need it.
// 
// TODO thread affinity
type FontFamilies struct {
	ff *C.uiDrawFontFamilies
}

// ListFontFamilies creates a new FontFamilies object ready for use.
func ListFontFamilies() *FontFamilies {
	return &FontFamilies{
		ff:	C.uiDrawListFontFamilies(),
	}
}

// Free destroys a FontFamilies. After calling Free, the FontFamilies
// cannot be used.
func (f *FontFamilies) Free() {
	C.uiDrawFreeFontFamilies(f.ff)
}

// NumFamilies returns the number of font families available.
func (f *FontFamilies) NumFamilies() int {
	return int(C.uiDrawFontFamiliesNumFamilies(f.ff))
}

// Family returns the name of the nth family in the list.
func (f *FontFamilies) Family(n int) string {
	cname := C.uiDrawFontFamiliesFamily(f.ff, C.uintmax_t(n))
	name := C.GoString(cname)
	C.uiFreeText(cname)
	return name
}

// TextWeight defines the various text weights, in order of
// increasing weight.
// 
// Note that if you leave this field unset, it will default to
// TextWeightThin. If you want the normal font weight, explicitly
// use the constant TextWeightNormal instead.
// TODO realign these?
// 
// TODO disclaimer
type TextWeight int
const (
	TextWeightThin TextWeight = iota
	TextWeightUltraLight
	TextWeightLight
	TextWeightBook
	TextWeightNormal
	TextWeightMedium
	TextWeightSemiBold
	TextWeightBold
	TextWeightUtraBold
	TextWeightHeavy
	TextWeightUltraHeavy
)

// TextItalic defines the various text italic modes.
// 
// TODO disclaimer
type TextItalic int
const (
	TextItalicNormal TextItalic = iota
	TextItalicOblique			// merely slanted text
	TextItalicItalic				// true italics
)

// TextStretch defines the various text stretches, in order of
// increasing wideness.
// 
// Note that if you leave this field unset, it will default to
// TextStretchUltraCondensed. If you want the normal font
// stretch, explicitly use the constant TextStretchNormal
// instead.
// TODO realign these?
// 
// TODO disclaimer
type TextStretch int
const (
	TextStretchUltraCondensed TextStretch = iota
	TextStretchExtraCondensed
	TextStretchCondensed
	TextStretchSemiCondensed
	TextStretchNormal
	TextStretchSemiExpanded
	TextStretchExpanded
	TextStretchExtraExpanded
	TextStretchUltraExpanded
)

// FontDescriptor describes a Font.
type FontDescriptor struct {
	Family		string
	Size			float64		// as a text size, for instance 12 for a 12-point font
	Weight		TextWeight
	Italic			TextItalic
	Stretch		TextStretch
}

// Font represents an actual font that can be drawn with.
type Font struct {
	f	*C.uiDrawTextFont
}

// LoadClosestFont loads a Font.
// 
// You pass the properties of the ideal font you want to load in the
// FontDescriptor you pass to this function. If the requested font
// is not available on the system, the closest matching font is used.
// This means that, for instance, if you specify a Weight of
// TextWeightUltraHeavy and the heaviest weight available for the
// chosen font family is actually TextWeightBold, that will be used
// instead. The specific details of font matching beyond this
// description are implementation defined. This also means that
// getting a descriptor back out of a Font may return a different
// desriptor.
// 
// TODO guarantee that passing *that* back into LoadClosestFont() returns the same font
func LoadClosestFont(desc *FontDescriptor) *Font {
	d := C.newFontDescriptor()		// both of these are freed by C.newFont()
	d.Family = C.CString(desc.Family)
	d.Size = C.double(desc.Size)
	d.Weight = C.uiDrawTextWeight(desc.Weight)
	d.Italic = C.uiDrawTextItalic(desc.Italic)
	d.Stretch = C.uiDrawTextStretch(desc.Stretch)
	return &Font{
		f:	C.newFont(d),
	}
}

// Free destroys a Font. After calling Free the Font cannot be used.
func (f *Font) Free() {
	C.uiDrawFreeTextFont(f.f)
}

// Handle returns the OS font object that backs this Font. On OSs
// that use reference counting for font objects, Handle does not
// increment the reference count; you are sharing package ui's
// reference.
// 
// On Windows this is a pointer to an IDWriteFont.
// 
// On Unix systems this is a pointer to a PangoFont.
// 
// On OS X this is a CTFontRef.
func (f *Font) Handle() uintptr {
	return uintptr(C.uiDrawTextFontHandle(f.f))
}

// Describe returns the FontDescriptor that most closely matches
// this Font.
// TODO guarantees about idempotency
// TODO rewrite that first sentence
func (f *Font) Describe() *FontDescriptor {
	panic("TODO unimplemented")
}

// FontMetrics holds various measurements about a Font.
// All metrics are in the same point units used for drawing.
type FontMetrics struct {
	// Ascent is the ascent of the font; that is, the distance from
	// the top of the character cell to the baseline.
	Ascent			float64

	// Descent is the descent of the font; that is, the distance from
	// the baseline to the bottom of the character cell. The sum of
	// Ascent and Descent is the height of the character cell (and
	// thus, the maximum height of a line of text).
	Descent			float64

	// Leading is the amount of space the font designer suggests
	// to have between lines (between the bottom of the first line's
	// character cell and the top of the second line's character cell).
	// This is a suggestion; it is chosen by the font designer to
	// improve legibility.
	Leading			float64

	// TODO figure out what these are
	UnderlinePos		float64
	UnderlineThickness	float64
}

// Metrics returns metrics about the given Font.
func (f *Font) Metrics() *FontMetrics {
	m := new(FontMetrics)
	mm := C.newFontMetrics()
	C.uiDrawTextFontGetMetrics(f.f, mm)
	m.Ascent = float64(mm.Ascent)
	m.Descent = float64(mm.Descent)
	m.Leading = float64(mm.Leading)
	m.UnderlinePos = float64(mm.UnderlinePos)
	m.UnderlineThickness = float64(mm.UnderlineThickness)
	C.freeFontMetrics(mm)
	return m
}

// TextLayout is the entry point for formatting a block of text to be
// drawn onto a DrawContext.
// 
// The block of text to lay out and the default font that is used if no
// font attributes are applied to a given character are provided
// at TextLayout creation time and cannot be changed later.
// However, you may add attributes to various points of the text
// at any time, even after drawing the text once (unlike a DrawPath).
// Some of these attributes also have initial values; refer to each
// method to see what they are.
// 
// The block of text can either be a single line or multiple
// word-wrapped lines, each with a given maximum width.
type TextLayout struct {
	l	*C.uiDrawTextLayout
}

// NewTextLayout creates a new TextLayout.
// For details on the width parameter, see SetWidth.
func NewTextLayout(text string, defaultFont *Font, width float64) *TextLayout {
	l := new(TextLayout)
	ctext := C.CString(text)		// freed by C.newTextLayout()
	l.l = C.newTextLayout(ctext, defaultFont.f, C.double(width))
	return l
}

// Free destroys a TextLayout. After calling Free the TextLayout
// cannot be used.
func (l *TextLayout) Free() {
	C.uiDrawFreeTextLayout(l.l)
}

// SetWidth sets the maximum width of the lines of text in a
// TextLayout. If the given width is negative, then the TextLayout
// will draw as a single line of text instead.
func (l *TextLayout) SetWidth(width float64) {
	C.uiDrawTextLayoutSetWidth(l.l, C.double(width))
}

// Extents returns the width and height that the TextLayout will
// actually take up when drawn. This measures full line allocations,
// even if no glyph reaches to the top of its ascent or bottom of its
// descent; it does not return a "best fit" rectnagle for the points that
// are actually drawn.
// 
// For a single-line TextLayout (where the width is negative), if there
// are no font changes throughout the TextLayout, then the height
// returned by TextLayout is equivalent to the sum of the ascent and
// descent of its default font's metrics. Or in other words, after
// 	f := ui.LoadClosestFont(...)
// 	l := ui.NewTextLayout("text", f, -1)
// 	metrics := f.Metrics()
// 	_, height := l.Extents()
// metrics.Ascent+metrics.Descent and height are equivalent.
func (l *TextLayout) Extents() (width float64, height float64) {
	cwidth := C.newDouble()
	cheight := C.newDouble()
	C.uiDrawTextLayoutExtents(l.l, cwidth, cheight)
	width = float64(*cwidth)
	height = float64(*cheight)
	C.freeDoubles(cwidth, cheight)
	return width, height
}

// Text draws the given TextLayout onto c at the given point.
// The point refers to the top-left corner of the text.
// (TODO bounding box or typographical extent?)
func (c *DrawContext) Text(x float64, y float64, layout *TextLayout) {
	C.uiDrawText(c.c, C.double(x), C.double(y), layout.l)
}
