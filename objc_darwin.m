// 15 may 2014

#include "objc_darwin.h"
#import <Foundation/NSString.h>

// see delegateuitask_darwin.m
// in this case, NSScrollView.h includes NSApplication.h
#import <AppKit/NSView.h>

#ifdef MAC_OS_X_VERSION_10_7
#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_7
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_7
#endif
#import <AppKit/NSApplication.h>
#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_6
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_6

#import <AppKit/NSView.h>
#import <AppKit/NSScrollView.h>

#define to(T, x) ((T *) (x))
#define _toNSString(x) to(NSString, (x))
#define toNSView(x) to(NSView, (x))
#define toNSScrollView(x) to(NSScrollView, (x))

// because the only way to make a new NSControl/NSView is with a frame (it gets overridden later)
NSRect dummyRect;

// this can be called before our NSApp is created, so keep a pool
id toNSString(char *str)
{
	NSAutoreleasePool *pool;
	NSString *s;

	pool = [NSAutoreleasePool new];
	s = [NSString stringWithUTF8String:str];
	[s retain];		// keep alive after releasing the pool
	[pool release];
	return s;
}

char *fromNSString(id str)
{
	return [_toNSString(str) UTF8String];
}

void display(id view)
{
	[toNSView(view) display];
}

struct xrect frame(id view)
{
	NSRect r;
	struct xrect s;

	r = [toNSView(view) frame];
	s.x = (intptr_t) r.origin.x;
	s.y = (intptr_t) r.origin.y;
	s.width = (intptr_t) r.size.width;
	s.height = (intptr_t) r.size.height;
	return s;
}

id makeScrollView(id content)
{
	NSScrollView *scrollview;

	scrollview = [[NSScrollView alloc]
		initWithFrame:dummyRect];
	[scrollview setHasHorizontalScroller:YES];
	[scrollview setHasVerticalScroller:YES];
	[scrollview setAutohidesScrollers:YES];
	// Interface Builder sets this for NSTableViews; we also want this on Areas
	[scrollview setDrawsBackground:YES];
	[scrollview setDocumentView:toNSView(content)];
	return scrollview;
}

void giveScrollViewBezelBorder(id scrollview)
{
	[toNSScrollView(scrollview) setBorderType:NSBezelBorder];
}

id scrollViewContent(id scrollview)
{
	return [toNSScrollView(scrollview) documentView];
}
