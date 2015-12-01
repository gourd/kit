package oauth2

import (
	"log"
	"net/http"

	gcontext "github.com/gorilla/context"
	gourdctx "github.com/gourd/kit/context"
	"github.com/gourd/kit/store"
)

type key int

const (
	storageKey key = iota
	accessKey  key = iota
)

// Middleware is a generic middleware
// to serve a Storage instance to
type Middleware struct {
	storage *Storage
}

// ServeHTTP implements http.Handler interface method.
// Attach a clone of the storage to context
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := gourdctx.GetRequestID(r)
	log.Printf("[%s] Clone storage into context", id)
	sc := m.storage.Clone()
	s := sc.(*Storage)
	s.SetRequest(r)
	gcontext.Set(r, storageKey, s)
}

// GetStorageOk returns oauth2 storage in context and a boolean flag.
// If process failed, boolean flag will be false
func GetStorageOk(r *http.Request) (s *Storage, ok bool) {
	raw := gcontext.Get(r, storageKey)
	s, ok = raw.(*Storage)
	return
}

// GetStorage returns oauth2 storage in context
// or nil if failed
func GetStorage(r *http.Request) *Storage {
	s, _ := GetStorageOk(r)
	return s
}

// GetRequestAccess returns oauth2 AccessData with token
// found in "Authority" header variable of the HTTP Request
func GetRequestAccess(r *http.Request) (d *AccessData, err error) {
	token := r.Header.Get("Authority")
	if token == "" {
		return // nothing
	}
	return GetTokenAccess(r, token)
}

// GetTokenAccess retrieves oauth2 AccessData of
// provided token
func GetTokenAccess(r *http.Request, token string) (d *AccessData, err error) {

	// retrieve context oauth2 storage
	s, ok := GetStorageOk(r)
	if !ok {
		log.Printf("Failed to retrieve storage from context")
		err = store.ErrorInternal
		return
	}

	sessid := gourdctx.GetRequestID(r)

	// get access by token
	od, err := s.LoadAccess(token)
	if err != nil {
		if err.Error() == "Not Found" {
			return nil, nil
		}

		log.Printf("[%s] Token: %s", sessid, token)
		log.Printf("[%s] Failed to load access: %s", sessid, err.Error())
		return
	}
	d = &AccessData{}
	d.ReadOsin(od)
	return
}
