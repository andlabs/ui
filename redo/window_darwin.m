// 8 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSWindow(x) ((NSWindow *) (x))

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

- (void)windowDidResize:(NSNotification *)n
{
	[self doWindowResize:[n object]];
}

- (void)doWindowResize:(id)win
{
	NSWindow *w;
	NSRect r;

	w = toNSWindow(win);
	r = [[w contentView] frame];
	windowResized(self->gowin, (uintptr_t) r.size.width, (uintptr_t) r.size.height);
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
	// calling the above the first time won't emit a size changed event (unlike on Windows and GTK+), so fake one to get the controls laid out properly
	windowRedraw(win);
}

void windowHide(id win)
{
	[toNSWindow(win) orderOut:toNSWindow(win)];
}

void windowClose(id win)
{
	[toNSWindow(win) close];
}

// fake a resize event under certain conditions; see each invocation for details
void windowRedraw(id win)
{
	goWindowDelegate *d;

	d = [toNSWindow(win) delegate];
	[d doWindowResize:win];
}