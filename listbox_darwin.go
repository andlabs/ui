// 2 march 2014
package ui

import (
	"runtime"
	"unsafe"
)

/*
The Cocoa API was not designed to be used directly in code; you were intended to build your user interfaces with Interface Builder. There is no dedicated listbox class; we have to synthesize it with a NSTableView. And this is difficult in code.

Under normal circumstances we would have to build our own data source class, as Cocoa doesn't provide premade data sources. Thankfully, Mac OS X 10.3 introduced the bindings system, which avoids all that. It's just not documented too well (again, because you're supposed to use Interface Builder). Bear with me here.

PERSONAL TODO - make a post somewhere that does all this in Objective-C itself, for the benefit of the programming community.
*/

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
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
	- (void)insertObject:(id)object atArrangedObjectsIndex:(NSInteger)index
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
	_insertObjectAtArrangedObjectsIndex = sel_getUid("insertObject:atArrangedObjectsIndex:")
	_removeObjectAtArrangedObjectsIndex = sel_getUid("removeObjectAtArrangedObjectsIndex:")
	_arrangedObjects = sel_getUid("arrangedObjects")
	_objectAtIndex = sel_getUid("objectAtIndex:")
)

func newListboxArray() C.id {
	array := objc_new(_NSArrayController)
	C.objc_msgSend_bool(array, _setAutomaticallyRearrangesObjects, C.BOOL(C.NO))
	return array
}

func appendListboxArray(array C.id, what string) {
	C.objc_msgSend_id(array, _addObject, toListboxItem(what))
}

func insertListboxArrayBefore(array C.id, what string, before int) {
	objc_msgSend_id_uint(array, _insertObjectAtArrangedObjectsIndex,
		toListboxItem(what), uintptr(before))
}

func deleteListboxArray(array C.id, index int) {
	objc_msgSend_id(array, _removeObjectAtArrangedObjectsIndex, uintptr(index))
}

func indexListboxArray(array C.id, index int) string {
	array = C.objc_msgSend_noargs(array, _arrangedObjects)
	dict := objc_msgSend_uint(array, _objectAtIndex, uintptr(index))
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
		C.id(C.nil))			// no options
}

func listboxArray(tableColumn C.id) C.id {
	dict := C.objc_msgSend_id(tableColumn, _infoForBinding, tableColumnBinding)
	return C.objc_msgSend_id(dict, _objectForKey, *C._NSObservedObjectKey)
}

/*
Now with all that done, we're ready to creat a table column.

Columns need string identifiers; we'll just reuse the item key.
*/

var (
	_NSTableColumn = objc_getClass("NSTableColumn")

	_initWithIdentifier = sel_getUid("initWithIdentifier:")
	_columnWithIdentifier = sel_getUid("columnWithIdentifier:")
)

func newListboxTableColumn() C.id {
	column := objc_alloc(_NSTableColumn)
	column = C.objc_msgSend_id(column, _initWithIdentifier, listboxItemKey)
	// TODO other properties?
	bindListboxArray(column, newListboxArray())
	return column
}

func listboxTableColumn(listbox C.id) C.id {
	return C.objc_msgSend_id(listbox, _columnWithIdentifier, listboxItemKey)
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
	_count = sel_getUid("_count")
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
	C.objc_msgSend_id(listbox, _setHeaderView, C.id(C.nil))
	// TODO others?
	windowView := C.objc_msgSend_noargs(parentWindow, _contentView)
	C.objc_msgSend_id(windowView, _addSubview, listbox)
	return listbox
}

func appendListbox(listbox C.id, what string) {
	array := listboxArray(listboxTableColumn(listbox))
	appendListboxArray(array, what)
}

func insertListboxBefore(listbox C.id, what string, before int) {
	array := listboxArray(listboxTableColumn(listbox))
	insertListboxArrayBefore(array, what, before)
}

// TODO this is inefficient!
// C.NSIndexSetEntries() makes two arrays of size count: one NSUInteger array and one C.uintptr_t array for returning; this makes a third of type []int for using
// if only NSUInteger was usable (see bleh_darwin.m)
func selectedListboxIndices(listbox C.id) (list []int) {
	var cindices []C.uintptr

	indices := C.objc_msgSend_noargs(listbox, _selectedRowIndexes)
	count := int(C.objc_msgSend_uintret(indices, _count))
	if count == 0 {
		ret
	}
	list = make([]int, count)
	cidx := C.NSIndexSetEntries(indices, C.uintptr_t(count))
	defer C.free(cidx)
	pcindices := (*reflect.SliceHeader)(unsafe.Pointer(&cindices))
	pcindices.Cap = count
	pcindices.Len = count
	pcindices.Data = uintptr(cidx)
	for i := 0; i < count; i++ {
		indices[i] = int(cidx[i])
	}
	return indices
}

func selectedListboxTexts(listbox C.id) (texts []string) {
	indices := selectedListboxIndices(listbox)
	if len(indices) == 0 {
		return nil
	}
	array := listboxArray(listboxTableColumn(listbox))
	texts = make([]string, len(indices))
	for i := 0; i < len(texts); i++ {
		texts[i] = indexListboxArray(array, indices[i])
	}
	return texts
}

func deleteListbox(listbox C.id, index int) {
	array := listboxArray(listboxTableColumn(listbox))
	deleteListboxArray(array, index)
}
