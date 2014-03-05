/* 28 february 2014 */

/*
This includes all Objective-C runtime headers for convenience. It also creates wrappers around objc_msgSend() out of necessity.

cgo doesn't support calling variable argument list C functions, so objc_msgSend() cannot be called directly.

Furthermore, Objective-C selectors work by basically sending the arguments to objc_msgSend() verbatim across the wire. This basically means we're stuck making wrapper functions for every possible argument list. What fun!

The format should be self-explanatory.
*/

/* for some reason I now have to use an include guard after commit [master 9b4e30c] ("Started to build a single global delegate object; now to fix issues.") */
#ifndef _OBJC_DARWIN_H_
#define _OBJC_DARWIN_H_

#include <objc/message.h>
#include <objc/objc.h>
#include <objc/runtime.h>

#include <stdint.h>

/* for delegate_darwin.go */
extern Class NilClass;

/* for listbox_darwin.go */
extern id *_NSObservedObjectKey;

static inline id objc_msgSend_noargs(id obj, SEL sel)
{
	return objc_msgSend(obj, sel);
}

struct xrect {
	int64_t x;
	int64_t y;
	int64_t width;
	int64_t height;
};

extern struct xrect objc_msgSend_stret_rect_noargs(id obj, SEL sel);

struct xsize {
	int64_t width;
	int64_t height;
};

extern struct xsize objc_msgSend_stret_size_noargs(id obj, SEL sel);

extern uintptr_t objc_msgSend_uintret_noargs(id objc, SEL sel);

extern intptr_t objc_msgSend_intret_noargs(id obj, SEL sel);

#define m1(name, type1) \
	static inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a) \
	{ \
		return objc_msgSend(obj, sel, a); \
	}

#define m2(name, type1, type2) \
	static inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a, type2 b) \
	{ \
		return objc_msgSend(obj, sel, a, b); \
	}

#define m3(name, type1, type2, type3) \
	static inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a, type2 b, type3 c) \
	{ \
		return objc_msgSend(obj, sel, a, b, c); \
	}

#define m4(name, type1, type2, type3, type4) \
	static inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a, type2 b, type3 c, type4 d) \
	{ \
		return objc_msgSend(obj, sel, a, b, c, d); \
	}

m1(str, char *)		/* TODO Go string? */
m1(id, id)
extern id _objc_msgSend_rect(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h);
m1(sel, SEL)
extern id _objc_msgSend_uint(id obj, SEL sel, uintptr_t a);
m1(ptr, void *)
m1(bool, BOOL)
extern id objc_msgSend_int(id obj, SEL sel, intptr_t a);
m1(double, double)

m2(id_id, id, id)
extern id _objc_msgSend_rect_bool(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h, BOOL b);
extern id objc_msgSend_id_int(id obj, SEL sel, id a, intptr_t b);
extern id objc_msgSend_id_uint(id obj, SEL sel, id a, uintptr_t b);

m3(id_id_id, id, id, id)
m3(sel_id_bool, SEL, id, BOOL)

extern id _objc_msgSend_rect_uint_uint_bool(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h, uintptr_t b, uintptr_t c, BOOL d);
m4(id_sel_id_id, id, SEL, id, id)
m4(id_id_id_id, id, id, id, id)

/* for listbox_darwin.go */
extern uintptr_t *NSIndexSetEntries(id, uintptr_t);

#endif
