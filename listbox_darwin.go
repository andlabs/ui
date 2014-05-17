// 2 march 2014

package ui

import (
	// ...
)

/*
The Cocoa API was not designed to be used directly in code; you were intended to build your user interfaces with Interface Builder. There is no dedicated listbox class; we have to synthesize it with a NSTableView. And this is difficult in code.

Under normal circumstances we would have to build our own data source class, as Cocoa doesn't provide premade data sources. Thankfully, Mac OS X 10.3 introduced the bindings system, which avoids all that. It's just not documented too well (again, because you're supposed to use Interface Builder). Bear with me here.

After switching from using the Objective-C runtime to using Objective-C directly, you will now need to look both here and in listbox_darwin.m to get what's going on.
*/

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
// #include "objc_darwin.h"
import "C"

/*
We bind our sole NSTableColumn to a NSArrayController.

NSArrayController is a subclass of NSObjectController, which handles key-value pairs. The object of a NSObjectController by default is a NSMutableDictionary of key-value pairs. The keys are the critical part here.

In effect, each object in our NSArrayController is a NSMutableDictionary with one item: a marker key and the actual string as the value.
*/

const (
	_listboxItemKey = "listboxitem"
)

var (
	listboxItemKey = toNSString(_listboxItemKey)
)

func toListboxItem(what string) C.id {
	return C.toListboxItem(listboxItemKey, toNSString(what))
}

func fromListboxItem(dict C.id) string {
	return fromNSString(C.fromListboxItem(dict, listboxItemKey))
}

/*
NSArrayController is what we bind.

This is what provides the actual list modification methods.
	- (void)addObject:(id)object
		adds object to the end of the list
	- (void)insertObject:(id)object atArrangedObjectIndex:(NSInteger)index
		adds object in the list before index
	- (void)removeObjectAtArrangedObjectIndex:(NSInteger)index
		removes the object at index
	- (id)arrangedObjects
		returns the underlying array; really a NSArray

But what is arrangedObjects? Why care about arranging objects? We don't have to arrange the objects; if we don't, they won't be arranged, and arrangedObjects just acts as the unarranged array.

Of course, Mac OS X 10.5 adds the ability to automatically arrange objects. So let's just turn that off to be safe.
*/

func makeListboxArray() C.id {
	return C.makeListboxArray()
}

func listboxArrayAppend(array C.id, what string) {
	C.listboxArrayAppend(array, toListboxItem(what))
}

func listboxArrayInsertBefore(array C.id, what string, before int) {
	C.listboxArrayInsertBefore(array, toListboxItem(what), C.uintptr_t(before))
}

func listboxArrayDelete(array C.id, index int) {
	C.listboxArrayDelete(array, C.uintptr_t(index))
}

func listboxArrayItemAt(array C.id, index int) string {
	dict := C.listboxArrayItemAt(array, C.uintptr_t(index))
	return fromListboxItem(dict)
}


/*
Now we have to establish the binding. To do this, we need the following things:
	- the object to bind (NSTableColumn)
	- the property of the object to bind (in this case, cell values, so the string @"value")
	- the object to bind to (NSArrayController)
	- the "key path" of the data to get from the array controller
	- any options for binding; we won't have any

The key path is easy: it's [the name of the list].[the key to grab]. The name of the list is arrangedObjects, as established above.

We will also need to be able to get the NSArrayController out to do things with it later; alas, the only real way to do it without storing it separately is to get the complete information set on our binding each time (a NSDictionary) and extract just the bound object.
*/

const (
	_listboxItemKeyPath = "arrangedObjects." + _listboxItemKey
)

var (
	tableColumnBinding = toNSString("value")
	listboxItemKeyPath = toNSString(_listboxItemKeyPath)
)

func bindListboxArray(tableColumn C.id, array C.id) {
	C.bindListboxArray(tableColumn, tableColumnBinding,
		array, listboxItemKeyPath)
}

func boundListboxArray(tableColumn C.id) C.id {
	return C.boundListboxArray(tableColumn, tableColumnBinding)
}

/*
Now with all that done, we're ready to creat a table column.

Columns need string identifiers; we'll just reuse the item key.

Editability is also handled here, as opposed to in NSTableView itself.
*/

func makeListboxTableColumn() C.id {
	column := C.makeListboxTableColumn(listboxItemKey)
	bindListboxArray(column, makeListboxArray())
	return column
}

func listboxTableColumn(listbox C.id) C.id {
	return C.listboxTableColumn(listbox, listboxItemKey)
}

/*
The NSTableViews don't draw their own scrollbars; we have to drop our NSTableViews in NSScrollViews for this. The NSScrollView is also what provides the Listbox's border.

The actual creation code was moved to objc_darwin.go.
*/

func makeListboxScrollView(listbox C.id) C.id {
	scrollview := makeScrollView(listbox)
	C.giveScrollViewBezelBorder(scrollview)		// this is what Interface Builder gives the scroll view
	return scrollview
}

func listboxInScrollView(scrollview C.id) C.id {
	return getScrollViewContent(scrollview)
}

/*
And now, a helper function that takes a scroll view and gets out the array.
*/

func listboxArray(listbox C.id) C.id {
	return boundListboxArray(listboxTableColumn(listboxInScrollView(listbox)))
}

/*
...and finally, we work with the NSTableView directly. These are the methods sysData calls.

We'll handle selections from the NSTableView too. The only trickery is dealing with the return value of -[NSTableView selectedRowIndexes]: NSIndexSet. The only way to get indices out of a NSIndexSet is to get them all out wholesale, and working with C arrays in Go is Not Fun.
*/

func makeListbox(parentWindow C.id, alternate bool, s *sysData) C.id {
	listbox := C.makeListbox(makeListboxTableColumn(), toBOOL(alternate))
	listbox = makeListboxScrollView(listbox)
	addControl(parentWindow, listbox)
	return listbox
}

func listboxAppend(listbox C.id, what string, alternate bool) {
	array := listboxArray(listbox)
	listboxArrayAppend(array, what)
}

func listboxInsertBefore(listbox C.id, what string, before int, alternate bool) {
	array := listboxArray(listbox)
	listboxArrayInsertBefore(array, what, before)
}

// technique from http://stackoverflow.com/questions/3773180/how-to-get-indexes-from-nsindexset-into-an-nsarray-in-cocoa
// we don't care that the indices were originally NSUInteger since by this point we have a problem anyway; Go programs generally use int indices anyway
// we also don't care about NSNotFound because we check the count first AND because NSIndexSet is always sorted (and NSNotFound can be a valid index if the list is large enough... since it's NSIntegerMax, not NSUIntegerMax)
func listboxSelectedIndices(listbox C.id) (list []int) {
	indices := C.listboxSelectedRowIndexes(listboxInScrollView(listbox))
	count := int(C.listboxIndexesCount(indices))
	if count == 0 {
		return nil
	}
	list = make([]int, count)
	list[0] = int(C.listboxIndexesFirst(indices))
	for i := 1; i < count; i++ {
		list[i] = int(C.listboxIndexesNext(indices, C.uintptr_t(list[i - 1])))
	}
	return list
}

func listboxSelectedTexts(listbox C.id) (texts []string) {
	indices := listboxSelectedIndices(listbox)
	if len(indices) == 0 {
		return nil
	}
	array := listboxArray(listbox)
	texts = make([]string, len(indices))
	for i := 0; i < len(texts); i++ {
		texts[i] = listboxArrayItemAt(array, indices[i])
	}
	return texts
}

func listboxDelete(listbox C.id, index int) {
	array := listboxArray(listbox)
	listboxArrayDelete(array, index)
}

func listboxLen(listbox C.id) int {
	return int(C.listboxLen(listboxInScrollView(listbox)))
}

func listboxSelectIndices(id C.id, indices []int) {
	listbox := listboxInScrollView(id)
	if len(indices) == 0 {
		C.listboxDeselectAll(listbox)
		return
	}
	panic("selectListboxIndices() > 0 not yet implemented (TODO)")
}
