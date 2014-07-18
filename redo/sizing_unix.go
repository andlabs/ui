// +build !windows,!darwin

// 23 february 2014

package ui

// #include "gtk_unix.h"
import "C"

type sizing struct {
	sizingbase

	// for size calculations
	// gtk+ needs nothing

	// for the actual resizing
	shouldVAlignTop	bool
}

const (
	gtkXMargin = 12
	gtkYMargin = 12
	gtkXPadding = 12
	gtkYPadding = 6
)

func (w *window) beginResize() (d *sizing) {
	d = new(sizing)
	if w.spaced {
		d.xmargin = gtkXMargin
		d.ymargin = gtkYMargin
		d.xpadding = gtkXPadding
		d.ypadding = gtkYPadding
	}
	return d
}

func (w *window) endResize(d *sizing) {
	C.gtk_widget_queue_draw(w.widget)
}

func (w *window) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// no need for coordinate conversion with gtk+
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
// TODO
/*
	if s.ctype == c_label && !s.alternate && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.shouldVAlignTop {
			// TODO should it be center-aligned to the first line or not
			gtk_misc_set_alignment(s.widget, 0, 0)
		} else {
			gtk_misc_set_alignment(s.widget, 0, 0.5)
		}
	}
*/
	// TODO this uses w.parentw directly; change?
	C.gtk_layout_move(w.parentw.layout, w.widget, C.gint(c.x), C.gint(c.y))
	C.gtk_widget_set_size_request(w.widget, C.gint(c.width), C.gint(c.height))
}

func (w *widgetbase) getAuxResizeInfo(d *sizing) {
//TODO
//	d.shouldVAlignTop = (s.ctype == c_listbox) || (s.ctype == c_area)
	d.shouldVAlignTop = false
}

// GTK+ 3 makes this easy: controls can tell us what their preferred size is!
// ...actually, it tells us two things: the "minimum size" and the "natural size".
// The "minimum size" is the smallest size we /can/ display /anything/. The "natural size" is the smallest size we would /prefer/ to display.
// The difference? Minimum size takes into account things like truncation with ellipses: the minimum size of a label can allot just the ellipses!
// So we use the natural size instead.
// There is a warning about height-for-width controls, but in my tests this isn't an issue.
// For Areas, we manually save the Area size and use that, just to be safe.

// We don't need to worry about y-offset because label alignment is "vertically center", which GtkLabel does for us.

func (w *widgetbase) preferredSize(d *sizing) (width int, height int) {
//TODO
/*
	if s.ctype == c_area {
		return s.areawidth, s.areaheight
	}
*/

	var r C.GtkRequisition

	C.gtk_widget_get_preferred_size(w.widget, nil, &r)
	return int(r.width), int(r.height)
}
