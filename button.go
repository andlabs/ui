// 12 february 2014

package ui

import (
	"sync"
)

// A Button represents a clickable button with some text.
type Button struct {
	// Clicked gets a message when the button is clicked.
	// You cannot change it once the Window containing the Button has been created.
	// If you do not respond to this signal, nothing will happen.
	Clicked chan struct{}

	lock     sync.Mutex
	created  bool
	sysData  *sysData
	initText string
}

// NewButton creates a new button with the specified text.
func NewButton(text string) (b *Button) {
	return &Button{
		sysData:  mksysdata(c_button),
		initText: text,
		Clicked:  newEvent(),
	}
}

// SetText sets the button's text.
func (b *Button) SetText(text string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.created {
		b.sysData.setText(text)
		return
	}
	b.initText = text
}

// Text returns the button's text.
func (b *Button) Text() string {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.created {
		return b.sysData.text()
	}
	return b.initText
}

func (b *Button) make(window *sysData) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.sysData.event = b.Clicked
	err := b.sysData.make(window)
	if err != nil {
		return err
	}
	b.sysData.setText(b.initText)
	b.created = true
	return nil
}

func (b *Button) allocate(x int, y int, width int, height int, d *sysSizeData) []*allocation {
	return []*allocation{&allocation{
		x:       x,
		y:       y,
		width:   width,
		height:  height,
		this:		b,
	}}
}

func (b *Button) preferredSize(d *sysSizeData) (width int, height int) {
	return b.sysData.preferredSize(d)
}

func (b *Button) commitResize(a *allocation, d *sysSizeData) {
	b.sysData.commitResize(a, d)
}

func (b *Button) getAuxResizeInfo(d *sysSizeData) {
	b.sysData.getAuxResizeInfo(d)
}
