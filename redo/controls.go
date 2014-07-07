// 7 july 2014

package ui

// Control represents a control.
// All Controls have event handlers that take a single argument (the Doer active during the event) and return nothing.
type Control interface {
	// TODO reparent (public)
	// TODO enable/disable (public)
	// TODO show/hide (public)
	// TODO sizing (likely private)
}

// Button is a clickable button that performs some task.
type Button interface {
	Control

	// OnClicked creates a Request to set the event handler for when the Button is clicked.
	OnClicked(func(d Doer)) *Request

	// Text and SetText creates a Request that get and set the Button's label text.
	Text() *Request
	SetText(text string) *Request
}

// NewButton creates a Request to create a new Button with the given label text.
func NewButton(text string) *Request {
	return newButton(text)
}

// GetNewButton is like NewButton but sends the Request along the given Doer and returns the resultant Button.
// Example:
// 	b := ui.GetNewButton(ui.Do, "OK")
func GetNewButton(c Doer, text string) Button {
	c <- newButton(text)
	return (<-c.resp).(Button)
}
