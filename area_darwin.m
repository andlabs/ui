// 13 may 2014

#include "objc_darwin.h"
#include "area_darwin.h"
#include "_cgo_export.h"
#include <AppKit/NSView.h>
#include <AppKit/NSTrackingArea.h>
#include <Foundation/NSGeometry.h>
#include <AppKit/NSEvent.h>

#define to(T, x) ((T *) (x))
#define toNSEvent(x) to(NSEvent, (x))

#define toNSInteger(x) ((NSInteger) (x))
#define fromNSInteger(x) ((intptr_t) (x))
#define toNSUInteger(x) ((NSUInteger) (x))
#define fromNSUInteger(x) ((uintptr_t) (x))

@interface areaView : NSView {
	NSTrackingArea *trackingArea;
}
@end

@implementation areaView

- (id)initWithFrame:(NSRect)r
{
	self = [super initWithFrame:r];
	if (self)
		[self retrack];
	// TODO other properties?
	return self;
}

- (void)drawRect:(NSRect)cliprect
{
	struct xrect rect;

	rect.x = (intptr_t) cliprect.origin.x;
	rect.y = (intptr_t) cliprect.origin.y;
	rect.width = (intptr_t) cliprect.size.width;
	rect.height = (intptr_t) cliprect.size.height;
	areaView_drawRect(self, rect);
}

- (BOOL)isFlipped
{
	return YES;
}

- (BOOL)acceptsFirstResponder
{
	return YES;
}

- (void)retrack
{
	trackingArea = makeTrackingArea(self);		// TODO make inline
	[self addTrackingArea:trackingArea];
}

- (void)updateTrackingAreas
{
	[self removeTrackingArea:trackingArea];
	[trackingArea release];
	[self retrack];
}

#define mouseEvent(m, f) \
	- (void)m:(NSEvent *)e \
	{ \
		f(self, e); \
	}
mouseEvent(mouseMoved, areaView_mouseMoved_mouseDragged)
mouseEvent(mouseDragged, areaView_mouseMoved_mouseDragged)
mouseEvent(rightMouseDragged, areaView_mouseMoved_mouseDragged)
mouseEvent(otherMouseDragged, areaView_mouseMoved_mouseDragged)
mouseEvent(mouseDown, areaView_mouseDown)
mouseEvent(rightMouseDown, areaView_mouseDown)
mouseEvent(otherMouseDown, areaView_mouseDown)
mouseEvent(mouseUp, areaView_mouseUp)
mouseEvent(rightMouseUp, areaView_mouseUp)
mouseEvent(otherMouseUp, areaView_mouseUp)

#define keyEvent(m, f) \
	- (void)m:(NSEvent *)e \
	{ \
		if (f(self, e) == NO) \
			[super m:e]; \
	}
keyEvent(keyDown, areaView_keyDown)
keyEvent(keyUp, areaView_keyUp)
keyEvent(flagsChanged, areaView_flagsChanged)

@end

id makeArea(void)
{
	return [[areaView alloc]
		initWithFrame:NSMakeRect(0, 0, 100, 100)];
}

uintptr_t modifierFlags(id e)
{
	return fromNSUInteger([toNSEvent(e) modifierFlags]);
}

// TODO move getTranslatedEventPoint() here

intptr_t buttonNumber(id e)
{
	return fromNSInteger([toNSEvent(e) buttonNumber]);
}

uintptr_t pressedMouseButtons(void)
{
	return fromNSUInteger([NSEvent pressedMouseButtons]);
}

uintptr_t keyCode(id e)
{
	return (uintptr_t) ([toNSEvent(e) keyCode]);
}
