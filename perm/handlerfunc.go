package perm

import (
	"golang.org/x/net/context"
)

// HandlerFunc is the basic unit for permission management
// It takes
type HandlerFunc func(ctx context.Context, perm string, info ...interface{}) error

// Allow calls f(r, pern, info...).
func (h HandlerFunc) Allow(ctx context.Context, perm string, info ...interface{}) error {
	return h(ctx, perm, info...)
}
