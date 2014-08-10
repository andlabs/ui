// 13 may 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSEvent(x) ((NSEvent *) (x))
// TODO rename to goAreaView
#define toAreaView(x) ((areaView *) (x))
#define toNSView(x) ((NSView *) (x))

#define toNSInteger(x) ((NSInteger) (x))
#define fromNSInteger(x) ((intptr_t) (x))
#define toNSUInteger(x) ((NSUInteger) (x))
#define fromNSUInteger(x) ((uintptr_t) (x))

@interface areaView : NSView {
@public
	void *goarea;
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
	areaView_drawRect(self, rect, self->goarea);
}

- (BOOL)isFlipped
{
	return YES;
}

- (BOOL)acceptsFirstResponder
{
	return YES;
}

// this will have the Area receive a click that switches to the Window it is in from another one
- (BOOL)acceptsFirstMouse:(NSEvent *)e
{
	return YES;
}

- (void)retrack
{
	self->trackingArea = [[NSTrackingArea alloc]
		initWithRect:[self bounds]
		// this bit mask (except for NSTrackingInVisibleRect, which was added later to prevent events from being triggered outside the visible area of the Area) comes from https://github.com/andlabs/misctestprogs/blob/master/cocoaviewmousetest.m (and I wrote this bit mask on 25 april 2014) and yes I know it includes enter/exit even though we don't watch those events; it probably won't really matter anyway but if it does I can change it easily
		options:(NSTrackingMouseEnteredAndExited | NSTrackingMouseMoved | NSTrackingActiveAlways | NSTrackingEnabledDuringMouseDrag | NSTrackingInVisibleRect)
		owner:self
		userInfo:nil];
	[self addTrackingArea:self->trackingArea];
}

- (void)updateTrackingAreas
{
	[self removeTrackingArea:self->trackingArea];
	[self->trackingArea release];
	[self retrack];
}

#define event(m, f) \
	- (void)m:(NSEvent *)e \
	{ \
		f(self, e, self->goarea); \
	}
event(mouseMoved, areaView_mouseMoved_mouseDragged)
event(mouseDragged, areaView_mouseMoved_mouseDragged)
event(rightMouseDragged, areaView_mouseMoved_mouseDragged)
event(otherMouseDragged, areaView_mouseMoved_mouseDragged)
event(mouseDown, areaView_mouseDown)
event(rightMouseDown, areaView_mouseDown)
event(otherMouseDown, areaView_mouseDown)
event(mouseUp, areaView_mouseUp)
event(rightMouseUp, areaView_mouseUp)
event(otherMouseUp, areaView_mouseUp)
event(keyDown, areaView_keyDown)
event(keyUp, areaView_keyUp)
event(flagsChanged, areaView_flagsChanged)

@end

Class getAreaClass(void)
{
	return [areaView class];
}

id newArea(void *goarea)
{
	areaView *a;

	a = [[areaView alloc] initWithFrame:NSZeroRect];
	a->goarea = goarea;
	return (id) a;
}

BOOL drawImage(void *pixels, intptr_t width, intptr_t height, intptr_t stride, intptr_t xdest, intptr_t ydest)
{
	unsigned char *planes[1];			// NSBitmapImageRep wants an array of planes; we have one plane
	NSBitmapImageRep *bitmap;
	BOOL success;

	planes[0] = (unsigned char *) pixels;
	bitmap = [[NSBitmapImageRep alloc]
		initWithBitmapDataPlanes:planes
		pixelsWide:toNSInteger(width)
		pixelsHigh:toNSInteger(height)
		bitsPerSample:8
		samplesPerPixel:4
		hasAlpha:YES
		isPlanar:NO
		// NSCalibratedRGBColorSpace changes the colors; let's not
		// thanks to JtRip in irc.freenode.net/#macdev
		colorSpaceName:NSDeviceRGBColorSpace
		bitmapFormat:0		// this is where the flag for placing alpha first would go if alpha came first; the default is alpha last, which is how we're doing things (otherwise the docs say "Color planes are arranged in the standard order—for example, red before green before blue for RGB color."); this is also where the flag for non-premultiplied colors would go if we used it (the default is alpha-premultiplied)
		bytesPerRow:toNSInteger(stride)
		bitsPerPixel:32];
	success = [bitmap drawInRect:NSMakeRect((CGFloat) xdest, (CGFloat) ydest, (CGFloat) width, (CGFloat) height)
		fromRect:NSZeroRect		// draw whole image
		operation:NSCompositeSourceOver
		fraction:1.0
		respectFlipped:YES
		hints:nil];
	[bitmap release];
	return success;
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

void areaRepaintAll(id view)
{
	[toNSView(view) display];
}
