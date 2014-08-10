// 30 july 2014

#import "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSView(x) ((NSView *) (x))
#define toNSControl(x) ((NSControl *) (x))

void parent(id control, id parentid)
{
	[toNSView(parentid) addSubview:toNSView(control)];
}

void controlSetHidden(id control, BOOL hidden)
{
	[toNSView(control) setHidden:hidden];
}

// also fine for NSCells
void setStandardControlFont(id control)
{
	[toNSControl(control) setFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSRegularControlSize]]];
}

// also fine for NSCells
void setSmallControlFont(id control)
{
	[toNSControl(control) setFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSSmallControlSize]]];
}

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

id newScrollView(id content, BOOL bordered)
{
	NSScrollView *sv;

	sv = [[NSScrollView alloc] initWithFrame:NSZeroRect];
	[sv setDocumentView:toNSView(content)];
	[sv setHasHorizontalScroller:YES];
	[sv setHasVerticalScroller:YES];
	[sv setAutohidesScrollers:YES];
	if (bordered)
		[sv setBorderType:NSBezelBorder];
	else
		[sv setBorderType:NSNoBorder];
	return (id) sv;
}
