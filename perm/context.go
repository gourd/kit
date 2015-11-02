package perm

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Use mux add a perm Mux to a context for later retrieve
func UseMux(m Mux) endpoint.Middleware {
	return func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			ctx = WithMux(ctx, m)
			return inner(ctx, request)
		}
	}
}

// GetMuxOk retrieve permission from current gorilla context
// and a boolean flag. If not found, it returns flase flag.
func GetMuxOk(ctx context.Context) (m Mux, ok bool) {
	// try to get current key
	mi := ctx.Value(contextKey)
	if mi == nil {
		ok = false
		return
	}

	m, ok = mi.(Mux)
	return
}

// GetMux retrieve permission mux from current gorilla context
func GetMux(ctx context.Context) (m Mux) {
	m, _ = GetMuxOk(ctx)
	return
}
