MAC OS X:
- NSComboBox scans the entered text to see if it matches one of the items and returns the index of that item if it does; find out how to suppress this so that it returns -1 unless the item was chosen from the list (like the other platforms)
	- asked: http://stackoverflow.com/questions/23046414/cocoa-how-do-i-get-nscombobox-indexofselecteditem-to-return-1-if-the-user-m
- make sure Areas get keyboard focus when clicking outside the actual Area space on Mac OS X
	- http://stackoverflow.com/questions/24102367/how-do-i-make-it-so-clicking-outside-the-actual-nsview-in-a-nsscrollview-but-wit

WINDOWS:
- windows: windows key handling is just wrong; figure out how to avoid (especially since Windows intercepts that key by default)
- redrawing controls after a window resize on Windows does not work properly
- when adding IsDialogMessage() find out if that make sthe area in the area bounds test automatically focused

UNIX:
- double-check to make sure MouseEvent.Held[] is sorted on Unix after we figure out how to detect buttons above button 5
- when resizing a GTK+ window smaller than a certain size, the controls inside will start clipping in bizarre ways (the horizontal scrollbar in Area will disappear smoothly; etc.)
- resizing seems to be completely and totally broken in the Wayland backend
	- https://github.com/andlabs/misctestprogs/blob/master/gtkwaylandtest.c http://imgur.com/qdxMFkY,7mnCXrL#0 wait for Jasper in IRC
- [12:55] <myklgo> pietro10: I meant to mention: 1073): Gtk-WARNING **: Theme parsing error: gtk.css:72:20: Not using units is deprecated. Assuming 'px'.    twice.
- figure out why Page Up/Page Down does tab stops

ALL PLATFORMS:
- make sure MouseEvent's documentation has dragging described correctly (both Windows and GTK+ do)
- make all widths and heights parameters in constructors in the same place (or drop the ones in Window entirely?)
