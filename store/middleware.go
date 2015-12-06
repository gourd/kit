package store

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Middleware takes a Defs and create a middleware that
// does WithStores and CloseAllIn
func Middleware(defs Defs) endpoint.Middleware {
	mware := func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			ctx = WithStores(ctx, defs)
			response, err = inner(ctx, request)
			CloseAllIn(ctx)
			return
		}
	}
	return mware
}
