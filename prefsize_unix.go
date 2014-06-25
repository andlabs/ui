// +build !windows,!darwin,!plan9

// 23 february 2014

package ui

// GTK+ 3 makes this easy: controls can tell us what their preferred size is!
// ...actually, it tells us two things: the "minimum size" and the "natural size".
// The "minimum size" is the smallest size we /can/ display /anything/. The "natural size" is the smallest size we would /prefer/ to display.
// The difference? Minimum size takes into account things like truncation with ellipses: the minimum size of a label can allot just the ellipses!
// So we use the natural size instead.
// There is a warning about height-for-width controls, but in my tests this isn't an issue.
// For Areas, we manually save the Area size and use that, just to be safe.

// We don't need to worry about y-offset because label alignment is "vertically center", which GtkLabel does for us.

func (s *sysData) preferredSize() (width int, height int, yoff int) {
	if s.ctype == c_area {
		return s.areawidth, s.areaheight, 0
	}

	_, _, width, height = gtk_widget_get_preferred_size(s.widget)
	return width, height, 0
}
