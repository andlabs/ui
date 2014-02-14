// 12 february 2014
package main

import (
	"sync"
)

// A Button represents a clickable button with some text.
type Button struct {
	// This channel gets a message when the button is clicked. Unlike other channels in this package, this channel is initialized to non-nil when creating a new button, and cannot be set to nil later.
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
		Clicked:	make(chan struct{}),
	}
}

// SetText sets the button's text.
func (b *Button) SetText(text string) (err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// TODO handle created
	b.initText = text
	return nil
}

func (b *Button) apply(window *sysData) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.sysData.event = b.Clicked
	return b.sysData.make(b.initText, 300, 300, window)
	// TODO size to parent size
}

func (b *Button) setRect(x int, y int, width int, height int) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.sysData.setRect(x, y, width, height)
}
