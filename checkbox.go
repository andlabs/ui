// 13 february 2014
//package ui
package main

import (
	"sync"
)

// A Checkbox is a clickable square with a label. The square can be either checked or unchecked.
type Checkbox struct {
	// TODO provide a channel for broadcasting check changes

	lock			sync.Mutex
	created		bool
	sysData		*sysData
	initText		string
	initCheck		bool
}

// NewCheckbox creates a new checkbox with the specified text.
func NewCheckbox(text string) (c *Checkbox) {
	return &Checkbox{
		sysData:	mksysdata(c_checkbox),
		initText:	text,
	}
}

// SetText sets the checkbox's text.
func (c *Checkbox) SetText(text string) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// TODO handle created
	c.initText = text
	return nil
}

// Checked() returns whether or not the checkbox has been checked.
func (c *Checkbox) Checked() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	check, err := c.sysData.isChecked()
	if err != nil {
		panic(err)		// TODO
	}
	return check
}

func (c *Checkbox) apply(window *sysData) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.sysData.make(c.initText, 300, 300, window)
	// TODO size to parent size
}

func (c *Checkbox) setRect(x int, y int, width int, height int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.sysData.setRect(x, y, width, height)
}
