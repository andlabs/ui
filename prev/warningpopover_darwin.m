// 26 august 2014

#include "objc_darwin.h"
#include <Cocoa/Cocoa.h>

// We would be able to just use plain old NSPopover here, but alas that steals focus.
// NSPopovers are intended for interactive content, and Apple seems to be diligent in enforcing this rule, as the known techniques for preventing a NSPopover from stealing focus no longer work in 10.9.
// Let's just fake it with a window.

// TODO better would be to use NSImageNameInvalidDataFreestandingTemplate somehow

@interface goWarningPopover : NSWindow {
@public
	id onBegin;
	id onEnd;
	id textfield;
	NSTextView *tv;
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
	if (self->onBegin != nil) {
		[[NSNotificationCenter defaultCenter] removeObserver:self->onBegin];
		self->onBegin = nil;
	}
	if (self->onEnd != nil) {
		[[NSNotificationCenter defaultCenter] removeObserver:self->onEnd];
		self->onEnd = nil;
	}
	if (self->tv != nil)
		[self->tv removeObserver:self forKeyPath:@"delegate"];
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

- (void)observeValueForKeyPath:(NSString *)keyPath ofObject:(id)object change:(NSDictionary *)change context:(void *)context
{
	if ([self->tv delegate] == self->textfield)
		[self orderFront:self];
	else
		[self orderOut:self];
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
	p->textfield = control;
	p->tv = (NSTextView *) [[v window] fieldEditor:NO forObject:nil];
	// thanks to http://stackoverflow.com/a/25562783/3408572 for suggesting KVO here
	[p->tv addObserver:p forKeyPath:@"delegate" options:NSKeyValueObservingOptionNew context:NULL];
	[p orderFront:p];
}
