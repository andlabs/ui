so I don't forget:
- describe thread-safety of Area.SetSize()
- should all instances of -1 as error returns from Windows functions be changed to ^0 or does the uintptr() conversion handle sign extension?

important things:
- I think Cocoa listbox item text is too low?
- NSPopUpButton does allow no initial selection ([b setSelectedIndex:-1]); use it
	- need to use it /after/ adding initial items, otherwise it won't work
	- find out if I can do the same with the ListBoxes
- NSComboBox scans the entered text to see if it matches one of the items and returns the index of that item if it does; find out how to suppress this so that it returns -1 unless the item was chosen from the list (like the other platforms)
- some Cocoa controls don't seem to resize correctly: Buttons have space around the edges
- LineEdit heights on Windows seem too big; either that or LineEdit, Button, and Label text is not vertically centered properly
	- are Checkboxes and Comboboxes too small?
	- Cocoa has similar margining issues (like Comboboxes having margins)
- make gcc (Unix)/clang (Mac OS X) pedantic about warnings/errors; also -Werror
	- problem: cgo-generated files trip -Werror up; I can't seem to turn off unused argument warnings with the -Wall/-Wextra/-pedantic options
- consolidate scroll view code in GTK+ and Mac OS X
- make sure mouse events trigger when we move the mouse over an Area with a button held on OS X
- area test time label weirdness
	- does not show anything past the date on windows
	- does not show initially on OS X; it shows up once you resize, and even shows up after you resize back to the original size

super ultra important things:
- window sizes are inconsistent and wrong: we need to see on which platforms the titlebars/borders/etc. count...
- formalize what happens if Modifiers by themselves are held
	- OS X: find out if multiple DIFFERENT modifiers released at once produces multiple events
	- in general, figure out what to do on multiple events, period
- OS X: handle Insert/Help key change in a sane and deterministic way
	- will need old and new Mac keyboards...
- should pressing modifier+key in the keyboard test mark the key alone as pressed as well? I'm leaning toward no, in which case make sure this behavior exists on all platforms
- make sure MouseEvent's documentation has dragging described correctly (both Windows and GTK+ do)
	- fix OS X so that it follows these rules
- cap click count to 2 on all platforms
	- cap mouse button count to 3? or should a function be used instead?
- the windows build appears to be unstable:
	- 64-bit crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 it works fine
	- 32-bit: it works, but if I save the class name converted to UTF-16 beforehand, wine indicates that the class name is replaced with the window title, so something there is wrong...
- david wendt is telling me he's getting frequent crashes on his end with the GTK+ amd64 build...
	TODO re-evaluate; I think I fixed them all ages ago now
- GTK+: stderr is flooded with
```
(test:17575): Gdk-CRITICAL **: gdk_device_ungrab: assertion 'GDK_IS_DEVICE (device)' failed

(test:17575): Gtk-CRITICAL **: gtk_device_grab_remove: assertion 'GDK_IS_DEVICE (device)' failed
```
	figure out why
	- I think it has to do with invalid list deletions; roll back the panics and check
- Cocoa: NSScrollView support is hacky at best
	- https://developer.apple.com/library/mac/documentation/cocoa/Conceptual/NSScrollViewGuide/Articles/Creating.html#//apple_ref/doc/uid/TP40003226-SW4 the warning about pixel alignment may or may not be heeded, not sure
	- frame sizes are a bit of a hack: the preferred size of a NSScrollView is the preferred size of its document view; the frameSize method described on the above link might be better but a real solution is optimal
- make sure the image drawn on an Area looks correct on all platforms (is not cropped incorrectly or blurred)
- when resizing a GTK+ window smaller than a certain size, the controls inside will start clipping in bizarre ways (progress bars/entry lines will just cut off; editable comboboxes will stretch slightly longer than noneditable ones; the horizontal scrollbar in Area will disappear smoothly; etc.)
- see update 18 March 2014 in README
- resizing seems to be completely and totally broken in the Wayland backend
	- TODO find out if this is a problem on the GTK+/Wayland side (no initial window-configure event?)
- redrawing Areas on Windows seems to be flaky: make the window small, scroll, then make it large again and watch the vertical corruption (alternatively "especially after changing the Area size to something larger than the window size and then resizing the window(???)")
	- redrawing controls after a window resize on Windows seems to be flaky
	- for Area again, the area bounds test's borders don't redraw correctly after certain series of resize operations
- clicking on Areas in GTK+ don't bring keyboard focus to them?
- make sure keyboard events on numpad off on all platforms don't switch between controls
- despite us explicitly clearing the clip area on Windows, Area still doesn't seem to draw alpha bits correctly... it appears as if we are drawing over the existing image each time
- on Windows, Shift+(num pad key) triggers the shifted key code when num lock is off; will need to reorder key code tests on all platforms to fix this
- pressing global keycodes (including kwin's zoom in/out) when running the keyboard test in wine causes the Area to lose keyboard focus; this doesn't happen on the GTK+ version (fix the Windows version to behave like the GTK+ version)
- GTK+ indefinite progress bar animation is choppy: make sure the speed we have now is the conventional speed for GTK+ programs (HIG doesn't list any) and that the choppiness is correct
- Message boxes are not application-modal on some platforms
- cast all objc_msgSend() direct invocations to the approrpiate types; this is how you're supposed to do things: https://developer.apple.com/library/ios/documentation/General/Conceptual/CocoaTouch64BitGuide/ConvertingYourAppto64-Bit/ConvertingYourAppto64-Bit.html http://lists.apple.com/archives/objc-language/2014/Jan/msg00011.html http://lists.apple.com/archives/cocoa-dev/2006/Feb/msg00753.html and many others

other things:
- on windows 7, progress bars seem to animate from 0 -> pos when you turn off marquee mode and set pos; see if that's documented or if I'm doing something wrong
	- intentional: http://social.msdn.microsoft.com/Forums/en-US/61350dc7-6584-4c4e-91b0-69d642c03dae/progressbar-disable-smooth-animation http://stackoverflow.com/questions/2217688/windows-7-aero-theme-progress-bar-bug http://discuss.joelonsoftware.com/default.asp?dotnet.12.600456.2 http://stackoverflow.com/questions/22469876/progressbar-lag-when-setting-position-with-pbm-setpos http://stackoverflow.com/questions/6128287/tprogressbar-never-fills-up-all-the-way-seems-to-be-updating-too-fast - these links have workarounds but blah; more proof that progressbars were programmatically intended to be incremented in steps
- clean up windows struct field names (holdover from when the intent was to make a wrapper lib first and then use it rather than using the windows API directly)
- make all widths and heights parameters in constructors in the same place (or drop the ones in Window entirely?)
