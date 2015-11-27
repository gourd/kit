package context

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

// UseGorilla implements go-kit http transport RequestFunc
func UseGorilla(parent context.Context, r *http.Request) context.Context {
	return WithGorilla(parent, r)
}

// ClearGorilla implements go-kit endpoint.Middleware that
// removes all values stored for a given request.
// Works like ClearHandler provided by gorilla
func ClearGorilla(inner endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = inner(ctx, request)
		if r := HTTPRequest(ctx); r != nil {
			gcontext.Clear(r)
		}
		return
	}
}
