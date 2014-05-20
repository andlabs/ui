// 15 may 2014

#include "objc_darwin.h"
#import <AppKit/NSControl.h>

// see delegateuitask_darwin.m
// in this case, NSScrollView.h, NSTableView.h, AND NSProgressIndicator.h all include NSApplication.h

#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_7
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_7
#import <AppKit/NSApplication.h>
#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_6
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_6

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