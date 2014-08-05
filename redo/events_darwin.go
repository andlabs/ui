// 30 march 2014

package ui

/*
Mac OS X uses its own set of hardware key codes that are different from PC keyboard scancodes, but are positional (like PC keyboard scancodes). These are defined in <HIToolbox/Events.h>, a Carbon header. As far as I can tell, there's no way to include this header without either using an absolute path or linking Carbon into the program, so the constant values are used here instead.

The Cocoa docs do guarantee that -[NSEvent keyCode] results in key codes that are the same as those returned by Carbon; that is, these codes.
*/

// use uintptr to be safe
var keycodeKeys = map[uintptr]byte{
	0x00: 'a',
	0x01: 's',
	0x02: 'd',
	0x03: 'f',
	0x04: 'h',
	0x05: 'g',
	0x06: 'z',
	0x07: 'x',
	0x08: 'c',
	0x09: 'v',
	0x0B: 'b',
	0x0C: 'q',
	0x0D: 'w',
	0x0E: 'e',
	0x0F: 'r',
	0x10: 'y',
	0x11: 't',
	0x12: '1',
	0x13: '2',
	0x14: '3',
	0x15: '4',
	0x16: '6',
	0x17: '5',
	0x18: '=',
	0x19: '9',
	0x1A: '7',
	0x1B: '-',
	0x1C: '8',
	0x1D: '0',
	0x1E: ']',
	0x1F: 'o',
	0x20: 'u',
	0x21: '[',
	0x22: 'i',
	0x23: 'p',
	0x25: 'l',
	0x26: 'j',
	0x27: '\'',
	0x28: 'k',
	0x29: ';',
	0x2A: '\\',
	0x2B: ',',
	0x2C: '/',
	0x2D: 'n',
	0x2E: 'm',
	0x2F: '.',
	0x32: '`',
	0x24: '\n',
	0x30: '\t',
	0x31: ' ',
	0x33: '\b',
}

var keycodeExtKeys = map[uintptr]ExtKey{
	0x41: NDot,
	0x43: NMultiply,
	0x45: NAdd,
	0x4B: NDivide,
	0x4C: NEnter,
	0x4E: NSubtract,
	0x52: N0,
	0x53: N1,
	0x54: N2,
	0x55: N3,
	0x56: N4,
	0x57: N5,
	0x58: N6,
	0x59: N7,
	0x5B: N8,
	0x5C: N9,
	0x35: Escape,
	0x60: F5,
	0x61: F6,
	0x62: F7,
	0x63: F3,
	0x64: F8,
	0x65: F9,
	0x67: F11,
	0x6D: F10,
	0x6F: F12,
	0x72: Insert, // listed as the Help key but it's in the same position on an Apple keyboard as the Insert key on a Windows keyboard; thanks to SeanieB from irc.badnik.net and Psy in irc.freenode.net/#macdev for confirming they have the same code
	0x73: Home,
	0x74: PageUp,
	0x75: Delete,
	0x76: F4,
	0x77: End,
	0x78: F2,
	0x79: PageDown,
	0x7A: F1,
	0x7B: Left,
	0x7C: Right,
	0x7D: Down,
	0x7E: Up,
}

var keycodeModifiers = map[uintptr]Modifiers{
	0x37: Super, // left command
	0x38: Shift, // left shift
	0x3A: Alt,   // left option
	0x3B: Ctrl,  // left control
	0x3C: Shift, // right shift
	0x3D: Alt,   // right alt
	0x3E: Ctrl,  // right control

	// the following is not in Events.h for some reason
	// thanks to Nicole and jedivulcan from irc.badnik.net
	0x36: Super, // right command
}

func fromKeycode(keycode uintptr) (ke KeyEvent, ok bool) {
	if key, ok := keycodeKeys[keycode]; ok {
		ke.Key = key
		return ke, true
	}
	if extkey, ok := keycodeExtKeys[keycode]; ok {
		ke.ExtKey = extkey
		return ke, true
	}
	return ke, false
}
