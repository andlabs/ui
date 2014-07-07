// 6 july 2014

package ui

// Doer is a channel that takes Requests returned by the various functions and methods of package ui.
// There are two main Doers: ui.Do, which is for outside event handlers, and the Doer passed into an event handler function.
// You should not create or use your own Doers; these are meaningless.
type Doer chan *Request

// Do is the main way to issue requests to package ui.
// Requests returned by the various functions and methods of package ui should be sent across Do to have them performed.
// When an event is dispatched to an event handler, that event handler will receive a new Doer which is active for the life of the event handler, and any requests made to Do will block until the event handler returns.
var Do = make(Doer)

// Request represents a request issued to the package.
// These are returned by the various functions and methods of package ui and are sent to either Do or to the currently active event handler channel.
// There are also several convenience functions that perfrom common operations with requests.
type Request struct {
	op		func()
	resp		chan interface{}
}

// Response returns a channel which is pulsed exactly once, then immeidately closed, with the response from the function that issued the request.
// If the function does not return a value, this value will be the zero value of struct{}.
// Otherwise, the type of this value depends on the function that created the Request.
func (r *Request) Response() <-chan interface{} {
	return r.resp
}

// Wait is a convenience function that performs a Request and waits for that request to be processed.
// If the request returns a value, it is discarded.
// You should generally use Wait on functions that do not return a value.
// See the documentation of Bool for an example.
func Wait(c Doer, r *Request) {
	c <- r
	<-r.resp
}

// Bool is a convenience function that performs a Request that returns a bool, waits for that request to be processed, and returns the result.
// For example:
// 	if ui.Bool(ui.Do, checkbox.Checked()) { /* do stuff */ }
func Bool(c Doer, r *Request) bool {
	c <- r
	return (<-r.resp).(bool)
}

// Int is like Bool, but for int.
func Int(c Doer, r *Request) int {
	c <- r
	return (<-r.resp).(int)
}

// String is like Bool, but for string.
func String(c Doer, r *Request) string {
	c <- r
	return (<-r.resp).(string)
}

// IntSlice is like Bool, but for []int.
func IntSlice(c Doer, r *Request) []int {
	c <- r
	return (<-r.resp).([]int)
}

// StringSlice is like Bool, but for []string.
func StringSlice(c Doer, r *Request) []string {
	c <- r
	return (<-r.resp).([]string)
}
