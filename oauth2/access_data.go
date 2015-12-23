//go:generate gourd gen store -type=AccessData -coll=oauth2_access $GOFILE
package oauth2

import (
	"fmt"
	"time"

	"github.com/RangelReale/osin"
)

// AccessData interfacing database to osin storage I/O of same name
type AccessData struct {

	// ID is the primary key of AccessData
	ID string `db:"id,omitempty" json:"id"`

	// ClientId is the client which this AccessData is linked to
	ClientID string `db:"client_id" json:"client_id"`

	// Client information
	Client *Client `db:"-" json:"-"`

	// Authorize data, for authorization code
	AuthorizeData *AuthorizeData `db:"-" json:"-"`

	// Authorize data, for authorization code
	AuthorizeDataJSON string `db:"auth_data_json,omitempty" json:"-"`

	// Previous access data, for refresh token
	AccessData *AccessData `db:"-" json:"-"`

	// AccessDataJSON stores the previous access data in JSON string
	AccessDataJSON string `db:"access_data_json,omitempty" json:"-"`

	// Access token
	AccessToken string `db:"access_token" json:"access_token"`

	// Refresh Token. Can be blank
	RefreshToken string `db:"refresh_token" json:"refresh_token"`

	// Token expiration in seconds
	ExpiresIn int32 `db:"expires_in" json:"expires_in"`

	// Requested scope
	Scope string `db:"scope" json:"scope"`

	// RedirectUri from request
	RedirectURI string `db:"redirect_uri" json:"redirect_uri"`

	// Date created
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// User Id the data is linked to
	UserID string `db:"user_id" json:"user_id"`

	// Data to be passed to storage. Not used by the osin library.
	UserData interface{} `db:"-"`
}

// ToOsin returns an osin version of the struct of osin I/O
func (d *AccessData) ToOsin() (od *osin.AccessData) {
	od = &osin.AccessData{}
	od.Client = d.Client
	od.AccessToken = d.AccessToken
	od.RefreshToken = d.RefreshToken
	od.ExpiresIn = d.ExpiresIn
	od.Scope = d.Scope
	od.RedirectUri = d.RedirectURI
	od.CreatedAt = d.CreatedAt
	od.UserData = d.UserData

	// TODO: do we need to do json.Unmarshal here for
	//       AuthorizeData and AccessData?? need to find out

	// indirect parameters
	if d.AuthorizeData != nil {
		od.AuthorizeData = d.AuthorizeData.ToOsin()
	}
	if d.AccessData != nil {
		od.AccessData = d.AccessData.ToOsin()
	}
	if d.UserData != nil {
		od.UserData = d.UserData
	}
	return
}

// ReadOsin reads an osin's AccessData into the AccessData instance
func (d *AccessData) ReadOsin(od *osin.AccessData) (err error) {

	// read parameters that could be directly read
	d.AccessToken = od.AccessToken
	d.RefreshToken = od.RefreshToken
	d.ExpiresIn = od.ExpiresIn
	d.Scope = od.Scope
	d.RedirectURI = od.RedirectUri
	d.CreatedAt = od.CreatedAt
	d.UserData = od.UserData

	// read indirect parameters
	if od.Client != nil {
		if c, ok := od.Client.(*Client); ok {
			d.Client = c
			d.ClientID = c.GetId()
		} else {
			err = fmt.Errorf("Failed to read client from osin.AccessData (%#v)", od.Client)
			return
		}
	}
	if od.AuthorizeData != nil {
		// read the AuthorizeData and store as JSON
		oaud := &AuthorizeData{}
		oaud.ReadOsin(od.AuthorizeData)
		d.AuthorizeData = oaud

		// remember the user_id
		d.UserID = oaud.UserID
	}
	if od.AccessData != nil {
		if *od == *od.AccessData {
			err = fmt.Errorf(".AccessData referencing itself")
			return
		}
		oacd := &AccessData{}
		oacd.ReadOsin(od.AccessData)
	}
	if od.UserData != nil {
		if d.UserID, err = UserDataID(od.UserData); err != nil {
			return
		}
	}

	return
}

// Scopes read the scope field into Scopes type
func (d *AccessData) Scopes() *Scopes {
	return ReadScopes(d.Scope)
}
