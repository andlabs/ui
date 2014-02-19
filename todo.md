so I don't forget:
- Window.SizeToFit() or WIndow.OptimalSize() (use: `Window.SetSize(Window.OptimalSize())`) for sizing a window to the control's interest
- Control.Show()/Control.Hide()
- Control.SetText()
- Groupbox
- determine if a selection in a non-editable combobox has been made
- see if we really need to track errors on Combobox.Selection()
	- in fact, see if we really need to track errors on a lot of things...
- password entry fields, character-limited entry fields, numeric entry fields, multiline entry fields
	- possible rename of LineEdit?
- more flexible size appropriation: allow a small button to be at the top of everything in the main() example here
- [Windows] should ListBox have a border style?
- padding and spacing in Stack; maybe a setting in Stack which keeps controls at their preferred size?
- change Stack/Combobox/Listbox constructors so that there's a separate constructor for each variant, rather than passing in parameters?
- allow Combobox to have initial settings
- Combobox and Listbox insertions and deletions should allow bulk (...string)
- Combobox/Listbox.DeleteAll
- Combobox/Listbox.Select (with Listbox.Select allowing bulk)
- Listbox.SelectAll
- have Combobox.InsertBefore, Listbox.InsertBefore, Combobox.Delete, and Listbox.Delete return an error on invalid index before creation
- make the Windows implementation of message boxes run on uitask

important things:
- there's no GTK+ error handling whatsoever; we need to figure out how it works
- make sure GTK+ documentation point differences don't matter

super ultra important things:
- for some reason events are now delayed on windows
- the windows build appears to be unstable:
	- 64-bit doesn't work, period: it crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 CreateWindowExW complains about an unregistered window class, yet the RegisterClassW appears to have succeeded and examining the stack in WinDbg indicates the correct class name is being sent (see below)
	- 32-bit: it works now, but if I save the class name converted to UTF-16 beforehand, wine indicates that the class name is replaced with the window title, so something there is wrong...
- handle in-library panics (internal errors) by reporting them to the user
- david wendt is telling me he's getting frequent crashes on his end with the GTK+ amd64 build...
	- I get soft deadlock if I mash the Click Me button repeatedly
	- occasionally maximizing/restoring a window will abort early and stay that way...?

important things:
- maybe make it so sysData doesn't need specialized info on every control type?
- write an implementation documentation.
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
- this:
	[16:27] <cespare> pietro10: depends what you mean by safe
	[16:27] <cespare> pietro10: sounds like you should move this functionality into a function though.
	[16:28] <cespare> (so the user can decide what to do with the error)
	[16:28] <cespare> typically people don't like their libraries calling exit :)
