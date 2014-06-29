In the new setup, Windows have WindowHandlers. A WindowHandler is defined as

``` go
type WindowHandler interface {
	Event(e Event, c interface{})
}
```

Whenever an event to the window or any control within the window comes in, the handler's `Event()` method is called with the event type and a context-specific value (usually the control that triggered the event) as an argument.

``` go
type Event int
const (
	Close Event = iota
	Clicked
	Checked
	Selected
	Dismissed		// for dialogs
	// ...
	CustomEvent = 5000		// arbitrary but high enough
)
```

The argument to `Close` is a pointer to a value that determnes whether to continue the closing of the window or not. The semantics of this value (type, possible values, and default; a special case in Cocoa means there could be three possible values) have yet to be defined.

The argument to all events `e` such that `Clicked` < `e` < `CustomEvent` is a pointer to the Control (or Dialog; see below) that triggered the event.

`CustomEvent` represents the first free ID that can be used by the program for whatever it wants, as a substitute for channels. The argument type is program-defened. To trigger a custom event, use the `Window.Send(e, data)` method. `Send()` panics if the event requested is not custom.

As an example, the timer from `wakeup` might be run on a goroutine:

``` go
func (w *MainWin) timerGoroutine() {
	for {
		select {
		case t := <-w.start:
			// set the timer up
		case <-w.timerChan:
			w.win.Send(CustomEvent, nil)
		case <-w.stop:
			// stop the timer
		}
	}
}
```

The underlying OS event handler is not existed until the event handling function returns.

With the exception of `Window.Create()`, `Window.Open()`, and `Window.Send()`, no objects and methods are safe for concurrent use anymore. They can only be used within an event handler. They can be used within `AreaHandler` methods as well as from the `WindowHandler` method.

`ui.Go()` no longer takes any arguments. Instead, when initiailization completes, it sends /and waits for the receipt of/ a semaphore value across the `ui.Started` channel, which is immeidately closed after first receipt. Programs should use this flag to know when it is safe to call `Window.Create()`, `Window.Open()`, and `Window.Send()`. A send of a semaphore value to `ui.Stop` will tell `ui.Go()` to return. This return is immediate; there is no opportunity for cleanup.

The semantics of dialogs will also need changing. It may be (I'm not sure yet) no longer possible to have "application-modal" dialogs. The standard dialog box methods on Window will still exist, but instead of returning a Control, they will return a new type Dialog which can be defined as

``` go
type Dialog interface {
	Result() Result
	Selection() interface{}	// string for file dialogs; some other type for other dialogs
	// TODO might contain hidden or unexported fields to prevent creating something that's compatible with Dialog but cannot be used as one for the sake of custom Dialogs; see below
	// TODO make it compatible with Control?
}
```

When the dialog is dismissed, a `Dismissed` event will be raised with that dialog as an argument; get the result code by calling `Result()`.

It might still be possible to have dialog boxes that do not return until the user takes an action and returns the result of that action. I do not know how these will work yet, or what names will be used for either type.

The Dialog specification above would still allow custom dialogs to be made. In fact, they could be built on top of Window perhaps (or even as a mode of Window), but they would need to be reusable somehow...
