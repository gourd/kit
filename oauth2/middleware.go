package oauth2

import (
	"log"
	"net/http"

	gcontext "github.com/gorilla/context"
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
	log.Printf("Clone storage into context")
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

	// get access by token
	od, err := s.LoadAccess(token)
	if err != nil {
		log.Printf("Token: %s", token)
		log.Printf("Failed to load access: %s", err.Error())
		return
	}
	d = &AccessData{}
	d.ReadOsin(od)
	return
}
