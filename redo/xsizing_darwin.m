// 15 may 2014

#include "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSControl(x) ((NSControl *) (x))
#define toNSTabView(x) ((NSTabView *) (x))
#define toNSScrollView(x) ((NSScrollView *) (x))
#define toNSView(x) ((NSView *) (x))

// TODO merge into control_darwin.m or sizing_darwin.m? really need to figure out what to do about the Go-side container struct...

// also good for NSTableView (TODO might not do what we want) and NSProgressIndicator
struct xsize controlPrefSize(id control)
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

struct xsize tabPrefSize(id control)
{
	NSTabView *tv;
	NSSize s;
	struct xsize t;

	tv = toNSTabView(control);
	s = [tv minimumSize];
	t.width = (intptr_t) s.width;
	t.height = (intptr_t) s.height;
	return t;
}

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
	// TODO set frame first?
	a.baseline = (intptr_t) [v baselineOffsetFromBottom];
	return a;
}

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
