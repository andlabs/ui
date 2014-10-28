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

- (void)setFrameSize:(NSSize)s
{
	[super setFrameSize:s];
	containerResized(self->gocontainer);
}

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
	NSView *v;
	NSRect frame;

	frame = NSMakeRect((CGFloat) x, (CGFloat) y, (CGFloat) width, (CGFloat) height);
	// mac os x coordinate system has (0,0) in the lower-left
	v = toNSView(c);
	frame.origin.y = ([[v superview] bounds].size.height - frame.size.height) - frame.origin.y;
	// here's the magic: what we specified was what we want the alignment rect to be; make it the actual frame
	frame = [v frameForAlignmentRect:frame];
	[v setFrame:frame];
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
