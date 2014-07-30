// 16 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSView(x) ((NSView *) (x))
#define toNSWindow(x) ((NSWindow *) (x))
#define toNSControl(x) ((NSControl *) (x))
#define toNSButton(x) ((NSButton *) (x))
#define toNSTextField(x) ((NSTextField *) (x))

// TODO move to control_darwin.m

void parent(id control, id parentid)
{
	[toNSView(parentid) addSubview:toNSView(control)];
}

void controlSetHidden(id control, BOOL hidden)
{
	[toNSView(control) setHidden:hidden];
}

// also fine for NSCells
void setStandardControlFont(id control)
{
	[toNSControl(control) setFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSRegularControlSize]]];
}

// also fine for NSCells
void setSmallControlFont(id control)
{
	[toNSControl(control) setFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSSmallControlSize]]];
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
	[b setButtonType:NSMomentaryPushInButton];
	[b setBordered:YES];
	[b setBezelStyle:NSRoundedBezelStyle];
	setStandardControlFont((id) b);
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
	[c setButtonType:NSSwitchButton];
	[c setBordered:NO];
	setStandardControlFont((id) c);
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

static id finishNewTextField(NSTextField *t, BOOL bordered)
{
	// same for text fields, password fields, and labels
	setStandardControlFont((id) t);
	// TODO text field/password field border (Interface Builder setting is confusing)
	if (!bordered)
		[t setBordered:NO];
	// smart quotes and other autocorrect features are handled by the window; see newWindow() in window_darwin.m for details
	// Interface Builder does this to make the text box behave properly
	// this disables both word wrap AND ellipsizing in one fell swoop
	// however, we need to send it to the control's cell, not to the control directly
	[[t cell] setLineBreakMode:NSLineBreakByClipping];
	// Interface Builder also sets this to allow horizontal scrolling
	// it also sets this for labels, despite those not being scrollable
	[[t cell] setScrollable:YES];
	return (id) t;
}

id newTextField(void)
{
	NSTextField *t;

	t = [[NSTextField alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	return finishNewTextField(t, YES);
}

id newPasswordField(void)
{
	NSSecureTextField *t;

	t = [[NSSecureTextField alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	return finishNewTextField(toNSTextField(t), YES);
}

const char *textFieldText(id t)
{
	return [[toNSTextField(t) stringValue] UTF8String];
}

void textFieldSetText(id t, char *text)
{
	[toNSTextField(t) setStringValue:[NSString stringWithUTF8String:text]];
}

id newLabel(void)
{
	NSTextField *l;

	l = [[NSTextField alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	[l setEditable:NO];
	[l setSelectable:NO];
	[l setDrawsBackground:NO];
	return finishNewTextField(l, NO);
}
