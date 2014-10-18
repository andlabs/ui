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
// thanks to mikeash and JtRip in irc.freenode.net/#macdev
@interface goContainerView : NSView {
@public
	void *gocontainer;
}
@end

@implementation goContainerView

@end

id newContainerView(void *gocontainer)
{
	goContainerView *c;

	c = [[goContainerView alloc] initWithFrame:NSZeroRect];
	c->gocontainer = gocontainer;
	return (id) c;
}

void moveControl(id c, intptr_t x, intptr_t y, intptr_t width, intptr_t height)
{
	[toNSView(c) setFrame:NSMakeRect((CGFloat) x, (CGFloat) y, (CGFloat) width, (CGFloat) height)];
}

struct xrect containerBounds(id view)
{
	NSRect b;
	struct xrect r;

	b = [toNSView(view) bounds];
	r.x = (intptr_t) b.origin.x;
	r.y = (intptr_t) b.origin.y;
	r.width = (intptr_t) b.size.width;
	r.height = (intptr_t) b.size.height;
	return r;
}
