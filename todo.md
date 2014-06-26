ALL:
- vertical alignment of labels is a botch right now
	- os x: very very wrong

MAC OS X:
- NSComboBox scans the entered text to see if it matches one of the items and returns the index of that item if it does; find out how to suppress this so that it returns -1 unless the item was chosen from the list (like the other platforms)
	- asked: http://stackoverflow.com/questions/23046414/cocoa-how-do-i-get-nscombobox-indexofselecteditem-to-return-1-if-the-user-m
- make sure Areas get keyboard focus when clicking outside the actual Area space on Mac OS X
	- http://stackoverflow.com/questions/24102367/how-do-i-make-it-so-clicking-outside-the-actual-nsview-in-a-nsscrollview-but-wit
- on initially starting the Area test, layout is totally wrong
- probably use fittingSize instead of sizeToFit

WINDOWS:
- there seems to be a caching issue: with the test program and `-dialog`, click one of the dialog buttons, then quickly tap one of the buttons in the main window. The dialog will pop up twice, and after both are closed the program aborts with a send on closed channel
	- appears to be a bug in my dialog code
	- appears to have *always* been a bug in dialog code...
- windows: windows key handling is just wrong; figure out how to avoid (especially since Windows intercepts that key by default)
- control sizing is a MAJOR pain
	- http://stackoverflow.com/questions/24130548/is-there-a-proper-way-to-get-the-preferred-size-of-windows-controls-there-are-s
- redrawing controls after a window resize on Windows does not work properly
- when adding IsDialogMessage() find out if that makes the area in the area bounds test automatically focused
- vertical alignment of labels can be wrong (I think; need to check again - TODO)

UNIX:
- figure out how to detect the alt key and mouse buttons above 5 properly for modifiers/Held[]
	- http://stackoverflow.com/questions/24053012/gdk-is-there-a-way-i-can-get-if-an-arbitrary-keyval-or-mouse-button-is-down-dur
- double-check to make sure MouseEvent.Held[] is sorted on Unix after we figure out how to detect buttons above button 5
- sizing with client-side decorations (Wayland) don't work
	- several people suggested connecting to size-allocate of the GtkLayout, but then I can wind up in a situation where there's extra padding or border space in the direction I resized
- [12:55] <myklgo> pietro10: I meant to mention: 1073): Gtk-WARNING **: Theme parsing error: gtk.css:72:20: Not using units is deprecated. Assuming 'px'.    twice.
- figure out why Page Up/Page Down does tab stops

ALL PLATFORMS:
- windows code presently wraps entire function bodies in uitask; make sure the other platforms do too
