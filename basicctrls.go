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

	// OnChanged is triggered when the text in a TextField is changed somehow.
	// Do not bother trying to figure out how the text was changed; instead, perform your validation and use Invalid to inform the user that the entered text is invalid instead.
	OnChanged(func())

	// Invalid throws a non-modal alert (whose nature is system-defined) on or near the TextField that alerts the user that input is invalid.
	// The string passed to Invalid will be displayed to the user to inform them of what specifically is wrong with the input.
	// Pass an empty string to remove the warning.
	Invalid(reason string)

	// ReadOnly and SetReadOnly get and set whether the TextField is read-only.
	// A read-only TextField cannot be changed by the user, but its text can still be manipulated in other ways (selecting, copying, etc.).
	ReadOnly() bool
	SetReadOnly(readonly bool)
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
// Labels are left-aligned. [FUTURE PLANS: For platform-specific horizontal alignment rules, use a Form.]
type Label interface {
	Control

	// Text and SetText get and set the Label's text.
	Text() string
	SetText(text string)
}

// NewLabel creates a new Label with the given text.
func NewLabel(text string) Label {
	return newLabel(text)
}

// Group is a Control that holds a single Control; if that Control also contains other Controls, then the Controls will appear visually grouped together.
// The appearance of a Group varies from system to system; for the most part a Group consists of a thin frame.
// All Groups have a text label indicating what the Group is for.
type Group interface {
	Control

	// Text and SetText get and set the Group's label text.
	Text() string
	SetText(text string)

	// Margined and SetMargined get and set whether the contents of the Group have a margin around them.
	// The size of the margin is platform-dependent.
	Margined() bool
	SetMargined(margined bool)
}

// NewGroup creates a new Group with the given text label and child Control.
func NewGroup(text string, control Control) Group {
	return newGroup(text, control)
}

// Textbox represents a multi-line text entry box.
// Text in a Textbox is unformatted, and scrollbars are applied automatically.
// TODO rename to TextBox? merge with TextField (but cannot use Invalid())? enable/disable line wrapping?
// TODO events
// TODO Tab key - insert horizontal tab or tab stop?
// TODO ReadOnly
// TODO line endings
type Textbox interface {
	Control

	// Text and SetText get and set the Textbox's text.
	Text() string
	SetText(text string)
}

// NewTextbox creates a new Textbox.
func NewTextbox() Textbox {
	return newTextbox()
}

// Spinbox is a Control that provides a text entry field that accepts integers and up and down buttons to increment and decrement those values.
// This control is in its preliminary state.
// TODO everything:
// - TODO set increment? (work on windows)
// - TODO set page step?
// - TODO wrapping
// - TODO negative values
type Spinbox interface {
	Control

	// Value and SetValue get and set the current value of the Spinbox, respectively.
	// For SetValue, if the new value is outside the current range of the Spinbox, it is set to the nearest extremity.
	Value() int
	SetValue(value int)

	// OnChanged sets the event handler for when the Spinbox's value is changed.
	// Under what conditions this event is raised when the user types into the Spinbox's edit field is platform-defined.
	OnChanged(func())
}

// NewSpinbox creates a new Spinbox with the given minimum and maximum.
// The initial value will be the minimum value.
// NewSpinbox() panics if min > max.
func NewSpinbox(min int, max int) Spinbox {
	if min > max {
		panic("min > max in NewSpinbox()")
	}
	return newSpinbox(min, max)
}

// ProgressBar is a Control that displays a horizontal bar which shows the level of completion of an operation.
// TODO indetermiante
type ProgressBar interface {
	Control

	// Percent and SetPrecent get and set the current percentage indicated by the ProgressBar, respectively.
	// This value must be between 0 and 100; all other values cause SetPercent to panic.
	// TODO rename to Progress/SetProgress?
	Percent() int
	SetPercent(percent int)
}

// NewProgressBar creates a new ProgressBar.
// It will initially show a progress of 0%.
func NewProgressBar() ProgressBar {
	return newProgressBar()
}
