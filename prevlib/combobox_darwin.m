// 17 may 2014

#include "objc_darwin.h"
#import <AppKit/NSPopUpButton.h>
#import <AppKit/NSComboBox.h>
#import <AppKit/NSArrayController.h>

/*
Cocoa doesn't have combo boxes in the sense that other systems do. NSPopUpButton is not editable and technically behaves like a menu on a menubar. NSComboBox is editable and is the more traditional combo box, but the edit field and list are even more separated than they are on other platforms.

Unfortunately, their default internal storage mechanisms exhibit the automatic selection behavior I DON'T want, so we're going to have to do that ourselves.

The NSArrayController we use in our Listboxes already behaves the way we want. Consequently, you'll notice a bunch of functions here call functions in listbox_darwin.m. How convenient =P (TODO separate into objc_darwin.m?)

TODO should we use NSComboBox's dataSource feature?
*/

extern NSRect dummyRect;

#define to(T, x) ((T *) (x))
#define toNSPopUpButton(x) to(NSPopUpButton, (x))
#define toNSComboBox(x) to(NSComboBox, (x))

#define toNSInteger(x) ((NSInteger) (x))
#define fromNSInteger(x) ((intptr_t) (x))

#define COMBOBOXKEY @"cbitem"
static NSString *comboboxKey = COMBOBOXKEY;
static NSString *comboboxBinding = @"contentValues";
static NSString *comboboxKeyPath = @"arrangedObjects." COMBOBOXKEY;

id makeCombobox(BOOL editable)
{
	NSArrayController *ac;

	ac = makeListboxArray();
#define BIND bind:comboboxBinding toObject:ac withKeyPath:comboboxKeyPath options:nil
// for NSPopUpButton, we need a little extra work to make it respect the NSArrayController's selection behavior properties
// thanks to stevesliva (http://stackoverflow.com/questions/23715275/cocoa-how-do-i-suppress-nspopupbutton-automatic-selection-synchronization-nsar)
// note: selectionIndex isn't listed in the Cocoa Bindings Reference for NSArrayController under exposed bindings, but is in the Cocoa Bindings Programming Topics under key-value observant properties, so we can still bind to it
#define BINDSEL bind:@"selectedIndex" toObject:ac withKeyPath:@"selectionIndex" options:nil

	if (!editable) {
		NSPopUpButton *pb;

		pb = [[NSPopUpButton alloc]
			initWithFrame:dummyRect
			pullsDown:NO];
		[pb BIND];
		[pb BINDSEL];
		return pb;
	}

	NSComboBox *cb;

	cb = [[NSComboBox alloc]
		initWithFrame:dummyRect];
	[cb setUsesDataSource:NO];
	[cb BIND];
	// no need to bind selection
	return cb;
}

id comboboxText(id c, BOOL editable)
{
	if (!editable)
		return [toNSPopUpButton(c) titleOfSelectedItem];
	return [toNSComboBox(c) stringValue];
}

void comboboxAppend(id c, BOOL editable, id str)
{
	id ac;

	ac = boundListboxArray(c, comboboxBinding);
	listboxArrayAppend(ac, toListboxItem(comboboxKey, str));
}

void comboboxInsertBefore(id c, BOOL editable, id str, intptr_t before)
{
	id ac;

	ac = boundListboxArray(c, comboboxBinding);
	listboxArrayInsertBefore(ac, toListboxItem(comboboxKey, str), before);
}

intptr_t comboboxSelectedIndex(id c)
{
	// both satisfy the selector
	return fromNSInteger([toNSPopUpButton(c) indexOfSelectedItem]);
}

void comboboxDelete(id c, intptr_t index)
{
	id ac;

	ac = boundListboxArray(c, comboboxBinding);
	listboxArrayDelete(ac, index);
}

intptr_t comboboxLen(id c)
{
	// both satisfy the selector
	return fromNSInteger([toNSPopUpButton(c) numberOfItems]);
}
