// 25 february 2014
package ui

import (
	"sync"
)

// A ProgressBar is a horizontal rectangle that fills up from left to right to indicate the progress of a long-running task.
// This progress is typically a percentage, so within the range [0,100].
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

// SetProgress sets the currently indicated progress amount on the ProgressBar. If this amount is outside the range [0,100] (ideally -1), the function will panic (it should allow indeterminate progress bars, alas those are not supported on Windows 2000).
func (p *ProgressBar) SetProgress(percent int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if percent < 0 || percent > 100 {
		panic("invalid percent")		// TODO
	}
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
