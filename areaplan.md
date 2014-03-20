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

For modifier keys pressed by themselves, neither `keyDown:` nor `keyUp:` appears to be sent; we need to handle `flagsChanged:` (if I'm reading this right, anyway). Whatever the case, `[e modifierFlags]` will always be valid.

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
