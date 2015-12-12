// 16 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSButton(x) ((NSButton *) (x))
#define toNSTextField(x) ((NSTextField *) (x))
#define toNSView(x) ((NSView *) (x))
#define toNSWindow(x) ((NSWindow *) (x))
#define toNSBox(x) ((NSBox *) (x))
#define toNSTextView(x) ((NSTextView *) (x))
#define toNSProgressIndicator(x) ((NSProgressIndicator *) (x))

@interface goControlDelegate : NSObject <NSTextFieldDelegate> {
@public
	void *gocontrol;
}
@end

@implementation goControlDelegate

- (IBAction)buttonClicked:(id)sender
{
	buttonClicked(self->gocontrol);
}

- (IBAction)checkboxToggled:(id)sender
{
	checkboxToggled(self->gocontrol);
}

- (void)controlTextDidChange:(NSNotification *)note
{
	textfieldChanged(self->gocontrol);
}

@end

id newButton(void)
{
	NSButton *b;

	b = [[NSButton alloc] initWithFrame:NSZeroRect];
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

	c = [[NSButton alloc] initWithFrame:NSZeroRect];
	[c setButtonType:NSSwitchButton];
	[c setBordered:NO];
	setStandardControlFont((id) c);
	return (id) c;
}

void checkboxSetDelegate(id checkbox, void *b)
{
	goControlDelegate *d;

	d = [goControlDelegate new];
	d->gocontrol = b;
	[toNSButton(checkbox) setTarget:d];
	[toNSButton(checkbox) setAction:@selector(checkboxToggled:)];
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

// also good for labels
// not static because area_darwin.m uses it
id finishNewTextField(id _t, BOOL bordered)
{
	NSTextField *t = toNSTextField(_t);

	// same for text fields, password fields, and labels
	setStandardControlFont((id) t);
	// these three are the same across text fields, password fields, and labels; the only difference is the setBezeled: value, and it's only different on labels
	// THE ORDER OF THESE CALLS IS IMPORTANT; CHANGE IT AND THE BORDERS WILL DISAPPEAR
	[t setBordered:NO];
	[t setBezelStyle:NSTextFieldSquareBezel];
	[t setBezeled:bordered];
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

	t = [[NSTextField alloc] initWithFrame:NSZeroRect];
	[t setSelectable:YES];		// otherwise the setting is masked by the editable default of YES
	return finishNewTextField((id) t, YES);
}

id newPasswordField(void)
{
	NSSecureTextField *t;

	t = [[NSSecureTextField alloc] initWithFrame:NSZeroRect];
	[t setSelectable:YES];		// otherwise the setting is masked by the editable default of YES
	return finishNewTextField((id) t, YES);
}

void textfieldSetDelegate(id textfield, void *t)
{
	goControlDelegate *d;

	d = [goControlDelegate new];
	d->gocontrol = t;
	[toNSTextField(textfield) setDelegate:d];
}

// also good for labels
const char *textfieldText(id t)
{
	return [[toNSTextField(t) stringValue] UTF8String];
}

// also good for labels
void textfieldSetText(id t, char *text)
{
	[toNSTextField(t) setStringValue:[NSString stringWithUTF8String:text]];
}

id textfieldOpenInvalidPopover(id textfield, char *reason)
{
	id popover;

	popover = newWarningPopover(reason);
	warningPopoverShow(popover, textfield);
	NSBeep();
	return (id) popover;
}

void textfieldCloseInvalidPopover(id popover)
{
	[toNSWindow(popover) close];
	// don't release; close does that already
}

BOOL textfieldEditable(id textfield)
{
	return [toNSTextField(textfield) isEditable];
}

void textfieldSetEditable(id textfield, BOOL editable)
{
	[toNSTextField(textfield) setEditable:editable];
}

id newLabel(void)
{
	NSTextField *l;

	l = [[NSTextField alloc] initWithFrame:NSZeroRect];
	[l setEditable:NO];
	[l setSelectable:NO];
	[l setDrawsBackground:NO];
	return finishNewTextField((id) l, NO);
}

id newGroup(id container)
{
	NSBox *group;

	group = [[NSBox alloc] initWithFrame:NSZeroRect];
	[group setBorderType:NSLineBorder];
	[group setBoxType:NSBoxPrimary];
	[group setTransparent:NO];
	// can't use setSmallControlFont() here because the selector is different
	[group setTitleFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSSmallControlSize]]];
	[group setTitlePosition:NSAtTop];
	[group setContentView:toNSView(container)];
	return (id) group;
}

const char *groupText(id group)
{
	return [[toNSBox(group) title] UTF8String];
}

void groupSetText(id group, char *text)
{
	[toNSBox(group) setTitle:[NSString stringWithUTF8String:text]];
}

id newTextbox(void)
{
	NSTextView *tv;

	tv = [[NSTextView alloc] initWithFrame:NSZeroRect];
	// verified against Interface Builder, except for rich text options
	[tv setAllowsDocumentBackgroundColorChange:NO];
	[tv setBackgroundColor:[NSColor textBackgroundColor]];
	[tv setTextColor:[NSColor textColor]];
	[tv setAllowsUndo:YES];
	[tv setEditable:YES];
	[tv setSelectable:YES];
	[tv setRichText:NO];
	[tv setImportsGraphics:NO];
	[tv setBaseWritingDirection:NSWritingDirectionNatural];
	// TODO default paragraph format
	[tv setAllowsImageEditing:NO];
	[tv setAutomaticQuoteSubstitutionEnabled:NO];
	[tv setAutomaticLinkDetectionEnabled:NO];
	[tv setUsesRuler:NO];
	[tv setRulerVisible:NO];
	[tv setUsesInspectorBar:NO];
	[tv setSelectionGranularity:NSSelectByCharacter];
//TODO	[tv setInsertionPointColor:[NSColor insertionColor]];
	[tv setContinuousSpellCheckingEnabled:NO];
	[tv setGrammarCheckingEnabled:NO];
	[tv setUsesFontPanel:NO];
	[tv setEnabledTextCheckingTypes:0];
	[tv setAutomaticDashSubstitutionEnabled:NO];
	[tv setAutomaticSpellingCorrectionEnabled:NO];
	[tv setAutomaticTextReplacementEnabled:NO];
	[tv setSmartInsertDeleteEnabled:NO];
	[tv setLayoutOrientation:NSTextLayoutOrientationHorizontal];
	// TODO default find panel behavior
	// now just to be safe; this will do some of the above but whatever
	disableAutocorrect((id) tv);
	// this option is complex; just set it to the Interface Builder default
	[[tv layoutManager] setAllowsNonContiguousLayout:YES];
	// this will work because it's the same selector
	setStandardControlFont((id) tv);
	return (id) tv;
}

char *textboxText(id tv)
{
	return [[toNSTextView(tv) string] UTF8String];
}

void textboxSetText(id tv, char *text)
{
	[toNSTextView(tv) setString:[NSString stringWithUTF8String:text]];
}

id newProgressBar(void)
{
	NSProgressIndicator *pi;

	pi = [[NSProgressIndicator alloc] initWithFrame:NSZeroRect];
	[pi setStyle:NSProgressIndicatorBarStyle];
	[pi setControlSize:NSRegularControlSize];
	[pi setControlTint:NSDefaultControlTint];
	[pi setBezeled:YES];
	[pi setDisplayedWhenStopped:YES];
	[pi setUsesThreadedAnimation:YES];
	[pi setIndeterminate:NO];
	[pi setMinValue:0];
	[pi setMaxValue:100];
	[pi setDoubleValue:0];
	return (id) pi;
}

intmax_t progressbarPercent(id pbar)
{
	return (intmax_t) [toNSProgressIndicator(pbar) doubleValue];
}

void progressbarSetPercent(id pbar, intmax_t percent)
{
	[toNSProgressIndicator(pbar) setDoubleValue:((double) percent)];
}
