// 1 march 2014
package ui

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
import "C"

// -[NSCell cellSize] is documented as determining the minimum size needed to draw its receiver. This will work for our case; it appears to be what GTK+ does as well.
// See also: http://stackoverflow.com/questions/1056079/is-there-a-way-to-programmatically-determine-the-proper-sizes-for-apples-built
// TODO figure out what to do if one of our controls returns the sentinel (10000, 10000) that indicates we can't use -[NSCell cellSize]

var (
	_cell = sel_getUid("cell")
	_cellSize = sel_getUid("cellSize")
)

func (s *sysData) preferredSize() (width int, height int) {
if classTypes[s.ctype].make == nil { return 0, 0 }	// prevent lockup during window resize
	cell := C.objc_msgSend_noargs(s.id, _cell)
	cs := C.objc_msgSend_stret_size_noargs(cell, _cellSize)
	return int(cs.width), int(cs.height)
}
