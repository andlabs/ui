// 13 may 2014

#include "objc_darwin.h"
#include "area_darwin.h"
#include "_cgo_export.h"
#include <AppKit/NSView.h>
#include <AppKit/NSTrackingArea.h>
#include <Foundation/NSGeometry.h>
#include <AppKit/NSEvent.h>
#include <AppKit/NSBitmapImageRep.h>

#define to(T, x) ((T *) (x))
#define toNSEvent(x) to(NSEvent, (x))
#define toAreaView(x) to(areaView, (x))

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
	trackingArea = [[NSTrackingArea alloc]
		initWithRect:[self bounds]
		// this bit mask (except for NSTrackingInVisibleRect, which was added later to prevent events from being triggered outside the visible area of the Area) comes from https://github.com/andlabs/misctestprogs/blob/master/cocoaviewmousetest.m (and I wrote this bit mask on 25 april 2014) and yes I know it includes enter/exit even though we don't watch those events; it probably won't really matter anyway but if it does I can change it easily
		options:(NSTrackingMouseEnteredAndExited | NSTrackingMouseMoved | NSTrackingActiveAlways | NSTrackingEnabledDuringMouseDrag | NSTrackingInVisibleRect)
		owner:self
		userInfo:nil];
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

void drawImage(void *pixels, intptr_t width, intptr_t height, intptr_t stride, intptr_t xdest, intptr_t ydest)
{
	unsigned char *planes[1];			// NSBitmapImageRep wants an array of planes; we have one plane
	NSBitmapImageRep *bitmap;

	planes[0] = (unsigned char *) pixels;
	bitmap = [[NSBitmapImageRep alloc]
		initWithBitmapDataPlanes:planes
		pixelsWide:toNSInteger(width)
		pixelsHigh:toNSInteger(height)
		bitsPerSample:8
		samplesPerPixel:4
		hasAlpha:YES
		isPlanar:NO
		colorSpaceName:NSCalibratedRGBColorSpace		// TODO NSDeviceRGBColorSpace?
		bitmapFormat:0		// this is where the flag for placing alpha first would go if alpha came first; the default is alpha last, which is how we're doing things (otherwise the docs say "Color planes are arranged in the standard order—for example, red before green before blue for RGB color."); this is also where the flag for non-premultiplied colors would go if we used it (the default is alpha-premultiplied)
		bytesPerRow:toNSInteger(stride)
		bitsPerPixel:32];
	// TODO this CAN fali; check error
	[bitmap drawInRect:NSMakeRect((CGFloat) xdest, (CGFloat) ydest, (CGFloat) width, (CGFloat) height)
		fromRect:NSZeroRect		// draw whole image
		operation:NSCompositeSourceOver
		fraction:1.0
		respectFlipped:YES
		hints:nil];
	[bitmap release];
}

uintptr_t modifierFlags(id e)
{
	return fromNSUInteger([toNSEvent(e) modifierFlags]);
}

struct xpoint getTranslatedEventPoint(id area, id e)
{
	NSPoint p;
	struct xpoint q;

	p = [toAreaView(area) convertPoint:[toNSEvent(e) locationInWindow] fromView:nil];
	q.x = (intptr_t) p.x;
	q.y = (intptr_t) p.y;
	return q;
}

intptr_t buttonNumber(id e)
{
	return fromNSInteger([toNSEvent(e) buttonNumber]);
}

intptr_t clickCount(id e)
{
	return fromNSInteger([toNSEvent(e) clickCount]);
}

uintptr_t pressedMouseButtons(void)
{
	return fromNSUInteger([NSEvent pressedMouseButtons]);
}

uintptr_t keyCode(id e)
{
	return (uintptr_t) ([toNSEvent(e) keyCode]);
}
