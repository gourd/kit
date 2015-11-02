package perm

import (
	"golang.org/x/net/context"
)

// Handler handle permission requests
type Handler interface {

	// Allow returns nil if permission is granted
	// or return an error if permission is denied
	Allow(ctx context.Context, perm string, info ...interface{}) error
}
