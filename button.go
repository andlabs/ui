// 12 february 2014
package ui

import (
	"sync"
)

// A Button represents a clickable button with some text.
type Button struct {
	// This channel gets a message when the button is clicked.
	// Unlike other channels in this package, this channel is initialized to non-nil when creating a new button, and cannot be set to nil later.
	Clicked	chan struct{}

	lock		sync.Mutex
	created	bool
	sysData	*sysData
	initText	string
}

// NewButton creates a new button with the specified text.
func NewButton(text string) (b *Button) {
	return &Button{
		sysData:	mksysdata(c_button),
		initText:	text,
		Clicked:	Event(),
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
	err := b.sysData.make(b.initText, window)
	if err != nil {
		return err
	}
	b.created = true
	return nil
}

func (b *Button) setRect(x int, y int, width int, height int, winheight int) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.sysData.setRect(x, y, width, height, winheight)
}

func (b *Button) preferredSize() (width int, height int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.sysData.preferredSize()
}
