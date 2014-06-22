// 12 may 2014

#include "objc_darwin.h"
#import <Foundation/NSGeometry.h>
#import <AppKit/NSWindow.h>
#import <AppKit/NSView.h>
#import <AppKit/NSFont.h>
#import <AppKit/NSControl.h>
#import <AppKit/NSButton.h>
#import <AppKit/NSTextField.h>
#import <AppKit/NSSecureTextField.h>
#import <AppKit/NSProgressIndicator.h>
#import <AppKit/NSScrollView.h>

// general TODO: go through all control constructors and their equivalent controls in Interface Builder to see if there's any qualities I'm missing

extern NSRect dummyRect;

#define to(T, x) ((T *) (x))
#define toNSWindow(x) to(NSWindow, (x))
#define toNSView(x) to(NSView, (x))
#define toNSControl(x) to(NSControl, (x))
#define toNSButton(x) to(NSButton, (x))
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

#define systemFontOfSize(s) ([NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:(s)]])

void applyStandardControlFont(id what)
{
	[toNSControl(what) setFont:systemFontOfSize(NSRegularControlSize)];
}

id makeWindow(id delegate)
{
	NSWindow *w;

	w = [[NSWindow alloc]
		initWithContentRect:dummyRect
		styleMask:(NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask)
		backing:NSBackingStoreBuffered
		defer:YES];	// defer creation of device until we show the window
	[w setDelegate:delegate];
	// we do not need setAcceptsMouseMovedEvents: here since we are using a tracking rect in Areas for that
	return w;
}

void windowShow(id window)
{
	[toNSWindow(window) makeKeyAndOrderFront:window];
}

void windowHide(id window)
{
	[toNSWindow(window) orderOut:window];
}

void windowSetTitle(id window, id title)
{
	[toNSWindow(window) setTitle:title];
}

id windowTitle(id window)
{
	return [toNSWindow(window) title];
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

void buttonSetText(id button, id text)
{
	[toNSButton(button) setTitle:text];
}

id buttonText(id button)
{
	return [toNSButton(button) title];
}

id makeCheckbox(void)
{
	NSButton *checkbox;

	checkbox = [[NSButton alloc]
		initWithFrame:dummyRect];
	[checkbox setButtonType:NSSwitchButton];
	return checkbox;
}

id makeLineEdit(BOOL password)
{
	id c;

	if (password)
		c = [[NSSecureTextField alloc]
			initWithFrame:dummyRect];
	else
		c = [[NSTextField alloc]
			initWithFrame:dummyRect];
	// Interface Builder does this to make the text box behave properly
	// see makeLabel() for other side effects
	[[toNSTextField(c) cell] setLineBreakMode:NSLineBreakByClipping];
	// Interface Builder also sets this to allow horizontal scrolling
	[[toNSTextField(c) cell] setScrollable:YES];
	return c;
}

void lineeditSetText(id lineedit, id text)
{
	[toNSTextField(lineedit) setStringValue:text];
}

id lineeditText(id lineedit)
{
	return [toNSTextField(lineedit) stringValue];
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
	// Interface Builder also sets this for its labels, so...
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

void center(id w)
{
	[toNSWindow(w) center];
}
