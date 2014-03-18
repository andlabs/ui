// 13 february 2014

package ui

import (
	"sync"
)

// A Checkbox is a clickable square with a label. The square can be either checked or unchecked. Checkboxes start out unchecked.
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
func (c *Checkbox) SetText(text string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		c.sysData.setText(text)
		return
	}
	c.initText = text
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

	if c.created {
		return c.sysData.isChecked()
	}
	return false
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

func (c *Checkbox) setRect(x int, y int, width int, height int, rr *[]resizerequest) {
	*rr = append(*rr, resizerequest{
		sysData:	c.sysData,
		x:		x,
		y:		y,
		width:	width,
		height:	height,
	})
}

func (c *Checkbox) preferredSize() (width int, height int) {
	return c.sysData.preferredSize()
}
