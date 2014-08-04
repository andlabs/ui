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

// TODO there has to be a better way
- (void)doWindowResize:(id)win
{
	NSWindow *w;
	NSRect r;

	w = toNSWindow(win);
	r = [[w contentView] frame];
	[[w contentView] setFrameSize:r.size];
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
	// TODO isolate into its own function when (if?) we add TextArea
	tv = (NSTextView *) [w fieldEditor:YES forObject:nil];
	[tv setEnabledTextCheckingTypes:0];
	[tv setAutomaticDashSubstitutionEnabled:NO];
	// don't worry about automatic data detection; it won't change stringValue (thanks pretty_function in irc.freenode.net/#macdev)
	[tv setAutomaticSpellingCorrectionEnabled:NO];
	[tv setAutomaticTextReplacementEnabled:NO];
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

id windowContentView(id win)
{
	return (id) [toNSWindow(win) contentView];
}

// fake a resize event under certain conditions; see each invocation for details
void windowRedraw(id win)
{
	goWindowDelegate *d;

	d = [toNSWindow(win) delegate];
	[d doWindowResize:win];
	// TODO new control sizes don't take effect properly, even with [toNSWindow(win) display];
}
