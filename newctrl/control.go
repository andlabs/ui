// 30 july 2014

package ui

// Control represents a control.
type Control interface {
	setParent(p *controlParent) // controlParent defined per-platform
//	nChildren() int		// TODO
	preferredSize(d *sizing) (width, height int)
	resize(x int, y int, width int, height int, d *sizing)
}

type controlbase struct {
	fsetParent			func(p *controlParent)
	fpreferredSize		func(d *sizing) (width, height int)
	fresize			func(x int, y int, width int, height int, d *sizing)
}

func (c *controlbase) setParent(p *controlParent) {
	c.fsetParent(p)
}

func (c *controlbase) preferredSize(d *sizing) (width, height int) {
	return c.fpreferredSize(d)
}

func (c *controlbase) resize(x int, y int, width int, height int, d *sizing) {
	c.fresize(x, y, width, height, d)
}
