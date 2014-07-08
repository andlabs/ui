// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

@interface appDelegateClass : NSObject
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

// thanks to mikeash in irc.freenode.net/#macdev for suggesting the use of Grand Dispatch and blocks for this
void issue(void *what)
{
	dispatch_async(dispatch_get_main_queue(), ^{
		doissue(what);
	});
}
