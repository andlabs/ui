/* 28 february 2014 */

/*
I wanted to avoid invoking Objective-C directly, preferring to do everything directly with the API. However, there are some things that simply cannot be done too well; for those situations, there's this. It does use the Objective-C runtime, eschewing the actual Objective-C part of this being an Objective-C file.

The main culprits are:
- data types listed as being defined in nonexistent headers
- 32-bit/64-bit type differences that are more than just a different typedef
- wrong documentation

Go wrapper functions (bleh_darwin.go) call these directly and take care of stdint.h -> Go type conversions.
*/

#include "objc_darwin.h"

#include <stdlib.h>

#include <Foundation/NSGeometry.h>
#include <AppKit/NSKeyValueBinding.h>

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

id _objc_msgSend_uint(id obj, SEL sel, uintptr_t a)
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

id _objc_msgSend_rect(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h)
{
	return objc_msgSend(obj, sel, OurRect());
}

id _objc_msgSend_rect_bool(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h, BOOL b)
{
	return objc_msgSend(obj, sel, OurRect(), b);
}

id _objc_msgSend_rect_uint_uint_bool(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h, uintptr_t b, uintptr_t c, BOOL d)
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
