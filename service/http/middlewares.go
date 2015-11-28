package http

import (
	"sort"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

// weight of different middlewares
//
// MWOuter is for middleware
// that runs before any protocol read
// or / and after any protocol output
//
// MWProtocol is for middleware
// that unwrap decoded request from a protocol
// or wraps the results in a proper line protocol
// before encode
//
// MWPrepare is for middleware
// to prepare the request before endpoint
// or prepare the response before protocol
//
// MWNormal is reserved for any
//
// MWInner is for middleware to run between
// normal middleware and the endpoint
//
const (
	MWOuter    = -10
	MWProtocol = -5
	MWPrepare  = -1
	MWNormal   = 0
	MWInner    = 10
)

// Middleware wrap a middleware with a weight parameter
type Middleware struct {
	endpoint.Middleware
	Weight int
}

// Middlewares contain middlewares to use in the service
type Middlewares []Middleware

// Len implements sort.Interface
func (wares Middlewares) Len() int {
	return len(wares)
}

// Less implements sort.Interface
func (wares Middlewares) Less(i, j int) bool {
	return wares[i].Weight < wares[j].Weight
}

// Swap implements sort.Interface
func (wares Middlewares) Swap(i, j int) {
	wares[i], wares[j] = wares[j], wares[i]
}

// Add appends a middleware to the set
func (wares *Middlewares) Add(weight int, waresToAdd ...endpoint.Middleware) {
	for _, ware := range waresToAdd {
		*wares = append(*wares, Middleware{ware, weight})
	}
}

// Slice returns a sorted []endpoint.Middleware
func (wares Middlewares) Slice() (slice []endpoint.Middleware) {
	slice = make([]endpoint.Middleware, 0, wares.Len())
	sort.Sort(wares)
	for _, ware := range wares {
		slice = append(slice, ware.Middleware)
	}
	return
}

// Chain is the helper function for composing middlewares
// in specific order
func (wares Middlewares) Chain() endpoint.Middleware {

	// if no middleware, return an passthrought middleware
	if wares.Len() == 0 {
		return func(inner endpoint.Endpoint) endpoint.Endpoint {
			return func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return inner(ctx, request)
			}
		}
	}

	// sort the middlewares
	slice := wares.Slice()
	return endpoint.Chain(slice[0], slice[1:]...)
}
