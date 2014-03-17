// 14 february 2014

package ui

import (
	"sync"
)

// A LineEdit is a control which allows you to enter a single line of text.
type LineEdit struct {
	// TODO Typing event

	lock			sync.Mutex
	created		bool
	sysData		*sysData
	initText		string
	password		bool
}

// NewLineEdit makes a new LineEdit with the specified text.
func NewLineEdit(text string) *LineEdit {
	return &LineEdit{
		sysData:	mksysdata(c_lineedit),
		initText:	text,
	}
}

// NewPasswordEdit makes a new LineEdit which allows the user to enter a password.
func NewPasswordEdit() *LineEdit {
	return &LineEdit{
		sysData:		mksysdata(c_lineedit),
		password:		true,
	}
}

// SetText sets the LineEdit's text.
func (l *LineEdit) SetText(text string) {
	if l.created {
		l.sysData.setText(text)
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.initText = text
}

// Text returns the LineEdit's text.
func (l *LineEdit) Text() string {
	if l.created {
		return l.sysData.text()
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.initText
}

func (l *LineEdit) make(window *sysData) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.sysData.alternate = l.password
	err := l.sysData.make(l.initText, window)
	if err != nil {
		return err
	}
	l.created = true
	return nil
}

func (l *LineEdit) setRect(x int, y int, width int, height int, winheight int) error {
	return l.sysData.setRect(x, y, width, height, winheight)
}

func (l *LineEdit) preferredSize() (width int, height int) {
	return l.sysData.preferredSize()
}
