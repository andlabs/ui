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

## Drawing and Scrolling

### Windows
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

For scrolling, the custom window class will come with scrollbars. We are reponsible for scrolling ourselves:
- we handle `WM_HSCROLL` and `WM_VSCROLL` messages, extrapolating the scroll data
	- we can use `GetScrollInfo` to get the current position, but the example code on MSDN adjusts it manually and then calls `ScrollWindow` then `UpdateWindow` (to accelerate redraw) and then `SetScrollInfo` (to update the scroll info)
- line size is 1, page size is visible dimension
- call `SetScrollInfo` on control resizes, passing in a `SCROLLINFO` which indicates the above, does not include `SIF_DISABLENOSCROLL` so scrollbars are auto-hidden, and does not change either thumb position (`nPos` and `nTrackPos`)
- the clipping rectangle must take scrolling into account; `GetScrollInfo` and add the position to the sent-out `cliprect` (only; still need regular `cliprect` for drawing) with `cliprect.Add()`
- we should probably cache the scroll position and window sizes so we wouldn't need to call those respective functions each `WM_PAINT` and `WM_HSCROLL`/`WM_VSCROLL`, respectively
	- TODO will resizing a window with built-in scrollbars/adjusting the page size set the thumb and signal repaint properly?

TODO is there a function to turn a `GpStatus` into a string?

TODO note http://msdn.microsoft.com/en-us/library/windows/desktop/bb775501%28v=vs.85%29.aspx#win_class for information on handling some key presses, tab switch, etc. (need to do this for the ones below too)

TODO standard scrollbars cannot be controlled by the keyboard; either we must provide an option for doing that or allow scrolling ourselves (the `myAreaGoroutine` would read the keyboard events and scroll manually, in the same way)

### GTK+
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

`GtkDrawingArea` is not natively scrollable, so we use `gtk_scrolled_window_add_with_viewport()` to add it to a `GtkScrolledWindow` with an implicit `GtkViewport` that handles scrolling for us. Otherwise, it's like what we did for Listbox.

TODO "Note that GDK automatically clears the exposed area to the background color before sending the expose event" decide what to do for the other platforms

### Cocoa
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

**IMPORTANT NOTE**: Before we move on to events, Cocoa requires that we override `acceptsFirstResponder` to return `YES` in order to accept events:
```objective-c
- (BOOL)acceptsFirstResponder
{
	return YES;
}
```

TODO erase clip rect?

## Mouse Events

TODO scroll wheel

### Windows
Back to our custom window prodcedure again. We receive:
```
WM_LBUTTONDBLCLK
WM_LBUTTONDOWN
WM_LBUTTONUP
WM_MBUTTONDBLCLK
WM_MBUTTONDOWN
WM_MBUTTONUP
WM_RBUTTONDBLCLK
WM_RBUTTONDOWN
WM_RBUTTONUP
WM_XBUTTONDBLCLK
WM_XBUTTONDOWN
WM_XBUTTONUP 
```
which specify the left, middle, right, and up to two additional mouse buttons.

Each of these returns the coordinates in the LPARAM and the modifier flags in the WPARAM:
```
MK_CONTROL
MK_LBUTTON
MK_MBUTTON
MK_RBUTTON
MK_SHIFT
MK_XBUTTON1
MK_XBUTTON2
```
where the button modifier flags allow handling simultaneous clicks. The XBUTTON messages also use WPARAM to encode which button was pressed.

In order to register double-clicks, we have to specify the `CS_DBLCLKS` style when calling `RegisterClass`. A mouse click event will always be sent before a double-click event.

That just leaves mouse moves. All mouse moves are handled with `WM_MOUSEMOVE`, which returns the same WPARAM and LPARAM format as above (so we use the WPARAM to see which mouse buttons were held during a move).

All of these messages expect us to return 0, except the XBUTTON messages, which expect us to return TRUE.

MSDN says to use macros to get the position and XBUTTON information:
```c
/* for all messages */
xPos = GET_X_LPARAM(lParam);
yPos = GET_Y_LPARAM(lParam);

/* for XBUTTON messages */
fwKeys = GET_KEYSTATE_WPARAM (wParam);
fwButton = GET_XBUTTON_WPARAM (wParam);
```
We will need to reimplement these macros ourselves.

All messages are supported on at least Windows 2000, so we're good using them all.

There does not seem to be an equivalent to the mouse entered signal provided by GTK+ and Cocoa. There *is* an equivalent to mouse left (`WM_MOUSELEAVE`), but it requires tracking support, which has to be set up in special ways.

Finally, the Alt key has to be retrieved a differnet way. [This](http://stackoverflow.com/questions/9205534/win32-mouse-and-keyboard-combination) says we can use `GetKeyState(VK_MENU)`.

### GTK+
- `"button-press-event"` for mouse button presses; needs `GDK_BUTTON_PRESS_MASK` and returns `GdkEventButton`
- `"button-release-event"` for mouse button releases; needs `GDK_BUTTON_RELEASE_MASK` and returns `GdkEventButton`
- `"enter-notify-event"` for when the mouse enters the widget; needs `GDK_ENTER_NOTIFY_MASK` and returns `GdkEventCrossing`
- `"leave-notify-event"` for when the mouse leaves the widget; needs `GDK_LEAVE_NOTIFY_MASK` and returns `GdkEventCrossing`
- `"motion-notify-event"` for when the mouse moves while inside the widget; needs `GDK_POINTER_MOTION_MASK` and returns `GdkEventMotion`

The following events may also be of use:
```
GDK_BUTTON_MOTION_MASK
	receive pointer motion events while any button is pressed

GDK_BUTTON1_MOTION_MASK
	receive pointer motion events while 1 button is pressed

GDK_BUTTON2_MOTION_MASK
	receive pointer motion events while 2 button is pressed

GDK_BUTTON3_MOTION_MASK
	receive pointer motion events while 3 button is pressed 
```

`GdkEventButton` tells us:
- event type: click, double-click, triple-click, release
	- a click event is always sent before a double-click and triple-click event
		- double-click: click, release, <u>click</u>, double-click, release
		- triple-click: C, R, C, DC, R, C, TC, R
			- this goes against other OSs which don't send both a click and double-click on the double-click
- x and y positions of event
- modifier keys and other mouse buttons held during event: see https://developer.gnome.org/gdk3/stable/gdk3-Windows.html#GdkModifierType
	- does not appear to have a way to differentiate left and right modifier keys
	- see note below about Alt/Meta
- button ID of event, with order 1 - left, 2 - middle, 3 - right

`GdkEventCrossing` tells us
- whether this was an enter or a leave
- x and y positions of event
- "crossing mode" and "notification type" [not sure if I'll need these - https://developer.gnome.org/gdk3/stable/gdk3-Event-Structures.html#GdkEventCrossing]
- modifier/mose button held flags (see above)

`GdkEventMotion` tells us
- the type of the event (I assume this is always going to be `GDK_MOTION_NOTIFY`)
- x and y positions of the event
- modifier keys/mouse buttons held (as above)

GDK by default doesn't map *all* the modifier keys away from their device-speicifc values into portable values; we have to tell it to do so:
```go
	C.gdk_keymap_add_virtual_modifiers(C.gdk_keymap_get_default(), &e.state)
```
(thanks to Daniel_S and daniels (two different people) in irc.gimp.net/#gtk+) (note: the GDK 3.4 documentation has a rather complex description of what `gdk_keymap_add_virtual_modifiers()` does; the latest version has a much better description)

### Cocoa
Our `NSView` subclass will override the following:
```
mouseDown:
mouseDragged:
mouseUp:
mouseMoved:
mouseEntered:
mouseExited:
rightMouseDragged:
rightMouseUp:
otherMouseDown:
otherMouseDragged:
otherMouseUp:
```
The `mouse...` selectors are for the left mouse button. Each of these selectors is of the form
```objective-c
- (void)selectorName:(NSEvent *)e
```
where `NSEvent` is a concrete type, not an abstract class, that contains all the information we need.

...almost. `NSEvent` doesn't record mouse position directly, but rather relative to the view's parent window. The NSView Programming Guide says we can do
```go
	nspoint := [self convertPoint:[e locationInWindow] fromView:nil]
```
to get the point we want. This *should* also obey `isFlipped:`, as that affects "the coordinate system of the receiver".

For the button number, there's `-[e buttonNumber]`. The exact number is described below. The reference also says "This method is intended for use with the NSOtherMouseDown, NSOtherMouseUp, and NSOtherMouseDragged events, but will return values for NSLeftMouse... and NSRightMouse... events also.", so since we build our class at runtime, we can just assign the same implementation function to each type of event (the `sel` argument will differ, but since we can just get the button number directly we don't have to worry).

The click count is specified in `-[e clickCount]`, so we can distinguish between single-click and double-click easily. Note "Returns 0 for a mouse-up event if a time threshold has passed since the corresponding mouse-down event. This is because if this time threshold passes before the mouse button is released, it is no longer considered a mouse click, but a mouse-down event followed by a mouse-up event.". The Event Programing Guide says "Find out how many mouse clicks occurred in quick succession (clickCount); multiple mouse clicks are conceptually treated as a single mouse-down event within a narrow time threshold (although they arrive in a series of mouseDown: messages). As with modifier keys, a double- or triple-click can change the significance of a mouse event for an application. (See Listing 4-3 for an example.)" which indicates that a click event is sent before a double-click.

`-[e modifierFlags]` gives us the modifier flags. The flag reference is https://developer.apple.com/library/mac/documentation/Cocoa/Reference/ApplicationKit/Classes/NSEvent_Class/Reference/Reference.html#//apple_ref/doc/uid/20000016-SW14 - no info on left/right keys seems to be provided.

The first held mouse button could be handled by the drag events. The rest can be grabbed with `+[NSEvent pressedMouseButtons]` (thanks to Psy| in irc.freenode.net/#macdev for confirming)

Also according to Psy|, the bit order of `pressedMouseButtons` corresponds to the `buttonNumber`, so 0 is the left button, 1 is the right button, 2 is the middle button, and so on.

TODO do we need to override `acceptsFirstMouse:` to return `YES` so a click event is sent when changing the current program to this one?

### Consensus
```go
// MouseEvent contains all the information for a mous event sent by Area.Mouse.
// Mouse button IDs start at 1, with 1 being the left mouse button, 2 being the middle mouse button, and 3 being the right mouse button.
// (TODO "If additional buttons are supported, they will be returned with 4 being the first additional button (XBUTTON1 on Windows), 5 being the second (XBUTTON2 on Windows), and so on."?) (TODO get the user-facing name for XBUTTON1/2; find out if there's a way to query available button count)
type MouseEvent struct {
	// Pos is the position of the mouse relative to the top-left of the area.
	Pos			image.Point

	// If the event was generated by a mouse button being pressed, Down contains the ID of that button.
	// Otherwise, Down contains 0.
	Down		uint

	// If the event was generated by a mouse button being released, Up contains the ID of that button.
	// Otherwise, Up contains 0.
	Up			uint

	// If Down is nonzero, Count indicates the number of clicks: 1 for single-click, 2 for double-click.
	// If Count == 2, AT LEAST one event with Count == 1 will have been sent prior.
	// (This is a platform-specific issue: some platforms send one, some send two.)
	Count		uint

	// Modifiers is a bit mask indicating the modifier keys being held during the event.
	Modifiers		Modifiers

	// Held is a slice of button IDs that indicate which mouse buttons are being held during the event.
	// (TODO "There is no guarantee that Held is sorted."?)
	// (TODO will this include or exclude Down and Up?)
	Held			[]uint
}

// HeldBits returns Held as a bit mask.
// Bit 0 maps to button 1, bit 1 maps to button 2, etc.
func (e MousEvent) HeldBits() (h uintptr) {
	for _, x := range e.Held {
		h |= uintptr(1) << (x - 1)
	}
	return h
}

// Modifiers indicates modifier keys being held during a mouse event.
// There is no way to differentiate between left and right modifier keys.
type Modifiers uintptr
const (
	Ctrl Modifiers = 1 << iota		// the canonical Ctrl keys ([TODO] on Mac OS X, Control on others)
	Alt						// the canonical Alt keys ([TODO] on Mac OS X, Meta on Unix systems, Alt on others)
	Shift						// the Shift keys
)
```

## Keyboard Events
You thought mouse events were vaguely compromise-y? Get ready... this is going to hurt. *Bad*.

### Windows
> Note: all messages here except `WM_UNICHAR` work on Windows 2000 and newer and require us to return 0 on handled. All messages (including `WM_UNICHAR` take the same parameter format.

**The good**: Windows keyboard message parameters are in a consistent, predictable format<br>
**The bad**: everything else is not

Windows distinguishes between typical user input and "system keys"; system keys constitute three conditions:
* Alt+(any key)
* F10 (in some cases? *TODO*)
* any key when there is no active window on screen
System keys are special: if we don't handle them explicitly, we have to send them up to the `DefWindowProc()`. If we don't, things like Alt+Tab (!) won't be handled.

The `TranslateMessage()` call that appears in the message loop takes key down events, and if possible, converts them into Unicode character requests, handling IME properly. The key down events are not removed; new character events are inserted instead. There does not seem to be a good way to tell if one (or more! different parts of the docs say different things about how many character events come in per key down event **TODO**) has been inserted except with `PeekMessage()`, as `TranslateMessage()` always returns nonzero if a key down event is passed in, regardless of whether or not it was converted.

At the end of the day, the messages:

 | Regular | System
----- | ----- | -----
Key down | `WM_KEYDOWN` | `WM_SYSKEYDOWN`
Key up | `WM_KEYUP` | `WM_SYSKEYUP`
Character | `WM_CHAR` | `WM_SYSCHAR`
Dead key (character; we can ignore these) | `WM_DEADCHAR` | `WM_SYSDEADCHAR`

The WPARAM is the key code or UTF-16 character value. [List of virtual key codes](http://msdn.microsoft.com/en-us/library/dd375731%28VS.85%29.aspx)

The low word of the LPARAM is the repeat count. Multiple key down events will be sent, but we can use this to hold a count for convenience.

There's a lot of information in the high word of the LPARAM, but none of that is really important (and is useless for the character messages because of multi-character codes; it'll just match the last key down event). GLFW does use bit 24, the "extended key" bit, to differentiate between left and right keys (which is documented0, however there are some catches that we'll get to in a bit. For reference, though, the docs say
> For enhanced 101- and 102-key keyboards, extended keys are the right ALT and the right CTRL keys on the main section of the keyboard; the INS, DEL, HOME, END, PAGE UP, PAGE DOWN and arrow keys in the clusters to the left of the numeric keypad; and the divide (/) and ENTER keys in the numeric keypad. Some other keyboards may support the extended-key bit in the lParam parameter. 

I'm not entirely sure if this is the case, but GLFW seems to think  Windows sends the base `VK_xxx` codes on keys that have both left and right equivalents, not the dedicated `VK_Lxxx`/`VK_Rxxx` codes. (Compatibility?) For most cases, the extended key bit mentioned above is sufficient to differentiate. There are two exceptions
* Shift requires [checking the hardware scancode](https://github.com/glfw/glfw/blob/master/src/win32_window.c#L182) (which, fortunately, Windows provides a way to find out from the virtual key code at runtime)
* left Control: Windows apparently does not send a single key code on the ["AltGr" key found on some keyboards](https://en.wikipedia.org/wiki/AltGr_key), but rather both a left Control and a right Alt at the same time; [fortunately we can check](https://github.com/glfw/glfw/blob/master/src/win32_window.c#L199)

Key release also has some snags that the GLFW sources point out:
* The Shift key differentiation in Windows is broken: [only one key up is sent, even if two Shift keys were released](https://github.com/glfw/glfw/blob/master/src/win32_window.c#L538)
* Print Screen never sends a key down event!

Finally, Windows XP adds `WM_UNICHAR`, another key down event. According to the GLFW sources, Windows itself doesn't use this, but some IME drivers do. WPARAM (which is 32-bit now) stores the UTF-32 representation of a requested character, and LPARAM works as usual.
* If WPARAM is the special constant `UNICODE_NOCHAR`, we return `TRUE` if we handle this event, and `FALSE` otherwise (`DefWindowProc()` returns `FALSE`). This is how drivers will tell if we support `WM_UNICHAR` at all.
* Otherwise, we handle the character and return `FALSE`.

I do not know if `WM_UNICHAR` follows the same rules as `WM_CHAR`. *TODO*

*TODO*
* do `WM_CHAR`/`WM_SYSCHAR`/`WM_UNICHAR` get sent on repeat?

### GTK+
> Note: GLFW doesn't really help here since we're using GDK for event handling and GLFW uses X11 directly. (I don't want to call out to X11 functions because of Wayland support; hell I don't know if GDK even provides the X11 key code!)

Before our `GtkDrawingArea` can take keyboard input, we need to turn its `can-focus` property on.

There are two events here: `key-press-event` and `key-release-event`. These take the same shared event function prototype as the mouse events above, with the `GdkEvent` actual type being `GdkEventKey`. This type tells us:
- the modifier flags, just like with mouse events (even mouse buttons!)
- the GDK virtual key code, which aren't explicitly listed in the documentation; the docs say to [check `<gdk/gdkkeysyms.h>` instead](https://git.gnome.org/browse/gtk+/tree/gdk/gdkkeysyms.h?h=gtk-3-4).

Repeats are not documented; it appears that we just get sent multiple `key-press-event`s in a row, with no way to tell if a key was repeated.

Character conversion is iffy. There doesn't really seem to be a way to handle character input properly...
- Originally the `GdkEventKey` had fields that gave you that information, but this is now deprecated because of GTK+ input methods.
- There's `gdk_keyval_to_unicode()`, but you can only handle one key at a time this way.
- The only way to properly handle IME with GTK+ is to use GTK+ input methods via `GtkIMContext`, however there does not seem to be a way to get available contexts, only make new ones or pull the context from an existing `GtkEntry`/`GtkTextView`. (You [can get a list of available context names](https://developer.gnome.org/gtk3/3.4/GtkSettings.html#GtkSettings--gtk-im-module), but that's as much as I could find, and even then this can be `NULL`.) *TODO*

### Cocoa
**Windows**: either virtual key codes or character codes<br>
**GTK+**: virtual key codes only<br>
**Cocoa**: take a wild guess -_-

Our `NSView` subclass has three selectors to override:
```objective-c
- (void)keyDown:(NSEvent *e)			// key down
- (void)keyUp:(NSEvent *e)			// key up
- (void)flagsChanged:(NSEvent *e)		// modifier key state changed
```

Now, as with GTK+, I lie: there is a way to get a raw key code: `[e keyCode]`. Unfortunately, unlike with Windows and GTK+, this key code table is [not device-independent and not keymap-independent](http://stackoverflow.com/questions/3202629/where-can-i-find-a-list-of-mac-virtual-key-codes): they are keyboard layout-specific, and do not appear to have been updated since the pre-Cocoa days. Indeed, neither Carbon nor Cocoa provide a definite list in their documentation, and the Xcode 5.0 version of Carbon's `HIToolbox/Events.h` header (`/System/Library/Frameworks/Carbon.framework/Versions/A/Frameworks/HIToolbox.framework/Versions/A/Headers/Events.h`) even says
```
 *    keyboard. Those constants with "ANSI" in the name are labeled
 *    according to the key position on an ANSI-standard US keyboard.
 *    For example, kVK_ANSI_A indicates the virtual keycode for the key
 *    with the letter 'A' in the US keyboard layout. Other keyboard
 *    layouts may have the 'A' key label on a different physical key;
 *    in this case, pressing 'A' will generate a different virtual
 *    keycode.
```
The only thing we can guarantee is that the Cocoa and Carbon codes are the same (as the documentation does guarantee this). (If you ever wondered why GLFW talks about US English keyboards in its virtual keycode docs... this is why. Yes, [GLFW interprets the raw key codes](https://github.com/glfw/glfw/blob/master/src/cocoa_window.m#L599).)

So. `NSEvent` provides two methods for getting character data:
- `[e characters]`, which just returns a string with characters if not dead
- `[e charactersIgnoringModifiers]` which bypasses Mac OS X's Option-key IME:
	- let's say we press Option-E to dead-key a tilde
	- `[e characters]` returns an empty `NSString`
	- `[e charactersIgnoringModifiers]` returns `@"E"`

Thankfully there IS a way to get keys that aren't printable characters! ...Mac OS X steals the private use area block of the Unicode BMP and reserves it for [its own virtual key codes](https://developer.apple.com/library/mac/documentation/cocoa/Reference/ApplicationKit/Classes/NSEvent_Class/Reference/Reference.html#//apple_ref/doc/uid/20000016-SW136); the first character of a one-character return from the `characters` methods will have a Unicode code point value equal to these. ([There's also these non-graphic characters constants.](https://developer.apple.com/library/mac/documentation/cocoa/Reference/ApplicationKit/Classes/NSText_Class/Reference/Reference.html#//apple_ref/doc/uid/20000367-SW46))

(Technically you're supposed to send incoming key events to `[self interpretKeyEvents:]`, which will generate a bunch of text-related method calls to make things easier, but we don't have to. Technically you're also supposed to use key equivalents, but that doesn't apply here...)

For modifier keys pressed by themselves, neither `keyDown:` nor `keyUp:` appears to be sent; we need to handle `flagsChanged:` (if I'm reading this right, anyway). Whatever the case, `[e modifierFlags]` will always be valid. In `flagsChanged:`, `characters` will **NOT** be valid and **WILL** throw an exception.

There's also `[e isARepeat]`, which tells us whether a key was repeated; it does not say how many times. (*TODO* does this mean `keyDown:` is sent multiple times?)

TODO
* Is there a way to differentiate Return and Enter? There doesn't seem to be a way to differentiate other left/right keys...
* Does `charactersIgnoringModifiers` always work, or only if `characters` would otherwise indicate a dead key?

### General TODOs
* What happens if I hold down a key, then switch programs with the mouse and release the key? If I decide to intercept modifiers and hand them out like with mouse events, this will be an issue. (Otherwise, I could just poll them each time, like with mouse events; on Windows and Cocoa this will work (GLFW seems to do this on Windows anyway) but I'm not sure about GTK+.)
* How is Shift handled in Windows character events?
* Figure out which keys we can provide and which we can't...

### Consensus??
```go
type KeyEvent struct {
	// TODO some key representation
	// maybe
	ApproximateKey	int
		// maps to a definite key on Windows and GTK+
		// on Mac OS X, charactersIgnoringModifiers is used instead
		// this means we ignore IME completely, which isn't optimal, but.
		// if this is 0, a Modifier was pressed by itself
		// there is no way to differentiate left and right keys

	Modifiers			Modifiers

	Handled			chan<- bool
		// must return on this channel due to Windows system keys
}
// also note: add Super to Modifiers
```

### Er wait oops
I forgot I wanted to make a tracker, whose input should in theory be layout independent; if we do the above we can't do this... we would need to use the key codes and hope key codes are keymap-dependent on both Windows and GTK+...

yeah

and guess what? they're keymap-independent on Windows and and likely so on GTK+ (thanks to tristan and LRN in irc.gimp.net/#gtk+ and exDM69 in irc.efnet.net/#winprog). Guess I'll need to write a quick test...<br>
A on keyboard; US English QWERTY - keyval:0x61<br>
A on keyboard; Georgian AZERTY Tskapo (A = ქ) - keyval:0x10010e5<br>
so.

So this leaves character-based input as the only real option. The only two questions that remain are:
- on Windows is there a reliable way to tell if a `WM_CHAR` DOES come up for the next code? `WM_UNICHAR`?
- related: will both a `WM_CHAR` and a `WM_UNCIHAR` ever come up for the same keystroke?
- how DO you load an existing `GtkIMContext`?
- related: will Cocoa's `charactersIgnoringModifiers` *always* ignore modifiers?

---------------------

<a name="sigh"></a>
Actually the real question is: is it possible to just get ONE domain of keyboard input on all platforms? GDK has constants for every possible language... so someone not using a Latin-based keyboard will wind up having their keystorkes rejected by the `GdkDrawingArea`...
>Keyboard input **MUST** be well defined, and it must be well defined **NOW**. As the author of the GUI library, I **MUST** guarantee that someone typing a character on the same physical machine on different operating systems each with the same keyboard layout gets the exact same response (with no unwanted side effects), and by extension that the programmer sees the same thing. But things are just different enough to screw this up.

Approach | Windows | GTK+ | Mac OS X
----- | ----- | ----- | -----
Virtual key code mapping | Adjusted by layout | Adjusted by layout | NOT adjusted by layout
Virtual key code range | Limited to physical keys on conventional keyboards; outside drivers do IME | NOT limited thus; virtual keycodes exist for languages | Limited to physical keys on conventional keyboards; OS provides IME facilities
Character translation at all | Provided by OS, but not sure about some behavioral details (*TODO*) | multiple; see the GTK+ section above; each problematic | Provided by OS; escape hatches available
Single-keystroke character translation | *TODO* | Constants exist for whatever keyboard layout you can imagine | *TODO*
Multi-keystroke character translation | `WM_DEADCHAR`/`WM_UNICHAR` | (see GTK+ s ectiona bove for issues) | Provided by OS; escape hatches avialable
Character translation ignoring input language (so the programmer can know that the A key was pressed regardless of language) | *TODO* | *TODO* | *maybe* `charactersIgnoringModifiers`? *TODO*

### OK new consensus
Windows: use virtual key codes<br>
OS X: use `charactersIgnoringModifiers`

```go
(on the Area comment)
// Do not use an Area if you intend to read text.
// Due to platform differences regarding text input,
// keyboard events have beem compromised in
// such a way that attempting to read Unicode data
// in platform-native ways is painful.
// [Use TextArea instead, providing a TextAreaHandler.]
(corollary: this means multi-line edits will need a different name; corollary 2: this means I will have to start providing font resource acquisition, something I didn't want to do (I wanted to relegate that to freetype-go or similar, assuming that was even capable of doing so))

// A KeyEvent represents a keypress in an Area.
// 
// KeyEvent has been designed to be predictable.
// As the different operating systems supported by package ui
// expose wildly different APIs and rules for reading keystrokes,
// this means that KeyEvent has certain rules and restrictions
// that you must mind. This makes KeyEvent unsuitable
// for reading text (as Area's comment will say).
// As another consequence, no KeyEvent will be generated if
// package ui cannot portably report a given key. Supported
// keys are described in the comments for the Rune field and
// the ExtKey and Modifiers types. Package ui will act as if
// false was sent on Handled, so these ignored keypresses are
// sent back to the operating system for handling.
type KeyArea struct {
	// Rune contains a lowercase rune specifying the name
	// of the key pressed that triggered the event.
	// Ideally, this would generally correspond to
	// the raw character pressed (so there would be two
	// events 'k', 'a' instead of 'か', if Japanese characters are
	// input that way on a given machine). This will hold true
	// on systems where IME returns are separate from
	// keypress codes. On other systems, an attempt has
	// been made to map backwards based on information
	// that can be provided in the most portable (if the
	// system provides multiple of its own backends) way.
	// See [TODO] for a list of Rune values that are guaranteed
	// to be available. There is no way to differentiate between
	// multiple differnet Keys with the same name (for instance,
	// there is no way to differentiate between '1' on the typewriter
	// section of a standard 101-key keyboard and '1' on the numeric
	// keypad section). Furthermore, note that Rune's value does not
	// [necessarily? TODO] indicate a physical position on the keyboard
	// (for instance, 'a' is returned when pressing A on both QWERTY
	// and AZERTY keyboards, not when pressing the key that would be
	// A on QWERTY keyboards on all layouts).
	// If this field is zero, see ExtKey.
	// [TODO: how do we handle numlock and capslock?]
	Rune			rune

	// If Rune is zero, ExtKey contains a predeclared identifier
	// naming an extended key. See ExtKey for details.
	// If both Rune and ExtKey are zero, a Modifier by itself
	// was pressed. Rune and ExtKey will not both be nonzero.
	ExtKey		ExtKey

	Modifiers		Modifiers

	// If Up is true, the key was released; if not, the key was pressed.
	// There is no guarantee that all pressed keys shall have
	// corresponding release events (for instance, if the user switches
	// programs while holding the key down, then releases the key).
	// Keys that have been held down are reported as multiple
	// key press events.
	Up			bool

	// When you are finished processing the incoming event,
	// send whether or not you did something in response
	// to the given keystroke over Handled. If you send false,
	// you indicate that you did not handle the keypress, and
	// that the system should handle it instead. (Some systems
	// will stop processing the keyboard event at all if you return
	// true unconditionally, which may result in unwanted behavior
	// like global task-switching keystrokes not being processed.)
	// Only one value may be sent on Handled. Do not close Handled;
	// the package will do it for you.
	Handled		chan<- bool
}

// ExtKey represents keys that do not have a Rune representation.
// There is no way to differentiate between left and right ExtKeys.
type ExtKey uintptr
const (
	keyname1 ExtKey = iota
	keyname2
	keyname3
)
(notes: this will have to be produced based on what's available on each platform (Mac OS X might be the biggest filter here); also as a personal favor-the-user decision Print Screen shall not be supported)
```

If I ever intend on providing alternate text-based widgets, I will need to use `GtkTextArea` and `NSTextArea` to make things work the most fluidly, so this will require another type. Woo...

Also this answers the what if a key has been held down and switches away from the program question: Windows does not send a key up.

This just leaves the GTK+ geometry mapping: there is a way to do it if X11 is the only supported backend, but Wayland exists...
    [12:29] <ebassi> yes, you can assume that they are the same
(irc.gimp.net/#gtk+) ok; that works too I guess
let's go!

### And finally... COMMON KEYS
I meant Windows might be the biggest filter because it has a concrete list of virtual key codes but meh

Windows | GDK | [Cocoa if known]
----- | ----- | -----
VK_BACK (0x08) - BACKSPACE key | GDK_KEY_BackSpace | either NSBackspaceCharacter or NSDeleteCharacter (TODO)
VK_TAB (0x09) - TAB key | GDK_KEY_Tab | NSTabCharacter; TODO do we also need to handle NSBackTabCharacter (Shift+Tab)?
VK_RETURN (0x0D) - ENTER key (note: also for numeric keypad) | [TODO either GDK_KEY_Linefeed or GDK_KEY_Return] and GDK_KEY_KP_Enter | either NSNewlineCharacter or NSCarriageReturnCharacter (TODO) (TODO also NSEnterCharacter?)
VK_SHIFT (0x10)/VK_LSHIFT (0xA0) - Left SHIFT key/VK_RSHIFT (0xA1) - Right SHIFT key - SHIFT key | [modifier] | [modifier]
VK_CONTROL (0x11)/VK_LCONTROL (0xA2) - Left CONTROL key/VK_RCONTROL (0xA3) - Right CONTROL key - CTRL key | [modifier] | [modifier]
VK_MENU (0x12)/VK_LMENU (0xA4) - Left MENU key/VK_RMENU (0xA5) - Right MENU key - ALT key | [modifier] | [modifier]
VK_PAUSE (0x13) - PAUSE key | GDK_KEY_Pause | NSPauseFunctionKey and NSBreakFunctionKey
VK_CAPITAL (0x14) - CAPS LOCK key | TODO either GDK_KEY_Caps_Lock or GDK_KEY_Shift_Lock | [modifier NSAlphaShiftKeyMask]
VK_ESCAPE (0x1B) - ESC key | GDK_KEY_Escape
VK_SPACE (0x20) - SPACEBAR | GDK_KEY_space
VK_PRIOR (0x21) - PAGE UP key | GDK_KEY_Page_Up | NSPageUpFunctionKey
VK_NEXT (0x22) - PAGE DOWN key | GDK_KEY_Page_Down | NSPageDownFunctionKey
VK_END (0x23) - END key | GDK_KEY_End | NSEndFunctionKey
VK_HOME (0x24) - HOME key | GDK_KEY_Home | NSHomeFunctionKey
VK_LEFT (0x25) - LEFT ARROW key | GDK_KEY_Left | NSLeftArrowFunctionKey
VK_UP (0x26) - UP ARROW key | GDK_KEY_Up | NSUpArrowFunctionKey
VK_RIGHT (0x27) - RIGHT ARROW key | GDK_KEY_Right | NSRightArrowFunctionKey
VK_DOWN (0x28) - DOWN ARROW key | GDK_KEY_Down | NSDownArrowFunctionKey
VK_INSERT (0x2D) - INS key | GDK_KEY_Insert | NSInsertFunctionKey (TODO also intercept the key that replaced it on newer Mac keyboards?)
VK_DELETE (0x2E) - DEL key | GDK_KEY_Delete (TODO really this one?) | either NSDeleteFunctionKey or NSDeleteCharacter (TODO)
0x30 - 0 key | GDK_KEY_0 and GDK_KEY_parenright (TODO how will we handle Rune for shifted symbols?)
0x31 - 1 key | GDK_KEY_1 and GDK_KEY_exclam
0x32 - 2 key | GDK_KEY_2 and GDK_KEY_at
0x33 - 3 key | GDK_KEY_3 and GDK_KEY_numbersign
0x34 - 4 key | GDK_KEY_4 and GDK_KEY_dollar (TODO how will this Rune be done on non-American setups?)
0x35 - 5 key | GDK_KEY_5 and GDK_KEY_percent
0x36 - 6 key | GDK_KEY_6 and GDK_KEY_asciicircum (TODO really this one?)
0x37 - 7 key | GDK_KEY_7 and GDK_KEY_ampersand
0x38 - 8 key | GDK_KEY_8 and GDK_KEY_asterisk
0x39 - 9 key | GDK_KEY_9 and GDK_KEY_parenleft
0x41 - A key | GDK_KEY_A and GDK_KEY_a
0x42 - B key | ...
0x43 - C key | ...
0x44 - D key | ...
0x45 - E key | ...
0x46 - F key | ...
0x47 - G key | ...
0x48 - H key | ...
0x49 - I key | ...
0x4A - J key | ...
0x4B - K key | ...
0x4C - L key | ...
0x4D - M key | ...
0x4E - N key | ...
0x4F - O key | ...
0x50 - P key | ...
0x51 - Q key | ...
0x52 - R key | ...
0x53 - S key | ...
0x54 - T key | ...
0x55 - U key | ...
0x56 - V key | ...
0x57 - W key | ...
0x58 - X key | ...
0x59 - Y key | ...
0x5A - Z key | ...
VK_LWIN (0x5B) - Left Windows key (Natural keyboard) | [modifier] | [modifier]
VK_RWIN (0x5C) - Right Windows key (Natural keyboard) | [modifier] | [modifier]
VK_APPS (0x5D) - Applications key (Natural keyboard) (TODO which key is this? the right click shortcut key?) | TODO | TODO if this is the right click key, is it NSMenuFunctionKey?
VK_NUMPAD0 (0x60) - Numeric keypad 0 key | GDK_KEY_KP_0 (TODO really this?)
VK_NUMPAD1 (0x61) - Numeric keypad 1 key | ...
VK_NUMPAD2 (0x62) - Numeric keypad 2 key | ...
VK_NUMPAD3 (0x63) - Numeric keypad 3 key | ...
VK_NUMPAD4 (0x64) - Numeric keypad 4 key | ...
VK_NUMPAD5 (0x65) - Numeric keypad 5 key | ...
VK_NUMPAD6 (0x66) - Numeric keypad 6 key | ...
VK_NUMPAD7 (0x67) - Numeric keypad 7 key | ...
VK_NUMPAD8 (0x68) - Numeric keypad 8 key | ...
VK_NUMPAD9 (0x69) - Numeric keypad 9 key | ...
VK_MULTIPLY (0x6A) - Multiply key | GDK_KEY_KP_Multiply (TODO really this one?)
VK_ADD (0x6B) - Add key | GDK_KEY_KP_Add (TODO really this one?)
VK_SUBTRACT (0x6D) - Subtract key | GDK_KEY_KP_Subtract (TODO really this one?)
VK_DECIMAL (0x6E) - Decimal key (. on numeric keypad) | GDK_KEY_KP_Decimal
VK_DIVIDE (0x6F) - Divide key | GDK_KEY_KP_Divide (TODO really this one?)
VK_F1 (0x70) - F1 key | GDK_KEY_F1 | NSF1FunctionKey
VK_F2 (0x71) - F2 key | ... | ...
VK_F3 (0x72) - F3 key | ... | ...
VK_F4 (0x73) - F4 key | ... | ...
VK_F5 (0x74) - F5 key | ... | ...
VK_F6 (0x75) - F6 key | ... | ...
VK_F7 (0x76) - F7 key | ... | ...
VK_F8 (0x77) - F8 key | ... | ...
VK_F9 (0x78) - F9 key | ... | ...
VK_F10 (0x79) - F10 key | ... | ...
VK_F11 (0x7A) - F11 key | ... | ...
VK_F12 (0x7B) - F12 key | ... | ...
VK_NUMLOCK (0x90) - NUM LOCK key | GDK_KEY_Num_Lock | NSClearLineFunctionKey (according to docs)
VK_SCROLL (0x91) - SCROLL LOCK key | GDK_KEY_Scroll_Lock | NSScrollLockFunctionKey
VK_OEM_1 (0xBA) - Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ';:' key (**TODO** see if the varying will hurt us somehow) | GDK_KEY_semicolon and GDK_KEY_colon
VK_OEM_PLUS (0xBB) - For any country/region, the '+' key | GDK_KEY_plus and GDK_KEY_equal
VK_OEM_COMMA (0xBC) - For any country/region, the ',' key | GDK_KEY_comma and GDK_KEY_less
VK_OEM_MINUS (0xBD) - For any country/region, the '-' key | GDK_KEY_minus and GDK_KEY_underscore
VK_OEM_PERIOD (0xBE) - For any country/region, the '.' key | GDK_KEY_period and GDK_KEY_greater
VK_OEM_2 (0xBF) - Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '/?' key | GDK_KEY_slash and GDK_KEY_question
VK_OEM_3 (0xC0) - Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '`~' key | GDK_KEY_quoteleft (TODO really?) and GDK_KEY_asciitilde
VK_OEM_4 (0xDB) - Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '[{' key | GDK_KEY_bracketleft and GDK_KEY_braceleft
VK_OEM_5 (0xDC) - Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '\|' key | GDK_KEY_backslash and GDK_KEY_bar (TODO really _bar?)
VK_OEM_6 (0xDD) - Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ']}' key | GDK_KEY_bracketright and GDK_KEY_braceright
VK_OEM_7 (0xDE) - Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the 'single-quote/double-quote' key | GDK_KEY_apostrophe and GDK_KEY_apostrophe (TODO really this one?)

Windows keys that will not be handled:
```
VK_LBUTTON (0x01) - Left mouse button
VK_RBUTTON (0x02) - Right mouse button
VK_CANCEL (0x03) - Control-break processing
	TODO - is this generated on Control-C?
	TODO if we handle Pause/Break, this will need to come back
VK_MBUTTON (0x04) - Middle mouse button (three-button mouse)
VK_XBUTTON1 (0x05) - X1 mouse button
VK_XBUTTON2 (0x06) - X2 mouse button
VK_CLEAR (0x0C) - CLEAR key (TODO what key is this) | GDK_KEY_Clear
	I have no idea what key this is...
VK_KANA/VK_HANGUEL/VK_HANGUL (0x15) - IME Kana mode/IME Hanguel mode (maintained for compatibility; use VK_HANGUL)/IME Hangul mode
VK_JUNJA (0x17) - IME Junja mode
VK_FINAL (0x18) - IME final mode
VK_HANJA/VK_KANJI (0x19) - IME Hanja mode/IME Kanji mode
VK_CONVERT (0x1C) - IME convert
VK_NONCONVERT (0x1D) - IME nonconvert
VK_ACCEPT (0x1E) - IME accept
VK_MODECHANGE (0x1F) - IME mode change request
VK_SELECT (0x29) - SELECT key (TODO what key is this?) | TODO
VK_PRINT (0x2A) - PRINT key (TODO what key is this?) | TODO
VK_EXECUTE (0x2B) - EXECUTE key (TODO what key is this?) | TODO
VK_SNAPSHOT (0x2C) - PRINT SCREEN key
VK_HELP (0x2F) - HELP key (TODO what key is this?) | TODO
VK_SLEEP (0x5F) - Computer Sleep key (TODO which key is this?) | TODO
VK_SEPARATOR (0x6C) - Separator key (Shift+Decimal key on some foreign layouts; inserts local thousands separator) | TODO
VK_F13 (0x7C) - F13 key
VK_F14 (0x7D) - F14 key
VK_F15 (0x7E) - F15 key
VK_F16 (0x7F) - F16 key
VK_F17 (0x80) - F17 key
VK_F18 (0x81) - F18 key
VK_F19 (0x82) - F19 key
VK_F20 (0x83) - F20 key
VK_F21 (0x84) - F21 key
VK_F22 (0x85) - F22 key
VK_F23 (0x86) - F23 key
VK_F24 (0x87) - F24 key
0x92 - OEM specific
0x93 - OEM specific
0x94 - OEM specific
0x95 - OEM specific
0x96 - OEM specific
VK_BROWSER_BACK (0xA6) - Browser Back key
VK_BROWSER_FORWARD (0xA7) - Browser Forward key
VK_BROWSER_REFRESH (0xA8) - Browser Refresh key
VK_BROWSER_STOP (0xA9) - Browser Stop key
VK_BROWSER_SEARCH (0xAA) - Browser Search key
VK_BROWSER_FAVORITES (0xAB) - Browser Favorites key
VK_BROWSER_HOME (0xAC) - Browser Start and Home key
VK_VOLUME_MUTE (0xAD) - Volume Mute key
VK_VOLUME_DOWN (0xAE) - Volume Down key
VK_VOLUME_UP (0xAF) - Volume Up key
VK_MEDIA_NEXT_TRACK (0xB0) - Next Track key
VK_MEDIA_PREV_TRACK (0xB1) - Previous Track key
VK_MEDIA_STOP (0xB2) - Stop Media key
VK_MEDIA_PLAY_PAUSE (0xB3) - Play/Pause Media key
VK_LAUNCH_MAIL (0xB4) - Start Mail key
VK_LAUNCH_MEDIA_SELECT (0xB5) - Select Media key
VK_LAUNCH_APP1 (0xB6) - Start Application 1 key
VK_LAUNCH_APP2 (0xB7) - Start Application 2 key
VK_OEM_8 (0xDF) - Used for miscellaneous characters; it can vary by keyboard.
0xE1 - OEM specific
VK_OEM_102 (0xE2) - Either the angle bracket key or the backslash key on the RT 102-key keyboard
	TODO actually SHOULD we handle this one?
0xE3 - OEM specific
0xE4 - OEM specific
VK_PROCESSKEY (0xE5) - IME PROCESS key
0xE6 - OEM specific
VK_PACKET (0xE7) - Used to pass Unicode characters as if they were keystrokes. The VK_PACKET key is the low word of a 32-bit Virtual Key value used for non-keyboard input methods. For more information, see Remark in KEYBDINPUT, SendInput, WM_KEYDOWN, and WM_KEYUP
0xE9 - OEM specific
0xEA - OEM specific
0xEB - OEM specific
0xEC - OEM specific
0xED - OEM specific
0xEE - OEM specific
0xEF - OEM specific
0xF0 - OEM specific
0xF1 - OEM specific
0xF2 - OEM specific
0xF3 - OEM specific
0xF4 - OEM specific
0xF5 - OEM specific
VK_ATTN (0xF6) - Attn key
VK_CRSEL (0xF7) - CrSel key
VK_EXSEL (0xF8) - ExSel key
VK_EREOF (0xF9) - Erase EOF key
VK_PLAY (0xFA) - Play key
VK_ZOOM (0xFB) - Zoom key
VK_NONAME (0xFC) - Reserved
VK_PA1 (0xFD) - PA1 key
VK_OEM_CLEAR (0xFE) - Clear key
```

BONUS OS X NOTE (which means we're right in ignoring numpad differences)
>`NSNumericPadKeyMask`
>Set if any key in the numeric keypad is pressed. The numeric keypad is generally on the right side of the keyboard. This is also set if any of the arrow keys are pressed (NSUpArrowFunctionKey, NSDownArrowFunctionKey, NSLeftArrowFunctionKey, and NSRightArrowFunctionKey). 

For the character keys on Windows, I suppose I could get away with checking all character keys for a WM_CHAR... will need to check, but first, need to find a keyboard layout that has different rules for symbols
* will also need to check what GTK+ and OS X do, in this case

For the locks: I'll need a way to get their state anyway, so simply checking keypresses isn't really effective
* **Caps Lock**: simply checking this and simulating Shift isn't good, because Windows and Unix have Caps Lock act like an inverse Shift, while Mac OS X has it act like a permanent Shift. I might have to add "The case of Rune is undefined; it does not necessarily correspond to the state of the Shift key. Correct code will need to check both uppercase and lowercase forms of letters."
* **Num Lock**: already handled; we're not differentiating between numpad and non-numpad anyway so knowing the Num Lock state is meaningless
* **Scroll Lock**: ...will need to figure out how to read lock state for this one.

Also undetermined:
* what happens on $ on other locales

### Consensus on characters
Mac OS X always sends shifted characters, even with `charactersIgnoringModifiers`. GTK+ has separate keycodes that we can use. This just leaves Windows. It appears what we need to do is process `WM_CHAR` and `WM_SYSCHAR` events, filtering out the characters that we can support.

I wonder if I should change Rune to ASCII because with the filtering we have applied these are guaranteed to be ASCII bytes...

...

```
65/GDK_KEY_A:
equiv 1/3: &ui._Ctype_GdkKeymapKey{keycode:0x26, group:0, level:1}
equiv 2/3: &ui._Ctype_GdkKeymapKey{keycode:0x1, group:38, level:1}
equiv 3/3: &ui._Ctype_GdkKeymapKey{keycode:0x1, group:1, level:38}
97/GDK_KEY_a:
equiv 1/3: &ui._Ctype_GdkKeymapKey{keycode:0x26, group:0, level:0}
equiv 2/3: &ui._Ctype_GdkKeymapKey{keycode:0x0, group:38, level:1}
equiv 3/3: &ui._Ctype_GdkKeymapKey{keycode:0x1, group:0, level:38}
```
oh you have got to be kidding me (Spanish keyboard layout)

OK so it looks like the only thing I can really do is just pretend that the keyboard input routines are perfect and handle the character events on Windows and OS X (and use `charactersIgnoringModifiers` on OS X because otherwise we can't handle the Option key) and HOPE TO GOD ALMIGHTY that all GDK users see the non-IME forms. SCREW THIS.

### Final Consensus, For Fuck's Sake.
```go
// A KeyEvent represents a keypress in an Area.
// 
// In a perfect world, KeyEvent would be 100% predictable.
// Despite my best efforts to do this, however, the various
// differences in input handling between each backend
// environment makes this completely impossible (I can
// work with two of the three identically, but not all three).
// Keep this in mind, and remember that Areas are not ideal
// for text.
// 
// If a key is pressed that is not supported by ASCII, ExtKey,
// or Modifiers, no KeyEvent will be produced and package
// ui will act as if false was sent on Handled.
type KeyEvent struct {
	// ASCII is a byte representing the character pressed.
	// Despite my best efforts, this cannot be trivialized
	// to produce predictable input rules on all OSs, even if
	// I try to handle physical keys instead of equivalent
	// characters. Therefore, what happens when the user
	// inserts a non-ASCII character is undefined (some systems
	// will give package ui the underlying ASCII key and we
	// return it; other systems do not). This is especially important
	// if the given input method uses Modifiers to enter characters.
	// If the parenthesized rule cannot be followed and the user
	// enters a non-ASCII character, it will be ignored (package ui
	// will act as above regarding keys it cannot handle).
	// In general, alphanumeric characters, ',', '.', '+', '-', and the
	// (space) should be available on all keyboards. Other ASCII
	// whitespace keys mentioned below may be available, but
	// mind layout differences.
	// Whether or not alphabetic characters are uppercase or
	// lowercase is undefined, and cannot be determined solely
	// by examining Modifiers for Shift. Correct code should handle
	// both uppercase and lowercase identically.
	// In addition, ASCII will contain
	// - ' ' (space) if the spacebar was pressed
	// - '\t' if Tab was pressed, regardless of Modifiers
	// - '\n' if any Enter/Return key was pressed, regardless of which
	// - '\b' if the typewriter Backspace key was pressed
	// If this value is 0, see ExtKey.
	ASCII	byte

	// ...
}
```

## Tweets mentioned in area.go's relevant comment
I tweeted a tl;dr of the whole debacle documented in detail above:
- https://twitter.com/pgandlabs/status/447790340251344896
- https://twitter.com/pgandlabs/status/447791528237596672
- https://twitter.com/pgandlabs/status/447791774749454336
- https://twitter.com/pgandlabs/status/447791890902286336
- https://twitter.com/pgandlabs/status/447791982992433152
- https://twitter.com/pgandlabs/status/447792201222066177
- https://twitter.com/pgandlabs/status/447792444311371777
- https://twitter.com/pgandlabs/status/447792609474654208
- https://twitter.com/pgandlabs/status/447792724969000960
- https://twitter.com/pgandlabs/status/447792787241832449
- https://twitter.com/pgandlabs/status/447792889620598784

# ...but wait!
When writing the Windows version of the above I got stuck (durr) and looked at the GLFW sources again; then I noticed **it was using scancodes for the typewriter keys**. But isn't that going to hurt when we try foreign keyboards?

[Nope.](http://www.quadibloc.com/comp/scan.htm)

I had found the above during my research above but never went down far enough to notice that it explains that yes, international keyboards DO use the same scancodes as American keyboards. So, given the scancodes on Windows and OS X's own positional key codes, we CAN do positional input after all! ...but we need a solution for GTK+.

We don't need to do the reverse lookup problem I had earlier with the Spanish keyboards as the `GdkEventKey` structure has a `hardware_keycode`. For that there's this (irc.gimp.net/#gtk+):
```
[Saturday, March 22, 2014] [12:22:04 PM] <andlabs> well I keep running into problems trying to just see how GdkKeymapKey works in GTK+/Wayland because I can't get Wayland working properly in VMs
[Saturday, March 22, 2014] [12:22:07 PM] Quit nkoep (~nik@koln-4db46460.pool.mediaWays.net) has left this server (Remote closed the connection).
[Saturday, March 22, 2014] [12:22:47 PM] <ebassi> andlabs: you should use a nested wayland inside x11, or run a wayland session
[Saturday, March 22, 2014] [12:22:48 PM] <jjavaholic> flashplugin > firefox > nvidia > libc > xorg ?
[Saturday, March 22, 2014] [12:23:11 PM] <andlabs> all right, thanks
[Saturday, March 22, 2014] [12:23:30 PM] <ebassi> jjavaholic: no. more like: flashplugin > xorg > nvidia; firefox > xorg > nvidia; firefox/flashplugin/xorg > libc
[Saturday, March 22, 2014] [12:23:43 PM] Join nkoep (~nik@koln-4db46460.pool.mediaWays.net) has joined this channel.
[Saturday, March 22, 2014] [12:28:41 PM] <andlabs> ebassi: so am I to assume that the hardware_keycodes in wayland are the smae as the hardware_keycodes in X? t hat's my main thing
[Saturday, March 22, 2014] [12:28:52 PM] <andlabs> if they are different I'd need to find out how
[Saturday, March 22, 2014] [12:29:00 PM] <ebassi> yes, you can assume that they are the same
[Saturday, March 22, 2014] [12:29:22 PM] <andlabs> ok, thanks
[Saturday, March 22, 2014] [12:30:35 PM] Quit nkoep (~nik@koln-4db46460.pool.mediaWays.net) has left this server (Remote closed the connection).
[Saturday, March 22, 2014] [12:31:51 PM] Quit gauteh (~gauteh@cD572BF51.dhcp.as2116.net) has left this server (Ping timeout: 600 seconds).
[Saturday, March 22, 2014] [12:31:56 PM] <jjavaholic> system profiler libc contents: http://i.imgur.com/herHBMU.png 
[Saturday, March 22, 2014] [12:32:36 PM] <jjavaholic> system profiler nvidia-304 contents: http://i.imgur.com/3NrkFNQ.png
[Saturday, March 22, 2014] [12:32:39 PM] <ebassi> andlabs: both x11 and wayland on linux get the events from the kernel evdev interface, so you get the same hw key codes
[Saturday, March 22, 2014] [12:32:59 PM] <ebassi> andlabs: also, both x11 and wayland use xkbcommon to handle key maps
[Saturday, March 22, 2014] [12:33:29 PM] <andlabs> cool, thanks
```
So I know that the codes are standard between X11 and Wayland, and that on Linux the codes are provided by evdev (which, as we will see later, uses standard codes). The question now is: is there a portable way to get the physical keys?

GLFW on Unix [uses the full XKB and its key names to do the job](https://github.com/glfw/glfw/blob/master/src/x11_init.c#L245). libxkbcommon doesn't seem to be able to do this, and we would need an X11 display to map back anyway, so that doesn't work. And when going through the GDK source to see how it fills in `hardware_keycode`, I notice it uses more XKB features that aren't part of libxkbcommon, like `struct xkb_desc` (even in [a Wayland patch](https://mail.gnome.org/archives/commits-list/2011-January/msg02064.html)). So these are useless.

There is one important thing that [we do get from both libxkbcommon](http://xkbcommon.org/doc/current/xkbcommon_8h.html#ac29aee92124c08d1953910ab28ee1997) and assorted other X docs and source code, though: X11 keycodes start at 8, so for the evdev case, the `hardware_keycode` is the evdev code + 8. (This latter part is in the libxkbcommon link in this paragraph.)

At this point I decided to figure out what `GdkEventKey.hardware_keycode` meant on each Unix variant, and ran a few tests to get some sample values of the `GdkEventKey.hardware_keycode`. Then I set up a FreeBSD 32-bit setup and tried the tests again... and got the same values. Huh?????

I decided to dig into the X.org source to see what was making key codes. There are two keyboard drivers of interest here:
<ul><li>xf86-input-keyboard, the generic keyboard driver: [relevant code](http://cgit.freedesktop.org/xorg/driver/xf86-input-keyboard/tree/src/kbd.c#n411)
<li>xf86-input-evdev, the Linux evdev client driver: [relevant code](http://cgit.freedesktop.org/xorg/driver/xf86-input-evdev/tree/src/evdev.c#n330)</ul>
The most that I could gather after going through this and jumping around the rest of the source tree a few times is that the xf86-input-keyboard driver simply passed scancodes up to the client program. (Indeed, pressing the 1 key on the keyboard produced 10, which is scancode 0x02 + 8, something I noticed during the FreeBSD tests.)

So the question is: are the scancodes the same as evdev's keyboard values for the typewriter keys? Those are the only ones we need; the other keys can be determined from their Windows virtual key codes and GDK key codes. I wrote a program to take [evdev's well-defined key equivalents](https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/input.h#n201) (this is also on your local Linux machine, inthe `<linux/input.h>` header file, under "Keys and buttons") and the same set of scancodes (from the Scan Codes Demystified article linked above) and... **they do match**!

<b><i><u>We're in the clear for positional keyboard input!</u></i></b> :D

There are a few loose ends:
* Can the OS X key codes be used to differentiate between numpad keys and non-numpad keys regardless of num lock state? If so, we can safely differentiate between the two, and can get rid of that arbitrary restriction.
	* TODO
* Can we also use scancodes for the numeric keypad, **including** the numeric keypad / key? GDK keysyms have Num Lock interpreted; we don't want that. This is just adding the scancodes for the numeric keypad to our test above...
	* We only need to worry about the number keys and ., it seems; everything else is unaffected by Num Lock (and uses extended scancodes, so.)
	* Okay, the numeric keys and . use the same scancodes and evdev key code values, so we can handle them with scancodes too.
* The GLFW source does not use the scancode 0x2B for \, claiming that it only exists on US keyboards (instead it uses one of the OEM virtual key codes on Windows). This goes against the Scan Codes Demystified page, which says that on international keyboards, that would be another key (with region-specific label) underneath and to the right of what would be the [ and ] keys on a US keyboard. This appears to be true in some cases; in others, the extra key is to the left of Backspace instead. Either way, this is close enough to the \ key's position on a US keyboard that we can just go ahead and use 0x2B anyway.
