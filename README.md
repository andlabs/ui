# Native UI library for Go
# THIS PACKAGE IS UNSTABLE AND PRELIMINARY. In fact, it is presently compiled as package `main` for ease of cross-platform testing and debugging. Once major issues are dealt with and the Mac OS X build working, I will likely move to packge `ui` and move `main()` to a test.

This is a simple library for building cross-platform GUI programs in Go. It targets Windows and all Unix variants (except Mac OS X until further notice) and provides a thread-safe, channel-based API.

There is documentation, but due to the note above, you won't be able to see it just yet. Refer to `main.go` for an example.
