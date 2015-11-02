package store

import (
	"fmt"
	"net/http"
)

// Provider is the default registry of providers
var Providers *providerDefs

// initialize the registry of providers
func init() {
	Providers = &providerDefs{
		defs: make(map[string]Provider),
	}
}

// ProviderDefs contains provider instances indexed by name
type providerDefs struct {
	defs map[string]Provider
}

// Define provider by name
func (d *providerDefs) Define(name string, p Provider) {
	d.defs[name] = p
}

// Define provider (with ProvideFunc) by name
func (d *providerDefs) DefineFunc(name string, f ProvideFunc) {
	d.defs[name] = f
}

// Get provider by name and error
func (d *providerDefs) Get(name string) (p Provider, err error) {
	var ok bool
	if p, ok = d.defs[name]; !ok {
		err = fmt.Errorf("Provider \"%s\" doesn't exists", name)
	}
	return
}

// MustGet retrieves provider by name or panic
func (d *providerDefs) MustGet(name string) (p Provider) {
	p, err := d.Get(name)
	if err != nil {
		panic(err.Error())
	}
	return
}

// Store provide named service from registered provider
// or return error if failed
func (d *providerDefs) Store(r *http.Request, name string) (s Store, err error) {
	p, err := d.Get(name)
	if err != nil {
		return
	}
	s, err = p.Store(r)
	return
}

// MustStore provide named service from registered provider
// or panic if failed
func (d *providerDefs) MustStore(r *http.Request, name string) (s Store) {
	s, err := d.MustGet(name).Store(r)
	if err != nil {
		panic(err)
	}
	return
}

// Provider defines service provider for web servers
type Provider interface {
	Store(r *http.Request) (Store, error)
}

// ProvideFunc is to simplify implementation of Provider
type ProvideFunc func(r *http.Request) (Store, error)

// Store implements Provider interface
func (f ProvideFunc) Store(r *http.Request) (Store, error) {
	return f(r)
}
