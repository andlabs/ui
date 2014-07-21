// 7 july 2014

package ui

// Control represents a control.
// All Controls have event handlers that take a single argument (the Doer active during the event) and return nothing.
type Control interface {
	unparent()
	parent(*window)
	// TODO enable/disable (public)
	// TODO show/hide (public)
	controlSizing
}

// Button is a clickable button that performs some task.
type Button interface {
	Control

	// OnClicked sets the event handler for when the Button is clicked.
	OnClicked(func())

	// Text and SetText get and set the Button's label text.
	Text() string
	SetText(text string)
}

// NewButton creates a new Button with the given label text.
func NewButton(text string) Button {
	return newButton(text)
}

// Checkbox is a clickable box that indicates some Boolean value.
type Checkbox interface {
	Control

	// OnClicked sets the event handler for when the Checkbox is clicked (to change its toggle state).
	// TODO change to OnCheckChanged or OnToggled?
	OnClicked(func())

	// Text and SetText are Requests that get and set the Checkbox's label text.
	Text() string
	SetText(text string)

	// Checked and SetChecked get and set the Checkbox's check state.
	Checked() bool
	SetChecked(checked bool)
}

// NewCheckbox creates a new Checkbox with the given label text.
// The Checkbox will be initially unchecked.
func NewCheckbox(text string) Checkbox {
//	return newCheckbox(text)
return nil
//TODO add checkbox after resolving other TODOs
}
