// +build !windows,!darwin,!plan9

// 23 february 2014
package ui

import (
	// ...
)

// GTK+ 3 makes this easy: controls can tell us what their preferred size is!
// ...actually, it tells us two things: the "minimum size" and the "natural size".
// The "minimum size" is the smallest size we /can/ display /anything/. The "natural size" is the smallest size we would /prefer/ to display.
// The difference? Minimum size takes into account things like truncation with ellipses: the minimum size of a label can allot just the ellipses!
// So we use the natural size instead, right?
// We could, but there's one snag: "Handle with care. Note that the natural height of a height-for-width widget will generally be a smaller size than the minimum height, since the required height for the natural width is generally smaller than the required height for the minimum width."
// This will have to be taken care of manually, so TODO; we'll just use the natural size for now

func (s *sysData) preferredSize() (width int, height int) {
	_, _, width, height = gtk_widget_get_preferred_size(s.widget)
	return width, height
}
