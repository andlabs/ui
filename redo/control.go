// 30 july 2014

package ui

// All Controls embed this structure, which provides the Control interface methods.
// If a Control needs to override one of these functions, it assigns to the function variables.
type controldefs struct {
	fsetParent func(p *controlParent)
	fcontainerShow func()
	fcontainerHide func()
	fallocate func(x int, y int, width int, height int, d *sizing) []*allocation
	fpreferredSize func(*sizing) (int, int)
	fcommitResize func(*allocation, *sizing)
	fgetAuxResizeInfo func(*sizing)
}

// There's no newcontroldefs() function; all defaults are set by controlbase.

func (w *controldefs) setParent(p *controlParent) {
	w.fsetParent(p)
}

func (w *controldefs) containerShow() {
	w.fcontainerShow()
}

func (w *controldefs) containerHide() {
	w.fcontainerHide()
}

func (w *controldefs) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return w.fallocate(x, y, width, height, d)
}

func (w *controldefs) preferredSize(d *sizing) (int, int) {
	return w.fpreferredSize(d)
}

func (w *controldefs) commitResize(c *allocation, d *sizing) {
	w.fcommitResize(c, d)
}

func (w *controldefs) getAuxResizeInfo(d *sizing) {
	w.fgetAuxResizeInfo(d)
}
