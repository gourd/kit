package oauth2

import (
	"log"
	"net/http"

	"golang.org/x/net/context"
)

type privateKey int

const (
	accessKey privateKey = iota
)

// ReadTokenAccess reads token information from header ("Authority")
// and, if AccessData found for the given token, add to context
func ReadTokenAccess(ctx context.Context, r *http.Request) context.Context {

	storage := &Storage{}
	storage.SetContext(ctx)

	token := r.Header.Get("Authority")
	if token == "" {
		return ctx
	}

	osinAccess, err := storage.LoadAccess(token)
	if err != nil {
		log.Printf("UseTokenAccess failed to load access. token=%#v err=%#v",
			token, err.Error())
		return ctx
	}

	ad := &AccessData{}
	ad.ReadOsin(osinAccess)
	return WithAccess(ctx, ad)
}

// WithAccess implements go-kit httptransport RequestFunc
// Adds the current HTTP Request to context.Context
func WithAccess(parent context.Context, ad *AccessData) context.Context {
	return context.WithValue(parent, accessKey, ad)
}

// GetAccess returns oauth2 AccessData stored in session
func GetAccess(ctx context.Context) (d *AccessData) {
	ptr := ctx.Value(accessKey)
	log.Printf("GetAccess(): %#v", ptr)
	if ptr == nil {
		return
	}
	log.Printf("ptr = %#v", ptr)
	if ad, ok := ptr.(*AccessData); ok {
		d = ad
	}
	return
}
