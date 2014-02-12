// 12 february 2014
package main

import (
	"fmt"
	"sync"
)

// A Button represents a clickable button with some text.
type Button struct {
	// This channel gets a message when the button is clicked. Unlike other channels in this package, this channel is initialized to non-nil when creating a new button, and cannot be set to nil later.
	Clicked	chan struct{}

	lock		sync.Mutex
	created	bool
	parent	Control
	pWin		*Window
	sysData	*sysData
	initText	string
}

// NewButton creates a new button with the specified text.
func NewButton(text string) (b *Button) {
	return &Button{
		sysData:	&sysData{
			cSysData:		cSysData{
				ctype:	c_button,
			},
		},
		initText:	text,
		Clicked:	make(chan struct{}),
	}
}

// SetText sets the button's text.
func (b *Button) SetText(text string) (err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.pWin != nil && b.pWin.created {
		panic("TODO")
	}
	b.initText = text
	return nil
}

func (b *Button) apply() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.pWin == nil {
		panic(fmt.Sprintf("button (initial text: %q) without parent window", b.initText))
	}
	if !b.pWin.created {
		b.sysData.clicked = b.Clicked
		b.sysData.parentWindow = b.pWin.sysData
		return b.sysData.make(b.initText, 300, 300)
		// TODO size to parent size
	}
	return b.sysData.show()
}

func (b *Button) setParent(c Control) {
	b.parent = c
	if w, ok := b.parent.(*Window); ok {
		b.pWin = w
	} else {
		b.pWin = c.parentWindow()
	}
}

func (b *Button) parentWindow() *Window {
	return b.pWin
}
