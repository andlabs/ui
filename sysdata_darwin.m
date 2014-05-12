// 12 may 2014

//#include "sysdata_darwin.h"
#include "objc_darwin.h"
#include <Foundation/NSGeometry.h>
#include <AppKit/NSWindow.h>
#include <AppKit/NSView.h>
#include <AppKit/NSCell.h>

static NSRect dummyRect;// = NSMakeRect(0, 0, 100, 100);

#define to(T, x) ((T *) x)
#define toNSWindow(x) to(NSWindow, x)
#define toNSView(x) to(NSView, x)

void addControl(id parentWindow, id control)
{
	[[toNSWindow(parentWindow) contentView] addSubview:control];
}

void controlShow(id what)
{
	[toNSView(what) setHidden:NO];
}

void controlHide(id what)
{
	[toNSView(what) setHidden:YES];
}

void applyStandardControlFont(id what)
{
	// TODO inline this
	objc_setFont(what, NSRegularControlSize);
}

id makeWindow(void)
{
	// TODO separate to initilaizer
	dummyRect = NSMakeRect(0, 0, 100, 100);
	return [[NSWindow alloc]
		initWithContentRect:dummyRect
		styleMask:(NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask)
		backing:NSBackingStoreBuffered
		defer:YES];	// defer creation of device until we show the window
}

void windowShow(id window)
{
	[toNSWindow(window) makeKeyAndOrderFront:window];
}

void windowHide(id window)
{
	[toNSWindow(window) orderOut:window];
}
