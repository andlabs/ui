// 28 october 2014

package ui

// #include "objc_darwin.h"
import "C"

// interface builder notes
// - the tops of the alignment rects should be identical
// - spinner properties: auto repeat
// - http://stackoverflow.com/questions/702829/integrate-nsstepper-with-nstextfield we'll need to bind the int value :S
// 	- TODO experiment with a dummy project
// - http://juliuspaintings.co.uk/cgi-bin/paint_css/animatedPaint/059-NSStepper-NSTextField.pl
// - http://www.youtube.com/watch?v=ZZSHU-O7HVo
// - http://andrehoffmann.wordpress.com/tag/nsstepper/ ?
// TODO
// - proper spacing between edit and spinner: Interface Builder isn't clear; NSDatePicker doesn't spill the beans
