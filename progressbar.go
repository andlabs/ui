// 25 february 2014

package ui

import (
	"sync"
)

// A ProgressBar is a horizontal rectangle that fills up from left to right to indicate the progress of a long-running task.
// This progress is represented by an integer within the range [0,100], representing a percentage.
// Alternatively, a progressbar can show an animation indicating that progress is being made but how much is indeterminate.
// Newly-created ProgressBars default to showing 0% progress.
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

// SetProgress sets the currently indicated progress amount on the ProgressBar.
// If percent is in the range [0,100], the progressBar shows that much percent complete.
// If percent is -1, the ProgressBar is made indeterminate.
// Otherwise, SetProgress panics.
// TODO what happens if you repeatedly call SetProgress(-1)?
func (p *ProgressBar) SetProgress(percent int) {
	p.lock.Lock()
	defer p.lock.Unlock()

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

func (p *ProgressBar) setRect(x int, y int, width int, height int, winheight int) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.sysData.setRect(x, y, width, height, winheight)
}

func (p *ProgressBar) preferredSize() (width int, height int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.sysData.preferredSize()
}
