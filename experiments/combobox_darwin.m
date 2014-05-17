// 17 may 2014

#include "objc_darwin.h"
#include <AppKit/NSPopUpButton.h>
#include <AppKit/NSComboBox.h>
#include <AppKit/NSArrayController.h>

/*
Cocoa doesn't have combo boxes in the sense that other systems do. NSPopUpButton is not editable and technically behaves like a menu on a menubar. NSComboBox is editable and is the more traditional combo box, but the edit field and list are even more separated than they are on other platforms.

Unfortunately, their default internal storage mechanisms exhibit the automatic selection behavior I DON'T want, so we're going to have to do that ourselves.

You will notice that a bunch of functions here call functions in listbox_darwin.m. How convenient =P
*/

@interface combobox : NSObject {
@public
	NSPopUpButton *pb;
	NSComboBox *cb;
	NSArrayController *ac;
}
// and these are so we can use a combobox like a reegular control
- (void)setHidden:(BOOL)hidden;
- (void)setFont:(NSFont *)font;
- (void)setFrame:(NSRect)r;
- (NSRect)frame;
- (void)sizeToFit;
@end

@implementation combobox

#define OVERRIDE(sig, msg) \
	sig \
	{ \
		if (pb != nil) { \
			[pb msg]; \
			return; \
		} \
		[cb msg]; \
	}

OVERRIDE(- (void)setHidden:(BOOL)hidden, setHidden:hidden)
OVERRIDE(- (void)setFont:(NSFont *)font, setFont:font)
OVERRIDE(- (void)setFrame:(NSRect)r, setFrame:r)

- (NSRect)frame
{
	if (pb != nil)
		return [pb frame];
	return [cb frame];
}

OVERRIDE(- (void)sizeToFit, sizeToFit)

@end

extern NSRect dummyRect;

#define to(T, x) ((T *) (x))
#define tocombobox(x) to(combobox, (x))

#define toNSInteger(x) ((NSInteger) (x))
#define fromNSInteger(x) ((intptr_t) (x))

#define COMBOBOXKEY @"cbitem"
static NSString *comboboxKey = COMBOBOXKEY;
static NSString *comboboxKeyPath = @"arrangedObjects." COMBOBOXKEY;

id makeCombobox(BOOL editable)
{
	combobox *c;

	c = [combobox new];
	c->pb = nil;
	c->cb = nil;
	c->ac = makeListboxArray();
#define BIND bind:@"contentValues" toObject:c->ac withKeyPath:comboboxKeyPath options:nil
	if (!editable) {
		c->pb = [[NSPopUpButton alloc]
			initWithFrame:dummyRect
			pullsDown:NO];
		[c->pb BIND];
	} else {
		c->cb = [[NSComboBox alloc]
			initWithFrame:dummyRect];
		[c->cb setUsesDataSource:NO];
		[c->cb BIND];
	}
	return c;
}

void addCombobox(id parentWindow, id cbox)
{
	combobox *c;

	c = tocombobox(cbox);
	if (c->pb != nil) {
		addControl(parentWindow, c->pb);
		return;
	}
	addControl(parentWindow, c->cb);
}

id comboboxText(id cbox, BOOL editable)
{
	combobox *c;

	c = tocombobox(cbox);
	if (!editable)
		return [c->pb titleOfSelectedItem];
	return [c->cb stringValue];
}

void comboboxAppend(id cbox, BOOL editable, id str)
{
	combobox *c;

	c = tocombobox(cbox);
	listboxArrayAppend(c->ac, toListboxItem(comboboxKey, str));
}

void comboboxInsertBefore(id cbox, BOOL editable, id str, intptr_t before)
{
	combobox *c;

	c = tocombobox(cbox);
	listboxArrayInsertBefore(c->ac, toListboxItem(comboboxKey, str), before);
}

intptr_t comboboxSelectedIndex(id cbox)
{
	combobox *c;

	c = tocombobox(cbox);
	if (c->pb != nil)
		return fromNSInteger([c->pb indexOfSelectedItem]);
	return fromNSInteger([c->cb indexOfSelectedItem]);
}

void comboboxDelete(id cbox, intptr_t index)
{
	combobox *c;

	c = tocombobox(cbox);
	listboxArrayDelete(c->ac, index);
}

intptr_t comboboxLen(id cbox)
{
	combobox *c;

	c = tocombobox(cbox);
	if (c->pb != nil)
		return fromNSInteger([c->pb numberOfItems]);
	return fromNSInteger([c->cb numberOfItems]);
}
