// 14 february 2014

package ui

import (
	"sync"
)

// A Label is a static line of text used to mark other controls.
// Label text is drawn on a single line; text that does not fit is truncated.
// TODO vertical alignment
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
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		l.sysData.setText(text)
		return
	}
	l.initText = text
}

// Text returns the Label's text.
func (l *Label) Text() string {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.text()
	}
	return l.initText
}

func (l *Label) make(window *sysData) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	err := l.sysData.make(window)
	if err != nil {
		return err
	}
	l.sysData.setText(l.initText)
	l.created = true
	return nil
}

func (l *Label) setRect(x int, y int, width int, height int, rr *[]resizerequest) {
	*rr = append(*rr, resizerequest{
		sysData:	l.sysData,
		x:		x,
		y:		y,
		width:	width,
		height:	height,
	})
}

func (l *Label) preferredSize() (width int, height int) {
	return l.sysData.preferredSize()
}
