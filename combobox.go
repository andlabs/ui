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
	c.sysData.editable = editable
	return c
}

// TODO Append, InsertBefore, Delete

// Selection returns the current selection.
func (c *Combobox) Selection() (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.sysData.text()
}

// TODO SelectedIndex

func (c *Combobox) make(window *sysData) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	err = c.sysData.make("", 300, 300, window)
	if err != nil {
		return err
	}
	for _, s := range c.initItems {
		err = c.sysData.append(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Combobox) setRect(x int, y int, width int, height int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.sysData.setRect(x, y, width, height)
}
