// 14 february 2014

package ui

import (
	"sync"
)

// A LineEdit is a control which allows you to enter a single line of text.
type LineEdit struct {
	lock     sync.Mutex
	created  bool
	sysData  *sysData
	initText string
	password bool
}

// NewLineEdit makes a new LineEdit with the specified text.
func NewLineEdit(text string) *LineEdit {
	return &LineEdit{
		sysData:  mksysdata(c_lineedit),
		initText: text,
	}
}

// NewPasswordEdit makes a new LineEdit which allows the user to enter a password.
func NewPasswordEdit() *LineEdit {
	return &LineEdit{
		sysData:  mksysdata(c_lineedit),
		password: true,
	}
}

// SetText sets the LineEdit's text.
func (l *LineEdit) SetText(text string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		l.sysData.setText(text)
		return
	}
	l.initText = text
}

// Text returns the LineEdit's text.
func (l *LineEdit) Text() string {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.text()
	}
	return l.initText
}

func (l *LineEdit) make(window *sysData) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.sysData.alternate = l.password
	err := l.sysData.make(window)
	if err != nil {
		return err
	}
	l.sysData.setText(l.initText)
	l.created = true
	return nil
}

func (l *LineEdit) allocate(x int, y int, width int, height int, d *sysSizeData) []*allocation {
	return []*allocation{&allocation{
		x:       x,
		y:       y,
		width:   width,
		height:  height,
		this:		l,
	}}
}

func (l *LineEdit) preferredSize(d *sysSizeData) (width int, height int) {
	return l.sysData.preferredSize(d)
}

func (l *LineEdit) commitResize(a *allocation, d *sysSizeData) {
	l.sysData.preferredSize(a, d)
}

func (l *LineEdit) getAuxResizeInfo(d *sysSizeData) {
	l.sysData.getAuxResizeInfo(d)
}
