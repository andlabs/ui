// 13 may 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#include <Foundation/NSObject.h>
#include <Foundation/NSValue.h>
#include <Foundation/NSNotification.h>
#include <AppKit/NSApplication.h>
#include <AppKit/NSWindow.h>
#include <Foundation/NSAutoreleasePool.h>
#include <AppKit/NSEvent.h>

@interface appDelegate : NSObject
@end

@implementation appDelegate

- (void)uitask:(NSValue *)fp
{
	appDelegate_uitask([fp pointerValue]);
}

- (BOOL)windowShouldClose:(id)win
{
	appDelegate_windowShouldClose(win);
	return NO;		// don't close
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
	appDelegate_applicationShouldTerminate();
	return NSTerminateCancel;
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
	[NSApplication sharedApplication];
	if ([NSApp setActivationPolicy:NSApplicationActivationPolicyRegular] != YES)
		return NO;
	[NSApp activateIgnoringOtherApps:YES];		// TODO actually do C.NO here? Russ Cox does YES in his devdraw; the docs say the Finder does NO
	[NSApp setDelegate:appDelegate];
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
