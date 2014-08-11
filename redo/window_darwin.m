// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSWindow(x) ((NSWindow *) (x))
#define toNSView(x) ((NSView *) (x))

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
	NSWindow *w;
	NSTextView *tv;

	w = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, (CGFloat) width, (CGFloat) height)
		styleMask:(NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask)
		backing:NSBackingStoreBuffered
		defer:YES];
	// we do not want substitutions
	// text fields, labels, etc. take their smart quotes and other autocorrect settings from their parent window, which provides a shared "field editor"
	// so we have to turn them off here
	// thanks akempgen in irc.freenode.net/#macdev
	// for some reason, this selector returns NSText but is documented to return NSTextView...
	disableAutocorrect((id) [w fieldEditor:YES forObject:nil]);
	return w;
}

void windowSetDelegate(id win, void *w)
{
	goWindowDelegate *d;

	d = [goWindowDelegate new];
	d->gowin = w;
	[toNSWindow(win) setDelegate:d];
}

void windowSetContentView(id win, id view)
{
	[toNSWindow(win) setContentView:toNSView(view)];
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
	// no need to worry about reshowing the window initially; that's handled by our container view (container_darwin.m)
}

void windowHide(id win)
{
	[toNSWindow(win) orderOut:toNSWindow(win)];
}

void windowClose(id win)
{
	[toNSWindow(win) close];
}

id windowContentView(id win)
{
	return (id) [toNSWindow(win) contentView];
}
