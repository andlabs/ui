// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type label struct {
	*controlSingleObject
}

func newLabel(text string) Label {
	l := &label{
		controlSingleObject:        newControlSingleObject(C.newLabel()),
	}
	l.SetText(text)
	return l
}

func (l *label) Text() string {
	return C.GoString(C.textfieldText(l.id))
}

func (l *label) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textfieldSetText(l.id, ctext)
}

/*TODO
func (l *label) commitResize(c *allocation, d *sizing) {
	if !l.standalone && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.neighborAlign.baseline != 0 { // no adjustment needed if the given control has no baseline
			// in order for the baseline value to be correct, the label MUST BE AT THE HEIGHT THAT OS X WANTS IT TO BE!
			// otherwise, the baseline calculation will be relative to the bottom of the control, and everything will be wrong
			origsize := C.controlPreferredSize(l._id)
			c.height = int(origsize.height)
			newrect := C.struct_xrect{
				x:      C.intptr_t(c.x),
				y:      C.intptr_t(c.y),
				width:  C.intptr_t(c.width),
				height: C.intptr_t(c.height),
			}
			ourAlign := C.alignmentInfo(l._id, newrect)
			// we need to find the exact Y positions of the baselines
			// fortunately, this is easy now that (x,y) is the bottom-left corner
			thisbasey := ourAlign.rect.y + ourAlign.baseline
			neighborbasey := d.neighborAlign.rect.y + d.neighborAlign.baseline
			// now the amount we have to move the label down by is easy to find
			yoff := neighborbasey - thisbasey
			// and we just add that
			c.y += int(yoff)
		}
		// in the other case, the most correct thing would be for Label to be aligned to the alignment rect, but I can't get this working, and it looks fine as it is anyway
	}
	basecommitResize(l, c, d)
}
*/
