// 25 february 2014
package ui

import (
	// ...
)

// A ProgressBar is a horizontal rectangle that fills up from left to right to indicate the progress of a long-running task.
// This progress is typically a percentage, so within the range [0,100].
// Alternatively, a progress bar can be "indeterminate": it indicates progress is being made, but is unclear as to how much.
// The presentation of indeterminate progress bars is system-specific (for instance, on Windows and many GTK+ skins, this is represented by a small chunk of progress going back and forth across the width of the bar).
type ProgressBar struct {
	lock		sync.Mutex
	created	bool
	sysData	*sysData
	initProg	int
}

// NewProgressBar creates a new ProgressBar.
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		sysData:	mksysdata(c_progressbar),
	}
}

// SetProgress sets the currently indicated progress amount on the ProgressBar. If this amount is outside the range [0,100] (ideally -1), the progress bar is indeterminate.
func (p *ProgressBar) SetProgress(percent int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.created {
		p.sysData.setProgress(percent)
		return
	}
	p.initProg = percent
}

func (p *ProgressBar) make(window *sysData) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.sysData.make("", window)
	if err != nil {
		return err
	}
	p.sysData.setProgress(p.initProg)
	p.created = true
	return nil
}

func (p *ProgressBar) setRect(x int, y int, width int, height int) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.sysData.setRect(x, y, width, height)
}

func (p *ProgressBar) preferredSize() (width int, height int, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	width, height = p.sysData.preferredSize()
	return
}
