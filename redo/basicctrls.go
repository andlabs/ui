// 7 july 2014

package ui

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

	// OnToggled sets the event handler for when the Checkbox is toggled.
	OnToggled(func())

	// Text and SetText get and set the Checkbox's label text.
	Text() string
	SetText(text string)

	// Checked and SetChecked get and set the Checkbox's check state.
	Checked() bool
	SetChecked(checked bool)
}

// NewCheckbox creates a new Checkbox with the given label text.
// The Checkbox will be initially unchecked.
func NewCheckbox(text string) Checkbox {
	return newCheckbox(text)
}

// TextField is a Control in which the user can enter a single line of text.
type TextField interface {
	Control

	// Text and SetText are Requests that get and set the TextField's text.
	Text() string
	SetText(text string)
}

// NewTextField creates a new TextField.
func NewTextField() TextField {
	return newTextField()
}

// NewPasswordField creates a new TextField for entering passwords; that is, it hides the text being entered.
func NewPasswordField() TextField {
	return newPasswordField()
}

// Tab is a Control that contains multiple pages of tabs, each containing a single Control.
// You can add and remove tabs from the Tab at any time.
// The appearance of a Tab with no tabs is implementation-defined.
type Tab interface {
	Control

	// Append adds a new tab to Tab.
	// The tab is added to the end of the current list of tabs.
	Append(name string, control Control)
}

// NewTab creates a new Tab with no tabs.
func NewTab() Tab {
	return newTab()
}

// Label is a Control that shows a static line of text.
// Label shows one line of text; any text that does not fit is truncated.
// A Label can either have smart vertical alignment relative to the control to its right or just be vertically aligned to the top (standalone).
// The effect of placing a non-standalone Label in any context other than to the immediate left of a Control is undefined.
// Both types of labels currently are left-aligned (TODO).
type Label interface {
	Control

	// Text and SetText get and set the Label's text.
	Text() string
	SetText(text string)
}

// NewLabel creates a new Label with the given text.
// The Label will smartly vertically position itself relative to the control to its immediate right.
// TODO Grids on GTK+ will not respect this unless SetFilling()
func NewLabel(text string) Label {
	return newLabel(text)
}

// NewStandaloneLabel creates a new Label with the given text.
// The Label will be vertically positioned at the top of its allocated space.
func NewStandaloneLabel(text string) Label {
	return newStandaloneLabel(text)
}
