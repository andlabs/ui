// 15 may 2014

#include "objc_darwin.h"
#import <AppKit/NSControl.h>
#import <AppKit/NSScrollView.h>
#import <AppKit/NSTableView.h>
#import <AppKit/NSProgressIndicator.h>
#import <AppKit/NSView.h>
// needed for the methods called by alignmentInfo()
#import <AppKit/NSLayoutConstraint.h>

#define to(T, x) ((T *) (x))
#define toNSControl(x) to(NSControl, (x))
#define toNSScrollView(x) to(NSScrollView, (x))
#define toNSTableView(x) to(NSTableView, (x))
#define toNSProgressIndicator(x) to(NSProgressIndicator, (x))
#define toNSView(x) to(NSView, (x))

#define inScrollView(x) ([toNSScrollView((x)) documentView])
#define listboxInScrollView(x) toNSTableView(inScrollView((x)))
#define areaInScrollView(x) inScrollView((x))

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

struct xsize listboxPrefSize(id scrollview)
{
	NSTableView *c;
	NSRect r;
	struct xsize s;

	c = listboxInScrollView(toNSScrollView(scrollview));
	[c sizeToFit];
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	return s;
}

struct xsize pbarPrefSize(id control)
{
	NSProgressIndicator *c;
	NSRect r;
	struct xsize s;

	c = toNSProgressIndicator(control);
	[c sizeToFit];
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	return s;
}

struct xsize areaPrefSize(id scrollview)
{
	NSView *c;
	NSRect r;
	struct xsize s;

	c = areaInScrollView(toNSScrollView(scrollview));
	// don't size to fit; the frame size is already the size we want
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	return s;
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
	a.alignmentRect.x = (intptr_t) r.origin.x;
	a.alignmentRect.y = (intptr_t) r.origin.y;
	a.alignmentRect.width = (intptr_t) r.size.width;
	a.alignmentRect.height = (intptr_t) r.size.height;
	a.baseline = (intptr_t) [v baselineOffsetFromBottom];
	return a;
}
