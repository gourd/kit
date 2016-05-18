package httpservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	gourdctx "github.com/gourd/kit/context"
	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
)

// jsonEncodeFunc encodes given response into JSON
func jsonEncodeFunc(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(response)
	return
}

// jsonErrorEncoder expands given error to StoreError then encode to JSON
func jsonErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")

	// quick fix for gokit bad request wrapping problem
	switch err.(type) {
	case httptransport.Error:
		err = err.(httptransport.Error).Err
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
		Weight:      0,
		Context:     gourdctx.NewEmpty(),
		Endpoint:    ep,
		Middlewares: &Middlewares{},
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
	Weight      int
	Context     context.Context
	Middlewares *Middlewares
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

// ServiceSlice attaches the method of sort.Interface to []*Service, sort in
// increasing order
type ServiceSlice []*Service

// Len implements sort.Interface
func (ss ServiceSlice) Len() int {
	return len(ss)
}

// Less implements sort.Interface
func (ss ServiceSlice) Less(i, j int) bool {
	return ss[i].Weight < ss[j].Weight
}

// Swap implements sort.Interface
func (ss ServiceSlice) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

// Sort is a short hand for sort.Sort(ServiceSlice)
func (ss ServiceSlice) Sort() {
	sort.Sort(ss)
}

// Services contain a group of named services
type Services map[string]*Service

// Patch takes a patch / patches and apply them to the group
func (services Services) Patch(patches ...ServicesPatch) {
	for _, patch := range patches {
		services = patch(services)
	}
}

// RouterFunc generalize router to route an http.Handler
type RouterFunc func(path string, methods []string, h http.Handler) error

// Each returns a channel that return services by weight order
func (services Services) Each() <-chan *Service {
	out := make(chan *Service)

	// turn services into a slices, sort it
	slice := ServiceSlice(make([]*Service, 0, len(services)))
	for _, s := range services {
		slice = append(slice, s)
	}
	sort.Sort(slice)

	// return *Service in the slice through output channel
	go func(out chan *Service, slice ServiceSlice) {
		defer close(out)
		for _, s := range slice {
			out <- s
		}
	}(out, slice)
	return out
}

// Route routes all services in the group
func (services Services) Route(rtr RouterFunc) (err error) {
	for service := range services.Each() {
		if err = service.Route(rtr); err != nil {
			err = fmt.Errorf("error routing %#v (method: %#v) (%#v)",
				service.Path, service.Methods, err.Error())
			return
		}
	}
	return
}

// ServicesPatch patches all children in a map[string]*Service
type ServicesPatch func(Services) Services
