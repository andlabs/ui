// 12 august 2018

package ui

// #include <stdlib.h>
// #include "ui.h"
// #include "util.h"
// typedef struct pkguiCColor pkguiCColor;
// struct pkguiCColor { double *r; double *g; double *b; double *a; };
// static inline pkguiCColor pkguiNewCColor(void)
// {
// 	pkguiCColor c;
// 
// 	c.r = (double *) pkguiAlloc(4 * sizeof (double));
// 	c.g = c.r + 1;
// 	c.b = c.g + 1;
// 	c.a = c.b + 1;
// 	return c;
// }
// static inline void pkguiFreeCColor(pkguiCColor c)
// {
// 	free(c.r);
// }
// static inline uiUnderlineColor *pkguiNewUnderlineColor(void)
// {
// 	return (uiUnderlineColor *) pkguiAlloc(sizeof (uiUnderlineColor));
// }
// static inline void pkguiFreeUnderlineColor(uiUnderlineColor *c)
// {
// 	free(c);
// }
import "C"

// Attribute stores information about an attribute in an
// AttributedString.
//
// The following types can be used as Attributes:
//
// 	- TextFamily
// 	- TextSize
// 	- TextWeight
// 	- TextItalic
// 	- TextStretch
// 	- TextColor
// 	- TextBackground
// 	- Underline
// 	- UnderlineColor
// 	- UnderlineColorCustom
// 	- OpenTypeFeatures
//
// For every Unicode codepoint in the AttributedString, at most one
// value of each attribute type can be applied.
type Attribute interface {
	toLibui() *C.uiAttribute
}

// TextFamily is an Attribute that changes the font family of the text
// it is applied to. Font family names are case-insensitive.
type TextFamily string

func (f TextFamily) toLibui() *C.uiAttribute {
	fstr := C.CString(string(f))
	defer freestr(fstr)
	return C.uiNewFamilyAttribute(fstr)
}

// TextSize is an Attribute that changes the size of the text it is
// applied to, in typographical points.
type TextSize float64

func (s TextSize) toLibui() *C.uiAttribute {
	return C.uiNewSizeAttribute(C.double(size))
}

// TextWeight is an Attribute that changes the weight of the text
// it is applied to. These roughly map to the OS/2 text weight field
// of TrueType and OpenType fonts, or to CSS weight numbers. The
// named constants are nominal values; the actual values may vary
// by font and by OS, though this isn't particularly likely. Any value
// between TextWeightMinimum and TextWeightMaximum,
// inclusive, is allowed.
//
// Note that due to restrictions in early versions of Windows, some
// fonts have "special" weights be exposed in many programs as
// separate font families. This is perhaps most notable with
// Arial Black. Package ui does not do this, even on Windows
// (because the DirectWrite API libui uses on Windows does not do
// this); to specify Arial Black, use family Arial and weight
// TextWeightBlack.
type TextWeight int
const (
	TextWeightMinimum TextWeight = 0
	TextWeightThin TextWeight = 100
	TextWeightUltraLight TextWeight = 200
	TextWeightLight TextWeight = 300
	TextWeightBook TextWeight = 350
	TextWeightNormal TextWeight = 400
	TextWeightMedium TextWeight = 500
	TextWeightSemiBold TextWeight = 600
	TextWeightBold TextWeight = 700
	TextWeightUltraBold TextWeight = 800
	TextWeightHeavy TextWeight = 900
	TextWeightUltraHeavy TextWeight = 950
	TextWeightMaximum TextWeight = 1000
)

func (w TextWeight) toLibui() *C.uiAttribute {
	return C.uiNewWeightAttribute(C.uiTextWeight(w))
}

// TextItalic is an Attribute that changes the italic mode of the text
// it is applied to. Italic represents "true" italics where the slanted
// glyphs have custom shapes, whereas oblique represents italics
// that are merely slanted versions of the normal glyphs. Most fonts
// usually have one or the other.
type TextItalic int
const (
	TextItalicNormal TextItalic = iota
	TextItalicOblique
	TextItalicItalic
)

func (i TextItalic) toLibui() *C.uiAttribute {
	return C.uiNewItalicAttribute(C.uiTextItalic(i))
}

// TextStretch is an Attribute that changes the stretch (also called
// "width") of the text it is applied to.
//
// Note that due to restrictions in early versions of Windows, some
// fonts have "special" stretches be exposed in many programs as
// separate font families. This is perhaps most notable with
// Arial Condensed. Package ui does not do this, even on Windows
// (because the DirectWrite API package ui uses on Windows does
// not do this); to specify Arial Condensed, use family Arial and
// stretch TextStretchCondensed.
type TextStretch int
const (
	TextStretchUltraCondensed TextStretch = iota
	TextStretchExtraCondensed
	TextStretchCondensed
	TextStretchSemiCondensed
	TextStretchNormal
	TextStretchSemiExpanded
	TextStretchExpanded
	TextStretchExtraExpanded
	TextStretchUltraExpanded
)

func (s TextStretch) toLibui() *C.uiAttribute {
	return C.uiNewStretchAttribute(C.uiTextStretch(s))
}

// TextColor is an Attribute that changes the color of the text it is
// applied to.
type TextColor struct {
	R	float64
	G	float64
	B	float64
	A	float64
}

func (c TextColor) toLibui() *C.uiAttribute {
	return C.uiNewColorAttribute(C.double(c.R), C.double(c.G), C.double(c.B), C.double(c.A))
}

// TextBackground is an Attribute that changes the background
// color of the text it is applied to.
type TextBackground struct {
	R	float64
	G	float64
	B	float64
	A	float64
}

func (b TextBackground) toLibui() *C.uiAttribute {
	return C.uiNewBackgroundAttribute(C.double(b.R), C.double(b.G), C.double(b.b), C.double(b.A))
}

// Underline is an Attribute that specifies a type of underline to use
// on text.
type Underline int
const (
	UnderlineNone Underline = iota
	UnderlineSingle
	UnderlineDouble
	UnderlineSuggestion		// wavy or dotted underlines used for spelling/grammar checkers
)

func (u Underline) toLibui() *C.uiAttribute {
	return C.uiNewUnderlineAttribute(C.uiUnderline(u))
}

// UnderlineColor is an Attribute that changes the color of any
// underline on the text it is applied to, regardless of the type of
// underline. In addition to being able to specify the
// platform-specific colors for suggestion underlines here, you can
// also use a custom color with UnderlineColorCustom.
// 
// To use the constants here correctly, pair them with
// UnderlineSuggestion (though they can be used on other types of
// underline as well).
// 
// If an underline type is applied but no underline color is
// specified, the text color is used instead. If an underline color
// is specified without an underline type, the underline color
// attribute is ignored, but not removed from the uiAttributedString.
type UnderlineColor int
const (
	UnderlineColorSpelling UnderlineColor = iota + 1
	UnderlineColorGrammar
	UnderlineColorAuxiliary		// for instance, the color used by smart replacements on macOS or in Microsoft Office
)

func (u UnderlineColor) toLibui() *C.uiAttribute {
	return C.uiNewUnderlineColorAttribute(C.uiUnderlineColor(u), 0, 0, 0, 0)
}

// UnderlineColorCustom is an Attribute like UnderlineColor, except
// it allows specifying a custom color.
type UnderlineColorCustom struct {
	R	float64
	G	float64
	B	float64
	A	float64
}

func (u UnderlineColorCustom) toLibui() *C.uiAttribute {
	return C.uiNewUnderlineColorAttribute(C.uiUnderlineColorCustom, C.double(u.R), C.double(u.G), C.double(u.B), C.double(u.A))
}

// OpenTypeFeatures is an Attribute that represents a set of
// OpenType feature tag-value pairs, for applying OpenType
// features to text. OpenType feature tags are four-character codes
// defined by OpenType that cover things from design features like
// small caps and swashes to language-specific glyph shapes and
// beyond. Each tag may only appear once in any given
// uiOpenTypeFeatures instance. Each value is a 32-bit integer,
// often used as a Boolean flag, but sometimes as an index to choose
// a glyph shape to use.
// 
// If a font does not support a certain feature, that feature will be
// ignored. (TODO verify this on all OSs)
// 
// See the OpenType specification at
// https://www.microsoft.com/typography/otspec/featuretags.htm
// for the complete list of available features, information on specific
// features, and how to use them.
// TODO invalid features
// 
// Note that if a feature is not present in a OpenTypeFeatures,
// the feature is NOT treated as if its value was zero, unlike in Go.
// Script-specific font shaping rules and font-specific feature
// settings may use a different default value for a feature. You
// should likewise NOT treat a missing feature as having a value of
// zero either. Instead, a missing feature should be treated as
// having some unspecified default value.
// 
// Note that despite OpenTypeFeatures being a map, its contents
// are copied by AttributedString. Modifying an OpenTypeFeatures
// after giving it to an AttributedString, or modifying one that comes
// out of an AttributedString, will have no effect.
type OpenTypeFeatures map[OpenTypeTag]uint32

func (o OpenTypeFeatures) toLibui() *C.uiAttribute {
	otf := C.uiNewOpenTypeFeatures()
	defer C.uiFreeOpenTypeFeatures(otf)
	for tag, value := range o {
		a := byte((tag >> 24) & 0xFF)
		b := byte((tag >> 16) & 0xFF)
		c := byte((tag >> 8) & 0xFF)
		d := byte(tag & 0xFF)
		C.uiOpenTypeFeaturesAdd(otf, C.char(a), C.char(b), C.char(c), C.char(d), C.uint32_t(value))
	}
	return C.uiNewFeaturesAttribute(otf)
}

// OpenTypeTag represents a four-byte OpenType feature tag.
type OpenTypeTag uint32

// ToOpenTypeTag converts the four characters a, b, c, and d into
// an OpenTypeTag.
func ToOpenTypeTag(a, b, c, d byte) OpenTypeTag {
	return (uint32(a) << 24) |
		(uint32(b) << 16) |
		(uint32(c) << 8) |
		uint32(d)
}

func attributeFromLibui(a *C.uiAttribute) Attribute {
	switch C.uiAttributeGetType(a) {
	case C.uiAttributeTypeFamily:
		cf := C.uiAttributeFamily(a)
		return TextFamily(C.GoString(cf))
	case C.uiAttributeTypeSize:
		return TextSize(C.uiAttributeSize(a))
	case C.uiAttributeTypeWeight:
		return TextWeight(C.uiAttributeWeight(a))
	case C.uiAttributeTypeItalic:
		return TextItalic(C.uiAttributeItalic(a))
	case C.uiAttributeTypeStretch:
		return TextStretch(C.uiAttributeStretch(a))
	case C.uiAttributeTypeColor:
		cc := C.pkguiNewCColor()
		defer C.pkguiFreeCColor(cc)
		C.uiAttributeColor(a, c.r, c.g, c.b, c.a)
		return TextColor{
			R:	float64(*(c.r)),
			G:	float64(*(c.g)),
			B:	float64(*(c.b)),
			A:	float64(*(c.a)),
		}
	case C.uiAttributeTypeBackground:
		cc := C.pkguiNewCColor()
		defer C.pkguiFreeCColor(cc)
		C.uiAttributeColor(a, c.r, c.g, c.b, c.a)
		return TextBackground{
			R:	float64(*(c.r)),
			G:	float64(*(c.g)),
			B:	float64(*(c.b)),
			A:	float64(*(c.a)),
		}
	case C.uiAttributeTypeUnderline:
		return Underline(C.uiAttributeUnderline(a))
	case C.uiAttributeTypeUnderlineColor:
		cu := C.pkguiNewUnderlineColor()
		defer C.pkguiFreeUnderlineColor(cu)
		cc := C.pkguiNewCColor()
		defer C.pkguiFreeCColor(cc)
		C.uiAttributeUnderlineColor(a, cu, c.r, c.g, c.b, c.a)
		if *cu == C.uiAttributeUnderlineColorCustom {
			return UnderlineColorCustom{
				R:	float64(*(c.r)),
				G:	float64(*(c.g)),
				B:	float64(*(c.b)),
				A:	float64(*(c.a)),
			}
		}
		return UnderlineColor(*cu)
	case C.uiAttributeTypeFeatures:
		// TODO
	}
	panic("unreachable")
}

///////// TODOTODO

// uiAttributedString represents a string of UTF-8 text that can
// optionally be embellished with formatting attributes. libui
// provides the list of formatting attributes, which cover common
// formatting traits like boldface and color as well as advanced
// typographical features provided by OpenType like superscripts
// and small caps. These attributes can be combined in a variety of
// ways.
//
// Attributes are applied to runs of Unicode codepoints in the string.
// Zero-length runs are elided. Consecutive runs that have the same
// attribute type and value are merged. Each attribute is independent
// of each other attribute; overlapping attributes of different types
// do not split each other apart, but different values of the same
// attribute type do.
//
// The empty string can also be represented by uiAttributedString,
// but because of the no-zero-length-attribute rule, it will not have
// attributes.
//
// A uiAttributedString takes ownership of all attributes given to
// it, as it may need to duplicate or delete uiAttribute objects at
// any time. By extension, when you free a uiAttributedString,
// all uiAttributes within will also be freed. Each method will
// describe its own rules in more details.
//
// In addition, uiAttributedString provides facilities for moving
// between grapheme clusters, which represent a character
// from the point of view of the end user. The cursor of a text editor
// is always placed on a grapheme boundary, so you can use these
// features to move the cursor left or right by one "character".
// TODO does uiAttributedString itself need this
//
// uiAttributedString does not provide enough information to be able
// to draw itself onto a uiDrawContext or respond to user actions.
// In order to do that, you'll need to use a uiDrawTextLayout, which
// is built from the combination of a uiAttributedString and a set of
// layout-specific properties.
typedef struct uiAttributedString uiAttributedString;

// uiAttributedStringForEachAttributeFunc is the type of the function
// invoked by uiAttributedStringForEachAttribute() for every
// attribute in s. Refer to that function's documentation for more
// details.
typedef uiForEach (*uiAttributedStringForEachAttributeFunc)(const uiAttributedString *s, const uiAttribute *a, size_t start, size_t end, void *data);

// @role uiAttributedString constructor
// uiNewAttributedString() creates a new uiAttributedString from
// initialString. The string will be entirely unattributed.
_UI_EXTERN uiAttributedString *uiNewAttributedString(const char *initialString);

// @role uiAttributedString destructor
// uiFreeAttributedString() destroys the uiAttributedString s.
// It will also free all uiAttributes within.
_UI_EXTERN void uiFreeAttributedString(uiAttributedString *s);

// uiAttributedStringString() returns the textual content of s as a
// '\0'-terminated UTF-8 string. The returned pointer is valid until
// the next change to the textual content of s.
_UI_EXTERN const char *uiAttributedStringString(const uiAttributedString *s);

// uiAttributedStringLength() returns the number of UTF-8 bytes in
// the textual content of s, excluding the terminating '\0'.
_UI_EXTERN size_t uiAttributedStringLen(const uiAttributedString *s);

// uiAttributedStringAppendUnattributed() adds the '\0'-terminated
// UTF-8 string str to the end of s. The new substring will be
// unattributed.
_UI_EXTERN void uiAttributedStringAppendUnattributed(uiAttributedString *s, const char *str);

// uiAttributedStringInsertAtUnattributed() adds the '\0'-terminated
// UTF-8 string str to s at the byte position specified by at. The new
// substring will be unattributed; existing attributes will be moved
// along with their text.
_UI_EXTERN void uiAttributedStringInsertAtUnattributed(uiAttributedString *s, const char *str, size_t at);

// TODO add the Append and InsertAtExtendingAttributes functions
// TODO and add functions that take a string + length

// uiAttributedStringDelete() deletes the characters and attributes of
// s in the byte range [start, end).
_UI_EXTERN void uiAttributedStringDelete(uiAttributedString *s, size_t start, size_t end);

// TODO add a function to uiAttributedString to get an attribute's value at a specific index or in a specific range, so we can edit

// uiAttributedStringSetAttribute() sets a in the byte range [start, end)
// of s. Any existing attributes in that byte range of the same type are
// removed. s takes ownership of a; you should not use it after
// uiAttributedStringSetAttribute() returns.
_UI_EXTERN void uiAttributedStringSetAttribute(uiAttributedString *s, uiAttribute *a, size_t start, size_t end);

// uiAttributedStringForEachAttribute() enumerates all the
// uiAttributes in s. It is an error to modify s in f. Within f, s still
// owns the attribute; you can neither free it nor save it for later
// use.
// TODO reword the above for consistency (TODO and find out what I meant by that)
// TODO define an enumeration order (or mark it as undefined); also define how consecutive runs of identical attributes are handled here and sync with the definition of uiAttributedString itself
_UI_EXTERN void uiAttributedStringForEachAttribute(const uiAttributedString *s, uiAttributedStringForEachAttributeFunc f, void *data);

// TODO const correct this somehow (the implementation needs to mutate the structure)
_UI_EXTERN size_t uiAttributedStringNumGraphemes(uiAttributedString *s);

// TODO const correct this somehow (the implementation needs to mutate the structure)
_UI_EXTERN size_t uiAttributedStringByteIndexToGrapheme(uiAttributedString *s, size_t pos);

// TODO const correct this somehow (the implementation needs to mutate the structure)
_UI_EXTERN size_t uiAttributedStringGraphemeToByteIndex(uiAttributedString *s, size_t pos);

// uiFontDescriptor provides a complete description of a font where
// one is needed. Currently, this means as the default font of a
// uiDrawTextLayout and as the data returned by uiFontButton.
// All the members operate like the respective uiAttributes.
typedef struct uiFontDescriptor uiFontDescriptor;

struct uiFontDescriptor {
	// TODO const-correct this or figure out how to deal with this when getting a value
	char *Family;
	double Size;
	uiTextWeight Weight;
	uiTextItalic Italic;
	uiTextStretch Stretch;
};

// uiDrawTextLayout is a concrete representation of a
// uiAttributedString that can be displayed in a uiDrawContext.
// It includes information important for the drawing of a block of
// text, including the bounding box to wrap the text within, the
// alignment of lines of text within that box, areas to mark as
// being selected, and other things.
//
// Unlike uiAttributedString, the content of a uiDrawTextLayout is
// immutable once it has been created.
//
// TODO talk about OS-specific differences with text drawing that libui can't account for...
typedef struct uiDrawTextLayout uiDrawTextLayout;

// uiDrawTextAlign specifies the alignment of lines of text in a
// uiDrawTextLayout.
// TODO should this really have Draw in the name?
_UI_ENUM(uiDrawTextAlign) {
	uiDrawTextAlignLeft,
	uiDrawTextAlignCenter,
	uiDrawTextAlignRight,
};

// uiDrawTextLayoutParams describes a uiDrawTextLayout.
// DefaultFont is used to render any text that is not attributed
// sufficiently in String. Width determines the width of the bounding
// box of the text; the height is determined automatically.
typedef struct uiDrawTextLayoutParams uiDrawTextLayoutParams;

// TODO const-correct this somehow
struct uiDrawTextLayoutParams {
	uiAttributedString *String;
	uiFontDescriptor *DefaultFont;
	double Width;
	uiDrawTextAlign Align;
};

// @role uiDrawTextLayout constructor
// uiDrawNewTextLayout() creates a new uiDrawTextLayout from
// the given parameters.
//
// TODO
// - allow creating a layout out of a substring
// - allow marking compositon strings
// - allow marking selections, even after creation
// - add the following functions:
// 	- uiDrawTextLayoutHeightForWidth() (returns the height that a layout would need to be to display the entire string at a given width)
// 	- uiDrawTextLayoutRangeForSize() (returns what substring would fit in a given size)
// 	- uiDrawTextLayoutNewWithHeight() (limits amount of string used by the height)
// - some function to fix up a range (for text editing)
_UI_EXTERN uiDrawTextLayout *uiDrawNewTextLayout(uiDrawTextLayoutParams *params);

// @role uiDrawFreeTextLayout destructor
// uiDrawFreeTextLayout() frees tl. The underlying
// uiAttributedString is not freed.
_UI_EXTERN void uiDrawFreeTextLayout(uiDrawTextLayout *tl);

// uiDrawText() draws tl in c with the top-left point of tl at (x, y).
_UI_EXTERN void uiDrawText(uiDrawContext *c, uiDrawTextLayout *tl, double x, double y);

// uiDrawTextLayoutExtents() returns the width and height of tl
// in width and height. The returned width may be smaller than
// the width passed into uiDrawNewTextLayout() depending on
// how the text in tl is wrapped. Therefore, you can use this
// function to get the actual size of the text layout.
_UI_EXTERN void uiDrawTextLayoutExtents(uiDrawTextLayout *tl, double *width, double *height);
