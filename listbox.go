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

func newListbox(multiple bool, items ...string) (l *Listbox) {
	l = &Listbox{
		sysData:		mksysdata(c_listbox),
		initItems:		items,
	}
	l.sysData.alternate = multiple
	return l
}

// NewListbox creates a new single-selection Listbox with the given items loaded initially.
func NewListbox(items ...string) *Listbox {
	return newListbox(false, items...)
}

// NewMultiSelListbox creates a new multiple-selection Listbox with the given items loaded initially.
func NewMultiSelListbox(items ...string) *Listbox {
	return newListbox(true, items...)
}

// Append adds items to the end of the Listbox's list.
// Append will panic if something goes wrong on platforms that do not abort themselves.
func (l *Listbox) Append(what ...string) {
	if l.created {
		for _, s := range what {
			l.sysData.append(s)
		}
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.initItems = append(l.initItems, what...)
}

// InsertBefore inserts a new item in the Listbox before the item at the given position. It panics if the given index is out of bounds.
// InsertBefore will also panic if something goes wrong on platforms that do not abort themselves.
func (l *Listbox) InsertBefore(what string, before int) {
	var m []string

	if l.created {
		if before < 0 || before >= l.sysData.len() {
			goto badrange
		}
		l.sysData.insertBefore(what, before)
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
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
func (l *Listbox) Delete(index int) {
	if l.created {
		if index < 0 || index >= l.sysData.len() {
			goto badrange
		}
		l.sysData.delete(index)
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if index < 0 || index >= len(l.initItems) {
		goto badrange
	}
	l.initItems = append(l.initItems[:index], l.initItems[index + 1:]...)
	return
badrange:
	panic(fmt.Errorf("index %d out of range in Listbox.Delete()", index))
}

// Selection returns a list of strings currently selected in the Listbox, or an empty list if none have been selected. This list will have at most one item on a single-selection Listbox.
func (l *Listbox) Selection() []string {
	if l.created {
		return l.sysData.selectedTexts()
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	return nil
}

// SelectedIndices returns a list of the currently selected indexes in the Listbox, or an empty list if none have been selected. This list will have at most one item on a single-selection Listbox.
func (l *Listbox) SelectedIndices() []int {
	if l.created {
		return l.sysData.selectedIndices()
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	return nil
}

// Len returns the number of items in the Listbox.
//
// On platforms for which this function may return an error, it panics if one is returned.
func (l *Listbox) Len() int {
	if l.created {
		return l.sysData.len()
	}
	l.lock.Lock()
	defer l.lock.Unlock()
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

func (l *Listbox) setRect(x int, y int, width int, height int, rr *[]resizerequest) {
	*rr = append(*rr, resizerequest{
		sysData:	l.sysData,
		x:		x,
		y:		y,
		width:	width,
		height:	height,
	})
}

func (l *Listbox) preferredSize() (width int, height int) {
	return l.sysData.preferredSize()
}
