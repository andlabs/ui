// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

@interface appDelegateClass : NSObject
@end

@implementation appDelegateClass

- (void)issue:(id)obj
{
	NSValue *what = (NSValue *) obj;

	doissue([what pointerValue]);
}

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

// Ideally we would have this work like on the other platforms and issue a NSEvent to the end of the event queue
// Unfortunately, there doesn't seem to be a way for NSEvents to hold pointer values, only (signed) NSIntegers
// So we'll have to do the performSelectorOnMainThread: approach
// [TODO]
void issue(void *what)
{
	NSAutoreleasePool *p;
	NSValue *v;

	p = [NSAutoreleasePool new];
	v = [NSValue valueWithPointer:what];
	[appDelegate performSelectorOnMainThread:@selector(issue:)
		withObject:v
		waitUntilDone:NO];
	[p release];
}
