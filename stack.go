// 13 february 2014
package ui

import (
	"fmt"
	"sync"
)

// Orientation defines the orientation of controls in a Stack.
type Orientation bool
const (
	Horizontal Orientation = false
	Vertical Orientation = true
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
	width, height	[]int		// caches to avoid reallocating these each time
}

// NewStack creates a new Stack with the specified orientation.
func NewStack(o Orientation, controls ...Control) *Stack {
	return &Stack{
		orientation:	o,
		controls:		controls,
		stretchy:		make([]bool, len(controls)),
		width:		make([]int, len(controls)),
		height:		make([]int, len(controls)),
	}
}

// SetStretchy marks a control in a Stack as stretchy.
func (s *Stack) SetStretchy(index int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.stretchy[index] = true			// TODO explicitly check for index out of bounds?
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
	var stretchywid, stretchyht int

	if len(s.controls) == 0 {		// do nothing if there's nothing to do
		return nil
	}
	// 1) get height and width of non-stretchy controls; figure out how much space is alloted to stretchy controls
	stretchywid = width
	stretchyht = height
	nStretchy := 0
	for i, c := range s.controls {
		if s.stretchy[i] {
			nStretchy++
			continue
		}
		w, h, err := c.preferredSize()
		if err != nil {
			return fmt.Errorf("error getting preferred size of control %d in Stack.setRect(): %v", i, err)
		}
		if s.orientation == Horizontal {			// all controls have same height
			s.width[i] = w
			s.height[i] = height
			stretchywid -= w
		} else {							// all controls have same width
			s.width[i] = width
			s.height[i] = h
			stretchyht -= h
		}
	}
	// 2) figure out size of stretchy controls
	if nStretchy != 0 {
		if s.orientation == Horizontal {			// split rest of width
			stretchywid /= nStretchy
		} else {							// split rest of height
			stretchyht /= nStretchy
		}
	}
	for i, c := range s.controls {
		if !s.stretchy[i] {
			continue
		}
		c.width[i] = stretchywid
		c.height[i] = stretchyht
	}
	// 3) now actually place controls
	for i, c := range s.controls {
		err := c.setRect(x, y, s.width[i], s.height[i])
		if err != nil {
			return fmt.Errorf("error setting size of control %d in Stack.setRect(): %v", i, err)
		}
		if s.orientation == Horizontal {
			x += s.width[i]
		} else {
			y += s.height[i]
		}
	}
	return nil
}
