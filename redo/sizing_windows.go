// 24 february 2014

package ui

import (
	"fmt"
)

type sizing struct {
	sizingbase

	// for size calculations
	baseX	int
	baseY	int

	// for the actual resizing
	// possibly the HDWP
}

const (
	marginDialogUnits = 7
	paddingDialogUnits = 4
)

func (w *window) beginResize() (d *sizing) {
	d = new(sizing)

	dc := getTextDC(w.hwnd)
	defer releaseTextDC(w.hwnd, dc)

	var tm s_TEXTMETRICW

	res, err := f_GetTextMetricsW(dc, &tm)
	if res == 0 {
		panic(fmt.Errorf("error getting text metrics for preferred size calculations: %v", err))
	}
	d.baseX = int(tm.tmAveCharWidth)		// TODO not optimal; third reference below has better way
	d.baseY = int(tm.tmHeight)

	if w.spaced {
		d.xmargin = f_MulDiv(marginDialogUnits, d.baseX, 4)
		d.ymargin = f_MulDiv(marginDialogUnits, d.baseY, 8)
		d.xpadding = f_MulDiv(paddingDialogUnits, d.baseX, 4)
		d.ypadding = f_MulDiv(paddingDialogUnits, d.baseY, 8)
	}

	return d
}

func (w *window) endResize(d *sizing) {
	// redraw
}

func (w *window) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// no translation needed on windows
}

// TODO move this to sizing.go
// TODO when doing so, account for margins
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
	yoff := stdDlgSizes[s.ctype].yoff
	if s.alternate {
		yoff = stdDlgSizes[s.ctype].yoffalt
	}
	if yoff != 0 {
		yoff = f_MulDiv(yoff, d.baseY, 8)
	}
	c.y += yoff
*/
	res, err := f_MoveWindow(w.hwnd, c.x, c.y, c.width, c.height, c_TRUE)
	if res == 0 {
		panic(fmt.Errorf("error setting window/control rect: %v", err))
	}
}

func (w *widgetbase) getAuxResizeInfo(d *sizing) {
	// do nothing
}

// For Windows, Microsoft just hands you a list of preferred control sizes as part of the MSDN documentation and tells you to roll with it.
// These sizes are given in "dialog units", which are independent of the font in use.
// We need to convert these into standard pixels, which requires we get the device context of the OS window.
// References:
// - http://msdn.microsoft.com/en-us/library/windows/desktop/aa511279.aspx#controlsizing for control sizes
// - http://msdn.microsoft.com/en-us/library/ms645502%28VS.85%29.aspx - the calculation needed
// - http://support.microsoft.com/kb/125681 - to get the base X and Y
// (thanks to http://stackoverflow.com/questions/58620/default-button-size)
// For push buttons, date/time pickers, links (which we don't use), toolbars, and rebars (another type of toolbar), Common Controls version 6 provides convenient methods to use instead, falling back to the old way if it fails.

// As we are left with incomplete data, an arbitrary size will be chosen
const (
	defaultWidth = 100 // 2 * preferred width of buttons
)

type dlgunits struct {
	width   int
	height  int
	longest bool // TODO actually use this
	getsize uintptr
	area    bool // use area sizes instead
	yoff		int
	yoffalt	int
}

// TODO
/*
var stdDlgSizes = [nctypes]dlgunits{
	c_button: dlgunits{
		width:   50,
		height:  14,
		getsize: _BCM_GETIDEALSIZE,
	},
	c_checkbox: dlgunits{
		// widtdh is not defined here so assume longest
		longest: true,
		height:  10,
	},
	c_combobox: dlgunits{
		// technically the height of a combobox has to include the drop-down list (this is a historical accident: originally comboboxes weren't drop-down)
		// but since we're forcing Common Controls version 6, we can take advantage of one of its mechanisms to automatically fix this mistake (bad practice but whatever)
		// see also: http://blogs.msdn.com/b/oldnewthing/archive/2006/03/10/548537.aspx
		// note that the Microsoft guidelines pages don't take the list size into account
		longest: true,
		height:  12, // from http://msdn.microsoft.com/en-us/library/windows/desktop/bb226818%28v=vs.85%29.aspx; the page linked above says 14
	},
	c_lineedit: dlgunits{
		longest: true,
		height:  14,
	},
	c_label: dlgunits{
		longest: true,
		height:  8,
		yoff:		3,
		yoffalt:	0,
	},
	c_listbox: dlgunits{
		longest: true,
		// height is not clearly defined here ("an integral number of items (3 items minimum)") so just use a three-line edit control
		height: 14 + 10 + 10,
	},
	c_progressbar: dlgunits{
		width:  237, // the first reference says 107 also works; TODO decide which to use
		height: 8,
	},
	c_area: dlgunits{
		area: true,
	},
}
*/

func getTextDC(hwnd uintptr) (dc uintptr) {
	dc, err := f_GetDC(hwnd)
	if dc == hNULL {
		panic(fmt.Errorf("error getting DC for preferred size calculations: %v", err))
	}
// TODO
/*
	// TODO save so it can be restored later
	res, err = f_SelectObject(dc, controlFont)
	if res == hNULL {
		panic(fmt.Errorf("error loading control font into device context for preferred size calculation: %v", err))
	}
*/
	return dc
}

func releaseTextDC(hwnd uintptr, dc uintptr) {
	res, err := f_ReleaseDC(hwnd, dc)
	if res == 0 {
		panic(fmt.Errorf("error releasing DC for preferred size calculations: %v", err))
	}
}

func (w *widgetbase) preferredSize(d *sizing) (width int, height int) {
// TODO
/*
	// the preferred size of an Area is its size
	if stdDlgSizes[s.ctype].area {
		return s.areawidth, s.areaheight
	}

	if msg := stdDlgSizes[s.ctype].getsize; msg != 0 {
		var size _SIZE

		r1, _, _ := _sendMessage.Call(
			uintptr(s.hwnd),
			msg,
			uintptr(0),
			uintptr(unsafe.Pointer(&size)))
		if r1 != uintptr(_FALSE) { // success
			return int(size.cx), int(size.cy)
		}
		// otherwise the message approach failed, so fall back to the regular approach
		println("message failed; falling back")
	}

	width = stdDlgSizes[s.ctype].width
	if width == 0 {
		width = defaultWidth
	}
	height = stdDlgSizes[s.ctype].height
	width = f_MulDiv(width, d.baseX, 4)   // equivalent to right of rect
	height = f_MulDiv(height, d.baseY, 8) // equivalent to bottom of rect
*/
	return width, height
}

// note on MulDiv():
// div will not be 0 in the usages above
// we also ignore overflow; that isn't likely to happen for our use case anytime soon
