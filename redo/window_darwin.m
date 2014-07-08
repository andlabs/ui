// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSWindow(x) ((NSWindow *) (x))

// TODO why do I need the explicit interface specification?
@interface goWindowDelegate : NSObject <NSWindowDelegate> {
@public
	void *gowin;
}
@end

@implementation goWindowDelegate

- (BOOL)windowShouldClose:(id)win
{
	return windowClosing(self->gowin);
}

@end

id newWindow(intptr_t width, intptr_t height)
{
	return [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, (CGFloat) width, (CGFloat) height)
		styleMask:(NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask)
		backing:NSBackingStoreBuffered
		defer:YES];
}

void windowSetDelegate(id win, void *w)
{
	goWindowDelegate *d;

	d = [goWindowDelegate new];
	d->gowin = w;
	[toNSWindow(win) setDelegate:d];
}

const char *windowTitle(id win)
{
	return [[toNSWindow(win) title] UTF8String];
}

void windowSetTitle(id win, const char * title)
{
	[toNSWindow(win) setTitle:[NSString stringWithUTF8String:title]];
}

void windowShow(id win)
{
	[toNSWindow(win) makeKeyAndOrderFront:toNSWindow(win)];
}

void windowHide(id win)
{
	[toNSWindow(win) orderOut:toNSWindow(win)];
}

void windowClose(id win)
{
	// TODO
}
