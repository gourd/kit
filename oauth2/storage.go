package oauth2

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/RangelReale/osin"
	"github.com/gourd/kit/store"
)

type storeKey int

// Keys for Storage to access different stores
// from provided context
const (
	KeyClient storeKey = iota
	KeyAuth
	KeyAccess
	KeyUser
)

// Storage implements osin.Storage
type Storage struct {
	ctx context.Context
}

// SetContext set the context of the storage clone
func (storage *Storage) SetContext(ctx context.Context) *Storage {
	storage.ctx = ctx
	return storage
}

// Clone implements osin.Storage.Clone
func (storage *Storage) Clone() (c osin.Storage) {
	c = &Storage{}
	return
}

// Close implements osin.Storage.Close
func (storage *Storage) Close() {
	// Close is not functional at All.
	// Should use store.CloseAllIn to wrap up
	// database connections.
}

// GetClient implements osin.Storage.GetClient
func (storage *Storage) GetClient(id string) (c osin.Client, err error) {

	log.Printf("GetClient %s", id)

	srv, err := store.Get(storage.ctx, KeyClient)
	if err != nil {
		log.Printf("Unable to get client store")
		return
	}
	defer srv.Close()

	e := &Client{}
	conds := store.NewConds()
	conds.Add("id", id)

	err = srv.One(conds, e)
	if err != nil {
		log.Printf("%#v", conds)
		log.Printf("Failed running One()")
		return
	} else if e == nil {
		log.Printf("Client not found for the id %#v", id)
		err = store.Error(http.StatusNotFound,
			"Client not found for the given id")
		return
	}

	c = e
	return
}

// SaveAuthorize saves authorize data.
func (storage *Storage) SaveAuthorize(d *osin.AuthorizeData) (err error) {

	log.Printf("SaveAuthorize %v", d)

	srv, err := store.Get(storage.ctx, KeyAuth)
	if err != nil {
		return
	}
	defer srv.Close()

	e := &AuthorizeData{}
	err = e.ReadOsin(d)
	if err != nil {
		return
	}

	// store client id with auth in database
	e.ClientId = e.Client.GetId()

	// create the auth data now
	err = srv.Create(store.NewConds(), e)
	return
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (storage *Storage) LoadAuthorize(code string) (d *osin.AuthorizeData, err error) {

	log.Printf("LoadAuthorize %s", code)

	// loading osin using osin storage
	srv, err := store.Get(storage.ctx, KeyAuth)
	if err != nil {
		return
	}
	defer srv.Close()

	e := &AuthorizeData{}
	conds := store.NewConds()
	conds.Add("code", code)

	err = srv.One(conds, e)
	if err != nil {
		return
	} else if e == nil {
		err = store.Error(http.StatusNotFound,
			"AuthorizeData not found for the code")
		return
	}

	// load client here
	var ok bool
	cli, err := storage.GetClient(e.ClientId)
	if err != nil {
		return
	} else if e.Client, ok = cli.(*Client); !ok {
		err = store.Error(http.StatusInternalServerError,
			"Internal Server Error")
		log.Printf("Unable to cast client into Client type: %#v", cli)
		return
	}

	// load user data here
	if e.UserId != "" {
		userStore, err := store.Get(storage.ctx, KeyUser)
		if err != nil {
			return d, err
		}
		user := &User{}
		userStore.One(store.NewConds().Add("id", e.UserId), user)
		e.UserData = user
	}

	d = e.ToOsin()
	return
}

// RemoveAuthorize revokes or deletes the authorization code.
func (storage *Storage) RemoveAuthorize(code string) (err error) {

	log.Printf("RemoveAuthorize %s", code)

	srv, err := store.Get(storage.ctx, KeyAuth)
	if err != nil {
		return
	}
	defer srv.Close()

	conds := store.NewConds()
	conds.Add("code", code)
	err = srv.Delete(conds)
	return
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (storage *Storage) SaveAccess(ad *osin.AccessData) (err error) {

	log.Printf("SaveAccess %#v", ad)

	srv, err := store.Get(storage.ctx, KeyAccess)
	if err != nil {
		return
	}
	defer srv.Close()

	// generate database access type
	e := &AccessData{}
	err = e.ReadOsin(ad)
	if err != nil {
		return
	}

	// store client id with access in database
	e.ClientId = e.Client.GetId()

	// if AuthorizeData is set, store as JSON
	if ad.AuthorizeData != nil {
		var b []byte
		authData := &AuthorizeData{}
		if err = authData.ReadOsin(ad.AuthorizeData); err != nil {
			return
		}
		if b, err = json.Marshal(authData); err != nil {
			return
		}
		e.AuthorizeDataJSON = string(b)
	}

	// if AccessData is set, store as JSON
	if ad.AccessData != nil {
		var b []byte
		accessData := &AccessData{}
		if err = accessData.ReadOsin(ad.AccessData); err != nil {
			return
		}
		if accessData.AccessData != nil {
			// forget data of too long ago
			accessData.AccessData = nil
		}
		if b, err = json.Marshal(accessData); err != nil {
			return
		}
		e.AccessDataJSON = string(b)
	}

	// create in database
	if err = srv.Create(store.NewConds(), e); err != nil {
		log.Printf("SaveAccess error: %#v", err.Error())
	}
	return
}

// loadAccessSupp loads supplementary data onto an *AccessData
func (storage *Storage) loadAccessSupp(e *AccessData) (err error) {

	// load client here
	var ok bool
	cli, err := storage.GetClient(e.ClientId)
	if err != nil {
		return
	} else if e.Client, ok = cli.(*Client); !ok {
		err = store.Error(http.StatusInternalServerError,
			"Internal Server Error")
		log.Printf("Unable to cast client into Client type: %#v", cli)
		return
	}
	e.ClientId = e.Client.GetId()

	// unserialize previous AuthorizeData here
	if e.AuthorizeDataJSON != "" {
		ad := &AuthorizeData{}
		json.Unmarshal([]byte(e.AuthorizeDataJSON), ad)
		e.AuthorizeData = ad
	}

	// unserialize previous AccessData here
	if e.AccessDataJSON != "" {
		ad := &AccessData{}
		json.Unmarshal([]byte(e.AccessDataJSON), ad)
		e.AccessData = ad
	}

	// load user data here
	if e.UserId != "" {
		userStore, err := store.Get(storage.ctx, KeyUser)
		if err != nil {
			return err
		}
		user := &User{}
		userStore.One(store.NewConds().Add("id", e.UserId), user)
		e.UserData = user
	}

	return

}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (storage *Storage) LoadAccess(token string) (d *osin.AccessData, err error) {

	log.Printf("LoadAccess %v", token)

	srv, err := store.Get(storage.ctx, KeyAccess)
	if err != nil {
		return
	}
	defer srv.Close()

	e := &AccessData{}
	conds := store.NewConds()
	conds.Add("access_token", token)

	err = srv.One(conds, e)
	if err != nil {
		return
	} else if e == nil {
		err = store.Error(http.StatusNotFound,
			"AccessData not found for the token")
		return
	}

	// load supplementary data
	if err = storage.loadAccessSupp(e); err != nil {
		return
	}

	d = e.ToOsin()
	return
}

// RemoveAccess revokes or deletes an AccessData.
func (storage *Storage) RemoveAccess(token string) (err error) {

	log.Printf("RemoveAccess %v", token)

	srv, err := store.Get(storage.ctx, KeyAccess)
	if err != nil {
		return
	}
	defer srv.Close()

	conds := store.NewConds()
	conds.Add("access_token", token)
	err = srv.Delete(conds)
	return
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (storage *Storage) LoadRefresh(token string) (d *osin.AccessData, err error) {

	log.Printf("LoadRefresh %v", token)

	srv, err := store.Get(storage.ctx, KeyAccess)
	if err != nil {
		return
	}
	defer srv.Close()

	e := &AccessData{}
	conds := store.NewConds()
	conds.Add("refresh_token", token)

	err = srv.One(conds, e)
	if err != nil {
		return
	} else if e == nil {
		err = store.Error(http.StatusNotFound,
			"AccessData not found for the refresh token")
		return
	}

	// load supplementary data
	if err = storage.loadAccessSupp(e); err != nil {
		return
	}

	d = e.ToOsin()
	return
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (storage *Storage) RemoveRefresh(token string) (err error) {

	log.Printf("RemoveRefresh %v", token)

	srv, err := store.Get(storage.ctx, KeyAccess)
	if err != nil {
		return
	}
	defer srv.Close()

	conds := store.NewConds()
	conds.Add("refresh_token", token)
	err = srv.Delete(conds)
	return
}
