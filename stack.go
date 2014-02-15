// 13 february 2014
package main

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

// A Stack stacks controls horizontally or vertically within the Stack's parent, alotting each the same size.
type Stack struct {
	lock			sync.Mutex
	created		bool
	orientation	Orientation
	controls		[]Control
}

// NewStack creates a new Stack with the specified orientation.
func NewStack(o Orientation, controls ...Control) *Stack {
	if o != Horizontal && o != Vertical {
		panic(fmt.Sprintf("invalid orientation %d given to NewStack()", o))
	}
	return &Stack{
		orientation:	o,
		controls:		controls,
	}
}

// TODO adorn errors with which stage failed
func (s *Stack) make(window *sysData) error {
	for _, c := range s.controls {
		err := c.make(window)
		if err != nil {
			return err
		}
	}
	s.created = true
	return nil
}

// TODO adorn errors with which stage failed
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
	for _, c := range s.controls {
		err := c.setRect(x, y, width, height)
		if err != nil {
			return err
		}
		x += dx
		y += dy
	}
	return nil
}
