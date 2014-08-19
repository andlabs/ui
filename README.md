Please don't use this package as it stands now. It is being rewritten. You can watch progress in redo/, but keep in mind that it can and will experience major API changes.

Hopefully the rewrite will complete before the end of August.

Note that anyone using this after Go 1.3 will experience intermittent crashes if their allocated objects don't escape to the heap. [Go issue 8310](https://code.google.com/p/go/issues/detail?id=8310) will make that worse as well, but until the Go team makes their proposal public, I don't have much of an alternative.

Note 2: due to missing header files, the Windows version of the rewrite requires [mingw-w64](http://mingw-w64.sourceforge.net/). Make sure your MinGW uses that version instead. If you're running on Windows and not sure what to download, get the mingw-builds distribution.
