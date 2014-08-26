// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSWindow(x) ((NSWindow *) (x))

static Class areaClass;

@interface goApplication : NSApplication
@end

@implementation goApplication

// by default, NSApplication eats some key events
// this prevents that from happening with Area
// see http://stackoverflow.com/questions/24099063/how-do-i-detect-keyup-in-my-nsview-with-the-command-key-held and http://lists.apple.com/archives/cocoa-dev/2003/Oct/msg00442.html
- (void)sendEvent:(NSEvent *)e
{
	NSEventType type;
	BOOL handled = NO;

	type = [e type];
	if (type == NSKeyDown || type == NSKeyUp || type == NSFlagsChanged) {
		id focused;

		focused = [[e window] firstResponder];
		if (focused != nil && [focused isKindOfClass:areaClass])
			switch (type) {
			case NSKeyDown:
				handled = [focused doKeyDown:e];
				break;
			case NSKeyUp:
				handled = [focused doKeyUp:e];
				break;
			case NSFlagsChanged:
				handled = [focused doFlagsChanged:e];
				break;
			}
	}
	if (!handled)
		[super sendEvent:e];
}

// ok AppKit, wanna play hardball? let's play hardball.
// because I can neither break out of the special version of the NSModalPanelRunLoopMode that the regular terminate: puts us in nor avoid the exit(0); call included, I'm taking control
// note that this is called AFTER applicationShouldTerminate:
- (void)terminate:(id)sender
{
	// DO ABSOLUTELY NOTHING
	// the magic is [NSApp run] will just... stop.
}

@end

@interface appDelegateClass : NSObject <NSApplicationDelegate>
@end

@implementation appDelegateClass

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)app
{
	NSArray *windows;
	NSUInteger i, n;

	windows = [NSApp windows];
	n = [windows count];
	for (i = 0; i < n; i++) {
		NSWindow *w;

		w = toNSWindow([windows objectAtIndex:i]);
		if (![[w delegate] windowShouldClose:w])
			// stop at the first rejection; thanks Lyle42 in irc.freenode.net/#macdev
			return NSTerminateCancel;
	}
	// all windows closed; stop gracefully for Go
	// note that this is designed for our special terminate: above
	return NSTerminateNow;
}

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)app
{
	return NO;
}

@end

appDelegateClass *appDelegate;

id getAppDelegate(void)
{
	return appDelegate;
}

void uiinit(char **errmsg)
{
	areaClass = getAreaClass();
	appDelegate = [appDelegateClass new];
	[goApplication sharedApplication];
	// don't check for a NO return; something (launch services?) causes running from application bundles to always return NO when asking to change activation policy, even if the change is to the same activation policy!
	// see https://github.com/andlabs/ui/issues/6
	[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
	[NSApp setDelegate:appDelegate];
}

void uimsgloop(void)
{
	[NSApp run];
//	NSLog(@"you shouldn't see this under normal circumstances, but screw the rules, I have SUBCLASSING");
}

// don't use [NSApp terminate:]; that quits the program
void uistop(void)
{
	NSEvent *e;

	[NSApp stop:NSApp];
	// stop: won't register until another event has passed; let's synthesize one
	e = [NSEvent otherEventWithType:NSApplicationDefined
		location:NSZeroPoint
		modifierFlags:0
		timestamp:[[NSProcessInfo processInfo] systemUptime]
		windowNumber:0
		context:[NSGraphicsContext currentContext]
		subtype:0
		data1:0
		data2:0];
	[NSApp postEvent:e atStart:NO];		// let pending events take priority
}

// thanks to mikeash in irc.freenode.net/#macdev for suggesting the use of Grand Central Dispatch and blocks for this
void issue(void *what)
{
	dispatch_async(dispatch_get_main_queue(), ^{
		doissue(what);
	});
}
