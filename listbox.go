// 14 february 2014
//package ui
package main

import (
	"sync"
)

// A Listbox is a vertical list of items, of which one or (optionally) more items can be selected at any given time.
type Listbox struct {
	// TODO Select event

	lock		sync.Mutex
	created	bool
	sysData	*sysData
	initItems	[]string
}

// NewCombobox makes a new combobox with the given items. If multiple is true, the listbox allows multiple selection.
func NewListbox(multiple bool, items ...string) (l *Listbox) {
	l = &Listbox{
		sysData:		mksysdata(c_listbox),
		initItems:		items,
	}
	l.sysData.alternate = multiple
	return l
}

// Append adds an item to the end of the Listbox's list.
func (l *Listbox) Append(what string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.append(what)
	}
	l.initItems = append(l.initItems, what)
	return nil
}

// InsertBefore inserts a new item in the Listbox before the item at the given position.
func (l *Listbox) InsertBefore(what string, before int) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.insertBefore(what, before)
	}
	m := make([]string, 0, len(l.initItems) + 1)
	m = append(m, l.initItems[:before]...)
	m = append(m, what)
	l.initItems = append(m, l.initItems[before:]...)
	return nil
}

// Delete removes the given item from the Listbox.
func (l *Listbox) Delete(index int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.delete(index)
	}
	l.initItems = append(l.initItems[:index], l.initItems[index + 1:]...)
	return nil
}

// TODO Selection

// TODO SelectedIndices

func (l *Listbox) make(window *sysData) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	err = l.sysData.make("", window)
	if err != nil {
		return err
	}
	for _, s := range l.initItems {
		err = l.sysData.append(s)
		if err != nil {
			return err
		}
	}
	l.created = true
	return nil
}

func (l *Listbox) setRect(x int, y int, width int, height int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.sysData.setRect(x, y, width, height)
}
