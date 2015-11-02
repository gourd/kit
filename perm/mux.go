package perm

import (
	"errors"
	"golang.org/x/net/context"
)

var HandlerNotFound error

func init() {
	HandlerNotFound = errors.New("Permission handler not found")
}

// Mux is the primary permission interface
// for user to obtain permission
type Mux interface {
	HandleFunc(perm string, h HandlerFunc)
	Handle(perm string, h Handler)
	Allow(ctx context.Context, perm string, info ...interface{}) error
	//ServeHTTP(w http.ResponseWriter, r *http.Request)
	Default(err error)
}

// NewMux returns a new DefaultMux
func NewMux() Mux {
	m := &DefaultMux{
		handlers: make(map[string]Handler),
	}
	m.Default(HandlerNotFound)
	return m
}

// DefaultMux route permission request to different
// permission Handler. The design mimics ServeMux
// in the core http pacakge
type DefaultMux struct {
	// unexported variables
	handlers map[string]Handler
	defaultH Handler
}

// Handle register a Handler to the DefaultMux.
// This Handler will be provided when calling ServePerm with
// the perm string equals perm.
// The design mimics *ServeMux.Handle
func (p *DefaultMux) Handle(perm string, h Handler) {
	// overwrite by default (until I figure something else)
	p.handlers[perm] = h
}

// HandleFunc register a HandlerFunc to the DefaultMux.
// This Handler will be provided when calling ServePerm with
// the perm string equals perm.
// The design mimics *ServeMux.HandleFunc
func (p *DefaultMux) HandleFunc(perm string, h HandlerFunc) {
	// overwrite by default (until I figure something else)
	p.handlers[perm] = h
}

// Allow dispatches the permission request to the registered
// handlers whose perm string matches / most close to the registered
// Handler
func (p *DefaultMux) Allow(ctx context.Context, perm string, info ...interface{}) error {
	if handler, ok := p.handlers[perm]; ok {
		return handler.Allow(ctx, perm, info...)
	}
	// TODO: find relevant permission string by pattern (i.e. `*`)

	return p.defaultH.Allow(ctx, perm, info...)
}

// Default set the default response of Mux
// if no handler is found. Default HandlerNotFound
// Set to nil if you want to pass by default
func (p *DefaultMux) Default(err error) {
	var h HandlerFunc = func(ctx context.Context, perm string, info ...interface{}) error {
		return err
	}
	p.defaultH = h
}
