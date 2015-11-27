package oauth2

import (
	"log"
	"net/http"

	"golang.org/x/net/context"
)

// WithAccess implements go-kit httptransport RequestFunc
// Adds the current HTTP Request to context.Context
func WithAccess(parent context.Context, r *http.Request) context.Context {
	access, err := GetRequestAccess(r)
	if err != nil {
		log.Printf("WithAccess error (%#v)", err.Error())
		return parent
	}
	log.Printf("WithAccess get access: %#v", access)
	return context.WithValue(parent, accessKey, access)
}

// GetAccess returns oauth2 AccessData stored in session
func GetAccess(ctx context.Context) (d *AccessData) {
	ptr := ctx.Value(accessKey)
	log.Printf("GetAccess(): %#v", ptr)
	if ptr == nil {
		return
	}
	d = ptr.(*AccessData)
	return
}
