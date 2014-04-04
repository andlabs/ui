/* 28 february 2014 */

/*
I wanted to avoid invoking Objective-C directly, preferring to do everything directly with the API. However, there are some things that simply cannot be done too well; for those situations, there's this. It does use the Objective-C runtime, eschewing the actual Objective-C part of this being an Objective-C file.

The main culprits are:
- data types listed as being defined in nonexistent headers
- 32-bit/64-bit type differences that are more than just a different typedef
- wrong documentation
though this is not always the case.

Go wrapper functions (bleh_darwin.go) call these directly and take care of stdint.h -> Go type conversions.
*/

#include "objc_darwin.h"

#include <stdlib.h>

#include <Foundation/NSGeometry.h>
#include <AppKit/NSKeyValueBinding.h>
#include <AppKit/NSEvent.h>
#include <AppKit/NSGraphics.h>
#include <AppKit/NSBitmapImageRep.h>
#include <AppKit/NSCell.h>

/* exception to the above: cgo doesn't like Nil and delegate_darwin.go has //export so I can't have this there */
Class NilClass = Nil;

/* used by listbox_darwin.go; requires NSString */
id *_NSObservedObjectKey = (id *) (&NSObservedObjectKey);

/*
NSUInteger is listed as being in <objc/NSObjCRuntime.h>... which doesn't exist. Rather than relying on undocumented header file locations or explicitly typedef-ing NSUInteger to the (documented) unsigned long, I'll just place things here for maximum safety. I use uintptr_t as that should encompass every possible unsigned long.
*/

uintptr_t objc_msgSend_uintret_noargs(id obj, SEL sel)
{
	return (uintptr_t) ((NSUInteger) objc_msgSend(obj, sel));
}

id objc_msgSend_uint(id obj, SEL sel, uintptr_t a)
{
	return objc_msgSend(obj, sel, (NSUInteger) a);
}

id objc_msgSend_id_uint(id obj, SEL sel, id a, uintptr_t b)
{
	return objc_msgSend(obj, sel, a, (NSUInteger) b);
}

/*
same as above, but for NSInteger
*/

intptr_t objc_msgSend_intret_noargs(id obj, SEL sel)
{
	return (intptr_t) ((NSInteger) objc_msgSend(obj, sel));
}

id objc_msgSend_int(id obj, SEL sel, intptr_t a)
{
	return objc_msgSend(obj, sel, (NSInteger) a);
}

id objc_msgSend_id_int(id obj, SEL sel, id a, intptr_t b)
{
	return objc_msgSend(obj, sel, a, (NSInteger) b);
}

/*
same as above, but for unsigned short
*/

uintptr_t objc_msgSend_ushortret_noargs(id obj, SEL sel)
{
	return (uintptr_t) ((unsigned short) objc_msgSend(obj, sel));
}

/*
These are the objc_msgSend() wrappers around NSRect. The problem is that while on 32-bit systems, NSRect is a concrete structure, on 64-bit systems it's just a typedef to CGRect. While in practice just using CGRect everywhere seems to work, better to be safe than sorry.

I use int64_t for maximum safety, as my coordinates are stored as Go ints and Go int -> C int (which is what is documented as happening) isn't reliable.
*/

/*
This is not documented in the docs, but is in various places on apple.com. In fact, the docs are actually WRONG: they say you pass a pointer to the structure as the first argument to objc_msgSend_stret()! And there might be some cases where we can't use stret because the struct is small enough...
*/
static NSRect (*objc_msgSend_stret_rect)(id, SEL, ...) =
	(NSRect (*)(id, SEL, ...)) objc_msgSend_stret;

struct xrect objc_msgSend_stret_rect_noargs(id obj, SEL sel)
{
	NSRect s;
	struct xrect t;

	s = objc_msgSend_stret_rect(obj, sel);
	t.x = (int64_t) s.origin.x;
	t.y = (int64_t) s.origin.y;
	t.width = (int64_t) s.size.width;
	t.height = (int64_t) s.size.height;
	return t;
}

#define OurRect() (NSMakeRect((CGFloat) x, (CGFloat) y, (CGFloat) w, (CGFloat) h))

id objc_msgSend_rect(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h)
{
	return objc_msgSend(obj, sel, OurRect());
}

id objc_msgSend_rect_bool(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h, BOOL b)
{
	return objc_msgSend(obj, sel, OurRect(), b);
}

id objc_msgSend_rect_uint_uint_bool(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h, uintptr_t b, uintptr_t c, BOOL d)
{
	return objc_msgSend(obj, sel, OurRect(), (NSUInteger) b, (NSUInteger) c, d);
}

/*
Same as NSRect above, but for NSSize now.
*/

/*
...like this one. (Note which function is being cast below.) This is an Intel-specific optimization; though this code won't run on PowerPC Macs (Go, and thus package ui, requires 10.6), if desktop ARM becomes a thing all bets are off. (tl;dr TODO)
*/
static NSSize (*objc_msgSend_stret_size)(id, SEL, ...) =
	(NSSize (*)(id, SEL, ...)) objc_msgSend;

struct xsize objc_msgSend_stret_size_noargs(id obj, SEL sel)
{
	NSSize s;
	struct xsize t;

	s = objc_msgSend_stret_size(obj, sel);
	t.width = (int64_t) s.width;
	t.height = (int64_t) s.height;
	return t;
}

/*
and again for NSPoint
*/

id objc_msgSend_point(id obj, SEL sel, int64_t x, int64_t y)
{
	return objc_msgSend(obj, sel, NSMakePoint((CGFloat) x, (CGFloat) y));
}

/*
This is a doozy: it deals with a NSUInteger array needed for this one selector, and converts them all into a uintptr_t array so we can use it from Go. The two arrays are created at runtime with malloc(); only the NSUInteger one is freed here, while Go frees the returned one. It's not optimal.
*/

static SEL getIndexes;
static BOOL getIndexes_init = NO;		/* because we can't initialize it out here */

uintptr_t *NSIndexSetEntries(id indexset, uintptr_t count)
{
	NSUInteger *nsuints;
	uintptr_t *ret;
	uintptr_t i;
	size_t countsize;

	if (getIndexes_init == NO) {
		getIndexes = sel_getUid("getIndexes:maxCount:inIndexRange:");
		getIndexes_init = YES;
	}
	countsize = (size_t) count;
	nsuints = (NSUInteger *) malloc(countsize * sizeof (NSUInteger));
	/* TODO check return value */
	objc_msgSend(indexset, getIndexes,
		nsuints, (NSUInteger) count, nil);
	ret = (uintptr_t *) malloc(countsize * sizeof (uintptr_t));
	for (i = 0; i < count; i++) {
		ret[i] = (uintptr_t) nsuints[i];
	}
	free(nsuints);
	return ret;
}

/*
See uitask_darwin.go: we need to synthesize a NSEvent so -[NSApplication stop:] will work. We cannot simply init the default NSEvent though (it throws an exception) so we must do it "the right way". This involves a very convoluted initializer; we'll just do it here to keep things clean on the Go side (this will only be run once anyway, on program exit).
*/

static id c_NSEvent;
static SEL s_newEvent;
static BOOL newEvent_init = NO;

id makeDummyEvent()
{
	if (newEvent_init == NO) {
		c_NSEvent = objc_getClass("NSEvent");
		s_newEvent = sel_getUid("otherEventWithType:location:modifierFlags:timestamp:windowNumber:context:subtype:data1:data2:");
		newEvent_init = YES;
	}
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

static void _areaView_drawRect(id self, SEL sel, NSRect r)
{
	struct xrect t;

	t.x = (int64_t) r.origin.x;
	t.y = (int64_t) r.origin.y;
	t.width = (int64_t) r.size.width;
	t.height = (int64_t) r.size.height;
	areaView_drawRect(self, t);
}

/* the only objective-c feature you'll see here */
/* TODO correct? "v@:" @encode(NSRect) complained about missing ; */
static char *avdrType = @encode(void(id, SEL, NSRect));

static SEL drawRect;
static BOOL drawRect_init = NO;

BOOL addAreaViewDrawMethod(Class what)
{
	if (drawRect_init == NO) {
		drawRect = sel_getUid("drawRect:");
		drawRect_init = YES;
	}
	return class_addMethod(what, drawRect, (IMP) _areaView_drawRect, avdrType);
}

/*
the NSBitmapImageRep constructor is complex; put it here
the only way to draw a NSBitmapImageRep in a flipped NSView is to use the most complex drawing method; put it here too
*/

static id c_NSBitmapImageRep;
static SEL s_alloc;
static SEL s_initWithBitmapDataPlanes;
static SEL s_drawInRect;
static SEL s_release;
static BOOL drawImage_init = NO;

void drawImage(void *pixels, int64_t width, int64_t height, int64_t stride, int64_t xdest, int64_t ydest)
{
	unsigned char *planes[1];			/* NSBitmapImageRep wants an array of planes; we have one plane */
	id bitmap;

	if (drawImage_init == NO) {
		c_NSBitmapImageRep = objc_getClass("NSBitmapImageRep");
		s_alloc = sel_getUid("alloc");
		s_initWithBitmapDataPlanes = sel_getUid("initWithBitmapDataPlanes:pixelsWide:pixelsHigh:bitsPerSample:samplesPerPixel:hasAlpha:isPlanar:colorSpaceName:bitmapFormat:bytesPerRow:bitsPerPixel:");
		s_drawInRect = sel_getUid("drawInRect:fromRect:operation:fraction:respectFlipped:hints:");
		s_release = sel_getUid("release");
		drawImage_init = YES;
	}
	bitmap = objc_msgSend(c_NSBitmapImageRep, s_alloc);
	planes[0] = (unsigned char *) pixels;
	bitmap = objc_msgSend(bitmap, s_initWithBitmapDataPlanes,
		planes,												/* initWithBitmapDataPlanes: */
		(NSInteger) width,										/* pixelsWide: */
		(NSInteger) height,										/* pixelsHigh: */
		(NSInteger) 8,											/* bitsPerSample: */
		(NSInteger) 4,											/* samplesPerPixel: */
		(BOOL) YES,											/* hasAlpha: */
		(BOOL) NO,											/* isPlanar: */
		NSCalibratedRGBColorSpace,								/* colorSpaceName: | TODO NSDeviceRGBColorSpace? */
		(NSBitmapFormat) NSAlphaNonpremultipliedBitmapFormat,		/* bitmapFormat: | this is where the flag for placing alpha first would go if alpha came first; the default is alpha last, which is how we're doing things (otherwise the docs say "Color planes are arranged in the standard order—for example, red before green before blue for RGB color.") */
		(NSInteger) stride,										/* bytesPerRow: */
		(NSInteger) 32);										/* bitsPerPixel: */
	objc_msgSend(bitmap, s_drawInRect,
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

static SEL s_locationInWindow;
static SEL s_convertPointFromView;
static BOOL getTranslatedEventPoint_init = NO;

static NSPoint (*objc_msgSend_stret_point)(id, SEL, ...) =
	(NSPoint (*)(id, SEL, ...)) objc_msgSend;

struct xpoint getTranslatedEventPoint(id self, id event)
{
	NSPoint p;
	struct xpoint ret;

	if (getTranslatedEventPoint_init == NO) {
		s_locationInWindow = sel_getUid("locationInWindow");
		s_convertPointFromView = sel_getUid("convertPoint:fromView:");
		getTranslatedEventPoint_init = YES;
	}
	p = objc_msgSend_stret_point(event, s_locationInWindow);
	p = objc_msgSend_stret_point(self, s_convertPointFromView,
		p,			/* convertPoint: */
		nil);			/* fromView: */
	ret.x = (int64_t) p.x;
	ret.y = (int64_t) p.y;
	return ret;
}

/*
we don't need this here technically — it can be done in Go just fine — but it's easier here
*/

static id c_NSFont;
static SEL s_setFont;
static SEL s_systemFontOfSize;
static SEL s_systemFontSizeForControlSize;
static BOOL setFont_init = NO;

static CGFloat (*objc_msgSend_cgfloatret)(id, SEL, ...) =
	(CGFloat (*)(id, SEL, ...)) objc_msgSend_fpret;

void objc_setFont(id what, unsigned int csize)
{
	CGFloat size;

	if (setFont_init == NO) {
		c_NSFont = objc_getClass("NSFont");
		s_setFont = sel_getUid("setFont:");
		s_systemFontOfSize = sel_getUid("systemFontOfSize:");
		s_systemFontSizeForControlSize = sel_getUid("systemFontSizeForControlSize:");
		setFont_init = YES;
	}

	size = objc_msgSend_cgfloatret(c_NSFont, s_systemFontSizeForControlSize, (NSControlSize) csize);
	objc_msgSend(what, s_setFont,
		objc_msgSend(c_NSFont, s_systemFontOfSize, size));
}
