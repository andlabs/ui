/* 28 february 2014 */

/*
I wanted to avoid invoking Objective-C directly, preferring to do everything directly with the API. However, there are some things that simply cannot be done too well; for those situations, there's this. It does use the Objective-C runtime, eschewing the actual Objective-C part of this being an Objective-C file.

The main culprits are:
- data types listed as being defined in nonexistent headers
- 32-bit/64-bit type differences that are more than just a different typedef
- wrong documentation
though this is not always the case.
*/

#include "objc_darwin.h"

#include <stdlib.h>

#include <Foundation/NSGeometry.h>
#include <AppKit/NSKeyValueBinding.h>
#include <AppKit/NSEvent.h>
#include <AppKit/NSGraphics.h>
#include <AppKit/NSBitmapImageRep.h>
#include <AppKit/NSCell.h>
#include <AppKit/NSApplication.h>
#include <AppKit/NSTrackingArea.h>

/*
These are all the selectors and class IDs used by the functions below.
*/

static id c_NSEvent;				/* makeDummyEvent() */
static SEL s_newEvent;
static id c_NSBitmapImageRep;	/* drawImage() */
static SEL s_alloc;
static SEL s_initWithBitmapDataPlanes;
static SEL s_drawInRect;
static SEL s_release;
static SEL s_locationInWindow;		/* getTranslatedEventPoint() */
static SEL s_convertPointFromView;
static id c_NSFont;
static SEL s_setFont;				/* objc_setFont() */
static SEL s_systemFontOfSize;
static SEL s_systemFontSizeForControlSize;
static id c_NSTrackingArea;
static SEL s_bounds;
static SEL s_initTrackingArea;

void initBleh()
{
	c_NSEvent = objc_getClass("NSEvent");
	s_newEvent = sel_getUid("otherEventWithType:location:modifierFlags:timestamp:windowNumber:context:subtype:data1:data2:");
	c_NSBitmapImageRep = objc_getClass("NSBitmapImageRep");
	s_alloc = sel_getUid("alloc");
	s_initWithBitmapDataPlanes = sel_getUid("initWithBitmapDataPlanes:pixelsWide:pixelsHigh:bitsPerSample:samplesPerPixel:hasAlpha:isPlanar:colorSpaceName:bitmapFormat:bytesPerRow:bitsPerPixel:");
	s_drawInRect = sel_getUid("drawInRect:fromRect:operation:fraction:respectFlipped:hints:");
	s_release = sel_getUid("release");
	s_locationInWindow = sel_getUid("locationInWindow");
	s_convertPointFromView = sel_getUid("convertPoint:fromView:");
	c_NSFont = objc_getClass("NSFont");
	s_setFont = sel_getUid("setFont:");
	s_systemFontOfSize = sel_getUid("systemFontOfSize:");
	s_systemFontSizeForControlSize = sel_getUid("systemFontSizeForControlSize:");
	c_NSTrackingArea = objc_getClass("NSTrackingArea");
	s_bounds = sel_getUid("bounds");
	s_initTrackingArea = sel_getUid("initWithRect:options:owner:userInfo:");
}

/*
See uitask_darwin.go: we need to synthesize a NSEvent so -[NSApplication stop:] will work. We cannot simply init the default NSEvent though (it throws an exception) so we must do it "the right way". This involves a very convoluted initializer; we'll just do it here to keep things clean on the Go side (this will only be run once anyway, on program exit).
*/

id makeDummyEvent()
{
	return objc_msgSend(c_NSEvent, s_newEvent,
		(NSUInteger) NSApplicationDefined,			/* otherEventWithType: */
		NSMakePoint(0, 0),						/* location: */
		(NSUInteger) 0,							/* modifierFlags: */
		(double) 0,							/* timestamp: */
		(NSInteger) 0,							/* windowNumber: */
		nil,									/* context: */
		(short) 0,								/* subtype: */
		(NSInteger) 0,							/* data1: */
		(NSInteger) 0);							/* data2: */
}

/*
[NSView drawRect:] needs to be overridden in our Area subclass. This takes a NSRect, which I'm not sure how to encode, so we're going to have to use @encode() and hope for the best for portability.
*/

extern void areaView_drawRect(id, struct xrect);

static void __areaView_drawRect(id self, SEL sel, NSRect r)
{
	struct xrect t;

	t.x = (int64_t) r.origin.x;
	t.y = (int64_t) r.origin.y;
	t.width = (int64_t) r.size.width;
	t.height = (int64_t) r.size.height;
	areaView_drawRect(self, t);
}

void *_areaView_drawRect = (void *) __areaView_drawRect;

/*
this and one below it are the only objective-c feature you'll see here

unfortunately NSRect both varies across architectures and is passed as just a structure, so its encoding has to be computed at compile time
because @encode() is NOT A LITERAL, we're going to just stick it all the way back in objc_darwin.go
see also: http://stackoverflow.com/questions/6812035/adding-methods-dynamically
*/
char *encodedNSRect = @encode(NSRect);

/*
the NSBitmapImageRep constructor is complex; put it here
the only way to draw a NSBitmapImageRep in a flipped NSView is to use the most complex drawing method; put it here too
*/

/*
hey guys you know what's fun? 32-bit ABI changes!
*/
static BOOL (*objc_msgSend_drawInRect)(id, SEL, NSRect, NSRect, NSCompositingOperation, CGFloat, BOOL, id) =
	(BOOL (*)(id, SEL, NSRect, NSRect, NSCompositingOperation, CGFloat, BOOL, id)) objc_msgSend;

void drawImage(void *pixels, int64_t width, int64_t height, int64_t stride, int64_t xdest, int64_t ydest)
{
	unsigned char *planes[1];			/* NSBitmapImageRep wants an array of planes; we have one plane */
	id bitmap;

	bitmap = objc_msgSend(c_NSBitmapImageRep, s_alloc);
	planes[0] = (unsigned char *) pixels;
	bitmap = objc_msgSend(bitmap, s_initWithBitmapDataPlanes,
		planes,								/* initWithBitmapDataPlanes: */
		(NSInteger) width,						/* pixelsWide: */
		(NSInteger) height,						/* pixelsHigh: */
		(NSInteger) 8,							/* bitsPerSample: */
		(NSInteger) 4,							/* samplesPerPixel: */
		(BOOL) YES,							/* hasAlpha: */
		(BOOL) NO,							/* isPlanar: */
		NSCalibratedRGBColorSpace,				/* colorSpaceName: | TODO NSDeviceRGBColorSpace? */
		(NSBitmapFormat) 0,					/* bitmapFormat: | this is where the flag for placing alpha first would go if alpha came first; the default is alpha last, which is how we're doing things (otherwise the docs say "Color planes are arranged in the standard order—for example, red before green before blue for RGB color."); this is also where the flag for non-premultiplied colors would go if we used it (the default is alpha-premultiplied) */
		(NSInteger) stride,						/* bytesPerRow: */
		(NSInteger) 32);						/* bitsPerPixel: */
	/* TODO this CAN fail; check error */
	objc_msgSend_drawInRect(bitmap, s_drawInRect,
		NSMakeRect((CGFloat) xdest, (CGFloat) ydest,
			(CGFloat) width, (CGFloat) height),					/* drawInRect: */
		NSZeroRect,										/* fromRect: | draw whole image */
		(NSCompositingOperation) NSCompositeSourceOver,		/* op: */
		(CGFloat) 1.0,										/* fraction: */
		(BOOL) YES,										/* respectFlipped: */
		nil);												/* hints: */
	objc_msgSend(bitmap, s_release);
}

/*
more NSPoint fumbling
*/

static NSPoint (*objc_msgSend_stret_point)(id, SEL, ...) =
	(NSPoint (*)(id, SEL, ...)) objc_msgSend;

struct xpoint getTranslatedEventPoint(id self, id event)
{
	NSPoint p;
	struct xpoint ret;

	p = objc_msgSend_stret_point(event, s_locationInWindow);
	p = objc_msgSend_stret_point(self, s_convertPointFromView,
		p,			/* convertPoint: */
		nil);			/* fromView: */
	ret.x = (int64_t) p.x;
	ret.y = (int64_t) p.y;
	return ret;
}
