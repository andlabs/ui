// 30 july 2014

#import "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSControl(x) ((NSControl *) (x))
#define toNSView(x) ((NSView *) (x))

// also good for NSTableView (TODO might not do what we want) and NSProgressIndicator
struct xsize controlPreferredSize(id control)
{
	NSControl *c;
	NSRect r;
	struct xsize s;

	c = toNSControl(control);
	[c sizeToFit];
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	return s;
}

// TODO verify this when we add more scrolling controls
// TODO no borders on Area
id newScrollView(id content)
{
	NSScrollView *sv;

	sv = [[NSScrollView alloc] initWithFrame:NSZeroRect];
	[sv setDocumentView:toNSView(content)];
	[sv setHasHorizontalScroller:YES];
	[sv setHasVerticalScroller:YES];
	[sv setAutohidesScrollers:YES];
	[sv setBorderType:NSBezelBorder];
	return (id) sv;
}
