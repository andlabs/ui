// 24 february 2014

package ui

// #include "winapi_windows.h"
import "C"

// For Windows, Microsoft just hands you a list of preferred control sizes as part of the MSDN documentation and tells you to roll with it.
// These sizes are given in "dialog units", which are independent of the font in use.
// We need to convert these into standard pixels, which requires we get the device context of the OS window.
// References:
// - http://msdn.microsoft.com/en-us/library/ms645502%28VS.85%29.aspx - the calculation needed
// - http://support.microsoft.com/kb/125681 - to get the base X and Y
// (thanks to http://stackoverflow.com/questions/58620/default-button-size)

type sizing struct {
	sizingbase

	// for size calculations
	baseX	C.int
	baseY	C.int

	// for the actual resizing
	// possibly the HDWP
}

// note on MulDiv():
// div will not be 0 in the usages below
// we also ignore overflow; that isn't likely to happen for our use case anytime soon

func fromdlgunitsX(du int, d *sizing) int {
	return int(C.MulDiv(C.int(du), d.baseX, 4))
}

func fromdlgunitsY(du int, d *sizing) int {
	return int(C.MulDiv(C.int(du), d.baseY, 8))
}

const (
	marginDialogUnits = 7
	paddingDialogUnits = 4
)

func (c *container) beginResize() (d *sizing) {
	d = new(sizing)

	d.baseX = C.baseX
	d.baseY = C.baseY

	if spaced {
		d.xmargin = fromdlgunitsX(marginDialogUnits, d)
		d.ymargin = fromdlgunitsY(marginDialogUnits, d)
		d.xpadding = fromdlgunitsX(paddingDialogUnits, d)
		d.ypadding = fromdlgunitsY(paddingDialogUnits, d)
	}

	return d
}

func (c *container) translateAllocationCoords(allocations []*allocation, winwidth, winheight int) {
	// no translation needed on windows
}

//TODO
/*
func (w *widgetbase) preferredSize(d *sizing) (width int, height int) {
	// the preferred size of an Area is its size
	if stdDlgSizes[s.ctype].area {
		return s.areawidth, s.areaheight
	}
}
*/
