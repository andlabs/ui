// 26 august 2014

#include "objc_darwin.h"
#include <Cocoa/Cocoa.h>

// We would be able to just use plain old NSPopover here, but alas that steals focus.
// NSPopovers are intended for interactive content, and Apple seems to be diligent in enforcing this rule, as the known techniques for preventing a NSPopover from stealing focus no longer work in 10.9.
// Let's just fake it with a window.

// TODO
// - doesn't get hidden properly when asked to order out
// - doesn't get hidden when changing first responders
// - doesn't get hidden when switching between programs/shown again
// - doesn't animate or have a transparent background; probably should

@interface goWarningPopover : NSWindow
@end

@implementation goWarningPopover

- (id)init
{
	self = [super initWithContentRect:NSZeroRect
		styleMask:NSBorderlessWindowMask
		backing:NSBackingStoreBuffered
		defer:YES];
	[self setOpaque:NO];
//	[self setAlphaValue:0.1];
	[self setHasShadow:YES];
	[self setExcludedFromWindowsMenu:YES];
	[self setMovableByWindowBackground:NO];
	[self setLevel:NSPopUpMenuWindowLevel];
	return self;
}

- (BOOL)canBecomeKeyWindow
{
	return NO;
}

- (BOOL)canBecomeMainWindow
{
	return NO;
}

@end

@interface goWarningView : NSView {
@public
	NSImageView *icon;
	NSTextField *label;
}
@end

@implementation goWarningView

- (void)sizeToFitAndArrange
{
	[self->label sizeToFit];

	CGFloat labelheight, imageheight;
	CGFloat targetwidth, imagewidth;

	labelheight = [self->label frame].size.height;
	imageheight = [[self->icon image] size].height;
	imagewidth = [[self->icon image] size].width;
	targetwidth = (imagewidth * labelheight) / imageheight;

	[self->icon setFrameSize:NSMakeSize(targetwidth, labelheight)];

	[self setFrameSize:NSMakeSize(targetwidth + [self->label frame].size.width, labelheight)];
	[self->icon setFrameOrigin:NSMakePoint(0, 0)];
	[self->label setFrameOrigin:NSMakePoint(targetwidth, 0)];
}

- (BOOL)acceptsFirstResponder
{
	return NO;
}

@end

id newWarningPopover(char *text)
{
	goWarningView *wv;

	wv = [[goWarningView alloc] initWithFrame:NSZeroRect];

	wv->icon = [[NSImageView alloc] initWithFrame:NSZeroRect];
	[wv->icon setImage:[NSImage imageNamed:NSImageNameCaution]];
	// TODO verify against Interface Builder
	[wv->icon setImageFrameStyle:NSImageFrameNone];
//	[wv->icon setImageAlignment:xxx];
	[wv->icon setImageScaling:NSImageScaleProportionallyUpOrDown];
	[wv->icon setEditable:NO];
	[wv->icon setAnimates:NO];
	[wv->icon setAllowsCutCopyPaste:NO];
	// TODO check other controls's values for this
	[wv->icon setRefusesFirstResponder:YES];

	wv->label = (NSTextField *) newLabel();
	// TODO rename to textfieldSetText
	textFieldSetText((id) wv->label, text);
	[wv->label setRefusesFirstResponder:YES];

	[wv addSubview:wv->icon];
	[wv addSubview:wv->label];
	[wv sizeToFitAndArrange];

	goWarningPopover *popover;

	popover = [[goWarningPopover alloc] init];		// explicitly use our initializer
	[[popover contentView] addSubview:wv];
	[popover setContentSize:[wv frame].size];

	return (id) popover;
}

void warningPopoverShow(id popover, id control)
{
	goWarningPopover *p = (goWarningPopover *) popover;
	NSView *v = (NSView *) control;
	NSRect vr;
	NSPoint vo;

	vr = [v convertRect:[v frame] toView:nil];
	vo = [[v window] convertRectToScreen:vr].origin;
	[p setFrameOrigin:NSMakePoint(vo.x, vo.y - [p frame].size.height)];
	[p orderFront:p];
}
