// 13 february 2014

package ui

import (
	"sync"
)

// A Checkbox is a clickable square with a label. The square can be either checked or unchecked. Checkboxes start out unchecked.
type Checkbox struct {
	lock      sync.Mutex
	created   bool
	sysData   *sysData
	initText  string
	initCheck		bool
}

// NewCheckbox creates a new checkbox with the specified text.
func NewCheckbox(text string) (c *Checkbox) {
	return &Checkbox{
		sysData:  mksysdata(c_checkbox),
		initText: text,
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

// SetChecked() changes the checked state of the Checkbox.
func (c *Checkbox) SetChecked(checked bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		c.sysData.setChecked(checked)
		return
	}
	c.initCheck = checked
}

// Checked() returns whether or not the Checkbox has been checked.
func (c *Checkbox) Checked() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.isChecked()
	}
	return c.initCheck
}

func (c *Checkbox) make(window *sysData) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.sysData.make(window)
	if err != nil {
		return err
	}
	c.sysData.setText(c.initText)
	c.sysData.setChecked(c.initCheck)
	c.created = true
	return nil
}

func (c *Checkbox) allocate(x int, y int, width int, height int, d *sysSizeData) []*allocation {
	return []*allocation{&allocation{
		x:       x,
		y:       y,
		width:   width,
		height:  height,
		this:		c,
	}}
}

func (c *Checkbox) preferredSize(d *sysSizeData) (width int, height int) {
	return c.sysData.preferredSize(d)
}

func (c *Checkbox) commitResize(a *allocation, d *sysSizeData) {
	c.sysData.commitResize(a, d)
}

func (c *Checkbox) getAuxResizeInfo(d *sysSizeData) {
	c.sysData.getAuxResizeInfo(d)
}
