package http

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Middlewares contain middlewares to use in the service
type Middlewares struct {
	Outer    []endpoint.Middleware
	Protocol endpoint.Middleware
	Prepare  endpoint.Middleware
	Inner    []endpoint.Middleware
}

// Chain is the helper function for composing middlewares
// in specific order
func (mws Middlewares) Chain() endpoint.Middleware {
	mwares := []endpoint.Middleware{}

	mwares = append(mwares, mws.Outer...)
	if mws.Protocol != nil {
		mwares = append(mwares, mws.Protocol)
	}
	if mws.Prepare != nil {
		mwares = append(mwares, mws.Prepare)
	}
	mwares = append(mwares, mws.Inner...)
	if len(mwares) == 0 {
		return func(inner endpoint.Endpoint) endpoint.Endpoint {
			return func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return inner(ctx, request)
			}
		}
	}

	return endpoint.Chain(mwares[0], mwares[1:]...)
}
