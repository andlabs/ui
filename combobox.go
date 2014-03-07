// 14 february 2014
package ui

import (
	"fmt"
	"sync"
)

// A Combobox is a drop-down list of items, of which only one can be selected at any given time. You may optionally make the combobox editable to allow custom items.
type Combobox struct {
	// TODO Select event

	lock		sync.Mutex
	created	bool
	sysData	*sysData
	initItems	[]string
}

func newCombobox(editable bool, items ...string) (c *Combobox) {
	c = &Combobox{
		sysData:		mksysdata(c_combobox),
		initItems:		items,
	}
	c.sysData.alternate = editable
	return c
}

// NewCombobox makes a new Combobox with the given items.
func NewCombobox(items ...string) *Combobox {
	return newCombobox(false, items...)
}

// NewEditableCombobox makes a new editable Combobox with the given items.
func NewEditableCombobox(items ...string) *Combobox {
	return newCombobox(true, items...)
}

// Append adds items to the end of the Combobox's list.
func (c *Combobox) Append(what ...string) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		for i, s := range what {
			err := c.sysData.append(s)
			if err != nil {
				return fmt.Errorf("error adding element %d in Combobox.Append() (%q): %v", i, s, err)
			}
		}
		return nil
	}
	c.initItems = append(c.initItems, what...)
	return nil
}

// InsertBefore inserts a new item in the Combobox before the item at the given position.
func (c *Combobox) InsertBefore(what string, before int) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.insertBefore(what, before)
	}
	m := make([]string, 0, len(c.initItems) + 1)
	m = append(m, c.initItems[:before]...)
	m = append(m, what)
	c.initItems = append(m, c.initItems[before:]...)
	return nil
}

// Delete removes the given item from the Combobox.
func (c *Combobox) Delete(index int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.delete(index)
	}
	c.initItems = append(c.initItems[:index], c.initItems[index + 1:]...)
	return nil
}

// Selection returns the current selection.
func (c *Combobox) Selection() string {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.text()
	}
	return ""
}

// SelectedIndex returns the index of the current selection in the Combobox. It returns -1 either if no selection was made or if text was manually entered in an editable Combobox.
func (c *Combobox) SelectedIndex() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.selectedIndex()
	}
	return -1
}

func (c *Combobox) make(window *sysData) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	err = c.sysData.make("", window)
	if err != nil {
		return err
	}
	for _, s := range c.initItems {
		err = c.sysData.append(s)
		if err != nil {
			return err
		}
	}
	c.created = true
	return nil
}

func (c *Combobox) setRect(x int, y int, width int, height int, winheight int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.sysData.setRect(x, y, width, height, winheight)
}

func (c *Combobox) preferredSize() (width int, height int, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	width, height = c.sysData.preferredSize()
	return
}
