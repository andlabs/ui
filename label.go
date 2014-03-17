// 14 february 2014

package ui

import (
	"sync"
)

// A Label is a static line of text used to mark other controls.
type Label struct {
	lock		sync.Mutex
	created	bool
	sysData	*sysData
	initText	string
}

// NewLabel creates a new Label with the specified text.
func NewLabel(text string) *Label {
	return &Label{
		sysData:	mksysdata(c_label),
		initText:	text,
	}
}

// SetText sets the Label's text.
func (l *Label) SetText(text string) {
	if l.created {
		l.sysData.setText(text)
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.initText = text
}

// Text returns the Label's text.
func (l *Label) Text() string {
	if l.created {
		return l.sysData.text()
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.initText
}

func (l *Label) make(window *sysData) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	err := l.sysData.make(l.initText, window)
	if err != nil {
		return err
	}
	l.created = true
	return nil
}

func (l *Label) setRect(x int, y int, width int, height int, winheight int) error {
	return l.sysData.setRect(x, y, width, height, winheight)
}

func (l *Label) preferredSize() (width int, height int) {
	return l.sysData.preferredSize()
}
