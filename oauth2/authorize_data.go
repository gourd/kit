//go:generate gourd gen store -type=AuthorizeData -coll=oauth2_auth $GOFILE
package oauth2

import (
	"fmt"
	"time"

	"github.com/RangelReale/osin"
)

// AuthorizeData interfacing database to osin storage I/O of same name
type AuthorizeData struct {

	// Authorize Data Id
	Id string `db:"id,omitempty" json:"id,omitempty"`

	// Client Id the data is linked to
	ClientId string `db:"client_id" json:"client_id"`

	// Client information
	Client *Client `db:"-" json:"-"`

	// Authorization code
	Code string `db:"code" json:"code"`

	// Token expiration in seconds
	ExpiresIn int32 `db:"expires_in" json:"expires_in"`

	// Requested scope
	Scope string `db:"scope" json:"scope"`

	// Redirect Uri from request
	RedirectUri string `db:"redirect_uri" json:"redirect_uri"`

	// State data from request
	State string `db:"state" json:"state"`

	// Date created
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// User Id the data is linked to
	UserId string `db:"user_id" json:"user_id"`

	// Data to be passed to storage. Not used by the osin library.
	UserData interface{} `db:"-"`
}

// ToOsin returns an osin version of the struct of osin I/O
func (d *AuthorizeData) ToOsin() (od *osin.AuthorizeData) {
	od = &osin.AuthorizeData{}
	od.Client = d.Client
	od.Code = d.Code
	od.ExpiresIn = d.ExpiresIn
	od.Scope = d.Scope
	od.RedirectUri = d.RedirectUri
	od.State = d.State
	od.CreatedAt = d.CreatedAt
	od.UserData = d.UserData
	return
}

// ReadOsin reads a *osin.AuthorizeData, takes its value
// then set to itself
func (d *AuthorizeData) ReadOsin(od *osin.AuthorizeData) (err error) {
	var ok bool
	if od.Client == nil {
		// skip for now
	} else if d.Client, ok = od.Client.(*Client); !ok {
		err = fmt.Errorf("osin client is not of Client type: %#v", od)
		return
	}
	d.ClientId = od.Client.GetId()
	d.Code = od.Code
	d.ExpiresIn = od.ExpiresIn
	d.Scope = od.Scope
	d.RedirectUri = od.RedirectUri
	d.State = od.State
	d.CreatedAt = od.CreatedAt
	d.UserData = od.UserData
	if d.UserData != nil {
		if d.UserId, err = UserDataID(d.UserData); err != nil {
			return
		}
	}
	return nil
}
