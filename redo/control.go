// 30 july 2014

package ui

// Control represents a control.
// All Controls have event handlers that take a single argument (the Doer active during the event) and return nothing.
type Control interface {
	setParent(p *controlParent)	// controlParent defined per-platform
	// TODO enable/disable (public)
	controlSizing
}

// this is the same across all platforms
func baseallocate(c Control, x int, y int, width int, height int, d *sizing) []*allocation {
	return []*allocation{&allocation{
		x:		x,
		y:		y,
		width:	width,
		height:	height,
		this:		c,
	}}
}
