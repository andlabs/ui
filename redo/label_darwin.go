// 16 july 2014

package ui

// #include "objc_darwin.h"
import "C"

// cheap trick
type label struct {
	*textField
	standalone			bool
	supercommitResize		func(c *allocation, d *sizing)
}

func finishNewLabel(text string, standalone bool) *label {
	l := &label{
		textField:		finishNewTextField(C.newLabel()),
		standalone:	standalone,
	}
	l.SetText(text)
	l.supercommitResize = l.fcommitResize
	l.fcommitResize = l.labelcommitResize
	return l
}

func newLabel(text string) Label {
	return finishNewLabel(text, false)
}

func newStandaloneLabel(text string) Label {
	return finishNewLabel(text, true)
}

func (l *label) labelcommitResize(c *allocation, d *sizing) {
	if !l.standalone && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.neighborAlign.baseline != 0 {		// no adjustment needed if the given control has no baseline
			// in order for the baseline value to be correct, the label MUST BE AT THE HEIGHT THAT OS X WANTS IT TO BE!
			// otherwise, the baseline calculation will be relative to the bottom of the control, and everything will be wrong
			origsize := C.controlPrefSize(l.id)
			c.height = int(origsize.height)
			newrect := C.struct_xrect{
				x:		C.intptr_t(c.x),
				y:		C.intptr_t(c.y),
				width:	C.intptr_t(c.width),
				height:	C.intptr_t(c.height),
			}
			ourAlign := C.alignmentInfo(l.id, newrect)
			// we need to find the exact Y positions of the baselines
			// fortunately, this is easy now that (x,y) is the bottom-left corner
			thisbasey := ourAlign.rect.y + ourAlign.baseline
			neighborbasey := d.neighborAlign.rect.y + d.neighborAlign.baseline
			// now the amount we have to move the label down by is easy to find
			yoff := neighborbasey - thisbasey
			// and we just add that
			c.y += int(yoff)
		}
		// TODO if there's no baseline, the alignment should be to the top /of the alignment rect/, not the frame
	}
	l.supercommitResize(c, d)
}
