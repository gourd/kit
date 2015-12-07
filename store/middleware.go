package store

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

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
