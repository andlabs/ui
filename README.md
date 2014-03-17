[![Build Status](https://travis-ci.org/andlabs/ui.png?branch=master)](https://travis-ci.org/andlabs/ui)
# Native UI library for Go
### CRITICAL UPDATE 17 March 2014: Due to deadlocks between resizing and setting label text once a second, Labels and Areas do not lock their internal mutex locks after their Window has been created. This seems to not have issues, and the Go race detector isn't saying anything, so IDK... Please help with this if you can. I will have to remove the mutexes from the other classes (later) to make this work properly.

### THIS PACKAGE IS UNDER ACTIVE DEVELOPMENT. It can be used; the API is stable enough at this point, but keep in mind there may still be crashes and API changes, as suggestions are always open. If you can help, please do! Run `./test` to build a test binary `test/test` which runs a (mostly) feature-complete UI test. Run `./d32 ./test` to build a 32-bit version (you will need a cgo-enabled 32-bit go environment, and I have only tested this on Mac OS X).

### UPDATE 12 March 2014: Windows 2000 is no longer supported [as it is no longer supported by Go](https://codereview.appspot.com/74790043).

This is a simple library for building cross-platform GUI programs in Go. It targets Windows, Mac OS X, Linux, and other Unixes, and provides a thread-safe, channel-based API. The API itself is minimal; it aims to provide only what is necessary for GUI program design. That being said, suggestions are welcome. Layout is done using various layout managers, and some effort is taken to conform to the target platform's UI guidelines. Otherwise, the library uses native toolkits.

ui aims to run on all supported versions of supported platforms. To be more precise, the system requirements are:

* Windows: Windows XP or newer. The Windows backend uses package `syscall` and calls Windows DLLs directly, so does not rely on cgo.
* Mac OS X: Mac OS X 10.6 (Snow Leopard) or newer. Objective-C dispatch is done by interfacing with libobjc directly, and thus this uses cgo.
	* Note: you will need Go 1.3 or newer (so until it is released, [go tip](http://tip.golang.org/doc/install/source#head)) for this verison, as it uses a single .m file due to technical restrictions (read the comments in `bleh_darwin.m` for details), and earlier versions of Go do not auto-build .m files.
* Other Unixes: The Unix backend uses GTK+, and thus cgo. It requires GTK+ 3.4 or newer; for Ubuntu this means 12.04 LTS (Precise Pangolin) at minimum. Check your distribution.

ui itself has no outside Go package dependencies; it is entirely self-contained.

To install, simply `go get` this package. On Mac OS X, make sure you have the Apple development headers. On other Unixes, make sure you have the GTK+ development files (for Ubuntu, `libgtk-3-dev` is sufficient).

Package documentation is available at http://godoc.org/github.com/andlabs/ui.

For an example of how ui is used, see https://github.com/andlabs/wakeup, which is a small program that implements a basic alarm clock.

## Known To Have Ever Been Built Matrices
For convenience's sake, here are matrices of builds that I have personally done at least once. Each cell represents the run status. These matrices represent builds that I have done at any point in development; it is not a guarantee that the current version works. (I built this list to answer questions of whether or not ui works with a specific configuration.) Only configurations marked with a * are tested during active development. "(invalid)" means the given OS/arch combination is not supported by Go.

   | 386 | amd64 | arm
----- | ----- | ----- | -----
windows | works on windows; works on wine* | works on windows; fails on wine | (invalid)
linux | see table below | see table below | Raspian: works
darwin (Mac OS X) | works* (cross-compiled from 64-bit) | works* | (invalid)
dragonfly | untested | untested | (invalid)
freebsd | untested (VM failure) | untested (VM failure) | untested
netbsd | untested | untested | untested
openbsd | untested | untested | (invalid)
solaris | (invalid) | Oracle Solaris 11: **GTK+ 3 not available from official repos** | (invalid)
plan9 | (not written yet; problems building Go) | (not written) | (invalid)
nacl | (not sure how to handle) | (not sure how to handle) | (invalid)

linux | 386 | amd64
----- | ----- | ----- | -----
Kubuntu (14.04) | works; cross-compiling on 64-bit Linux fails due to nonexistent .so symlinks | works*
Fedora | untested | untested
openSUSE | untested | untested
Arch Linux | untested | untested
Mandriva (TODO choose between PCLinuxOS and Mageia - it appears the original Mandriva is either dead or nonfree and I would rather choose the fork that structures packages identically for parity; do they both?) | untested | untested
Slackware | untested | untested
Gentoo | untested | untested

(The above list should cover all the bases of major Linux distributions and variants thereof; I might add a dedicated Debian test later but other than that... suggestions welcome. Kubuntu 64-bit is my main system and the main development platform; the Windows builds are cross-compiled from here. And yes, this also implies I seriously consider a Plan 9 port of the library using [libcontrol](http://plan9.bell-labs.com/magic/man2html/2/control), though I'm guessing this will blow up in my face due to any possible conflicts between libthread and Go's runtime (I need to see how the Go runtime implements OS threads on Plan 9).)

## Contributing
Contributions are welcome. File issues, pull requests, approach me on IRC (pietro10 in #go-nuts; andlabs elsewhere), etc. Even suggestions are welcome: while I'm mainly drawing from my own GUI programming experience, everyone is different. I have received emails, however I am not likely to see those right away, so I don't suggest contacting me by email if your communication is urgent.

If you want to dive in, read implementation.md: this is a description of how the library works. (Feel free to suggest improvements to this as well.) The other .md files in this repository contain various development notes.

Please suggest documentation improvements as well.
