// 30 july 2014

package ui

// Control represents a control.
type Control interface {
	setParent(p *controlParent) // controlParent defined per-platform
	preferredSize(d *sizing) (width, height int)
	resize(x int, y int, width int, height int, d *sizing)
	nTabStops() int		// used by the Windows backend
}

type controlbase struct {
	fsetParent			func(p *controlParent)
	fpreferredSize		func(d *sizing) (width, height int)
	fresize			func(x int, y int, width int, height int, d *sizing)
	fnTabStops		func() int
}

// children should not use the same name as these, otherwise weird things will happen

func (c *controlbase) setParent(p *controlParent) {
	c.fsetParent(p)
}

func (c *controlbase) preferredSize(d *sizing) (width, height int) {
	return c.fpreferredSize(d)
}

func (c *controlbase) resize(x int, y int, width int, height int, d *sizing) {
	c.fresize(x, y, width, height, d)
}

func (c *controlbase) nTabStops() int {
	return c.fnTabStops()
}
