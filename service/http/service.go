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
