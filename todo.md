so I don't forget:
- Window.SizeToFit() or WIndow.OptimalSize() (use: `Window.SetOptimalSize())`) for sizing a window to the control's interest
- Control.Show()/Control.Hide()
- Groupbox
- see if we really need to track errors on a lot of places that report errors
	- it appears GTK+ and Cocoa both either don't provide a convenient way to grab errors or you're not supposed to; I assume you're supposed to just assume everything works... but on Windows we check errors for functions that return errors, and there's no guarantee that only certian errors will be returned...
- character-limited entry fields, numeric entry fields, multiline entry fields
	- possible rename of LineEdit?
		- especially for password fields - NewPasswordEntry()?
- [Windows, Mac OS X] should ListBox have a border style?
	- [Windows] a different border on LineEdits?
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
- change sysData.make() so it does not take the initial window text as an argument and instead have the respective Control/Window.make() call sysData.setText() expressly; this would allow me to remove the "no such concept of text" checks from the GTK+ and Mac OS X backends

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
- there's no GTK+ or Cocoa error handling whatsoever; we need to figure out how it works
	- I know how Cocoa error handling works: it uses Cocoa exceptions; need to figure out how to catch and handle them somehow
- make sure GTK+ documentation version point differences (x in 4.3.x) don't matter
- button sizes and LineEdit sizes on Windows seem too big; Comboboxes have margins
	- Cocoa has similar margining issues (like on Comboboxes)
- sometimes the size of the drop-down part of a Combobox becomes 0 or 1 or some other impossibly small value on Windows
- make gcc (Unix)/clang (Mac OS X) pedantic about warnings/errors; also -Werror
- make sure scrollbars in Listbox work identically on all platforms (specifically the existence and autohiding of both horizontal and vertical scrollbars)
	- pin down this behavior; also note non-editability
- GTK+ windows cannot be resized smaller than their controls's current sizes in their current positions; find out how to overrule that so they can be freely resized

super ultra important things:
- the windows build appears to be unstable:
	- 64-bit crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 it works fine
	- 32-bit: it works, but if I save the class name converted to UTF-16 beforehand, wine indicates that the class name is replaced with the window title, so something there is wrong...
- handle in-library panics (internal errors) by reporting them to the user
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
- the user can still [NSApp terminate:] from the Dock icon, bypassing Go itself
	- ideally we need a QuitItem() function for this case if/when we add menus

important things:
- make specific wording in documentation consistent (make/create, etc.)
	- document minor details like wha thappens on specific events so that they are guaranteed to work the same on all platforms (are there any left?)
		- what happens when the user clicks and drags on a listbox
- make passing of parameters and type conversions of parameters to uitask on Windows consistent: explicit _WPARAM(xxx)/_LPARAM(xxx)/uintptr(xxx), for example
	- do this for type signatures in exported functions: (err error) or just error?
	- do this for the names of GTK+ helper functions (gtkXXX or gXXX)

far off:
- localization
- strip unused constants from the Windows files
- combine more Windows files; rename some?
- tab stops

maybe:
- rename Stack to Box?
- make Combobox and Listbox satisfy sort.Interface?
- indeterminate progress bars (not supported on Windows 2000)
- should a noneditable Combobox be allowed to return to unselected mode by the user?
- since all events are dispatched without blocking uitask, don't bother requiring explicit dispatch? remove ui.Event() and make Window.Closing initialized by default; if we don't listen on the channel, nothing will happen
