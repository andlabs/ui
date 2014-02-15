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

// TODO Append, InsertBefore, Delete

// TODO Selection

// TODO SelectedIndices

func (l *Listbox) make(window *sysData) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	err = l.sysData.make("", 300, 300, window)
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
