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

TODO for all of the following: verify API call data types before moving code

## Windows
We create another custom window class that does `WM_PAINT` and handles input events thereof.

For this mockup, I'll extract the message handling into its own function and assume I can call Windows API functions and use their types and constants as normal. For `WM_PAINT` both `wparam` and `lparam` are unused.
```go
func repaint(s *sysData) HRESULT {
	var xrect RECT
	var ps PAINTSTRUCT

	// TODO send TRUE if we want to erase the clip area
	if GetUpdateRect(s.hwnd, &xrect, FALSE) == 0 {
		// no update rect, so we're done
		return 0
	}
	hdc, err := BeginPaint(s.hwnd, &ps)
	if hdc == 0 {		// failure
		panic(fmt.Errorf("error beginning Area repaint: %v", err))
	}

	cliprect := image.Rect(int(xrect.Left), int(xrect.Top), int(xrect.Right), int(xrect.Bottom))
	imgret := make(chan *image.NRGBA)
	defer close(imgret)
	s.paint <- PaintRequest{
		Rect:		cliprect,
		Out:		imgret,
	}
	i := <-imgret

	// drawing code here; see below

	EndPaint(s.hwnd, &ps)
	return 0
}
```

We can use GDI+ (gdiplus.dll) and its flat API for drawing...
```c
GpStatus WINGDIPAPI GdipCreateBitmapFromScan0(INT width, INT height, INT stride, PixelFormat format, BYTE* scan0, GpBitmap** bitmap);
GpStatus WINGDIPAPI GdipCreateFromHDC(HDC hdc, GpGraphics **graphics);
GpStatus WINGDIPAPI GdipDrawImageI(GpGraphics *graphics, GpImage *image, INT x, INT y);
GpStatus WINGDIPAPI GdipDeleteGraphics(GpGraphics *graphics);
GpStatus WINGDIPAPI GdipDisposeImage(GpImage *image);
```
(`GpBitmap` extends `GpImage`.) The only problem is the pixel format: the most appropriate one is `PixelFormat32bppARGB`, which is not premultiplied, but the components are in the wrong order... (specifically in BGRA order) (there is no RGBA pixel format in any bit width) (TODO `GdipDisposeImage` seems wrong since it bypasses `~Bitmap()` and goes right for `~Image()` but I don't see an explicit `~Bitmap()`...)

Disregarding the RGBA issue, the draw code would be
```go
	var bitmap, graphics uintptr

	status := GdipCreateBitmapFromScan0(
		i.Rect.Dx(),
		i.Rect.Dy(),
		i.Stride,
		PixelFormat32bppARGB,
		(*byte)(unsafe.Pointer(&i.Pix[0])),
		&bitmap)
	if status != 0 {		// failure
		panic(fmt.Errorf("error creating GDI+ bitmap to blit (GDI+ error code %d)", status))
	}
	status = GdipCreateFromHDC(hdc, &graphics)
	if status != 0 {		// failure
		panic(fmt.Errorf("error creating GDI+ graphics context to blit to (GDI+ error code %d)", status))
	}
	status = GdipDrawImageI(graphics, bitmap, cliprect.Min.X, cliprect.Min.Y)
	if status != 0 {		// failure
		panic(fmt.Errorf("error blitting GDI+ bitmap (GDI+ error code %d)", status))
	}
	status = GdipDeleteGraphics(graphics)
	if status != 0 {		// failure
		panic(fmt.Errorf("error freeing GDI+ graphics context to blit to (GDI+ error code %d)", status))
	}
	status = GdipDisposeImage(bitmap)
	if status != 0 {		// failure
		panic(fmt.Errorf("error freeing GDI+ bitmap to blit (GDI+ error code %d)", status))
	}
```

Upon further review, there really doesn't seem to be any way around it: we have to shuffle the image data around. We seem to be in good company: [go.wde needs to do so as well](https://github.com/skelterjohn/go.wde/blob/master/win/dib_windows.go). But you can't be too sure...
```go
	realbits := make([]byte, 4 * i.Rect.Dx() * I.Rect.Dy())
	q := 0
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		k := i.Pix[y * i.Stride:]
		for x := i.Rect.Min.X; x < i.Rect.Max.X; x += 4 {
			realbits[q + 0] = byte(k[y + x + 2])	// B
			realbits[q + 1] = byte(k[y + x + 1])	// G
			realbits[q + 2] = byte(k[y + x + 0])	// R
			realbits[q + 3] = byte(k[y + x + 3])	// A
			q += 4
		}
	}

	var bitmap, graphics uintptr

	status := GdipCreateBitmapFromScan0(
		i.Rect.Dx(),
		i.Rect.Dy(),
		i.Rect.Dy() * 4,			// got rid of extra stride
		PixelFormat32bppARGB,
		&realbits[0],
		&bitmap)
	// rest of code
```

We must also initialize and shut down GDI+ in uitask:
```go
var (
	gdiplustoken uintptr
)

	// init
	startupinfo := &GdiplusStartupInput{
		GdiplusVersion:	1,
	}
	status := GdiplusStartup(&gdiplustoken, startupinfo, nil)
	if status != 0 {		// failure
		return fmt.Errorf("error initializing GDI+ (GDI+ error code %d)", status)
	}

	// shutdown
	GdiplusShutdown(gdiplustoken)
```

TODO is there a function to turn a `GpStatus` into a string?

TODO note http://msdn.microsoft.com/en-us/library/windows/desktop/bb775501%28v=vs.85%29.aspx#win_class for information on handling some key presses, tab switch, etc. (need to do this for the ones below too)

TODO figure out scrolling
	- windows can have either standard or custom scrollbars; need to figure out how the standard scrollbars work
	- standard scrollbars cannot be controlled by the keyboard; either we must provide an option for doing that or allow scrolling ourselves (the `myAreaGoroutine` would read the keyboard events and scroll manually, in the same way)

## GTK+
We can use `GtkDrawingArea`. We hook into the `draw` signal; it does something equivalent to

```go
func draw_callback(widget *C.GtkWidget, cr *C.cairo_t, data C.gpointer) C.gboolean {
	var x, y, w, h C.double

	s := (*sysData)(unsafe.Pointer(data))
	// thanks to desrt in irc.gimp.net/#gtk+
	C.cairo_clip_extents(cr, &x, &y, &w, &h)
	cliprect := image.Rect(int(x), int(y), int(w), int(h))
	imgret := make(chan *image.NRGBA)
	defer close(imgret)
	s.paint <- PaintRequest{
		Rect:		cliprect,
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
		C.gdouble(cliprect.Min.X),
		C.gdouble(cliprect.Min.Y))
	C.g_object_unref((C.gpointer)(unsafe.Pointer(pixbuf)))		// free pixbuf
	return C.FALSE		// TODO what does this return value mean? docs don't say
}
```

[Example 1 on this page](https://developer.gnome.org/gdk-pixbuf/2.26/gdk-pixbuf-The-GdkPixbuf-Structure.html) indicates the pixels are in RGBA order, which is good.

On alpha premultiplication:
```
12:27	andlabs	Hi. Is the pixel data fed to gdk-pixbuf alpha premultiplied, not alpha premultiplied, or is that settable? I need to feed it data from a source that doesn't know about the underlying rendering system. Thanks.
12:29		*** KaL_out is now known as KaL
12:29	desrt	andlabs: pixbuf is non-premultiplied
12:30	mclasen	sad that this information is not obvious in the docs
12:30	andlabs	there is no information about premultiplied in any of the GTK+ documentation, period
12:30	desrt	andlabs: we have a utility function to copy it to a cairo surface that does the multiply for you...
12:30	andlabs	(in versions compatible with ubuntu 12.04, at least)
12:31	andlabs	good to know, thanks
12:31	desrt	andlabs: i think it's because gdkpixbuf existed before premultiplication was a wide practice
12:31	desrt	so at the time nobody would have asked the question
12:31	andlabs	huh
```

`GtkDrawingArea` is not natively scrollable, so we use `gtk_scrolled_window_add_with_viewport()` to add it to a `GtkScrollingWindow` with an implicit `GtkViewport` that handles scrolling for us. Otherwise, it's like what we did for Listbox.

TODO "Note that GDK automatically clears the exposed area to the background color before sending the expose event" decide what to do for the other platforms

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
var (
	// for later
	initWithBitmapDataPlanes = sel_getUid("initWithBitmapDataPlanes:pixelsWide:pixelsHigh:bitsPerSample:samplesPerPixel:hasAlpha:isPlanar:colorSpaceName:bitmapFormat:bytesPerRow:bitsPerPixel:")
)

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
		bitmapFormat:NSAlphaNonpremultipliedBitmapFormat		// this is where the flag for placing alpha first would go if alpha came first; the default is alpha last, which is how we're doing things (otherwise the docs say "Color planes are arranged in the standard order—for example, red before green before blue for RGB color.")
		bytesPerRow:i.Stride
		bitsPerPixel:32]
	[bitmap drawAtPoint:NSMakePoint(cliprect.Min.X, cliprect.Min.Y)]
	[bitmap release]
}
```
Due to the size of the `NSBitmapImageRep` constructor, I might just have another C function that performs the `NSBitmapImageRep` constructor using the `image.NRGBA` fields.

Finally, we need to override `-[NSView isFlipped]` since we want to keep (0,0) at the top-left:
```go
//export our_isFlipped
func our_isFlipped(self C.id, sel C.SEL) C.BOOL {
	return C.BOOL(C.YES)
}
```

For scrolling, we simply wrap our view in a `NSScrollView` just as we did with Listbox; Cocoa handles all the details for us.

TODO erase clip rect?
