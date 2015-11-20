package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	gourdctx "github.com/gourd/kit/context"
	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
)

// jsonEncodeFunc encodes given response into JSON
func jsonEncodeFunc(w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(response)
	return
}

// jsonErrorEncoder expands given error to StoreError then encode to JSON
func jsonErrorEncoder(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")

	// quick fix for gokit bad request wrapping problem
	switch err.(type) {
	case httptransport.BadRequestError:
		err = err.(httptransport.BadRequestError).Err
	}

	serr := store.ExpandError(err)
	log.Printf("error: %#v", serr.ServerMsg)
	json.NewEncoder(w).Encode(serr)
}

// NewJSONService creates a service descriptor
// with defaults for a simple JSON service
func NewJSONService(path string, ep endpoint.Endpoint) *Service {
	return &Service{
		Path:        path,
		Methods:     []string{"GET"},
		Context:     gourdctx.NewEmpty(),
		Endpoint:    ep,
		Middlewares: Middlewares{[]endpoint.Middleware{}, nil, nil, []endpoint.Middleware{}},
		EncodeFunc:  jsonEncodeFunc,
		Options: []httptransport.ServerOption{
			httptransport.ServerBefore(gourdctx.UseGorilla),
			httptransport.ServerErrorEncoder(jsonErrorEncoder),
		},
	}
}

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

// Service contains all parameters needed to call
// httptransport.NewServer
type Service struct {
	Path        string
	Methods     []string
	Context     context.Context
	Middlewares Middlewares
	Endpoint    endpoint.Endpoint
	DecodeFunc  httptransport.DecodeRequestFunc
	EncodeFunc  httptransport.EncodeResponseFunc
	Options     []httptransport.ServerOption
}

// Handler returns go-kit http transport server
// of the given definition
func (s Service) Handler() http.Handler {
	ep := s.Middlewares.Chain()(s.Endpoint)
	return httptransport.NewServer(
		s.Context,
		ep,
		s.DecodeFunc,
		s.EncodeFunc,
		s.Options...)
}

// Route add the given service to router with RouterFunc
func (s Service) Route(rtr RouterFunc) error {
	return rtr(s.Path, s.Methods, s.Handler())
}

// RouterFunc generalize router to route an http.Handler
type RouterFunc func(path string, methods []string, h http.Handler) error
