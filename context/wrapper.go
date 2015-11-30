package context

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

type key int

const (
	reqKey key = iota
	idKey  key = iota
)

// WithHTTPRequest adds the current HTTP Request to context.Context
func WithHTTPRequest(parent context.Context, r *http.Request) context.Context {
	return context.WithValue(parent, reqKey, r)
}

// HTTPRequest returns the *http.Request associated with ctx using NewContext,
// if any.
func HTTPRequest(ctx context.Context) *http.Request {
	reqItf := ctx.Value(reqKey)
	if reqItf == nil {
		return nil // if not found, return nil
	}

	// return stored request
	return reqItf.(*http.Request)
}

// WithGorilla wraps a given context.Context with our wrapper context.
// It also runs WithHTTPRequest inside.
func WithGorilla(parent context.Context, r *http.Request) context.Context {
	return &wrapper{WithHTTPRequest(parent, r)}
}

// wrapper is based on gorilla wrapper in
// Golang blog: // https://blog.golang.org/context/gorilla/gorilla.go
type wrapper struct {
	context.Context
}

// Value returns Gorilla's context package's value for this Context's request
// and key. It delegates to the parent Context if there is no such value.
func (ctx *wrapper) Value(key interface{}) interface{} {
	if key == reqKey {
		// do nothing, fall to Context.Value
	} else if val, ok := gcontext.GetOk(HTTPRequest(ctx.Context), key); ok {
		return val
	}
	return ctx.Context.Value(key)
}

// WithID add a string ID to the context (for session tracking)
func WithID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, idKey, id)
}

// GetID get the string ID from request
func GetID(ctx context.Context) string {
	if v := ctx.Value(idKey); v != nil {
		return v.(string)
	}
	return ""
}
