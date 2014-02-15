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
- scrollbars on listboxes (shouldn't they be automatic? or is that just wine being dumb?)
- [Windows] should ListBox have a border style?
- padding and spacing in Stack; maybe a setting in Stack which keeps controls at their preferred size?
- change Stack/Combobox/Listbox constructors so that there's a separate constructor for each variant, rather than passing in parameters?
- allow Combobox to have initial settings

super ultra important things:
- the windows build appears to be unstable:
	- 64-bit doesn't work, period: it crashes in malloc in wine with heap corruption warnings aplenty during DLL loading; in windows 7 CreateWindowExW complains about an unregistered window class, yet the RegisterClassW appears to have succeeded and examining the stack in WinDbg indicates the correct class name is being sent (see below)
	- 32-bit: it works now, but if I save the class name converted to UTF-16 beforehand, wine indicates that the class name is replaced with the window title, so something there is wrong...

important things:
- maybe make it so sysData doesn't need specialized info on every control type?
- write an implementation documentation.
- Control.preferredSize() (definitely needed for Grid and Form)

far off:
- localization
- strip unused constants from the Windows files
- combine more Windows files; rename some?
- normalize error handling to adorn errors with function call information

maybe:
- rename Stack to Box?
