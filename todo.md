so I don't forget:
- Window.SizeToFit() or WIndow.OptimalSize() (use: `Window.SetOptimalSize())`) for sizing a window to the control's interest
- Control.Show()/Control.Hide()
- Groupbox
- see if we really need to track errors on a lot of places that report errors
	- Window.Show()/Window.Hide() report errors due to UpdateWindow(), which can fail, but that is only called when the window is first opened: split that functionality out
- character-limited entry fields, numeric entry fields, multiline entry fields
	- possible rename of LineEdit?
		- especially for password fields - NewPasswordEntry()?
- [Windows, Mac OS X] should ListBox have a border style?
	- [Windows] a different border on LineEdits?
- padding and spacing in Stack
- change Listbox constructor so that there's a separate constructor for each variant, rather than passing in parameters
- allow Combobox to have initial settings
- Combobox and Listbox insertions and deletions should allow bulk (...string)
- Combobox/Listbox.DeleteAll
- Combobox/Listbox.Select (with Listbox.Select allowing bulk)
- Listbox.SelectAll
- have Combobox.InsertBefore, Listbox.InsertBefore, Combobox.Delete, and Listbox.Delete return an error on invalid index before creation, or have them panic like an invalid array index, etc.; decide which to do as these do differnet things on different platforms by default
	- same for other methods that take indices, like the Stack and Grid stretchy methods
- make the Windows implementation of message boxes run on uitask
	- ensure MsgBoxError can run if initialization failed if things change ever
- should Labels be selectable?
	- should message box text be selectable on all platforms or only on those that make it the default?

important things:
- because the main event loop is not called if initialization fails, it is presently impossible for MsgBoxError() to work if UI initialization fails; this basically means we cannot allow initializiation to fail on Mac OS X if we want to be able to report UI init failures to the user with one
- figure out where to auto-place windows in Cocoa (also window coordinates are still not flipped properly so (0,0) on screen is the bottom-left)
	- also provide a method to center windows; Cocoa provides one for us but
- I think Cocoa NSButton text is not vertically aligned properly...?
- NSPopUpButton does allow no initial selection ([b setSelectedIndex:-1]); use it
	- find out if I can do the same with the ListBoxes
- NSComboBox scans the entered text to see if it matches one of the items and returns the index of that item if it does; find out how to suppress this so that it returns -1 unless the item was chosen from the list (like the other platforms)
- some Cocoa controls don't seem to resize correctly: Buttons have space around the edges and don't satisfy stretchiness
- there's no GTK+ or Cocoa error handling whatsoever; we need to figure out how it works
	- I know how Cocoa error handling works: it uses Cocoa exceptions; need to figure out how to catch and handle them somehow
- make sure GTK+ documentation version point differences (x in 4.3.x) don't matter
- button sizes and LineEdit sizes on Windows seem too big; Comboboxes have margins
	- Cocoa has similar margining issues (like on Comboboxes)
- sometimes the size of the drop-down part of a Combobox becomes 0 or 1 or some other impossibly small value on Windows
- make gcc (Unix)/clang (Mac OS X) pedantic about warnings/errors; also -Werror
- make sure scrollbars in Listbox work identically on all platforms (specifically the existence and autohiding of both horizontal and vertical scrollbars)
- GTK+ windows cannot be resized smaller than their controls's current sizes in their current positions; find out how to overrule that so they can be freely resized

super ultra important things:
- the windows build appears to be unstable:
	- 64-bit doesn't work, period: it crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 it works fine
	- 32-bit: it works now, but if I save the class name converted to UTF-16 beforehand, wine indicates that the class name is replaced with the window title, so something there is wrong...
- handle in-library panics (internal errors) by reporting them to the user
- david wendt is telling me he's getting frequent crashes on his end with the GTK+ amd64 build...
	TODO re-evaluate; I think I fixed them all ages ago now
- occasionally I get
		panic: error sending message to message loop to call function: Invalid thread ID.
	when starting up the windows/386 build; race in ui()/msgloop()?
	- happens the first time I run a new build in wine; also if my computer is running too slowly when running in wine

important things:
- make specific wording in documentation consistent (make/create, etc.)
	- document minor details like wha thappens on specific events so that they are guaranteed to work the same on all platforms
		- for instance, initial selection state of Combobox and Listbox
			- related: should a noneditable Combobox be allowed to return to unselected mode by the user?
- make passing of parameters and type conversions of parameters to uitask consistent
	- TODO figure out what I meant by this; I don't remember

far off:
- localization
- strip unused constants from the Windows files
- combine more Windows files; rename some?
- tab stops

maybe:
- rename Stack to Box?
- make Combobox and Listbox satisfy sort.Interface?
- indeterminate progress bars (not supported on Windows 2000)
