// 26 august 2014

#include "objc_darwin.h"
#include <Cocoa/Cocoa.h>

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

	NSPopover *popover;
	NSViewController *vc;

	vc = [NSViewController new];
	[vc setView:wv];
	popover = [NSPopover new];
	[popover setContentViewController:vc];
	[popover setContentSize:[wv frame].size];

	return (id) popover;
}
