/* 28 february 2014 */

/*
This includes all Objective-C runtime headers for convenience. It also creates wrappers around objc_msgSend() out of necessity.

cgo doesn't support calling variable argument list C functions, so objc_msgSend() cannot be called directly.

Furthermore, Objective-C selectors work by basically sending the arguments to objc_msgSend() verbatim across the wire. This basically means we're stuck making wrapper functions for every possible argument list. What fun!

The format should be self-explanatory.
*/

#include <objc/message.h>
#include <objc/objc.h>
#include <objc/runtime.h>

#include <stdint.h>

inline id objc_msgSend_noargs(id obj, SEL sel)
{
	return objc_msgSend(obj, sel);
}

struct xsize {
	int64_t width;
	int64_t height;
};

extern struct xsize objc_msgsSend_stret_size_noargs(id obj, SEL sel);

#define m1(name, type1) \
	inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a) \
	{ \
		return objc_msgSend(obj, sel, a); \
	}

#define m2(name, type1, type2) \
	inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a, type2 b) \
	{ \
		return objc_msgSend(obj, sel, a, b); \
	}

#define m3(name, type1, type2, type3) \
	inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a, type2 b, type3 c) \
	{ \
		return objc_msgSend(obj, sel, a, b, c); \
	}

#define m4(name, type1, type2, type3, type4) \
	inline id objc_msgSend_ ## name (id obj, SEL sel, type1 a, type2 b, type3 c, type4 d) \
	{ \
		return objc_msgSend(obj, sel, a, b, c, d); \
	}

m1(str, char *)		/* TODO Go string? */
m1(id, id)
extern id _objc_msgSend_rect(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h);
m1(sel, SEL)
extern id _objc_msgSend_uint(id obj, SEL sel, uintptr_t a);

m2(id_id, id, id)

m3(id_id_id, id, id, id)
m3(sel_id_bool, SEL, id, BOOL)

extern id _objc_msgSend_rect_uint_uint_bool(id obj, SEL sel, int64_t x, int64_t y, int64_t w, int64_t h, uintptr_t b, uintptr_t c, BOOL d);
m4(id_sel_id_id, id, SEL, id, id)
