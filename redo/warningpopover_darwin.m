// 26 august 2014

#include "objc_darwin.h"
#include <Cocoa/Cocoa.h>

// We would be able to just use plain old NSPopover here, but alas that steals focus.
// NSPopovers are intended for interactive content, and Apple seems to be diligent in enforcing this rule, as the known techniques for preventing a NSPopover from stealing focus no longer work in 10.9.
// Let's just fake it with a window.

@interface goWarningPopover : NSWindow {
@public
	id onBegin;
	id onEnd;
}
@end

@implementation goWarningPopover

- (id)init
{
	self = [super initWithContentRect:NSZeroRect
		styleMask:NSBorderlessWindowMask
		backing:NSBackingStoreBuffered
		defer:YES];
	[self setOpaque:NO];
	[self setHasShadow:YES];
	[self setExcludedFromWindowsMenu:YES];
	[self setMovableByWindowBackground:NO];
	[self setLevel:NSPopUpMenuWindowLevel];
	[self setHidesOnDeactivate:YES];
	self->onBegin = nil;
	self->onEnd = nil;
	return self;
}

- (void)close
{
	NSLog(@"disposing");
	if (self->onBegin != nil) {
		[[NSNotificationCenter defaultCenter] removeObserver:self->onBegin];
		self->onBegin = nil;
	}
	if (self->onEnd != nil) {
		[[NSNotificationCenter defaultCenter] removeObserver:self->onEnd];
		self->onEnd = nil;
	}
	[super close];
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
	[wv->icon setImageFrameStyle:NSImageFrameNone];
	[wv->icon setImageAlignment:NSImageAlignCenter];
	[wv->icon setImageScaling:NSImageScaleProportionallyUpOrDown];
	[wv->icon setEditable:NO];
	[wv->icon setAnimates:NO];
	[wv->icon setAllowsCutCopyPaste:NO];
	[wv->icon setRefusesFirstResponder:YES];

	wv->label = (NSTextField *) newLabel();
	textfieldSetText((id) wv->label, text);
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

	// note that the frame is a rect of the superview
	vr = [[v superview] convertRect:[v frame] toView:nil];
	vo = [[v window] convertRectToScreen:vr].origin;
	[p setFrameOrigin:NSMakePoint(vo.x, vo.y - [p frame].size.height)];
	[p orderFront:p];

	// auto-show/hide when control gains/loses focus
	// TODO this notification is only sent when a character is pressed
	p->onBegin = [[NSNotificationCenter defaultCenter] addObserverForName:NSControlTextDidBeginEditingNotification
		object:v
		queue:nil
		usingBlock:^(NSNotification *note){
			[p orderFront:p];
		}];
	p->onEnd = [[NSNotificationCenter defaultCenter] addObserverForName:NSControlTextDidEndEditingNotification
		object:v
		queue:nil
		usingBlock:^(NSNotification *note){
			[p orderOut:p];
		}];
}
