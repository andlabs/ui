// 17 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSView(x) ((NSView *) (x))

void moveControl(id c, intptr_t x, intptr_t y, intptr_t width, intptr_t height)
{
	// TODO appropriate casts
	[toNSView(c) setFrame:NSMakeRect(x, y, width, height)];
}
