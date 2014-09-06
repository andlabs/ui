# ui: platform-native GUI library for Go

## Feature requests wanted! (Really; IDK what to add next!)

This is a library that aims to provide simple GUI software development in Go. It runs on/requires:

* Windows: cgo, mingw-w64, Windows XP and newer
* Mac OS X: cgo, Mac OS X 10.7 and newer
* other Unixes: cgo, GTK+ 3.4 and newer

Go 1.3 is required. Note that vanilla 1.3 has a bug in Mac OS X cgo; the next release will fix it.

(this README needs some work)

Be sure to have at least each outermost Window escaping to the heap until a good resolution to Go issue 8310 comes out.

prevlib.tar contains the previous version of the library as it stood when I restarted; don't bother using it.

# Updates

**3 September 2014**<br>The new GtkGrid-style Grid is now implemented! See its documentation for more details. Also, debugging spew has been removed.

**31 August 2014**<br>Grid is now renamed SimpleGrid in preparation for implementing a more [GtkGrid](https://developer.gnome.org/gtk3/unstable/GtkGrid.html)-like Grid. Mind the change.

# Screenshots
The example widget gallery on GTK+ in the Adwaita theme (3.13/master):

![widget gallery example](https://raw.githubusercontent.com/andlabs/ui/master/examples/widgetgallery/widgetgallery.png)
