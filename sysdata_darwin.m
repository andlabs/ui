// 12 may 2014

//#include "sysdata_darwin.h"
#include "objc_darwin.h"
#include <Foundation/NSGeometry.h>
#include <AppKit/NSWindow.h>
#include <AppKit/NSView.h>
#include <AppKit/NSCell.h>
#include <AppKit/NSButton.h>
#include <AppKit/NSPopUpButton.h>
#include <AppKit/NSComboBox.h>
#include <AppKit/NSTextField.h>
#include <AppKit/NSSecureTextField.h>
#include <AppKit/NSProgressIndicator.h>
#include <AppKit/NSScrollView.h>

static NSRect dummyRect;// = NSMakeRect(0, 0, 100, 100);

#define to(T, x) ((T *) (x))
#define toNSWindow(x) to(NSWindow, (x))
#define toNSView(x) to(NSView, (x))
#define toNSButton(x) to(NSButton, (x))
#define toNSPopUpButton(x) to(NSPopUpButton, (x))
#define toNSComboBox(x) to(NSComboBox, (x))
#define toNSTextField(x) to(NSTextField, (x))
#define toNSProgressIndicator(x) to(NSProgressIndicator, (x))
#define toNSScrollView(x) to(NSScrollView, (x))

#define toNSInteger(x) ((NSInteger) (x))
#define fromNSInteger(x) ((intptr_t) (x))

#define inScrollView(x) ([toNSScrollView((x)) documentView])
#define areaInScrollView(x) inScrollView((x))

void addControl(id parentWindow, id control)
{
	[[toNSWindow(parentWindow) contentView] addSubview:control];
}

void controlShow(id what)
{
	[toNSView(what) setHidden:NO];
}

void controlHide(id what)
{
	[toNSView(what) setHidden:YES];
}

void applyStandardControlFont(id what)
{
	// TODO inline this
	objc_setFont(what, NSRegularControlSize);
}

id makeWindow(void)
{
	// TODO separate to initilaizer
	dummyRect = NSMakeRect(0, 0, 100, 100);
	return [[NSWindow alloc]
		initWithContentRect:dummyRect
		styleMask:(NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask)
		backing:NSBackingStoreBuffered
		defer:YES];	// defer creation of device until we show the window
}

void windowShow(id window)
{
	[toNSWindow(window) makeKeyAndOrderFront:window];
}

void windowHide(id window)
{
	[toNSWindow(window) orderOut:window];
}

id makeButton(void)
{
	NSButton *button;

	button = [[NSButton alloc]
		initWithFrame:dummyRect];
	[button setBezelStyle:NSRoundedBezelStyle];
	return button;
}

void buttonSetTargetAction(id button, id delegate)
{
	[toNSButton(button) setTarget:delegate];
	[toNSButton(button) setAction:@selector(buttonClicked:)];
}

id makeCheckbox(void)
{
	NSButton *checkbox;

	checkbox = [[NSButton alloc]
		initWithFrame:dummyRect];
	[checkbox setButtonType:NSSwitchButton];
	return checkbox;
}

id makeCombobox(BOOL editable)
{
	if (!editable) {
		NSPopUpButton *combobox;

		combobox = [[NSPopUpButton alloc]
			initWithFrame:dummyRect
			pullsDown:NO];
		return combobox;
	}

	NSComboBox *combobox;

	combobox = [[NSComboBox alloc]
		initWithFrame:dummyRect];
	[combobox setUsesDataSource:NO];
	return combobox;
}

void comboboxAppend(id combobox, BOOL editable, id str)
{
	if (!editable) {
		[toNSPopUpButton(combobox) addItemWithTitle:str];
		return;
	}
	[toNSComboBox(combobox) addItemWithObjectValue:str];
}

void comboboxInsertBefore(id combobox, BOOL editable, id str, intptr_t before)
{
	if (!editable) {
		[toNSPopUpButton(combobox) insertItemWithTitle:str atIndex:toNSInteger(before)];
		return;
	}
	[toNSComboBox(combobox) insertItemWithObjectValue:str atIndex:toNSInteger(before)];
}

intptr_t comboboxSelectedIndex(id combobox)
{
	// works the same for both NSPopUpButton and NSComboBox
	return fromNSInteger([toNSPopUpButton(combobox) indexOfSelectedItem]);
}

void comboboxDelete(id combobox, intptr_t index)
{
	// works the same for both NSPopUpButton and NSComboBox
	[toNSPopUpButton(combobox) removeItemAtIndex:toNSInteger(index)];
}

intptr_t comboboxLen(id combobox)
{
	// works the same for both NSPopUpButton and NSComboBox
	return fromNSInteger([toNSPopUpButton(combobox) numberOfItems]);
}

void comboboxSelectIndex(id combobox, BOOL editable, intptr_t index)
{
	NSInteger i;
	NSInteger selected;
	NSComboBox *c;

	i = toNSInteger(index);
	// NSPopUpButton documents -1 as deselecting, so we can just use selectItemAtindex: directly
	if (!editable) {
		[toNSPopUpButton(combobox) selectItemAtIndex:i];
		return;
	}
	// NSComboBox, on the other hand, does not, so to be safe, we do things the long way
	c = toNSComboBox(combobox);
	if (i == -1) {		// deselect
		selected = [c indexOfSelectedItem];
		if (selected != -1)
			[c deselectItemAtIndex:selected];
		return;
	}
	[c selectItemAtIndex:i];
}

id makeLineEdit(BOOL password)
{
	if (password)
		return [[NSSecureTextField alloc]
			initWithFrame:dummyRect];
	return [[NSTextField alloc]
		initWithFrame:dummyRect];
}

id makeLabel(void)
{
	NSTextField *label;

	label = [[NSTextField alloc]
		initWithFrame:dummyRect];
	[label setEditable:NO];
	[label setBordered:NO];
	[label setDrawsBackground:NO];
	// this disables both word wrap AND ellipsizing in one fell swoop
	// we have to send to the control's cell for this
	[[label cell] setLineBreakMode:NSLineBreakByClipping];
	// for a multiline label, we either use WordWrapping and send setTruncatesLastVisibleLine: to disable ellipsizing OR use one of those ellipsizing styles
	return label;
}

id makeProgressBar(void)
{
	NSProgressIndicator *pbar;

	pbar = [[NSProgressIndicator alloc]
		initWithFrame:dummyRect];
	[pbar setStyle:NSProgressIndicatorBarStyle];
	[pbar setIndeterminate:NO];
	[pbar stopAnimation:pbar];
	return pbar;
}

void setRect(id what, intptr_t x, intptr_t y, intptr_t width, intptr_t height)
{
	[toNSView(what) setFrame:NSMakeRect((CGFloat) x, (CGFloat) y, (CGFloat) width, (CGFloat) height)];
}

BOOL isCheckboxChecked(id checkbox)
{
	return [toNSButton(checkbox) state] == NSOnState;
}

void windowSetContentSize(id window, intptr_t width, intptr_t height)
{
	NSWindow *win;

	win = toNSWindow(window);
	// use -[NSWindow setContentSize:], which will resize the window without taking the titlebar as part of the given size and without needing us to consider the window's position (the function takes care of both for us)
	[win setContentSize:NSMakeSize((CGFloat) width, (CGFloat) height)];
	[win display];			// TODO needed?
}

void setProgress(id pbar, intptr_t percent)
{
	NSProgressIndicator *p;

	p = toNSProgressIndicator(pbar);
	if (percent == -1) {
		[p setIndeterminate:YES];
		[p startAnimation:p];
		return;
	}
	[p stopAnimation:p];			// will have no effect if we were already determinate
	[p setIndeterminate:NO];
	[p setDoubleValue:((double) percent)];
}

void setAreaSize(id scrollview, intptr_t width, intptr_t height)
{
	NSView *area;

	area = areaInScrollView(scrollview);
	[area setFrame:NSMakeRect(0, 0, (CGFloat) width, (CGFloat) height)];
	[area display];			// and redraw
}
