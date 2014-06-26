package ui

// Recursively removes border margins and padding from controls, replaces
// nil values with stretchy spaces and reorients nested stack to have
// opposing orientations.
func resetControls(parent *Stack) {
	for i, control := range parent.controls {
		switch control.(type) {
		case *Stack:
			stack := control.(*Stack)
			stack.borderMargin = 0
			stack.orientation = !parent.orientation
			resetControls(stack)
		case nil:
			emptySpace := newStack(horizontal)
			parent.controls[i] = emptySpace
			parent.stretchy[i] = true
		}
	}
}

// Creates a new Stack from the given controls. The topmost Stack will have
// vertical orientation and margin borders, with each nested stack being
// oriented oppositely. Controls are displayed with a default padding
// between them.
func Layout(controls ...Control) *Stack {
	stack := &Stack{
		orientation:  vertical,
		controls:     controls,
		stretchy:     make([]bool, len(controls)),
		width:        make([]int, len(controls)),
		height:       make([]int, len(controls)),
		padding:      10,
		borderMargin: 15,
	}

	resetControls(stack)

	return stack
}
