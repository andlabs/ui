// 13 february 2014
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

	if c.created {
		return c.sysData.setText(text)
	}
	c.initText = text
	return nil
}

// Text returns the checkbox's text.
func (c *Checkbox) Text() string {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.text()
	}
	return c.initText
}

// Checked() returns whether or not the checkbox has been checked.
func (c *Checkbox) Checked() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	// TODO if not created
	return c.sysData.isChecked()
}

func (c *Checkbox) make(window *sysData) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.sysData.make(c.initText, window)
	if err != nil {
		return err
	}
	c.created = true
	return nil
}

func (c *Checkbox) setRect(x int, y int, width int, height int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.sysData.setRect(x, y, width, height)
}
