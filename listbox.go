// 14 february 2014
package ui

import (
	"fmt"
	"sync"
)

// A Listbox is a vertical list of items, of which either at most one or any number of items can be selected at any given time.
// On creation, no item is selected.
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

// Append adds items to the end of the Listbox's list.
// Append will panic if something goes wrong on platforms that do not abort themselves.
func (l *Listbox) Append(what ...string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		for _, s := range what {
			l.sysData.append(s)
		}
		return
	}
	l.initItems = append(l.initItems, what...)
}

// InsertBefore inserts a new item in the Listbox before the item at the given position. It panics if the given index is out of bounds.
// InsertBefore will also panic if something goes wrong on platforms that do not abort themselves.
func (l *Listbox) InsertBefore(what string, before int) {
	l.lock.Lock()
	defer l.lock.Unlock()

	var m []string

	if l.created {
		if before < 0 || before >= l.sysData.len() {
			goto badrange
		}
		l.sysData.insertBefore(what, before)
		return
	}
	if before < 0 || before >= len(l.initItems) {
		goto badrange
	}
	m = make([]string, 0, len(l.initItems) + 1)
	m = append(m, l.initItems[:before]...)
	m = append(m, what)
	l.initItems = append(m, l.initItems[before:]...)
	return
badrange:
	panic(fmt.Errorf("index %d out of range in Listbox.InsertBefore()", before))
}

// Delete removes the given item from the Listbox. It panics if the given index is out of bounds.
func (l *Listbox) Delete(index int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		if index < 0 || index >= l.sysData.len() {
			goto badrange
		}
		return l.sysData.delete(index)
	}
	if index < 0 || index >= len(l.initItems) {
		goto badrange
	}
	l.initItems = append(l.initItems[:index], l.initItems[index + 1:]...)
	return nil
badrange:
	panic(fmt.Errorf("index %d out of range in Listbox.Delete()", index))
}

// Selection returns a list of strings currently selected in the Listbox, or an empty list if none have been selected. This list will have at most one item on a single-selection Listbox.
func (l *Listbox) Selection() []string {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.selectedTexts()
	}
	return nil
}

// SelectedIndices returns a list of the currently selected indexes in the Listbox, or an empty list if none have been selected. This list will have at most one item on a single-selection Listbox.
func (l *Listbox) SelectedIndices() []int {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.selectedIndices()
	}
	return nil
}

// Len returns the number of items in the Listbox.
//
// On platforms for which this function may return an error, it panics if one is returned.
func (l *Listbox) Len() int {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.created {
		return l.sysData.len()
	}
	return len(l.initItems)
}

func (l *Listbox) make(window *sysData) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	err = l.sysData.make("", window)
	if err != nil {
		return err
	}
	for _, s := range l.initItems {
		l.sysData.append(s)
	}
	l.created = true
	return nil
}

func (l *Listbox) setRect(x int, y int, width int, height int, winheight int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.sysData.setRect(x, y, width, height, winheight)
}

func (l *Listbox) preferredSize() (width int, height int) {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.sysData.preferredSize()
}
