package context

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

type key int

const reqKey key = 0

// WithGorilla wraps a given context.Context with our wrapper context
func WithGorilla(parent context.Context, r *http.Request) context.Context {
	return &wrapper{parent, r}
}

// wrapper is based on gorilla wrapper in
// Golang blog: // https://blog.golang.org/context/gorilla/gorilla.go
type wrapper struct {
	context.Context
	req *http.Request
}

// Value returns Gorilla's context package's value for this Context's request
// and key. It delegates to the parent Context if there is no such value.
func (ctx *wrapper) Value(key interface{}) interface{} {
	if key == reqKey {
		return ctx.req
	}
	if val, ok := gcontext.GetOk(ctx.req, key); ok {
		return val
	}
	return ctx.Context.Value(key)
}

// HTTPRequest returns the *http.Request associated with ctx using NewContext,
// if any.
func HTTPRequest(ctx context.Context) (*http.Request, bool) {
	// We cannot use ctx.(*wrapper).req to get the request because ctx may
	// be a Context derived from a *wrapper. Instead, we use Value to
	// access the request if it is anywhere up the Context tree.
	req, ok := ctx.Value(reqKey).(*http.Request)
	return req, ok
}
