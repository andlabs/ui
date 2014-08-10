// 15 may 2014

#include "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSControl(x) ((NSControl *) (x))
#define toNSTabView(x) ((NSTabView *) (x))
#define toNSScrollView(x) ((NSScrollView *) (x))
#define toNSView(x) ((NSView *) (x))

// TODO figure out where these should go

// this function is safe to call on Areas; it'll just return the frame and a baseline of 0 since it uses the default NSView implementations
struct xalignment alignmentInfo(id c, struct xrect newrect)
{
	NSView *v;
	struct xalignment a;
	NSRect r;

	v = toNSView(c);
	r = NSMakeRect((CGFloat) newrect.x,
		(CGFloat) newrect.y,
		(CGFloat) newrect.width,
		(CGFloat) newrect.height);
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

// TODO remove?
struct xrect frame(id c)
{
	NSRect r;
	struct xrect s;

	r = [toNSView(c) frame];
	s.x = (intptr_t) r.origin.x;
	s.y = (intptr_t) r.origin.y;
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	return s;
}
