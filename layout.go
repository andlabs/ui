package ui

// Recursively replaces nils with stretchy empty spaces and changes the orientation
// of inner stack so they are perpenticular to each other.
func resetControls(parent *Stack) {
	for i, control := range parent.controls {
		switch control.(type) {
		case *Stack:
			stack := control.(*Stack)
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
	}

	resetControls(stack)

	return stack
}
