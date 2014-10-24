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

// also fine for NSCells and NSTexts (NSTextViews)
void setStandardControlFont(id control)
{
	[toNSControl(control) setFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSRegularControlSize]]];
}

// also fine for NSCells
void setSmallControlFont(id control)
{
	[toNSControl(control) setFont:[NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:NSSmallControlSize]]];
}

// also good for NSBox and NSProgressIndicator
struct xsize controlPreferredSize(id control)
{
	NSControl *c;
	NSRect r;
	struct xsize s;

	c = toNSControl(control);
	[c sizeToFit];
	// use alignmentRect here instead of frame because we'll be resizing based on that
	r = [c alignmentRectForFrame:[c frame]];
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

// these function are safe to call on Areas; they'll just return the frame and a baseline of 0 since they use the default NSView implementations

static struct xalignment doAlignmentInfo(NSView *v, NSRect r)
{
	struct xalignment a;

	r = [v alignmentRectForFrame:r];
	a.rect.x = (intptr_t) r.origin.x;
	a.rect.y = (intptr_t) r.origin.y;
	a.rect.width = (intptr_t) r.size.width;
	a.rect.height = (intptr_t) r.size.height;
	// I'm not sure if we need to set the frame for -[NSView baselineOffsetFromBottom], but let's do it just to be safe
	[v setFrame:r];
	a.baseline = (intptr_t) [v baselineOffsetFromBottom];
	return a;
}

struct xalignment alignmentInfo(id c, struct xrect newrect)
{
	NSView *v;
	NSRect r;

	v = toNSView(c);
	r = NSMakeRect((CGFloat) newrect.x,
		(CGFloat) newrect.y,
		(CGFloat) newrect.width,
		(CGFloat) newrect.height);
	return doAlignmentInfo(v, r);
}

struct xalignment alignmentInfoFrame(id c)
{
	NSView *v;

	v = toNSView(c);
	return doAlignmentInfo(v, [v frame]);
}
