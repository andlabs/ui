// 25 february 2014

package ui

// A ProgressBar is a horizontal rectangle that fills up from left to right to indicate the progress of a long-running task.
// This progress is represented by an integer within the range [0,100], representing a percentage.
// Alternatively, a progressbar can show an animation indicating that progress is being made but how much is indeterminate.
// Newly-created ProgressBars default to showing 0% progress.
type ProgressBar struct {
	created  bool
	sysData  *sysData
	initProg int
}

// NewProgressBar creates a new ProgressBar.
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		sysData: mksysdata(c_progressbar),
	}
}

// SetProgress sets the currently indicated progress amount on the ProgressBar.
// If percent is in the range [0,100], the progressBar shows that much percent complete.
// If percent is -1, the ProgressBar is made indeterminate.
// Otherwise, SetProgress panics.
// Calling SetProgress(-1) repeatedly will neither leave indeterminate mode nor stop any animation involved in indeterminate mode indefinitely; any other side-effect of doing so is implementation-defined.
func (p *ProgressBar) SetProgress(percent int) {
	if percent < -1 || percent > 100 {
		panic("percent value out of range")
	}
	if p.created {
		p.sysData.setProgress(percent)
		return
	}
	p.initProg = percent
}

func (p *ProgressBar) make(window *sysData) error {
	err := p.sysData.make(window)
	if err != nil {
		return err
	}
	p.sysData.setProgress(p.initProg)
	p.created = true
	return nil
}

func (p *ProgressBar) allocate(x int, y int, width int, height int, d *sysSizeData) []*allocation {
	return []*allocation{&allocation{
		x:       x,
		y:       y,
		width:   width,
		height:  height,
		this:		p,
	}}
}

func (p *ProgressBar) preferredSize(d *sysSizeData) (width int, height int) {
	return p.sysData.preferredSize(d)
}

func (p *ProgressBar) commitResize(a *allocation, d *sysSizeData) {
	p.sysData.commitResize(a, d)
}

func (p *ProgressBar) getAuxResizeInfo(d *sysSizeData) {
	p.sysData.getAuxResizeInfo(d)
}
