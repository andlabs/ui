This is a file to keep track of API restrictions that simplify the implementation of the package. I would like to eliminate them, but...

- Once you open a window, the controls are finalized: you cannot change the window's control or add/remove controls to layouts.
- Once you open a window, you cannot change its event channels or its controls's event channels.
- [Windows] At most 65535 controls can be made, period. This is because child window IDs are alloted by the UI library application-global, not window-local, and BN_CLICKED only stores the control ID in a word (and I start counting at 1 to be safe). If I keep the first restriction and amend it such that you can only set the control of a window at the time of first open (somehow; split create and open?), I can easily make them window-local.
