// 30 july 2014

#import "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSView(x) ((NSView *) (x))

// TODO verify this when we add more scrolling controls
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
