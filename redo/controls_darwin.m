// 16 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSView(x) ((NSView *) (x))
#define toNSWindow(x) ((NSWindow *) (x))
#define toNSControl(x) ((NSControl *) (x))
#define toNSButton(x) ((NSButton *) (x))

void parent(id control, id parentid)
{
	[[toNSWindow(parentid) contentView] addSubview:toNSView(control)];
}

void controlSetHidden(id control, BOOL hidden)
{
	[toNSView(control) setHidden:hidden];
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

id newButton(void)
{
	NSButton *b;

	// TODO cache the initial rect?
	b = [[NSButton alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	// TODO verify all of these against Interface Builder
	[b setButtonType:NSMomentaryLightButton];
	[b setBordered:YES];
	[b setBezelStyle:NSRoundedBezelStyle];
	setStandardControlFont(b);
	return (id) b;
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

id newCheckbox(void)
{
	NSButton *c;

	c = [[NSButton alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	// TODO verify all of these against Interface Builder
	[c setButtonType:NSSwitchButton];
	[c setBordered:NO];
	setStandardControlFont(c);
	return (id) c;
}

BOOL checkboxChecked(id c)
{
	if ([toNSButton(c) state] == NSOnState)
		return YES;
	return NO;
}

void checkboxSetChecked(id c, BOOL checked)
{
	NSInteger state;

	state = NSOnState;
	if (checked == NO)
		state = NSOffState;
	[toNSButton(c) setState:state];
}
