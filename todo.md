so I don't forget:
- Window.SizeToFit() or WIndow.OptimalSize() (use: `Window.SetSize(Window.OptimalSize())`) for sizing a window to the control's interest
- Control.Show()/Control.Hide()

important things:
- maybe make it so I don't need to expose Window.sysData to controls? I need a way to get the window HWND for the Windows one...
- maybe make it so sysData doesn't need specialized info on every control type?
- write an implementation documentation.

far off:
- localization
- strip unused constants from the Windows files
- combine more Windows files; rename some?
