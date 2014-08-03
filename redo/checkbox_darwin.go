// 16 july 2014

package ui

// #include "objc_darwin.h"
import "C"

type checkbox struct {
	*button
}

func newCheckbox(text string) *checkbox {
	return &checkbox{
		button:	finishNewButton(C.newCheckbox(), text),
	}
}

// we don't need to define our own event here; we can just reuse Button's
// (it's all target-action anyway)

func (c *checkbox) Checked() bool {
	return fromBOOL(C.checkboxChecked(c.id))
}

func (c *checkbox) SetChecked(checked bool) {
	C.checkboxSetChecked(c.id, toBOOL(checked))
}
