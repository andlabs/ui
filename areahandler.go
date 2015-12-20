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
// 	if (ah == NULL)
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
// from an Area.
type AreaHandler interface {
	// TODO document all these
	Draw(a *Area, dp *AreaDrawParams)
	MouseEvent(a *Area, me *AreaMouseEvent)
	MouseCrossed(a *Area, left bool)
	DragBroken(a *Area)
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

// AreaDrawParams defines the TODO.
type AreaDrawParams struct {
	// TODO document all these
	Context		*DrawContext
	AreaWidth	float64
	AreaHeight	float64
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
type AreaMouseEvent struct {
	X			float64
	Y			float64
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
