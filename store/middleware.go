package store

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// RequestFuncWithFactory returns a go-kit/kit/transpoirt/http.RequestFunc
// that will add a given factory to the context
func RequestFuncWithFactory(factory Factory) httptransport.RequestFunc {
	return func(parent context.Context, r *http.Request) context.Context {
		return WithFactory(parent, factory)
	}
}

// ClearFactory cleans up after RequestFunc.
// If you called RequestFunc in your http transport,
// you should use this
func ClearFactory(inner endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = inner(ctx, request)
		CloseAllIn(ctx)
		return
	}
}

// Middleware takes a Factory and create a middleware that
// does WithStores and CloseAllIn
func Middleware(factory Factory) endpoint.Middleware {
	mware := func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			ctx = WithFactory(ctx, factory)
			response, err = inner(ctx, request)
			CloseAllIn(ctx)
			return
		}
	}
	return mware
}
