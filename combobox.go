// 14 february 2014
//package ui
package main

import (
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

// NewCombobox makes a new combobox with the given items. If editable is true, the combobox is editable.
func NewCombobox(editable bool, items ...string) (c *Combobox) {
	c = &Combobox{
		sysData:		mksysdata(c_combobox),
		initItems:		items,
	}
	c.sysData.alternate = editable
	return c
}

// Append adds an item to the end of the Combobox's list.
func (c *Combobox) Append(what string) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.append(what)
	}
	c.initItems = append(c.initItems, what)
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

// TODO Delete

// Selection returns the current selection.
func (c *Combobox) Selection() string {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.text()
	}
	return ""
}

// TODO SelectedIndex

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

func (c *Combobox) setRect(x int, y int, width int, height int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.sysData.setRect(x, y, width, height)
}
