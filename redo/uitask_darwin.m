// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

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
	appDelegate = [appDelegateClass new];
	[NSApplication sharedApplication];
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
