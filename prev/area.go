// 14 march 2014

package ui

import (
	"fmt"
	"image"
	"image/draw"
	"reflect"
	"unsafe"
)

// Area represents a blank canvas upon which programs may draw anything and receive arbitrary events from the user.
// An Area has an explicit size, represented in pixels, that may be different from the size shown in its Window.
// For information on scrollbars, see "Scrollbars" in the Overview.
// The coordinate system of an Area always has an origin of (0,0) which maps to the top-left corner; all image.Points and image.Rectangles sent across Area's channels conform to this.
// The size of an Area must be at least 1x1 (that is, neither its width nor its height may be zero or negative).
// For control layout purposes, an Area prefers to be at the size you set it to (so if an Area is not stretchy in its layout, it will ask to have that size).
//
// To handle events to the Area, an Area must be paired with an AreaHandler.
// See AreaHandler for details.
//
// Area will accept keyboard focus if tabbed into, but will refuse to relinquish keyboard focus if tabbed out.
//
// Do not use an Area if you intend to read text.
// Area reads keys based on their position on a standard
// 101-key keyboard, and does no character processing.
// Character processing methods differ across operating
// systems; trying ot recreate these yourself is only going
// to lead to trouble.
// If you absolutely need to enter text somehow, use OpenTextFieldAt() and its related methods.
type Area interface {
	Control

	// SetSize sets the Area's internal drawing size.
	// It has no effect on the actual control size.
	// SetSize will also signal the entirety of the Area to be redrawn as in RepaintAll.
	// It panics if width or height is zero or negative.
	SetSize(width int, height int)

	// Repaint marks the given rectangle of the Area as needing to be redrawn.
	// The given rectangle is clipped to the Area's size.
	// If, after clipping, the rectangle is empty, Repaint does nothing.
	Repaint(r image.Rectangle)

	// RepaintAll marks the entirety of the Area as needing to be redrawn.
	RepaintAll()

	// OpenTextFieldAt opens a TextField with the top-left corner at the given coordinates of the Area.
	// It panics if the coordinates fall outside the Area.
	// Any text previously in the TextField (be it by the user or by a call to SetTextFieldText()) is retained.
	// The TextField receives the input focus so the user can type things; when the TextField loses the input focus, it hides itself and signals the event set by OnTextFieldDismissed.
	// The TextField will also dismiss itself on some platforms when the user "completes editing"; the exact meaning of this is platform-specific.
	OpenTextFieldAt(x int, y int)

	// TextFieldText and TextFieldSetText get and set the OpenTextFieldAt TextField's text, respectively.
	TextFieldText() string
	SetTextFieldText(text string)

	// OnTextFieldDismissed is an event that is fired when the OpenTextFieldAt TextField is dismissed.
	OnTextFieldDismissed(f func())
}

type areabase struct {
	width   int
	height  int
	handler AreaHandler
}

// AreaHandler represents the events that an Area should respond to.
// These methods are all executed on the main goroutine, not necessarily the same one that you created the AreaHandler in; you are responsible for the thread safety of any members of the actual type that implements ths interface.
// (Having to use this interface does not strike me as being particularly Go-like, but the nature of Paint makes channel-based event handling a non-option; in practice, deadlocks occur.)
type AreaHandler interface {
	// Paint is called when the Area needs to be redrawn.
	// The part of the Area that needs to be redrawn is stored in cliprect.
	// Before Paint() is called, this region is cleared with a system-defined background color.
	// You MUST handle this event, and you MUST return a valid image, otherwise deadlocks and panicking will occur.
	// The image returned must have the same size as rect (but does not have to have the same origin points).
	// Example:
	// 	imgFromFile, _, err := image.Decode(file)
	// 	if err != nil { panic(err) }
	// 	img := image.NewRGBA(imgFromFile.Rect)
	// 	draw.Draw(img, img.Rect, imgFromFile, image.ZP, draw.Over)
	// 	// ...
	// 	func (h *myAreaHandler) Paint(rect image.Rectangle) *image.RGBA {
	// 		return img.SubImage(rect).(*image.RGBA)
	// 	}
	Paint(cliprect image.Rectangle) *image.RGBA

	// Mouse is called when the Area receives a mouse event.
	// You are allowed to do nothing in this handler (to ignore mouse events).
	// See MouseEvent for details.
	// After handling the mouse event, package ui will decide whether to perform platform-dependent event chain continuation based on that platform's designated action (so it is not possible to override global mouse events this way).
	Mouse(e MouseEvent)

	// Key is called when the Area receives a keyboard event.
	// Return true to indicate that you handled the event; return false to indicate that you did not and let the system handle the event.
	// You are allowed to do nothing in this handler (to ignore keyboard events); in this case, return false.
	// See KeyEvent for details.
	Key(e KeyEvent) (handled bool)
}

// MouseEvent contains all the information for a mous event sent by Area.Mouse.
// Mouse button IDs start at 1, with 1 being the left mouse button, 2 being the middle mouse button, and 3 being the right mouse button.
// If additional buttons are supported, they will be returned with 4 being the first additional button.
// For example, on Unix systems where mouse buttons 4 through 7 are pseudobuttons for the scroll wheel directions, the next button, button 8, will be returned as 4, 9 as 5, etc.
// The association between button numbers and physical buttons are system-defined.
// For example, on Windows, buttons 4 and 5 are mapped to what are internally referred to as "XBUTTON1" and "XBUTTON2", which often correspond to the dedicated back/forward navigation buttons on the sides of many mice.
// The examples here are NOT a guarantee as to how many buttons maximum will be available on a given system.
//
// If the user clicked on the Area to switch to the Window it is contained in from another window in the OS, the Area will receive a MouseEvent for that click.
type MouseEvent struct {
	// Pos is the position of the mouse in the Area at the time of the event.
	Pos image.Point

	// If the event was generated by a mouse button being pressed, Down contains the ID of that button.
	// Otherwise, Down contains 0.
	// If Down contains nonzero, the Area will also receive keyboard focus.
	Down uint

	// If the event was generated by a mouse button being released, Up contains the ID of that button.
	// Otherwise, Up contains 0.
	// If both Down and Up are 0, the event represents mouse movement (with optional held buttons for dragging; see below).
	// Down and Up shall not both be nonzero.
	Up uint

	// If Down is nonzero, Count indicates the number of clicks: 1 for single-click, 2 for double-click, 3 for triple-click, and so on.
	// The order of events will be Down:Count=1 -> Up -> Down:Count=2 -> Up -> Down:Count=3 -> Up -> ...
	Count uint

	// Modifiers is a bit mask indicating the modifier keys being held during the event.
	Modifiers Modifiers

	// Held is a slice of button IDs that indicate which mouse buttons are being held during the event.
	// Held will not include Down and Up.
	// Held will be sorted.
	// Only buttons 1, 2, and 3 are guaranteed to be detected by Held properly; whether or not any others are is implementation-defined.
	//
	// If Held is non-empty but Up and Down are both zero, the mouse is being dragged, with all the buttons in Held being held.
	// Whether or not a drag into an Area generates MouseEvents is implementation-defined.
	// Whether or not a drag over an Area when the program is inactive generates MouseEvents is also implementation-defined.
	// Moving the mouse over an Area when the program is inactive and no buttons are held will, however, generate MouseEvents.
	Held []uint
}

// HeldBits returns Held as a bit mask.
// Bit 0 maps to button 1, bit 1 maps to button 2, etc.
func (e MouseEvent) HeldBits() (h uintptr) {
	for _, x := range e.Held {
		h |= uintptr(1) << (x - 1)
	}
	return h
}

// A KeyEvent represents a keypress in an Area.
//
// Key presses are based on their positions on a standard
// 101-key keyboard found on most computers. The
// names chosen for keys here are based on their names
// on US English QWERTY keyboards; see Key for details.
//
// If a key is pressed that is not supported by Key, ExtKey,
// or Modifiers, no KeyEvent will be produced and package ui will behave as if false was returned from the event handler.
type KeyEvent struct {
	// Key is a byte representing a character pressed
	// in the typewriter section of the keyboard.
	// The value, which is independent of whether the
	// Shift key is held, is a constant with one of the
	// following (case-sensitive) values, drawn according
	// to the key's position on the keyboard.
	//    ` 1 2 3 4 5 6 7 8 9 0 - =
	//     q w e r t y u i o p [ ] \
	//      a s d f g h j k l ; '
	//       z x c v b n m , . /
	// The actual key entered will be the key at the respective
	// position on the user's keyboard, regardless of the actual
	// layout. (Some keyboards move \ to either the row above
	// or the row below but in roughly the same spot; this is
	// accounted for. Some keyboards have an additonal key
	// to the left of 'z' or additional keys to the right of '='; these
	// cannot be read.)
	// In addition, Key will contain
	// - ' ' (space) if the spacebar was pressed
	// - '\t' if Tab was pressed, regardless of Modifiers
	// - '\n' if the typewriter Enter key was pressed
	// - '\b' if the typewriter Backspace key was pressed
	// If this value is zero, see ExtKey.
	Key byte

	// If Key is zero, ExtKey contains a predeclared identifier
	// naming an extended key. See ExtKey for details.
	// If both Key and ExtKey are zero, a Modifier by itself
	// was pressed. Key and ExtKey will not both be nonzero.
	ExtKey ExtKey

	// If both Key and ExtKey are zero, Modifier will contain exactly one of its bits set, indicating which Modifier was pressed or released.
	// As with Modifiers itself, there is no way to differentiate between left and right modifier keys.
	// As such, the result of pressing and/or releasing both left and right of the same Modifier is system-defined.
	// Furthermore, the result of holding down a Key or ExtKey, then pressing a Modifier, and then releasing the original key is system-defined.
	// Under no condition shall Key, ExtKey, AND Modifier all be zero.
	Modifier Modifiers

	// Modifiers contains all the modifier keys currently being held at the time of the KeyEvent.
	// If Modifier is nonzero, Modifiers will not contain Modifier itself.
	Modifiers Modifiers

	// If Up is true, the key was released; if not, the key was pressed.
	// There is no guarantee that all pressed keys shall have
	// corresponding release events (for instance, if the user switches
	// programs while holding the key down, then releases the key).
	// Keys that have been held down are reported as multiple
	// key press events.
	Up bool
}

// ExtKey represents keys that are not in the typewriter section of the keyboard.
type ExtKey uintptr

const (
	Escape ExtKey = iota + 1
	Insert        // equivalent to "Help" on Apple keyboards
	Delete
	Home
	End
	PageUp
	PageDown
	Up
	Down
	Left
	Right
	F1 // F1..F12 are guaranteed to be consecutive
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	N0 // numpad keys; independent of Num Lock state
	N1 // N0..N9 are guaranteed to be consecutive
	N2
	N3
	N4
	N5
	N6
	N7
	N8
	N9
	NDot
	NEnter
	NAdd
	NSubtract
	NMultiply
	NDivide
	_nextkeys // for sanity check
)

// EffectiveKey returns e.Key if it is set.
// Otherwise, if e.ExtKey denotes a numpad key,
// EffectiveKey returns the equivalent e.Key value
// ('0'..'9', '.', '\n', '+', '-', '*', or '/').
// Otherwise, EffectiveKey returns zero.
func (e KeyEvent) EffectiveKey() byte {
	if e.Key != 0 {
		return e.Key
	}
	k := e.ExtKey
	switch {
	case k >= N0 && k <= N9:
		return byte(k-N0) + '0'
	case k == NDot:
		return '.'
	case k == NEnter:
		return '\n'
	case k == NAdd:
		return '+'
	case k == NSubtract:
		return '-'
	case k == NMultiply:
		return '*'
	case k == NDivide:
		return '/'
	}
	return 0
}

// Modifiers indicates modifier keys being held during an event.
// There is no way to differentiate between left and right modifier keys.
// As such, what KeyEvents get sent if the user does something unusual with both of a certain modifier key at once is undefined.
type Modifiers uintptr

const (
	Ctrl  Modifiers = 1 << iota // the keys labelled Ctrl or Control on all platforms
	Alt                         // the keys labelled Alt or Option or Meta on all platforms
	Shift                       // the Shift keys
	Super                       // the Super keys on platforms that have one, or the Windows keys on Windows, or the Command keys on Mac OS X
)

func checkAreaSize(width int, height int, which string) {
	if width <= 0 || height <= 0 {
		panic(fmt.Errorf("invalid size %dx%d in %s", width, height, which))
	}
}

// NewArea creates a new Area with the given size and handler.
// It panics if handler is nil or if width or height is zero or negative.
func NewArea(width int, height int, handler AreaHandler) Area {
	checkAreaSize(width, height, "NewArea()")
	if handler == nil {
		panic("handler passed to NewArea() must not be nil")
	}
	return newArea(&areabase{
		width:   width,
		height:  height,
		handler: handler,
	})
}

// internal function, but shared by all system implementations: &img.Pix[0] is not necessarily the first pixel in the image
func pixelDataPos(img *image.RGBA) int {
	return img.PixOffset(img.Rect.Min.X, img.Rect.Min.Y)
}

func pixelData(img *image.RGBA) *uint8 {
	return &img.Pix[pixelDataPos(img)]
}

// some platforms require pixels in ARGB order in their native endianness (because they treat the pixel array as an array of uint32s)
// this does the conversion
// you need to convert somewhere (Windows and cairo give us memory to use; Windows has stride==width but cairo might not)
func toARGB(i *image.RGBA, memory uintptr, memstride int, toNRGBA bool) {
	var realbits []byte

	rbs := (*reflect.SliceHeader)(unsafe.Pointer(&realbits))
	rbs.Data = memory
	rbs.Len = 4 * i.Rect.Dx() * i.Rect.Dy()
	rbs.Cap = rbs.Len
	p := pixelDataPos(i)
	q := 0
	iPix := i.Pix
	if toNRGBA { // for Windows image lists
		j := image.NewNRGBA(i.Rect)
		draw.Draw(j, j.Rect, i, i.Rect.Min, draw.Src)
		iPix = j.Pix
	}
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		nextp := p + i.Stride
		nextq := q + memstride
		for x := i.Rect.Min.X; x < i.Rect.Max.X; x++ {
			argb := uint32(iPix[p+3]) << 24 // A
			argb |= uint32(iPix[p+0]) << 16 // R
			argb |= uint32(iPix[p+1]) << 8  // G
			argb |= uint32(iPix[p+2])       // B
			// the magic of conversion
			native := (*[4]byte)(unsafe.Pointer(&argb))
			realbits[q+0] = native[0]
			realbits[q+1] = native[1]
			realbits[q+2] = native[2]
			realbits[q+3] = native[3]
			p += 4
			q += 4
		}
		p = nextp
		q = nextq
	}
}
