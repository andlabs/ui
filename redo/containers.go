// 25 july 2014

package ui

// Tab is a Control that contains multiple pages of tabs, each containing a single Control.
// You can add and remove tabs from the Tab at any time.
// TODO rename?
// TODO implement containerShow()/containerHide() on this
type Tab interface {
	Control

	// Append adds a new tab to Tab.
	// The tab is added to the end of the current list of tabs.
	Append(name string, control Control)

	// Delete removes the given tab.
	// It panics if index is out of range.
//	Delete(index int)
//TODO
}

// NewTab creates a new Tab with no tabs.
func NewTab() Tab {
	return newTab()
}
