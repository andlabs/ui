so I don't forget:
- Window.SizeToFit() or WIndow.OptimalSize() (use: `Window.SetOptimalSize())`) for sizing a window to the control's interest
- Control.Show()/Control.Hide()
- Groupbox
- character-limited entry fields, numeric entry fields, multiline entry fields
	- possible rename of LineEdit?
		- especially for password fields - NewPasswordEntry()?
- [all platforms] Listbox should have a border style
	- [Windows] a different border on LineEdits and Listboxes
- padding and spacing in Stack
- allow Combobox to have initial settings
- Combobox and Listbox insertions and deletions should allow bulk (...string)
- Combobox/Listbox.DeleteAll
- Combobox/Listbox.Select (with Listbox.Select allowing bulk)
	- Checkbox.Check or Checkbox.SetChecked
- Listbox.SelectAll
- should Labels be selectable?
	- should message box text be selectable on all platforms or only on those that make it the default?
- Listbox/Combobox.Index(n)
	- Index(n) is the name used by reflect.Value; use a different one?
- Message boxes should not show secondary text if none is specified.
- note that you can change event channels before opening the window; this allows unifying menus/toolbars/etc.
	- will probably want to bring back Event() for this (but as NewEvent())
- add bounds checking to Area's sizing methods
- describe thread-safety of Area.SetSize()

important things:
- because the main event loop is not called if initialization fails, it is presently impossible for MsgBoxError() to work if UI initialization fails; this basically means we cannot allow initializiation to fail on Mac OS X if we want to be able to report UI init failures to the user with one (which would be desirable, maybe (would violate Windows HIG?))
- figure out where to auto-place windows in Cocoa (also window coordinates are still not flipped properly so (0,0) on screen is the bottom-left)
	- also provide a method to center windows; Cocoa provides one for us but
- I think Cocoa NSButton text is not vertically aligned properly...?
	- and listbox item text is too low?
- NSPopUpButton does allow no initial selection ([b setSelectedIndex:-1]); use it
	- need to use it /after/ adding initial items, otherwise it won't work
	- find out if I can do the same with the ListBoxes
- NSComboBox scans the entered text to see if it matches one of the items and returns the index of that item if it does; find out how to suppress this so that it returns -1 unless the item was chosen from the list (like the other platforms)
- some Cocoa controls don't seem to resize correctly: Buttons have space around the edges and don't satisfy stretchiness
- make sure GTK+ documentation version point differences (x in 4.3.x) don't matter
- LineEdit heights on Windows seem too big; either that or LineEdit, Button, and Label text is not vertically centered properly
	- are Checkboxes too small?
	- Cocoa has similar margining issues (like Comboboxes having margins)
- sometimes the size of the drop-down part of a Combobox becomes 0 or 1 or some other impossibly small value on Windows
- make gcc (Unix)/clang (Mac OS X) pedantic about warnings/errors; also -Werror
- make sure scrollbars in Listbox work identically on all platforms (specifically the existence and autohiding of both horizontal and vertical scrollbars)
	- pin down this behavior; also note non-editability
- listboxes spanning the vertical height of the window don't always align with the bottom border of the edit control attached to the bottom of the window...
- make sure mouse events don't trigger if the control size is larger than the Area size and the mouse event happens outside the Area range on all platforms

super ultra important things:
- formalize what happens if Modifiers by themselves are held
	- OS X: find out if multiple DIFFERENT modifiers released at once produces multiple events
	- in general, figure out what to do on multiple events, period
- OS X: handle Insert/Help key change in a sane and deterministic way
	- will need old and new Mac keyboards...
- should pressing modifier+key in the keyboard test mark the key alone as pressed as well? I'm leaning toward no, in which case make sure this behavior exists on all platforms
- formalize dragging
	- implement dragging on windows
	- may need to drop Held depending on weirdness I see in OS X
- cap click count to 2 on all platforms
	- cap mouse button count to 3? or should a function be used instead?
- the windows build appears to be unstable:
	- 64-bit crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 it works fine
	- 32-bit: it works, but if I save the class name converted to UTF-16 beforehand, wine indicates that the class name is replaced with the window title, so something there is wrong...
- david wendt is telling me he's getting frequent crashes on his end with the GTK+ amd64 build...
	TODO re-evaluate; I think I fixed them all ages ago now
- occasionally I get
		panic: error sending message to message loop to call function: Invalid thread ID.
	when starting up the windows/386 build; race in ui()/msgloop()?
	- happens the first time I run a new build in wine; also if my computer is running too slowly when running in wine
- GTK+: stderr is flooded with
```
(test:17575): Gdk-CRITICAL **: gdk_device_ungrab: assertion 'GDK_IS_DEVICE (device)' failed

(test:17575): Gtk-CRITICAL **: gtk_device_grab_remove: assertion 'GDK_IS_DEVICE (device)' failed
```
	figure out why
	- I think it has to do with invalid list deletions; roll back the panics and check
- the user can still [NSApp terminate:] from the Dock icon, bypassing Go itself
	- ideally we need a QuitItem() function for this case if/when we add menus
	- check this on all platforms
- Cocoa: NSScrollView support is hacky at best
	- https://developer.apple.com/library/mac/documentation/cocoa/Conceptual/NSScrollViewGuide/Articles/Creating.html#//apple_ref/doc/uid/TP40003226-SW4 the warning about pixel alignment may or may not be heeded, not sure
	- frame sizes are a bit of a hack: the preferred size of a NSScrollView is the preferred size of its document view; the frameSize method described on the above link might be better but a real solution is optimal
- make sure the image drawn on an Area looks correct on all platforms (is not cropped incorrectly or blurred)
- when resizing a GTK+ window smaller than a certain size, the controls inside will start clipping in bizarre ways (progress bars/entry lines will just cut off; editable comboboxes will stretch slightly longer than noneditable ones; the horizontal scrollbar in Area will disappear smoothly; etc.)
- the window background of a GTK+ window seems to be... off - I think it has to do with the GtkLayout
- see update 18 March 2014 in README
- resizing seems to be completely and totally broken in the Wayland backend
- redrawing Areas on Windows seems to be flaky
- make sure the first and last rows and columns of an Area are being drawn on Windows
- clicking on Areas in GTK+ don't bring keyboard focus to them?
- make sure GTK+ keyboard events on numpad off don't switch between controls
- for our two custom window classes, we should allocate extra space in the window class's info structure and then use SetWindowLongPtrW() during WM_CREATE to store the sysData and not have to make a new window class each time; this might also fix the s != nil && s.hwnd != 0 special cases in the Area WndProc if done right
	- references: https://github.com/glfw/glfw/blob/master/src/win32_window.c#L182, http://www.catch22.net/tuts/custom-controls
- Area redraw on Windows is still a bit flaky, especially after changing the Area size to something larger than the window size and then resizing the window(???)
- despite us explicitly clearing the clip area on Windows, Area still doesn't seem to draw alpha bits correctly... it appears as if we are drawing over the existing image each time
- on Windows, Shift+(num pad key) triggers the shifted key code when num lock is off; will need to reorder key code tests on all platforms to fix this
- pressing global keycodes (including kwin's zoom in/out) when running the keyboard test in wine causes the Area to lose keyboard focus; this doesn't happen on the GTK+ version (fix the Windows version to behave like the GTK+ version)

important things:
- make specific wording in documentation consistent (make/create, etc.)
	- document minor details like wha thappens on specific events so that they are guaranteed to work the same on all platforms (are there any left?)
		- what happens when the user clicks and drags on a listbox
	- should field descriptions in method comments include the receiver name? (for instance e.Held vs. Held) - see what Go's own documentation does
- make passing of parameters and type conversions of parameters to uitask on Windows consistent: explicit _WPARAM(xxx)/_LPARAM(xxx)/uintptr(xxx), for example
	- do this for type signatures in exported functions: (err error) or just error?
	- do this for the names of GTK+ helper functions (gtkXXX or gXXX)
- on windows 7, progress bars seem to animate from 0 -> pos when you turn off marquee mode and set pos; see if that's documented or if I'm doing something wrong
- clean up windows struct field names (holdover from when the intent was to make a wrapper lib first and then use it rather than using the windows API directly)
- make all widths and heights parameters in constructors in the same place (or drop the ones in Window entirely?)

far off:
- localization
- strip unused constants from the Windows files
- combine more Windows files; rename some?
- tab stops

maybe:
- rename Stack to Box?
- make Combobox and Listbox satisfy sort.Interface?
- should a noneditable Combobox be allowed to return to unselected mode by the user?
- provide a way for MouseEvent/KeyEvent to signal that the keypress caused the Area to gain focus
	- provide an event for leaving focus so a focus rectangle can be drawn
- change the Windows code to use extra class space (as in http://www.catch22.net/tuts/custom-controls)
	- this is a bit flakier as SetWindowLongPtr() can fail, and it can also succeed in such a way that the last error is unreliable
