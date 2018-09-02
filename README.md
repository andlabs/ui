# ui: platform-native GUI library for Go

This is a library that aims to provide simple GUI software development in Go. It is based on my [libui](https://github.com/andlabs/libui), a simple cross-platform library that does the same thing, but written in C.

It runs on/requires:

- Windows: cgo, Windows Vista SP2 with Platform Update and newer
- Mac OS X: cgo, Mac OS X 10.8 and newer
- other Unixes: cgo, GTK+ 3.10 and newer
	- Debian, Ubuntu, etc.: `sudo apt-get install libgtk-3-dev`
	- Red Hat/Fedora, etc.: `sudo dnf install gtk3-devel`

It also requires Go 1.8 or newer.

It currently aligns to libui's Alpha 4.1, with only a small handful of functions not available.

# Status

Package ui is currently **mid-alpha** software. Much of what is currently present runs stabily enough for the examples and perhaps some small programs to work, but the stability is still a work-in-progress, much of what is already there is not feature-complete, some of it will be buggy on certain platforms, and there's a lot of stuff missing. The libui README has more information.

# Installation

Once you have the dependencies installed, a simple

```
go get github.com/andlabs/ui/...
```

should suffice.

# Documentation

The in-code documentation is sufficient to get started, but needs improvement.

Some simple example programs are in the `examples` directory. You can `go build` each of them individually.

## Windows manifests

Package ui requires a manifest that specifies Common Controls v6 to run on Windows. It should at least also state as supported Windows Vista and Windows 7, though to avoid surprises with other packages (or with Go itself; see [this issue](https://github.com/golang/go/issues/17835)) you should state compatibility with higher versions of Windows too.

The simplest option is provided as a subpackage `winmanifest`; you can simply import it without a name, and it'll set things up properly:

```go
import _ "github.com/andlabs/ui/winmanifest"
```

You do not have to worry about importing this in non-Windows-only files; it does nothing on non-Windows platforms.

If you wish to use your own manifest instead, you can use the one in `winmanifest` as a template to see what's required and how. You'll need to specify the template in a `.rc` file and use `windres` in MinGW-w64 to generate a `.syso` file as follows:

```
windres -i resources.rc -o winmanifest_windows_GOARCH.syso -O coff
```

You may also be interested in the `github.com/akavel/rsrc` and `github.com/josephspurrier/goversioninfo` packages, which provide other Go-like options for embedding the manifest.

Note that if you choose to ship a manifest as a separate `.exe.manifest` file instead of embedding it in your binary, and you use Cygwin or MSYS2 as the source of your MinGW-w64, Cygwin and MSYS2 instruct gcc to embed a default manifest of its own if none is specified. **This default will override your manifest file!** See [this issue](https://github.com/Alexpux/MSYS2-packages/issues/454) for more details, including workaround instructions.

## macOS program execution

If you run a macOS program binary directly from the command line, it will start in the background. This is intentional; see [this](https://github.com/andlabs/libui#why-does-my-program-start-in-the-background-on-os-x-if-i-run-from-the-command-line) for more details.
