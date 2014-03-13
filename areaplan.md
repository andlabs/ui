```go
type Area struct {		// implements Control
	// Paint receives requests to redraw from the window system.
	Paint		chan PaintRequest

	// Keyboard receives keyboard events.
	Key		chan KeyEvent			// not covered here

	// Mouse receives mouse events.
	Mouse	chan MouseEvent		// not covered here
}

// PaintRequest represents a request to redraw an Area.
// It is sent across Area.Paint.
type PaintRequest struct {
	// Rect is the clipping rectangle that needs redraw.
	Rect		image.Rect

	// Out is a channel on which you send the image to redraw.
	Out		chan<- *image.NRGBA
}
```

and an example of intended use:

```go
func myAreaGoroutine(area *ui.Area, start <-chan bool) {
	var img *image.NRGBA

	// initialize img here
	<-start		// sent after calling Window.Open()
	area.SetSize(img.Rect.Dx(), img.Rect.Dy())	// sets the internal size; scrollbars and scrolling is handled automatically
	for {
		select {
		case req := <-area.Paint:
			req.Out <- img.SubImage(req.Rect).(*image.NRGBA)
		case e := <-area.Mouse:
			// draw on a mouse click, for instance
		}
	}
}
```

TODO is there a race on `area.SetSize()`?

## Windows
TODO

## GTK+
We can use `GtkDrawingArea`. We hook into the `draw` signal; it does something equivalent to

```go
func draw_callback(widget *C.GtkWidget, cr *C.cairo_t, data C.gpointer) C.gboolean {
	s := (*sysData)(unsafe.Pointer(data))
	// TODO get clip rectangle that needs drawing
	imgret := make(chan *image.NRGBA)
	defer close(imgret)
	s.paint <- PaintRequest{
		Rect:		/* clip rect */,
		Out:		imgret,
	}
	i := <-imgret
	pixbuf := C.gdk_pixbuf_new_from_data(
		(*C.guchar)(unsafe.Pointer(&i.Pix[0])),
		C.GDK_COLORSPACE_RGB,
		C.TRUE,			// has alpha channel
		8,				// bits per sample
		C.int(i.Rect.Dx()),
		C.int(i.Rect.Dy()),
		C.int(i.Stride),
		nil, nil)			// do not free data
	C.gdk_cairo_set_source_pixbuf(cr,
		pixbuf,
		/* x of clip rect */,
		/* y of clip rect */)
	return C.FALSE		// TODO what does this return value mean? docs don't say
}
```

TODO figure out how scrolling plays into this

## Cocoa
For this one we **must** create a subclass of `NSView` that overrides the drawing and keyboard/mouse event messages.

The drawing message is `-[NSView drawRect:]`, which just takes the `NSRect` as an argument. So we already need to use `bleh_darwin.m` to grab the actual `NSRect` and convert it into something with a predictable data type before passing it back to Go. If we do this:
```go
//export our_drawRect
func our_drawRect(self C.id, rect C.struct_xrect) {
```
we can call `our_drawRect()` from this C wrapper:
```objective-c
extern void our_drawRect(id, struct xrect);

void _our_drawRect(id self, SEL sel, NSRect r)
{
	struct xrect t;

	t.x = (int64_t) s.origin.x;
	t.y = (int64_t) s.origin.y;
	t.width = (int64_t) s.size.width;
	t.height = (int64_t) s.size.height;
	our_drawRect(self, t);
}
```
This just leaves `our_drawRect` itself. For this mockup, I will use "Objective-Go":
```go
//export our_drawRect
func our_drawRect(self C.id, rect C.struct_xrect) {
	s := getSysData(self)
	cliprect := image.Rect(int(rect.x), int(rect.y), int(rect.width), int(rect.height))
	imgret := make(chan *image.NRGBA)
	defer close(imgret)
	s.paint <- PaintRequest{
		Rect:		cliprect,
		Out:		imgret,
	}
	i := <-imgret
	// the NSBitmapImageRep constructor requires a list of pointers
	_bitmapData := [1]*uint8{&i.Pix[0]}
	bitmapData := (**C.uchar)(unsafe.Pointer(&bitmapData))
	bitmap := [[NSBitmapImageRep alloc]
		initWithBitmapDataPlanes:bitmapData
		pixelsWide:i.Rect.Dx()
		pixelsHigh:i.Rect.Dy()
		bitsPerSample:8
		samplesPerPixel:4
		hasAlpha:YES
		isPlanar:NO
		colorSpaceName:NSCalibratedRGBColorSpace		// TODO NSDeviceRGBColorSpace?
		bitmapFormat:NSAlphaNonpremultipliedBitmapFormat
		bytesPerRow:i.Stride
		bitsPerPixel:32]
	[bitmap drawAtPoint:NSMakePoint(cliprect.x, cliprect.y)]
	[bitmap release]
}
```
Due to the utter complexity of all that `NSImage` stuff, I might just have another C function that performs the `NSBitmapImageRep` constructor using the `image.NRGBA` fields.

Finally, we need to override `-[NSView isFlipped]` since we want to keep (0,0) at the top-left:
```go
//export our_isFlipped
func our_isFlipped(self C.id, sel C.SEL) C.BOOL {
	return C.BOOL(C.YES)
}
```

TODO figure out scrolling
