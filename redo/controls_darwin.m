// 16 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSView(x) ((NSView *) (x))
#define toNSWindow(x) ((NSWindow *) (x))
#define toNSControl(x) ((NSControl *) (x))
#define toNSButton(x) ((NSButton *) (x))

void unparent(id control)
{
	NSWindow *old;

	[toNSView(control) retain];		// save from being freed when released by the removal selector below
	old = [toNSView(control) window];
	[toNSView(control) removeFromSuperview];
	// redraw since we changed controls
	windowRedraw((id) old);
}

void parent(id control, id parentid, BOOL floating)
{
	[[toNSWindow(parentid) contentView] addSubview:toNSView(control)];
	if (floating)		// previously unparented
		[toNSView(control) release];
	// redraw since we changed controls
	windowRedraw(parentid);
}

static inline void setStandardControlFont(id control)
{
	[toNSControl(control) setFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSRegularControlSize]]];
}

@interface goControlDelegate : NSObject {
@public
	void *gocontrol;
}
@end

@implementation goControlDelegate

- (IBAction)buttonClicked:(id)sender
{
	buttonClicked(self->gocontrol);
}

@end

id newButton(char *text)
{
	NSButton *b;

	// TODO cache the initial rect?
	b = [[NSButton alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	// TODO verify all of these against Interface Builder
	[b setButtonType:NSMomentaryLightButton];
	[b setTitle:[NSString stringWithUTF8String:text]];
	[b setBordered:YES];
	[b setBezelStyle:NSRoundedBezelStyle];
	setStandardControlFont(b);
	return b;
}

void buttonSetDelegate(id button, void *b)
{
	goControlDelegate *d;

	d = [goControlDelegate new];
	d->gocontrol = b;
	[toNSButton(button) setTarget:d];
	[toNSButton(button) setAction:@selector(buttonClicked:)];
}

const char *buttonText(id button)
{
	return [[toNSButton(button) title] UTF8String];
}

void buttonSetText(id button, char *text)
{
	[toNSButton(button) setTitle:[NSString stringWithUTF8String:text]];
}
