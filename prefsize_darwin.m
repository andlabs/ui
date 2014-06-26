// 15 may 2014

#include "objc_darwin.h"
#import <AppKit/NSControl.h>
#import <AppKit/NSView.h>
#import <AppKit/NSScrollView.h>
#import <AppKit/NSTableView.h>
#import <AppKit/NSProgressIndicator.h>

#define to(T, x) ((T *) (x))
#define toNSControl(x) to(NSControl, (x))
#define toNSScrollView(x) to(NSScrollView, (x))
#define toNSTableView(x) to(NSTableView, (x))
#define toNSProgressIndicator(x) to(NSProgressIndicator, (x))

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
