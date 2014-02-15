// 14 february 2014
//package ui
package main

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

// TODO SetText()/Text()?

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

func (l *Label) setRect(x int, y int, width int, height int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.sysData.setRect(x, y, width, height)
}
