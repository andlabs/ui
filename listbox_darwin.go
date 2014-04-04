// 2 march 2014

package ui

import (
	"reflect"
	"unsafe"
)

/*
The Cocoa API was not designed to be used directly in code; you were intended to build your user interfaces with Interface Builder. There is no dedicated listbox class; we have to synthesize it with a NSTableView. And this is difficult in code.

Under normal circumstances we would have to build our own data source class, as Cocoa doesn't provide premade data sources. Thankfully, Mac OS X 10.3 introduced the bindings system, which avoids all that. It's just not documented too well (again, because you're supposed to use Interface Builder). Bear with me here.

PERSONAL TODO - make a post somewhere that does all this in Objective-C itself, for the benefit of the programming community.

TODO - change the name of some of these functions? specifically the functions that get data about the NSTableView?
*/

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
// #include "objc_darwin.h"
// /* cgo doesn't like nil */
// id nilid = nil;
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
	_NSMutableDictionary = objc_getClass("NSMutableDictionary")

	_dictionaryWithObjectForKey = sel_getUid("dictionaryWithObject:forKey:")
	_objectForKey = sel_getUid("objectForKey:")

	listboxItemKey = toNSString(_listboxItemKey)
)

func toListboxItem(what string) C.id {
	return C.objc_msgSend_id_id(_NSMutableDictionary,
		_dictionaryWithObjectForKey,
		toNSString(what), listboxItemKey)
}

func fromListboxItem(dict C.id) string {
	return fromNSString(C.objc_msgSend_id(dict, _objectForKey, listboxItemKey))
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

var (
	_NSArrayController = objc_getClass("NSArrayController")

	_setAutomaticallyRearrangesObjects = sel_getUid("setAutomaticallyRearrangesObjects:")
	_addObject = sel_getUid("addObject:")
	_insertObjectAtArrangedObjectIndex = sel_getUid("insertObject:atArrangedObjectIndex:")
	_removeObjectAtArrangedObjectIndex = sel_getUid("removeObjectAtArrangedObjectIndex:")
	_arrangedObjects = sel_getUid("arrangedObjects")
	_objectAtIndex = sel_getUid("objectAtIndex:")
)

func newListboxArray() C.id {
	array := C.objc_msgSend_noargs(_NSArrayController, _new)
	C.objc_msgSend_bool(array, _setAutomaticallyRearrangesObjects, C.BOOL(C.NO))
	return array
}

func appendListboxArray(array C.id, what string) {
	C.objc_msgSend_id(array, _addObject, toListboxItem(what))
}

func insertListboxArrayBefore(array C.id, what string, before int) {
	C.objc_msgSend_id_uint(array, _insertObjectAtArrangedObjectIndex,
		toListboxItem(what), C.uintptr_t(before))
}

func deleteListboxArray(array C.id, index int) {
	C.objc_msgSend_uint(array, _removeObjectAtArrangedObjectIndex,
		C.uintptr_t(index))
}

func indexListboxArray(array C.id, index int) string {
	array = C.objc_msgSend_noargs(array, _arrangedObjects)
	dict := C.objc_msgSend_uint(array, _objectAtIndex, C.uintptr_t(index))
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
	_bindToObjectWithKeyPathOptions = sel_getUid("bind:toObject:withKeyPath:options:")
	_infoForBinding = sel_getUid("infoForBinding:")

	tableColumnBinding = toNSString("value")
	listboxItemKeyPath = toNSString(_listboxItemKeyPath)
)

func bindListboxArray(tableColumn C.id, array C.id) {
	C.objc_msgSend_id_id_id_id(tableColumn, _bindToObjectWithKeyPathOptions,
		tableColumnBinding,
		array, listboxItemKeyPath,
		C.nilid)				// no options
}

func listboxArrayController(tableColumn C.id) C.id {
	dict := C.objc_msgSend_id(tableColumn, _infoForBinding, tableColumnBinding)
	return C.objc_msgSend_id(dict, _objectForKey, *C._NSObservedObjectKey)
}

/*
Now with all that done, we're ready to creat a table column.

Columns need string identifiers; we'll just reuse the item key.

Editability is also handled here, as opposed to in NSTableView itself.
*/

var (
	_NSTableColumn = objc_getClass("NSTableColumn")

	_initWithIdentifier = sel_getUid("initWithIdentifier:")
	_tableColumnWithIdentifier = sel_getUid("tableColumnWithIdentifier:")
	// _setEditable in sysdata_darwin.go
)

func newListboxTableColumn() C.id {
	column := objc_alloc(_NSTableColumn)
	column = C.objc_msgSend_id(column, _initWithIdentifier, listboxItemKey)
	C.objc_msgSend_bool(column, _setEditable, C.BOOL(C.NO))
	// TODO other properties?
	bindListboxArray(column, newListboxArray())
	return column
}

func listboxTableColumn(listbox C.id) C.id {
	return C.objc_msgSend_id(listbox, _tableColumnWithIdentifier, listboxItemKey)
}

/*
The NSTableViews don't draw their own scrollbars; we have to drop our NSTableViews in NSScrollViews for this.
*/

var (
	_NSScrollView = objc_getClass("NSScrollView")

	_setHasHorizontalScroller = sel_getUid("setHasHorizontalScroller:")
	_setHasVerticalScroller = sel_getUid("setHasVerticalScroller:")
	_setAutohidesScrollers = sel_getUid("setAutohidesScrollers:")
	_setDocumentView = sel_getUid("setDocumentView:")
	_documentView = sel_getUid("documentView")
)

func newListboxScrollView(listbox C.id) C.id {
	scrollview := objc_alloc(_NSScrollView)
	scrollview = objc_msgSend_rect(scrollview, _initWithFrame,
		0, 0, 100, 100)
	C.objc_msgSend_bool(scrollview, _setHasHorizontalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setHasVerticalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setAutohidesScrollers, C.BOOL(C.YES))
	C.objc_msgSend_id(scrollview, _setDocumentView, listbox)
	return scrollview
}

func listboxInScrollView(scrollview C.id) C.id {
	return C.objc_msgSend_noargs(scrollview, _documentView)
}

/*
And now, a helper function that takes a scroll view and gets out the array.
*/

func listboxArray(listbox C.id) C.id {
	return listboxArrayController(listboxTableColumn(listboxInScrollView(listbox)))
}

/*
...and finally, we work with the NSTableView directly. These are the methods sysData calls.

We'll handle selections from the NSTableView too. The only trickery is dealing with the return value of -[NSTableView selectedRowIndexes]: NSIndexSet. The only way to get indices out of a NSIndexSet is to get them all out wholesale, and working with C arrays in Go is Not Fun.
*/

var (
	_NSTableView = objc_getClass("NSTableView")

	_addTableColumn = sel_getUid("addTableColumn:")
	_setAllowsMultipleSelection = sel_getUid("setAllowsMultipleSelection:")
	_setAllowsEmptySelection = sel_getUid("setAllowsEmptySelection:")
	_setHeaderView = sel_getUid("setHeaderView:")
	_selectedRowIndexes = sel_getUid("selectedRowIndexes")
	_count = sel_getUid("count")
	_numberOfRows = sel_getUid("numberOfRows")
)

func makeListbox(parentWindow C.id, alternate bool) C.id {
	listbox := objc_alloc(_NSTableView)
	listbox = objc_msgSend_rect(listbox, _initWithFrame,
		0, 0, 100, 100)
	C.objc_msgSend_id(listbox, _addTableColumn, newListboxTableColumn())
	multi := C.BOOL(C.NO)
	if alternate {
		multi = C.BOOL(C.YES)
	}
	C.objc_msgSend_bool(listbox, _setAllowsMultipleSelection, multi)
	C.objc_msgSend_bool(listbox, _setAllowsEmptySelection, C.BOOL(C.YES))
	C.objc_msgSend_id(listbox, _setHeaderView, C.nilid)
	// TODO others?
	listbox = newListboxScrollView(listbox)
	addControl(parentWindow, listbox)
	return listbox
}

func appendListbox(listbox C.id, what string, alternate bool) {
	array := listboxArray(listbox)
	appendListboxArray(array, what)
}

func insertListboxBefore(listbox C.id, what string, before int, alternate bool) {
	array := listboxArray(listbox)
	insertListboxArrayBefore(array, what, before)
}

// TODO this is inefficient!
// C.NSIndexSetEntries() makes two arrays of size count: one NSUInteger array and one C.uintptr_t array for returning; this makes a third of type []int for using
// if only NSUInteger was usable (see bleh_darwin.m)
func selectedListboxIndices(listbox C.id) (list []int) {
	var cindices []C.uintptr_t

	indices := C.objc_msgSend_noargs(listboxInScrollView(listbox), _selectedRowIndexes)
	count := int(C.objc_msgSend_uintret_noargs(indices, _count))
	if count == 0 {
		return nil
	}
	list = make([]int, count)
	cidx := C.NSIndexSetEntries(indices, C.uintptr_t(count))
	defer C.free(unsafe.Pointer(cidx))
	pcindices := (*reflect.SliceHeader)(unsafe.Pointer(&cindices))
	pcindices.Cap = count
	pcindices.Len = count
	pcindices.Data = uintptr(unsafe.Pointer(cidx))
	for i := 0; i < count; i++ {
		list[i] = int(cindices[i])
	}
	return list
}

func selectedListboxTexts(listbox C.id) (texts []string) {
	indices := selectedListboxIndices(listbox)
	if len(indices) == 0 {
		return nil
	}
	array := listboxArray(listbox)
	texts = make([]string, len(indices))
	for i := 0; i < len(texts); i++ {
		texts[i] = indexListboxArray(array, indices[i])
	}
	return texts
}

func deleteListbox(listbox C.id, index int) {
	array := listboxArray(listbox)
	deleteListboxArray(array, index)
}

func listboxLen(listbox C.id) int {
	return int(C.objc_msgSend_intret_noargs(listboxInScrollView(listbox), _numberOfRows))
}
