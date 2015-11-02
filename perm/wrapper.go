package perm

import (
	"golang.org/x/net/context"
)

type keyType int

const contextKey keyType = 0

// WithMux adds a mux
func WithMux(parent context.Context, m Mux) context.Context {
	return &wrapper{parent, m}
}

//
type wrapper struct {
	context.Context
	mux Mux
}

// Value returns perm Mux
func (ctx *wrapper) Value(key interface{}) interface{} {
	if key == contextKey {
		return ctx.mux
	}
	return ctx.Context.Value(key)
}
