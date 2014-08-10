// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

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

	type = [e type];
	if (type == NSKeyDown || type == NSKeyUp || type == NSFlagsChanged) {
		id focused;

		focused = [[e window] firstResponder];
		if (focused != nil && [focused isKindOfClass:areaClass])
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

@interface appDelegateClass : NSObject <NSApplicationDelegate>
@end

@implementation appDelegateClass
@end

appDelegateClass *appDelegate;

id getAppDelegate(void)
{
	return appDelegate;
}

BOOL uiinit(void)
{
	areaClass = getAreaClass();
	appDelegate = [appDelegateClass new];
	[goApplication sharedApplication];
	// don't check for a NO return; something (launch services?) causes running from application bundles to always return NO when asking to change activation policy, even if the change is to the same activation policy!
	// see https://github.com/andlabs/ui/issues/6
	[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
	[NSApp activateIgnoringOtherApps:YES];		// TODO rsc does this; finder says NO?
	[NSApp setDelegate:appDelegate];
	return YES;
}

void uimsgloop(void)
{
	[NSApp run];
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
