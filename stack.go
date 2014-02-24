// 13 february 2014
package ui

import (
	"fmt"
	"sync"
)

// Orientation defines the orientation of controls in a Stack.
type Orientation int
const (
	Horizontal Orientation = iota
	Vertical
)

// A Stack stacks controls horizontally or vertically within the Stack's parent.
// A horizontal Stack gives all controls the same height and their preferred widths.
// A vertical Stack gives all controls the same width and their preferred heights.
// Some controls may be marked as "stretchy": when the Window they are in changes size, stretchy controls resize to take up the remaining space after non-stretchy controls are laid out. If multiple controls are marked stretchy, they are alloted equal distribution of the remaining space.
type Stack struct {
	lock			sync.Mutex
	created		bool
	orientation	Orientation
	controls		[]Control
	stretchy		[]bool
	xpos, ypos	[]int		// caches to avoid reallocating these each time
	width, height	[]int
}

// NewStack creates a new Stack with the specified orientation.
func NewStack(o Orientation, controls ...Control) *Stack {
	if o != Horizontal && o != Vertical {
		panic(fmt.Sprintf("invalid orientation %d given to NewStack()", o))
	}
	return &Stack{
		orientation:	o,
		controls:		controls,
		stretchy:		make([]bool, len(controls)),
		xpos:		make([]int, len(controls)),
		ypos:		make([]int, len(controls)),
		width:		make([]int, len(controls)),
		height:		make([]int, len(controls)),
	}
}

// SetStretchy marks a control in a Stack as stretchy.
func (s *Stack) SetStretchy(index int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s[index] = true			// TODO explicitly check for index out of bounds?
}

func (s *Stack) make(window *sysData) error {
	for i, c := range s.controls {
		err := c.make(window)
		if err != nil {
			return fmt.Errorf("error adding control %d to Stack: %v", i, err)
		}
	}
	s.created = true
	return nil
}

func (s *Stack) setRect(x int, y int, width int, height int) error {
	var dx, dy int

	if len(s.controls) == 0 {		// do nothing if there's nothing to do
		return nil
	}
	switch s.orientation {
	case Horizontal:
		dx = width / len(s.controls)
		width = dx
	case Vertical:
		dy = height / len(s.controls)
		height = dy
	default:
		panic(fmt.Sprintf("invalid orientation %d given to Stack.setRect()", s.orientation))
	}
	for i, c := range s.controls {
		err := c.setRect(x, y, width, height)
		if err != nil {
			return fmt.Errorf("error setting size of control %d: %v", i, err)
		}
		x += dx
		y += dy
	}
	return nil
}
