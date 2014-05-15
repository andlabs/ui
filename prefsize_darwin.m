// 15 may 2014

#include "objc_darwin.h"
#include "prefsize_darwin.h"
#include <AppKit/NSControl.h>
#include <AppKit/NSScrollView.h>
#include <AppKit/NSTableView.h>
#include <AppKit/NSProgressIndicator.h>

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