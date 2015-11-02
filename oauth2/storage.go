package oauth2

import (
	"github.com/RangelReale/osin"
	"github.com/gourd/kit/store"
	"log"
	"net/http"
)

// Storage implements osin.Storage
type Storage struct {
	r      *http.Request
	Client store.Provider
	Auth   store.Provider
	Access store.Provider
	User   store.Provider
}

// SetRequest set the request
func (storage *Storage) SetRequest(r *http.Request) *Storage {
	storage.r = r
	return storage
}

// UseClientFrom set the Client provider
func (storage *Storage) UseClientFrom(p store.Provider) *Storage {
	storage.Client = p
	return storage
}

// UseAuthFrom set the Auth provider
func (storage *Storage) UseAuthFrom(p store.Provider) *Storage {
	storage.Auth = p
	return storage
}

// UseAccessFrom set the Access provider
func (storage *Storage) UseAccessFrom(p store.Provider) *Storage {
	storage.Access = p
	return storage
}

// UseUserFrom set the User provider
func (storage *Storage) UseUserFrom(p store.Provider) *Storage {
	storage.User = p
	return storage
}

// Clone the storage
func (storage *Storage) Clone() (c osin.Storage) {
	c = &Storage{
		Client: storage.Client,
		Auth:   storage.Auth,
		Access: storage.Access,
		User:   storage.User,
	}
	return
}

// Close the connection to the storage
func (storage *Storage) Close() {
	// placeholder now, will revisit when doing mongodb
}

// GetClient loads the client by id (client_id)
func (storage *Storage) GetClient(id string) (c osin.Client, err error) {

	log.Printf("GetClient %s", id)

	srv, err := storage.Client.Store(storage.r)
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

	srv, err := storage.Auth.Store(storage.r)
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
	srv, err := storage.Auth.Store(storage.r)
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

	d = e.ToOsin()
	return
}

// RemoveAuthorize revokes or deletes the authorization code.
func (storage *Storage) RemoveAuthorize(code string) (err error) {

	log.Printf("RemoveAuthorize %s", code)

	srv, err := storage.Auth.Store(storage.r)
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

	srv, err := storage.Access.Store(storage.r)
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

	// store authorize id with access in database
	if ad.AuthorizeData != nil {
		e.AuthorizeCode = ad.AuthorizeData.Code
	}

	// store previous access id with access in database
	if ad.AccessData != nil {
		e.PrevAccessToken = ad.AccessData.AccessToken
	}

	// create in database
	err = srv.Create(store.NewConds(), e)
	log.Printf("SaveAccess last error: %#v", err)
	return
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (storage *Storage) LoadAccess(token string) (d *osin.AccessData, err error) {

	log.Printf("LoadAccess %v", token)

	srv, err := storage.Access.Store(storage.r)
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
	err = func(e *AccessData) (err error) {

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

		// load authdata here
		if e.AuthorizeCode != "" {
			a, err := storage.LoadAuthorize(e.AuthorizeCode)
			if err != nil {
				// ignore "Not Found"
				code, msg := store.ParseError(err)
				if code == 404 {
					log.Printf("Failed to load Auth: %#v. Ignore", msg)
				} else {
					log.Printf("Failed to load Auth: %#v", msg)
					return err
				}
			} else {
				log.Printf("Auth data found")
				ad := &AuthorizeData{}
				if err = ad.ReadOsin(a); err != nil {
					return err
				}
				e.AuthorizeData = ad
			}
		}

		// load previous access here
		if e.PrevAccessToken != "" {
			a, err := storage.LoadAccess(e.PrevAccessToken)
			if err != nil {
				return err
			}
			ad := &AccessData{}
			if err = ad.ReadOsin(a); err != nil {
				return err
			}
			e.AccessData = ad
		}

		return
	}(e)

	if err != nil {
		return
	}

	d = e.ToOsin()
	return
}

// RemoveAccess revokes or deletes an AccessData.
func (storage *Storage) RemoveAccess(token string) (err error) {

	log.Printf("RemoveAccess %v", token)

	srv, err := storage.Access.Store(storage.r)
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

	srv, err := storage.Access.Store(storage.r)
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

	d = e.ToOsin()
	return
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (storage *Storage) RemoveRefresh(token string) (err error) {

	log.Printf("RemoveRefresh %v", token)

	srv, err := storage.Access.Store(storage.r)
	if err != nil {
		return
	}
	defer srv.Close()

	conds := store.NewConds()
	conds.Add("refresh_token", token)
	err = srv.Delete(conds)
	return
}
