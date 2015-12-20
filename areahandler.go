// 13 december 2015

package ui

// #include <stdlib.h>
// #include "ui.h"
// extern void doAreaHandlerDraw(uiAreaHandler *, uiArea *, uiAreaDrawParams *);
// extern void doAreaHandlerMouseEvent(uiAreaHandler *, uiArea *, uiAreaMouseEvent *);
// extern void doAreaHandlerMouseCrossed(uiAreaHandler *, uiArea *, int);
// extern void doAreaHandlerDragBroken(uiAreaHandler *, uiArea *);
// extern int doAreaHandlerKeyEvent(uiAreaHandler *, uiArea *, uiAreaKeyEvent *);
// static inline uiAreaHandler *allocAreaHandler(void)
// {
// 	uiAreaHandler *ah;
// 
// 	ah = (uiAreaHandler *) malloc(sizeof (uiAreaHandler));
// 	if (ah == NULL)		// TODO
// 		return NULL;
// 	ah->Draw = doAreaHandlerDraw;
// 	ah->MouseEvent = doAreaHandlerMouseEvent;
// 	ah->MouseCrossed = doAreaHandlerMouseCrossed;
// 	ah->DragBroken = doAreaHandlerDragBroken;
// 	ah->KeyEvent = doAreaHandlerKeyEvent;
// 	return ah;
// }
// static inline void freeAreaHandler(uiAreaHandler *ah)
// {
// 	free(ah);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var areahandlers = make(map[*C.uiAreaHandler]AreaHandler)

// AreaHandler defines the functionality needed for handling events
// from an Area. Each of the methods on AreaHandler is called from
// the GUI thread, and every parameter (other than the Area itself)
// should be assumed to only be valid during the life of the method
// call (so for instance, do not save AreaDrawParams.AreaWidth, as
// that might change without generating an event).
// 
// Coordinates to Draw and MouseEvent are given in points. Points
// are generic, floating-point, device-independent coordinates with
// (0,0) at the top left corner. You never have to worry about the
// mapping between points and pixels; simply draw everything using
// points and you get nice effects like looking sharp on high-DPI
// monitors for free. Proper documentation on the matter is being
// written. In the meantime, there are several referenes to this kind of
// drawing, most notably on Apple's website: https://developer.apple.com/library/mac/documentation/GraphicsAnimation/Conceptual/HighResolutionOSX/Explained/Explained.html#//apple_ref/doc/uid/TP40012302-CH4-SW1
// 
// For a scrolling Area, points are automatically offset by the scroll
// position. So if the mouse moves to position (5,5) while the
// horizontal scrollbar is at position 10 and the horizontal scrollbar is
// at position 20, the coordinate stored in the AreaMouseEvent
// structure is (15,25). The same applies to drawing.
type AreaHandler interface {
	// Draw is sent when a part of the Area needs to be drawn.
	// dp will contain a drawing context to draw on, the rectangle
	// that needs to be drawn in, and (for a non-scrolling area) the
	// size of the area. The rectangle that needs to be drawn will
	// have been cleared by the system prior to drawing, so you are
	// always working on a clean slate.
	// 
	// If you call Save on the drawing context, you must call Release
	// before returning from Draw, and the number of calls to Save
	// and Release must match. Failure to do so results in undefined
	// behavior.
	Draw(a *Area, dp *AreaDrawParams)

	// MouseEvent is called when the mouse moves over the Area
	// or when a mouse button is pressed or released. See
	// AreaMouseEvent for more details.
	// 
	// If a mouse button is being held, MouseEvents will continue to
	// be generated, even if the mouse is not within the area. On
	// some systems, the system can interrupt this behavior;
	// see DragBroken.
	MouseEvent(a *Area, me *AreaMouseEvent)

	// MouseCrossed is called when the mouse either enters or
	// leaves the Area. It is called even if the mouse buttons are being
	// held (see MouseEvent above). If the mouse has entered the
	// Area, left is false; if it has left the Area, left is true.
	// 
	// If, when the Area is first shown, the mouse is already inside
	// the Area, MouseCrossed will be called with left=false.
	// TODO what about future shows?
	MouseCrossed(a *Area, left bool)

	// DragBroken is called if a mouse drag is interrupted by the
	// system. As noted above, when a mouse button is held,
	// MouseEvent will continue to be called, even if the mouse is
	// outside the Area. On some systems, this behavior can be
	// stopped by the system itself for a variety of reasons. This
	// method is provided to allow your program to cope with the
	// loss of the mouse in this case. You should cope by cancelling
	// whatever drag-related operation you were doing.
	// 
	// Note that this is only generated on some systems under
	// specific conditions. Do not implement behavior that only
	// takes effect when DragBroken is called.
	DragBroken(a *Area)

	// KeyEvent is called when a key is pressed while the Area has
	// keyboard focus (if the Area has been tabbed into or if the
	// mouse has been clicked on it). See AreaKeyEvent for specifics.
	// 
	// Because some keyboard events are handled by the system
	// (for instance, menu accelerators and global hotkeys), you
	// must return whether you handled the key event; return true
	// if you did or false if you did not. If you wish to ignore the
	// keyboard outright, the correct implementation of KeyEvent is
	// 	func (h *MyHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// 		return false
	// 	}
	// DO NOT RETURN TRUE UNCONDITIONALLY FROM THIS
	// METHOD. BAD THINGS WILL HAPPEN IF YOU DO.
	KeyEvent(a *Area, ke *AreaKeyEvent) (handled bool)
}

func registerAreaHandler(ah AreaHandler) *C.uiAreaHandler {
	uah := C.allocAreaHandler()
	areahandlers[uah] = ah
	return uah
}

func unregisterAreaHandler(uah *C.uiAreaHandler) {
	delete(areahandlers, uah)
	C.freeAreaHandler(uah)
}

// AreaDrawParams provides a drawing context that can be used
// to draw on an Area and tells you where to draw. See AreaHandler
// for introductory information.
type AreaDrawParams struct {
	// Context is the drawing context to draw on. See DrawContext
	// for how to draw.
	Context		*DrawContext

	// AreaWidth and AreaHeight provide the size of the Area for
	// non-scrolling Areas. For scrolling Areas both values are zero.
	// 
	// To reiterate the AreaHandler documentation, do NOT save
	// these values for later; they can change without generating
	// an event.
	AreaWidth	float64
	AreaHeight	float64

	// These four fields define the rectangle that needs to be
	// redrawn. The system will not draw anything outside this
	// rectangle, but you can make your drawing faster if you
	// also stay within the lines.
	ClipX		float64
	ClipY		float64
	ClipWidth		float64
	ClipHeight	float64
}

//export doAreaHandlerDraw
func doAreaHandlerDraw(uah *C.uiAreaHandler, ua *C.uiArea, udp *C.uiAreaDrawParams) {
	ah := areahandlers[uah]
	a := areas[ua]
	dp := &AreaDrawParams{
		Context:		&DrawContext{udp.Context},
		AreaWidth:	float64(udp.AreaWidth),
		AreaHeight:	float64(udp.AreaHeight),
		ClipX:		float64(udp.ClipX),
		ClipY:		float64(udp.ClipY),
		ClipWidth:		float64(udp.ClipWidth),
		ClipHeight:	float64(udp.ClipHeight),
	}
	ah.Draw(a, dp)
}

// TODO document all these
// 
// TODO note that in the case of a drag, X and Y can be out of bounds, or in the event of a scrolling area, in places that are not visible
type AreaMouseEvent struct {
	X			float64
	Y			float64

	// AreaWidth and AreaHeight provide the size of the Area for
	// non-scrolling Areas. For scrolling Areas both values are zero.
	// 
	// To reiterate the AreaHandler documentation, do NOT save
	// these values for later; they can change without generating
	// an event.
	AreaWidth	float64
	AreaHeight	float64

	Down		uint
	Up			uint
	Count		uint
	Modifiers		Modifiers
	Held			[]uint
}

func appendBits(out []uint, held C.uint64_t) []uint {
	n := uint(1)
	for i := 0; i < 64; i++ {
		if held & 1 != 0 {
			out = append(out, n)
		}
		held >>= 1
		n++
	}
	return out
}

//export doAreaHandlerMouseEvent
func doAreaHandlerMouseEvent(uah *C.uiAreaHandler, ua *C.uiArea, ume *C.uiAreaMouseEvent) {
	ah := areahandlers[uah]
	a := areas[ua]
	me := &AreaMouseEvent{
		X:			float64(ume.X),
		Y:			float64(ume.Y),
		AreaWidth:	float64(ume.AreaWidth),
		AreaHeight:	float64(ume.AreaHeight),
		Down:		uint(ume.Down),
		Up:			uint(ume.Up),
		Count:		uint(ume.Count),
		Modifiers:		Modifiers(ume.Modifiers),
		Held:		make([]uint, 0, 64),
	}
	me.Held = appendBits(me.Held, ume.Held1To64)
	ah.MouseEvent(a, me)
}

//export doAreaHandlerMouseCrossed
func doAreaHandlerMouseCrossed(uah *C.uiAreaHandler, ua *C.uiArea, left C.int) {
	ah := areahandlers[uah]
	a := areas[ua]
	ah.MouseCrossed(a, tobool(left))
}

//export doAreaHandlerDragBroken
func doAreaHandlerDragBroken(uah *C.uiAreaHandler, ua *C.uiArea) {
	ah := areahandlers[uah]
	a := areas[ua]
	ah.DragBroken(a)
}

// TODO document all these
type AreaKeyEvent struct {
	Key		rune
	ExtKey	ExtKey
	Modifier	Modifiers
	Modifiers	Modifiers
	Up		bool
}

//export doAreaHandlerKeyEvent
func doAreaHandlerKeyEvent(uah *C.uiAreaHandler, ua *C.uiArea, uke *C.uiAreaKeyEvent) C.int {
	ah := areahandlers[uah]
	a := areas[ua]
	ke := &AreaKeyEvent{
		Key:			rune(uke.Key),
		ExtKey:		ExtKey(uke.ExtKey),
		Modifier:		Modifiers(uke.Modifier),
		Modifiers:		Modifiers(uke.Modifiers),
		Up:			tobool(uke.Up),
	}
	return frombool(ah.KeyEvent(a, ke))
}

// TODO document
// 
// Note: these must be numerically identical to their libui equivalents.
type Modifiers uint
const (
	Ctrl Modifiers = 1 << iota
	Alt
	Shift
	Super
)

// TODO document
// 
// Note: these must be numerically identical to their libui equivalents.
type ExtKey int
const (
	Escape ExtKey = iota + 1
	Insert			// equivalent to "Help" on Apple keyboards
	Delete
	Home
	End
	PageUp
	PageDown
	Up
	Down
	Left
	Right
	F1			// F1..F12 are guaranteed to be consecutive
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	N0			// numpad keys; independent of Num Lock state
	N1			// N0..N9 are guaranteed to be consecutive
	N2
	N3
	N4
	N5
	N6
	N7
	N8
	N9
	NDot
	NEnter
	NAdd
	NSubtract
	NMultiply
	NDivide
)
