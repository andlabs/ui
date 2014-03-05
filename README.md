[![Build Status](https://travis-ci.org/tompao/ui.png?branch=travis-ci)](https://travis-ci.org/tompao/ui)

# Native UI library for Go
### THIS PACKAGE IS UNDER ACTIVE DEVELOPMENT. Feel free to start using it, but mind: it's far from feature-complete, it's still in need of testing and crash-fixing, and the API can (and will) change. If you can help, please do! Run `./test` to build a test binary `test/test` which runs a (mostly) feature-complete UI test. Run `./d32 ./test` to build a 32-bit version (you will need a cgo-enabled 32-bit go environment, and I have only tested this on Mac OS X).

This is a simple library for building cross-platform GUI programs in Go. It targets Windows, Mac OS X, Linux, and other Unixes, and provides a thread-safe, channel-based API. The API itself is minimal; it aims to provide only what is necessary for GUI program design. That being said, suggestions are welcome. Layout is done using various layout managers, and some effort is taken to conform to the target platform's UI guidelines. Otherwise, the library uses native toolkits.

ui aims to run on all supported versions of supported platforms. To be more precise, the system requirements are:

* Windows: Windows 2000 or newer. The Windows backend uses package `syscall` and calls Windows DLLs directly, so does not rely on cgo.
* Mac OS X: Mac OS X 10.6 (Snow Leopard) or newer. Objective-C dispatch is done by interfacing with libobjc directly, and thus this uses cgo.
* Other Unixes: The Unix backend uses GTK+, and thus cgo. It requires GTK+ 3.4 or newer; for Ubuntu this means 12.04 LTS (Precise Pangolin) at minimum. Check your distribution.

ui itself has no outside Go package dependencies; it is entirely self-contained.

To install, simply `go get` this package. On Mac OS X, make sure you have the Apple development headers. On other Unixes, make sure you have the GTK+ development files (for Ubuntu, `libgtk-3-dev` is sufficient).

Package documentation is available at http://godoc.org/github.com/andlabs/ui.

The following is an example program to illustrate what programming with ui is like:
```go
(see test/main.go)
```

## Contributing
Contributions are welcome. File issues, pull requests, approach me on IRC (pietro10 in #go-nuts; andlabs elsewhere), etc. Even suggestions are welcome: while I'm mainly drawing from my own GUI programming experience, everyone is different.

If you want to dive in, read implementation.md: this is a description of how the library works. (Feel free to suggest improvements to this as well.) The other .md files in this repository contain various development notes.

Please suggest documentation improvements as well.
