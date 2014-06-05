// 13 may 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#import <Foundation/NSObject.h>
#import <Foundation/NSValue.h>
#import <Foundation/NSNotification.h>

// see the hack below; we'll include everything first just in case some other headers get included below; if the hack ever gets resolved, we can use the ones below instead
#import <Foundation/NSAutoreleasePool.h>
#import <AppKit/NSEvent.h>
#import <AppKit/NSAlert.h>

// HACK.
// Apple's header files are bugged: there's an enum that was introduced in 10.7 with new values added in 10.8, but instead of wrapping the whole enum in a version check, they wrap just the fields. This means that on 10.6 that enum will be empty, which is illegal.
// As only one other person on the entire internet has had this problem (TODO get link) and no one ever replied to his report, we're on our own here. This is dumb and will break compile-time availability and deprecation checks, but we don't have many other options.
// I could use SDKs here, but on 10.6 itself Xcode 4.3, which changed the location of SDKs, is only available to people with a paid Apple developer account, and Beelsebob on irc.freenode.net/#macdev told me that any other configuration is likely to have a differnet directory entirely, so...
// Of course, if Go were ever to drop 10.6 support, this problem would go away (hopefully).
// Oh, and Xcode 4.2 for Snow Leopard comes with headers that don't include MAC_OS_X_VERSION_10_7 and thus won't have this problem, so we need to watch for that too...
#ifdef MAC_OS_X_VERSION_10_7
#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_7
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_7
#endif
#import <AppKit/NSApplication.h>
#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_6
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_6

// NSWindow.h is below because it includes NSApplication.h
// we'll use the other headers's include guards so if the above is resolved and I forget to uncomment anything below it won't matter
#import <AppKit/NSWindow.h>
#import <Foundation/NSAutoreleasePool.h>
#import <AppKit/NSEvent.h>
#import <AppKit/NSAlert.h>

extern NSRect dummyRect;

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

- (void)alertDidEnd:(NSAlert *)alert returnCode:(NSInteger)returnCode contextInfo:(void *)contextInfo
{
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
	dummyRect = NSMakeRect(0, 0, 100, 100);
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
