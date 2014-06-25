// 15 may 2014

#include "objc_darwin.h"
#import <AppKit/NSControl.h>
#import <AppKit/NSView.h>
#import <AppKit/NSTextField.h>
#import <AppKit/NSScrollView.h>
#import <AppKit/NSTableView.h>
#import <AppKit/NSProgressIndicator.h>

#define to(T, x) ((T *) (x))
#define toNSControl(x) to(NSControl, (x))
#define toNSTextField(x) to(NSTextField, (x))
#define toNSScrollView(x) to(NSScrollView, (x))
#define toNSTableView(x) to(NSTableView, (x))
#define toNSProgressIndicator(x) to(NSProgressIndicator, (x))

#define inScrollView(x) ([toNSScrollView((x)) documentView])
#define listboxInScrollView(x) toNSTableView(inScrollView((x)))
#define areaInScrollView(x) inScrollView((x))

struct xprefsize controlPrefSize(id control, BOOL alternate)
{
	NSControl *c;
	NSRect r;
	struct xprefsize s;

	c = toNSControl(control);
	[c sizeToFit];
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	s.yoff = 0;		// no yoff for most controls
	return s;
}

struct xprefsize labelPrefSize(id control, BOOL alternate)
{
	NSControl *c;
	NSRect r;
	struct xprefsize s;

	c = toNSControl(control);
	[c sizeToFit];
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	s.yoff = 0;		// no yoff for standalone labels
	if (!alternate) {
		// TODO this seems really hacky
		// temporarily enable the border, compute its height, and take the height difference
		[toNSTextField(c) setBordered:YES];
		[c sizeToFit];
		r = [c frame];
		[toNSTextField(c) setBordered:NO];
		// - 1 since the sizes are exclusive (????? TODO)
		s.yoff = ((intptr_t) r.size.height) - s.height - 1;
	}
	return s;
}

struct xprefsize listboxPrefSize(id scrollview, BOOL altenrate)
{
	NSTableView *c;
	NSRect r;
	struct xprefsize s;

	c = listboxInScrollView(toNSScrollView(scrollview));
	[c sizeToFit];
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	s.yoff = 0;		// no yoff for listboxes
	return s;
}

struct xprefsize pbarPrefSize(id control, BOOL alternate)
{
	NSProgressIndicator *c;
	NSRect r;
	struct xprefsize s;

	c = toNSProgressIndicator(control);
	[c sizeToFit];
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	s.yoff = 0;		// no yoff for progress bars
	return s;
}

struct xprefsize areaPrefSize(id scrollview, BOOL alternate)
{
	NSView *c;
	NSRect r;
	struct xprefsize s;

	c = areaInScrollView(toNSScrollView(scrollview));
	// don't size to fit; the frame size is already the size we want
	r = [c frame];
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	s.yoff = 0;		// no yoff for areas
	return s;
}
