# Area
An Area is a blank-slate control. It will have the following capabilities when all is said and done:
- Events:
	- Paint ( cliprect image.Rectangle -- img *image.NRGBA )<br>Called when a paint event comes in; expects to be fed an image to draw back. cliprect is used to constrain the drawing to only what's needed.
	- Mouse ( e MouseEvent -- )<br>Called when a mouse event (motion/click/drag/etc.) comes in.
	- Key ( e KeyEvent -- )<br>Called when a keyboard event comes in.
- Accessors/Mutators
	- SetSize(width, height int) / Size() (width, height int)<br>Determines the internal size of the Area. The actual size of the control may be different; scrolling is handled automatically.
	- ScrollTo(x, y int)<br>Scrolls the Area such that the given point is visible.
- Must satisfy the Control interface, which contains unexported methods.

The brick in the window here is Paint, which stops the world until it gets something back.

## Attempt 1: Event Channels
The original implementation of Area was an opaque object which implemented Control and used channels to transmit events:
```
type Area struct {
	Paint		chan PaintRequest
	Mouse	chan MouseEvent
}
type PaintRequest {
	Rect		image.Rectangle
	Out		chan<- *image.NRGBA
}
```
and the intended usage was
```
for {
	select {
	case req := <-area.Paint:
		req.Out <- img.SubImage(req.Rect).(*image.NRGBA)
	case e := <-area.Mouse:
		// handle mouse event
	// other events
	}
}
```
Alas, in practice, this failed miserably: though it worked just fine when the only control being worked on was the Area, as soon as we added another control:
```
for {
	select {
	case req := <-area.Paint:
		req.Out <- img.SubImage(req.Rect).(*image.NRGBA)
	case e := <-area.Mouse:
		// handle mouse event
	case time := <-ticker:
		label.SetText(time.String())
	}
}
```
a deadlock occurred: a Paint event would be issued while label.SetText() was being called. The Paint event would not return to the UI main loop until it got something on req.Out, and label.SetText() would not return until the UI main loop fully rexecuted the text change and signaled back that it had done so.

This would not be an issue if the drawing tasks could be detached from the main thread; alas, no OS seems to let you do this. (Mac OS X's NSApplication has a method named detachDrawingThread:toTarget:withObject: that sounds like it does this, but in reality this is a misnomer and it really should be called runSelectorOnNewThread:receiver:withArgument:.)

## Attempt 2: Embedding
This one I couldn't figure out how to do right.

In this case, you would have a custom type that would embed Area and where the events would be overloaded methods:
```
type MyArea struct {
	*ui.Area
	img		*image.NRGBA
}
func newMyArea() *MyArea {
	return &MyArea{
		Area:	ui.NewArea(),
		img:		/* new image here */,
	}
}
func (a *MyArea) Paint(rect image.Rectangle) *image.NRGBA {
	return a.img.SubImage(rect).(*image.NRGBA)
}
```
and unimplemented events (other than Paint, which must be implemented) would do nothing.

When I got this compiled, however, Paint() panicked, as it was still using the ui.Area.Paint() method!

So then I tried the interface approach instead:
```
type Area interface {
	Control			// must be a control
	Paint(image.Rectangle) *image.NRGBA
	Mouse(MouseEvent)
}
```
but I don't know how I'm going to handle the Control embed, as that uses unexported methods. The unexported methods are needed to do the acutal OS-dependent work, and I figured exporting these to the public would result in misuse.

So I'm stuck. Unless there's a way I can get either attempt working, or I can figure out some other Go-like approach... I personally felt weird about Attempt 2, but go.wde does it. I personally don't want to use a delegate type that provides the event functions and have Area require one, as that sounds both not-Go-like and more OOP-like than I want this to be. I already silently reject requests for the callback approach as it is (I should respond to them)...

## Attempt 3: Delegates
Well I did it anyway
```
type Area {
	handler	AreaHandler
	// ...
}
type AreaHandler interface {
	Paint(image.Rectangle) *image.NRGBA
	Mouse(MouseEvent)
}
func NewArea(handler AreaHandler) *Area {
	// ...
}
```
This seems to have fixed the deadlocks!... but I can still deadlock elsewhere, so something /else/ is wrong, sigh...
