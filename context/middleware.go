package context

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	gcontext "github.com/gorilla/context"
	"github.com/satori/go.uuid"
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

// generate a new string ID with UUID
func newID() string {
	uid := uuid.NewV4()
	return strings.TrimRight(base64.URLEncoding.EncodeToString(uid[:]), "=")
}

// UseID add a string id to http request header and context
func UseID(parent context.Context, r *http.Request) context.Context {
	if prevID := r.Header.Get("X-GOURD-ID"); prevID == "" {
		id := newID()
		r.Header.Set("X-GOURD-ID", id)
		return WithID(parent, id)
	}
	return parent
}

// GetRequestID get string id from http request
func GetRequestID(r *http.Request) string {
	return r.Header.Get("X-GOURD-ID")
}
