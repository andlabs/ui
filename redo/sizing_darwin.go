// 1 march 2014

package ui

// #include "objc_darwin.h"
import "C"

type sizing struct {
	sizingbase

	// for size calculations
	// nothing for mac

	// for the actual resizing
//	neighborAlign		C.struct_xalignment
}

// THIS IS A GUESS. TODO.
// The only indication that this is remotely correct is the Auto Layout Guide implying that 12 pixels is the "Aqua space".
const (
	macXMargin = 12
	macYMargin = 12
	macXPadding = 12
	macYPadding = 12
)

func (w *window) beginResize() (d *sizing) {
	d = new(sizing)
	if w.spaced {
		d.xmargin = macXMargin
		d.ymargin = macYMargin
		d.xpadding = macXPadding
		d.ypadding = macYPadding
	}
	return d
}

func (w *window) endResize(d *sizing) {
	// redraw
}

func (w *window) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	for _, a := range allocations {
		// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
		// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
		a.y = (winheight - a.y) - a.height
	}
}

func (w *widgetbase) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
        return []*allocation{&allocation{
                x:      x,
                y:      y,
                width:  width,
                height: height,
                this:   w,
        }}
}

func (w *widgetbase) commitResize(c *allocation, d *sizing) {
//TODO
/*
	if s.ctype == c_label && !s.alternate && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.neighborAlign.baseline != 0 {		// no adjustment needed if the given control has no baseline
			// in order for the baseline value to be correct, the label MUST BE AT THE HEIGHT THAT OS X WANTS IT TO BE!
			// otherwise, the baseline calculation will be relative to the bottom of the control, and everything will be wrong
			origsize := C.controlPrefSize(s.id)
			c.height = int(origsize.height)
			newrect := C.struct_xrect{
				x:		C.intptr_t(c.x),
				y:		C.intptr_t(c.y),
				width:	C.intptr_t(c.width),
				height:	C.intptr_t(c.height),
			}
			ourAlign := C.alignmentInfo(s.id, newrect)
			// we need to find the exact Y positions of the baselines
			// fortunately, this is easy now that (x,y) is the bottom-left corner
			thisbasey := ourAlign.alignmentRect.y + ourAlign.baseline
			neighborbasey := d.neighborAlign.alignmentRect.y + d.neighborAlign.baseline
			// now the amount we have to move the label down by is easy to find
			yoff := neighborbasey - thisbasey
			// and we just add that
			c.y += int(yoff)
		}
		// TODO if there's no baseline, the alignment should be to the top /of the alignment rect/, not the frame
	}
*/
	C.moveControl(w.id, C.intptr_t(c.x), C.intptr_t(c.y), C.intptr_t(c.width), C.intptr_t(c.height))
}

func (w *widgetbase) getAuxResizeInfo(d *sizing) {
// TODO
//	d.neighborAlign = C.alignmentInfo(s.id, C.frame(s.id))
}

/*
Cocoa doesn't provide a reliable way to get the preferred size of a control (you're supposed to use Interface Builder and have it set up autoresizing for you). The best we can do is call [control sizeToFit] (which is defined for NSControls and has a custom implementation for the other types here) and read the preferred size. Though this changes the size, we're immediately overriding the change on return from sysData.preferredSize(), so no harm done. (This is similar to what we are doing with GTK+, except GTK+ does not actually change the size.)
*/

//TODO
/*
// standard case: control immediately passed in
func controlPrefSize(control C.id) (width int, height int) {
	r := C.controlPrefSize(control)
	return int(r.width), int(r.height)
}

// NSTableView is actually in a NSScrollView so we have to get it out first
func listboxPrefSize(control C.id) (width int, height int) {
	r := C.listboxPrefSize(control)
	return int(r.width), int(r.height)
}

// and for type checking reasons, progress bars are separate despite responding to -[sizeToFit]
func pbarPrefSize(control C.id) (width int, height int) {
	r := C.pbarPrefSize(control)
	return int(r.width), int(r.height)
}

// Areas know their own preferred size
func areaPrefSize(control C.id) (width int, height int) {
	r := C.areaPrefSize(control)
	return int(r.width), int(r.height)
}

var prefsizefuncs = [nctypes]func(C.id) (int, int){
	c_button:      controlPrefSize,
	c_checkbox:    controlPrefSize,
	c_combobox:    controlPrefSize,
	c_lineedit:    controlPrefSize,
	c_label:       controlPrefSize,
	c_listbox:     listboxPrefSize,
	c_progressbar: pbarPrefSize,
	c_area:        areaPrefSize,
}
*/

func (w *widgetbase) preferredSize(d *sizing) (width int, height int) {
//TODO
//	return prefsizefuncs[s.ctype](s.id)
return 0,0
}
