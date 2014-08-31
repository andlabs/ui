# ui: platform-native GUI library for Go

This is a library that aims to provide simple GUI software development in Go. It runs on/requires:

* Windows: cgo, mingw-w64, Windows XP and newer
* Mac OS X: cgo, Mac OS X 10.7 and newer
* other Unixes: cgo, GTK+ 3.4 and newer

Go 1.3 is required. Note that vanilla 1.3 has a bug in Mac OS X cgo; the next release will fix it.

(this README needs some work)

prevlib.tar contains the previous version of the library as it stood when I restarted; don't bother using it.

# Screenshots
The example widget gallery on GTK+ in the Adwaita theme (3.13/master):

!(https://raw.githubusercontent.com/andlabs/ui/master/examples/widgetgallery/widgetgallery.png)
