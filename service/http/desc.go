package http

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Desc is the descriptor of a RESTful service
type Desc interface {
	Paths() Paths
	SetMiddleware(name string, mware endpoint.Middleware)
	GetMiddleware(name string) endpoint.Middleware
	SetDecodeFunc(name string, dec httptransport.DecodeRequestFunc)
	GetDecodeFunc(name string) httptransport.DecodeRequestFunc
}

// desc is the default implementation of Desc interface
type desc struct {
	paths       Paths
	middlewares map[string]endpoint.Middleware
	decodeFuncs map[string]httptransport.DecodeRequestFunc
}

// Path implements Desc interface
func (d desc) Paths() Paths {
	return d.paths
}

// SetMethod implements Desc inteface
func (d *desc) SetMiddleware(name string, mware endpoint.Middleware) {
	d.middlewares[name] = mware
}

// GetMethod implements Desc interface
func (d desc) GetMiddleware(name string) endpoint.Middleware {
	if mware, ok := d.middlewares[name]; ok {
		return mware
	}
	return nil
}

// SetDecodeFunc implements Desc interface
func (d desc) SetDecodeFunc(name string, dec httptransport.DecodeRequestFunc) {
	d.decodeFuncs[name] = dec
}

// GetDecodeFunc implements Desc interface
func (d desc) GetDecodeFunc(name string) httptransport.DecodeRequestFunc {
	if dec, ok := d.decodeFuncs[name]; ok {
		return dec
	}
	return nil
}

// NewDesc returns the default implementation of Desc
func NewDesc(paths Paths) Desc {
	return &desc{
		paths:       paths,
		middlewares: make(map[string]endpoint.Middleware),
		decodeFuncs: make(map[string]httptransport.DecodeRequestFunc),
	}
}
