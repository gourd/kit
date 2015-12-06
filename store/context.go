package store

import (
	"fmt"

	"golang.org/x/net/context"
)

type keys int

// context keys
const (
	storesKey keys = iota
	DefaultSrc
)

// WithStores create a Stores and add to the context
func WithStores(parent context.Context, defs Defs) context.Context {

	return context.WithValue(parent,
		storesKey, &stores{defs, make(map[interface{}]Conn)})
}

// Get try to connect to a store with provided source
// and provider definition. If fail, return nil and error
func Get(ctx context.Context,
	key interface{}) (s Store, err error) {

	v := ctx.Value(storesKey)
	if v == nil {
		err = fmt.Errorf("Stores not in context")
		return
	}

	stores := v.(Stores)
	s, err = stores.Get(key)
	return
}

// CloseAllIn close all Store connections in the context
func CloseAllIn(ctx context.Context) {

	v := ctx.Value(storesKey)
	if v == nil {
		return
	}

	stores := v.(Stores)
	stores.Close()
}

// Stores is an interface for store
// with connection pool management.
//
// Each HTTP request should have its
// own Stores instance in the context
type Stores interface {

	// Connect connects a provider at a source
	// and return connection and, if any, connection error
	Get(key interface{}) (s Store, err error)

	// Close close all Conn in the set
	Close()
}

type stores struct {
	defs  Defs
	conns map[interface{}]Conn
}

// Connect connects gets a connection to the key
func (sts *stores) Get(key interface{}) (s Store, err error) {

	// find provider
	srcKey, provider := sts.defs.Get(key)
	if srcKey == nil && provider == nil {
		err = fmt.Errorf("Store provider not found")
		return
	}

	// find existing connection
	var conn Conn
	var ok bool
	if conn, ok = sts.conns[srcKey]; !ok {
		source := sts.defs.GetSource(srcKey)
		conn, err = source()
	}
	if err != nil {
		return
	}

	sts.conns[srcKey] = conn
	s, err = provider(conn.Raw())
	return
}

// Close close all the Conn in the set
func (sts *stores) Close() {
	for _, conn := range sts.conns {
		conn.Close()
	}
}
