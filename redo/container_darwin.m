// 4 august 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#include <Cocoa/Cocoa.h>

#define toNSView(x) ((NSView *) (x))

// calling -[className] on the content views of NSWindow, NSTabItem, and NSBox all return NSView, so I'm assuming I just need to override these
// fornunately:
// - NSWindow resizing calls -[setFrameSize:] (but not -[setFrame:])
// - NSTab resizing calls both -[setFrame:] and -[setFrameSIze:] on the current tab
// - NSTab switching tabs calls both -[setFrame:] and -[setFrameSize:] on the new tab
// so we just override setFrameSize:
// (TODO NSBox)
// thanks to mikeash and JtRip in irc.freenode.net/#macdev
@interface goContainerView : NSView {
@public
	void *gocontainer;
}
@end

@implementation goContainerView

- (void)setFrameSize:(NSSize)s
{
	[super setFrameSize:s];
	containerResized(self->gocontainer, (intptr_t) s.width, (intptr_t) s.height);
}

@end

id newContainerView(void *gocontainer)
{
	goContainerView *c;

	c = [[goContainerView alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	c->gocontainer = gocontainer;
	return (id) c;
}

void moveControl(id c, intptr_t x, intptr_t y, intptr_t width, intptr_t height)
{
	[toNSView(c) setFrame:NSMakeRect((CGFloat) x, (CGFloat) y, (CGFloat) width, (CGFloat) height)];
}
