// 14 february 2014
package ui

import (
	"fmt"
	"sync"
)

// A Combobox is a drop-down list of items, of which at most one can be selected at any given time. You may optionally make the combobox editable to allow custom items. Initially, no item will be selected (and no text entered in an editable Combobox's entry field).
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
// Append will panic if something goes wrong on platforms that do not abort themselves.
func (c *Combobox) Append(what ...string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		for i, s := range what {
			c.sysData.append(s)
		}
		return
	}
	c.initItems = append(c.initItems, what...)
}

// InsertBefore inserts a new item in the Combobox before the item at the given position. It panics if the given index is out of bounds.
// InsertBefore will also panic if something goes wrong on platforms that do not abort themselves.
func (c *Combobox) InsertBefore(what string, before int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var m []string

	if c.created {
		if before < 0 || before >= c.sysData.len() {
			goto badrange
		}
		c.sysData.insertBefore(what, before)
		return
	}
	if before < 0 || before >= len(c.initItems) {
		goto badrange
	}
	m = make([]string, 0, len(c.initItems) + 1)
	m = append(m, c.initItems[:before]...)
	m = append(m, what)
	c.initItems = append(m, c.initItems[before:]...)
	return
badrange:
	panic(fmt.Errorf("index %d out of range in Combobox.InsertBefore()", before))
}

// Delete removes the given item from the Combobox. It panics if the given index is out of bounds.
func (c *Combobox) Delete(index int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		if index < 0 || index >= c.sysData.len() {
			goto badrange
		}
		return c.sysData.delete(index)
	}
	if index < 0 || index >= len(c.initItems) {
		goto badrange
	}
	c.initItems = append(c.initItems[:index], c.initItems[index + 1:]...)
	return nil
badrange:
	panic(fmt.Errorf("index %d out of range in Combobox.Delete()", index))
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

// Len returns the number of items in the Combobox.
//
// On platforms for which this function may return an error, it panics if one is returned.
func (c *Combobox) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.created {
		return c.sysData.len()
	}
	return len(c.initItems)
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
