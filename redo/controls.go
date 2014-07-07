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

	// OnClicked sets the event handler for when the Button is clicked.
	OnClicked(func(d Doer))

	// Text and SetText are Requests that get and set the Button's label text.
	Text() *Request
	SetText(text string) *Request
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
	OnClicked(func(d Doer))

	// Text and SetText are Requests that get and set the Checkbox's label text.
	Text() *Request
	SetText(text string) *Request

	// Checked and SetChecked are Requests that get and set the Checkbox's check state.
	Checked() *Request
	SetChecked(checked bool) *Request
}

// NewCheckbox creates a new Checkbox with the given label text.
// The Checkbox will be initially unchecked.
func NewCheckbox(text string) Checkbox {
	return newCheckbox(text)
}

// Combobox is a drop-down list from which one item can be selected.
// Each item of a Combobox is a text string.
// The Combobox can optionally be editable, in which case the user can type in a selection not in the list.
// [TODO If an item is selected in an editable Combobox, the edit field will be changed ot reflect the selection.]
type Combobox interface {
	Control

	// TODO events

	// Append, InsertBefore, and Delete are Requests that change the Combobox's list.
	// InsertBefore and Delete panic if the index passed in is out of range.
	Append(item string) *Request
	InsertBefore(item string, before int) *Request
	Delete(index int) *Request

	// SelectedIndex and SelectedText are Requests that return the current Combobox selection, either as the index into the list or as its label.
	// SelectedIndex returns -1 and SelectedText returns an empty string if no selection has been made.
	// If the Combobox is editable, SelectedIndex returns -1 if the user has entered their own string, in which case SelectedText will return that string.
	SelectedIndex() *Request
	SelectedText() *Request

	// SelectIndex is a Request that selects an index from the list.
	// SelectIndex panics if the given index is out of range.
	// [TODO SelectText or SetCustomText]
	SelectIndex(index int) *Request

	// Len is a Request that returns the number of items in the list.
	// At is a Request that returns a given item's text.
	// At panics if the given index is out of range.
	Len() *Request
	At(index int) *Request
}

// NewCombobox creates a new Combobox with the given items.
// The Checkbox will have nothing selected initially.
func NewCombobox(items ...string) Combobox {
	return newCombobox(items)
}

// NewEditableCombobox creates a new editable Combobox with the given items.
// The Combobox will have nothing selected initially and no custom text initially.
func NewEditableCombobox(items ...string) Combobox {
	return newEditableCombobox(items)
}

// LineEdit
// Label
// Listox
// ProgressBar
// (Area, Stack, and Grid will remain in their own file)
