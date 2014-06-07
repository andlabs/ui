MAC OS X:
- NSComboBox scans the entered text to see if it matches one of the items and returns the index of that item if it does; find out how to suppress this so that it returns -1 unless the item was chosen from the list (like the other platforms)
	- asked: http://stackoverflow.com/questions/23046414/cocoa-how-do-i-get-nscombobox-indexofselecteditem-to-return-1-if-the-user-m
- make sure Areas get keyboard focus when clicking outside the actual Area space on Mac OS X

WINDOWS:
- windows: windows key handling is just wrong; figure out how to avoid (especially since Windows intercepts that key by default)
- the windows build appears to be unstable:
	- 64-bit crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 it works fine
- redrawing controls after a window resize on Windows does not work properly
- on windows 7, progress bars seem to animate from 0 -> pos when you turn off marquee mode and set pos; see if that's documented or if I'm doing something wrong
	- intentional: http://social.msdn.microsoft.com/Forums/en-US/61350dc7-6584-4c4e-91b0-69d642c03dae/progressbar-disable-smooth-animation http://stackoverflow.com/questions/2217688/windows-7-aero-theme-progress-bar-bug http://discuss.joelonsoftware.com/default.asp?dotnet.12.600456.2 http://stackoverflow.com/questions/22469876/progressbar-lag-when-setting-position-with-pbm-setpos http://stackoverflow.com/questions/6128287/tprogressbar-never-fills-up-all-the-way-seems-to-be-updating-too-fast - these links have workarounds but blah; more proof that progressbars were programmatically intended to be incremented in steps
	- related: in wine,
		- set progress to 0, indeterminate, dec - frozen indeterminate animation
		- set progress to 100, indeterminate, inc - frozen indetemrinate animation
		- need to see if this is a wine bug or not
- check all uses of RECT.right/.bottom in Windows that don't have an accompanying -RECT.left/.top to make sure they're correct
- when adding IsDialogMessage() find out if that make sthe area in the area bounds test automatically focused

UNIX:
- double-check to make sure MouseEvent.Held[] is sorted on Unix after we figure out how to detect buttons above button 5
- david wendt is telling me he's getting frequent crashes on his end with the GTK+ amd64 build...
	TODO re-evaluate; I think I fixed them all ages ago now
- when resizing a GTK+ window smaller than a certain size, the controls inside will start clipping in bizarre ways (the horizontal scrollbar in Area will disappear smoothly; etc.)
- resizing seems to be completely and totally broken in the Wayland backend
	- TODO find out if this is a problem on the GTK+/Wayland side (no initial window-configure event?)
- [12:55] <myklgo> pietro10: I meant to mention: 1073): Gtk-WARNING **: Theme parsing error: gtk.css:72:20: Not using units is deprecated. Assuming 'px'.    twice.
- figure out why Page Up/Page Down does tab stops

ALL PLATFORMS:
- make sure MouseEvent's documentation has dragging described correctly (both Windows and GTK+ do)
- make sure the preferred size of a Listbox is the minimum size needed to display everything on all platforms (capped at the screen height, of course?)
- make sure the image drawn on an Area looks correct on all platforms (is not cropped incorrectly or blurred)
- make all widths and heights parameters in constructors in the same place (or drop the ones in Window entirely?)
- Message boxes that belong to agiven parent are still application-modal on all platforms except Mac OS X because the whole system waits... we'll need to use a channel for this, I guess :S
