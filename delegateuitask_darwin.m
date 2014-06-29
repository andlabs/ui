// 13 may 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#import <Foundation/NSObject.h>
#import <Foundation/NSValue.h>
#import <Foundation/NSNotification.h>
#import <AppKit/NSApplication.h>
#import <AppKit/NSWindow.h>
#import <Foundation/NSAutoreleasePool.h>
#import <AppKit/NSEvent.h>
#import <AppKit/NSAlert.h>

extern NSRect dummyRect;

@interface ourApplication : NSApplication
@end

@implementation ourApplication

// by default, NSApplication eats some key events
// this prevents that from happening with Area
// see http://stackoverflow.com/questions/24099063/how-do-i-detect-keyup-in-my-nsview-with-the-command-key-held and http://lists.apple.com/archives/cocoa-dev/2003/Oct/msg00442.html
- (void)sendEvent:(NSEvent *)e
{
	NSEventType type;

	type = [e type];
	if (type == NSKeyDown || type == NSKeyUp || type == NSFlagsChanged) {
		id focused;

		focused = [[e window] firstResponder];
		// TODO can focused be nil? the isKindOfClass: docs don't say if it handles nil receivers
		if ([focused isKindOfClass:areaClass])
			switch (type) {
			case NSKeyDown:
				[focused keyDown:e];
				return;
			case NSKeyUp:
				[focused keyUp:e];
				return;
			case NSFlagsChanged:
				[focused flagsChanged:e];
				return;
			}
		// else fall through
	}
	// otherwise, let NSApplication do it
	[super sendEvent:e];
}

@end

@interface appDelegate : NSObject
@end

@implementation appDelegate

- (void)uitask:(NSValue *)fp
{
	appDelegate_uitask([fp pointerValue]);
}

- (BOOL)windowShouldClose:(id)win
{
	return appDelegate_windowShouldClose(win);
}

- (void)windowDidResize:(NSNotification *)n
{
	appDelegate_windowDidResize([n object]);
}

- (void)buttonClicked:(id)button
{
	appDelegate_buttonClicked(button);
}

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)app
{
	NSArray *windows;
	NSUInteger i;

	// try to close all windows
	windows = [NSApp windows];
	for (i = 0; i < [windows count]; i++)
		[[windows objectAtIndex:i] performClose:self];
	// if any windows are left, cancel
	if ([[NSApp windows] count] != 0)
		return NSTerminateCancel;
	// no windows are left; we're good
	return NSTerminateNow;
}

- (void)alertDidEnd:(NSAlert *)alert returnCode:(NSInteger)returnCode contextInfo:(void *)data
{
	NSInteger *ret = (NSInteger *) data;

	*ret = returnCode;
}

@end

id makeAppDelegate(void)
{
	return [appDelegate new];
}

id windowGetContentView(id window)
{
	return [((NSWindow *) window) contentView];
}

BOOL initCocoa(id appDelegate)
{
	// on 10.6 the -[NSApplication setDelegate:] method complains if we don't have one
	NSAutoreleasePool *pool;

	pool = [NSAutoreleasePool new];
	dummyRect = NSMakeRect(0, 0, 100, 100);
	initAreaClass();
	[ourApplication sharedApplication];			// makes NSApp an object of type ourApplication
	if ([NSApp setActivationPolicy:NSApplicationActivationPolicyRegular] != YES)
		return NO;
	[NSApp activateIgnoringOtherApps:YES];		// TODO actually do C.NO here? Russ Cox does YES in his devdraw; the docs say the Finder does NO
	[NSApp setDelegate:appDelegate];
	[pool release];
	return YES;
}

void douitask(id appDelegate, void *p)
{
	NSAutoreleasePool *pool;
	NSValue *fp;

	// we need to make an NSAutoreleasePool, otherwise we get leak warnings on stderr
	pool = [NSAutoreleasePool new];
	fp = [NSValue valueWithPointer:p];
	[appDelegate performSelectorOnMainThread:@selector(uitask:)
		withObject:fp
		waitUntilDone:YES];			// wait so we can properly drain the autorelease pool; on other platforms we wind up waiting anyway (since the main thread can only handle one thing at a time) so
	[pool release];
}

void breakMainLoop(void)
{
	NSEvent *e;

	// -[NSApplication stop:] stops the event loop; it won't do a clean termination, but we're not too concerned with that (at least not on the other platforms either so)
	// we can't call -[NSApplication terminate:] because that will just quit the program, ensuring we never leave ui.Go()
	[NSApp stop:NSApp];
	// simply calling -[NSApplication stop:] is not good enough, as the stop flag is only checked when an event comes in
	// we have to create a "proper" event; a blank event will just throw an exception
	e = [NSEvent otherEventWithType:NSApplicationDefined
		location:NSZeroPoint
		modifierFlags:0
		timestamp:0
		windowNumber:0
		context:nil
		subtype:0
		data1:0
		data2:0];
	[NSApp postEvent:e atStart:NO];			// not at start, just in case there are other events pending (TODO is this correct?)
}

void cocoaMainLoop(void)
{
	[NSApp run];
}
