package context

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// New returns a context.Context that also
// return gorilla/context values
func New(r *http.Request) context.Context {
	return WithGorilla(NewEmpty(), r)
}

// NewEmpty returns a basic implementation of
// context.Context that has no value at all
func NewEmpty() context.Context {
	return new(emptyContext)
}

var epoch time.Time = time.Unix(0, 0)

// emptyContext
type emptyContext int

// Deadline returns the time when work done on behalf of this context
// should be canceled.  Deadline returns ok==false when no deadline is
// set.  Successive calls to Deadline return the same results.
func (ctx *emptyContext) Deadline() (deadline time.Time, ok bool) {
	return epoch, false
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled.  Done may return nil if this context can
// never be canceled.  Successive calls to Done return the same valu
func (ctx *emptyContext) Done() <-chan struct{} {
	return nil // emptyContext cannot be canceled
}

// Err returns a non-nil error value after Done is closed.  Err returns
// Canceled if the context was canceled or DeadlineExceeded if the
// context's deadline passed.  No other values for Err are defined.
// After Done is closed, successive calls to Err return the same value.
func (ctx *emptyContext) Err() error {
	return nil // since emptyContext cannot be canceled, it has no error
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.  Successive calls to Value with
// the same key returns the same result.
func (ctx *emptyContext) Value(interface{}) interface{} {
	return nil // emptyContext contains no value and will always return nil
}
