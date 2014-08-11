// 15 may 2014

#include "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSControl(x) ((NSControl *) (x))
#define toNSTabView(x) ((NSTabView *) (x))
#define toNSScrollView(x) ((NSScrollView *) (x))
#define toNSView(x) ((NSView *) (x))

// TODO figure out where this should go

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
