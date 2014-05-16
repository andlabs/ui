A Stack is a stack of controls arranged horizontally or vertically.
	If horizontal, all controls take the same heights, and use their preferred widths.
	If vertical, all controls take the same widths, and use their preferred heights.

PROBLEM
	We want to have controls that stretch with the window
SOLUTION
	Delegate one control as the "stretchy" control.
		s := NewVerticalStack(c1, c2, c3)
		s.SetStretchy(1)	// c2
	When drawing, all other controls have their positions and sizes determiend, then the stretchy one's is.

PROBLEM
	We want to have a control that's anchored (resizes with) the top left and one that's anchored with the top right.
SOLUTION
	Allow stretchiness for arbitrary controls
		s.SetStretchy(1)
		s.SetStretchy(2)
	When drawing, the sizes of non-stretchy controls are determined first, then all stretchy controls are given equal amounts of the rest.
	The preferred size of the stack is the preferred size of non-stretchy controls + (the number of stretchy controls * the largest preferred size of the stretchy controls).

PROBLEM
	Non-equal size distribution of stretchy controls: for instance, in a horizontal stack, a navigation bar is usually a fixed size and the content area fills the rest.
SOLUTION
	I'm not entirely sure how to fix this one yet; I do know for a navigation pane the user is usually in control of the size, so... will figure out later.
