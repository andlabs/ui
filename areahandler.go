// 13 december 2015

package ui

// #include "ui.h"
// extern void doAreaHandlerDraw(uiAreaHandler *, uiArea *, uiAreaDrawParams *);
// static inline void uiAreaHandler *allocAreaHandler(void)
// {
// 	uiAreaHandler *ah;
// 
// 	ah = (uiAreaHandler *) malloc(sizeof (uiAreaHandler));
// 	if (ah == NULL)
// 		return NULL;
// 	ah->Draw = doAreaHandlerDraw;
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
}

func registerHandler(ah AreaHandler) *C.uiAreaHandler {
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
	ClientWidth	float64
	ClientHeight	float64
	ClipX		float64
	ClipY		float64
	ClipWidth		float64
	ClipHeight	float64
	HScrollPos	int
	VScrollPos	int
}

// export doAreaHandlerDraw
func doAreaHandlerDraw(uah *C.uiAreaHandler, ua *C.uiArea, udp *C.uiAreaDrawParams) {
	ah := areahandlers[uah]
	a := areas[ua]
	dp := &AreaDrawParams{
		Context:		&DrawContext{udp.Context},
		ClientWidth:	float64(udp.ClientWidth),
		ClientHeight:	float64(udp.ClientHeight),
		ClipX:		float64(udp.ClipX),
		ClipY:		float64(udp.ClipY),
		ClipWidth:		float64(udp.ClipWidth),
		ClipHeight:	float64(udp.ClipHeight),
		HScrollPos:	int(udp.HScrollPos),
		VScrollPos:	int(udp.VScrollPos),
	}
	ah.Draw(a, dp)
}
