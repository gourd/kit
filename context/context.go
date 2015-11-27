package context

import (
	"net/http"

	"golang.org/x/net/context"
)

// New returns a context.Context that also
// return gorilla/context values
func New(r *http.Request) context.Context {
	return WithGorilla(context.Background(), r)
}

// NewEmpty returns a basic implementation of
// context.Context that has no value at all
func NewEmpty() context.Context {
	return context.Background()
}
