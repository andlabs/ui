so I don't forget:
- Window.SizeToFit() or WIndow.OptimalSize() (use: `Window.SetSize(Window.OptimalSize())`) for sizing a window to the control's interest
- Control.Show()/Control.Hide()
- Groupbox
- determine if a selection in a non-editable combobox has been made
- see if we really need to track errors on Combobox.Selection()
	- in fact, see if we really need to track errors on a lot of things...
- character-limited entry fields, numeric entry fields, multiline entry fields
	- possible rename of LineEdit?
		- especially for password fields - NewPasswordEntry()?
- more flexible size appropriation: allow a small button to be at the top of everything in the main() example here
- [Windows] should ListBox have a border style?
	- a different border on LineEdits?
- padding and spacing in Stack; maybe a setting in Stack which keeps controls at their preferred size?
- change Stack/Combobox/Listbox constructors so that there's a separate constructor for each variant, rather than passing in parameters?
- allow Combobox to have initial settings
- Combobox and Listbox insertions and deletions should allow bulk (...string)
- Combobox/Listbox.DeleteAll
- Combobox/Listbox.Select (with Listbox.Select allowing bulk)
- Listbox.SelectAll
- have Combobox.InsertBefore, Listbox.InsertBefore, Combobox.Delete, and Listbox.Delete return an error on invalid index before creation
- make the Windows implementation of message boxes run on uitask
	- ensure MsgBoxError can run if initialization failed if things change ever
- should Labels be selectable?
	- should message box text be selectable on all platforms or only on those that make it the default?

important things:
- ui.Go() should exit when the main() you pass in exits
- because the main event loop is not called if initialization fails, it is presently impossible for MsgBoxError() to work if UI initialization fails; this basically means we cannot allow initializiation to fail on Mac OS X if we want to be able to report UI init failures to the user with one
- Cocoa coordinates have (0,0) at the bottom left: need to fix this somehow
- I think Cocoa NSButton text is not vertically aligned properly...?
- NSPopUpButton doesn't seem to allow no initial selection? need to be sure
- NSComboBox scans the entered text to see if it matches one of the items and returns the index of that item if it does; find out how to suppress this so that it returns -1 unless the item was chosen from the list (like the other platforms)
- figure out what to do about deleting a nonexistent item; each backend responds differently by default
- some Cocoa controls don't seem to resize correctly: Buttons have space around the edges and don't satisfy stretchiness; Labels in the grid test window get cropped
- find out how to deselect the first item in the NSTableView when first creating the Listbox (to make it work like on other platforms)
- there's no GTK+ error handling whatsoever; we need to figure out how it works
- make sure GTK+ documentation point differences don't matter
- button sizes and LineEdit sizes on Windows seem too big; Comboboxes have margins
- sometimes the size of the drop-down part of a Combobox becomes 0 or 1 or some other impossibly small value on Windows
- make gcc (Unix)/clang (Mac OS X) pedantic about warnings/errors; also -Werror

super ultra important things:
- for some reason events are now delayed on windows
- the windows build appears to be unstable:
	- 64-bit doesn't work, period: it crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 it works fine
	- 32-bit: it works now, but if I save the class name converted to UTF-16 beforehand, wine indicates that the class name is replaced with the window title, so something there is wrong...
- on 64-bit windows 7 comboboxes don't show their lists
- handle in-library panics (internal errors) by reporting them to the user
- david wendt is telling me he's getting frequent crashes on his end with the GTK+ amd64 build...
	TODO re-evaluate; I think I fixed them all ages ago now
- occasionally I get
		panic: error sending message to message loop to call function: Invalid thread ID.
	when starting up the windows/386 build; race in ui()/msgloop()?

important things:
- Control.preferredSize() (definitely needed for Grid and Form)
- make specific wording in documentation consistent (make/create, etc.)
- make passing of parameters and type conversions of parameters to uitask consistent

far off:
- localization
- strip unused constants from the Windows files
- combine more Windows files; rename some?

maybe:
- rename Stack to Box?
- make Combobox and Listbox satisfy sort.Interface?
- indeterminate progress bars (not supported on Windows 2000)
