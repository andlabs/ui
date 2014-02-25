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
// Any extra space at the end of a Stack is left blank.
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

// SetStretchy marks a control in a Stack as stretchy. This cannot be called once the Window containing the Stack has been opened.
func (s *Stack) SetStretchy(index int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.created {
		panic("call to Stack.SetStretchy() after Stack has been created")
	}
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
	for i := range s.controls {
		if !s.stretchy[i] {
			continue
		}
		s.width[i] = stretchywid
		s.height[i] = stretchyht
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

// The preferred size of a Stack is the sum of the preferred sizes of non-stretchy controls + (the number of stretchy controls * the largest preferred size among all stretchy controls).
func (s *Stack) preferredSize() (width int, height int, err error) {
	max := func(a int, b int) int {
		if a > b {
			return a
		}
		return b
	}

	var nStretchy int
	var maxswid, maxsht int

	if len(s.controls) == 0 {		// no controls, so return emptiness
		return 0, 0, nil
	}
	for i, c := range s.controls {
		w, h, err := c.preferredSize()
		if err != nil {
			return 0, 0, fmt.Errorf("error getting preferred size of control %d in Stack.preferredSize(): %v", i, err)
		}
		if s.stretchy[i] {
			nStretchy++
			maxswid = max(maxswid, w)
			maxsht = max(maxsht, h)
		}
		if s.orientation == Horizontal {			// max vertical size
			if !s.stretchy[i] {
				width += w
			}
			height = max(height, h)
		} else {
			width = max(width, w)
			if !s.stretchy[i] {
				height += h
			}
		}
	}
	if s.orientation == Horizontal {
		width += nStretchy * maxswid
	} else {
		height += nStretchy * maxsht
	}
	return
}

// Space returns a null control intended for padding layouts with blank space where otherwise impossible (for instance, at the beginning or in the middle of a Stack).
// In order for a Space to work, it must be marked as stretchy in its parent layout; otherwise its size is undefined.
func Space() Control {
	// As above, a Stack with no controls draws nothing and reports no errors; its parent will still size it properly if made stretchy.
	return NewStack(Horizontal)
}
