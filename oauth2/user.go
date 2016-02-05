//go:generate gourd gen store -type=User -coll=user $GOFILE
//go:generate gourd gen endpoints -type=User -store=UserStore -storekey=KeyUser user_store.go
//go:generate gourd gen rest -type=User -store=UserStore -storekey=KeyUser user_store.go
package oauth2

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// OAuth2User is the generic user interface
// for OAuth2 login check
type OAuth2User interface {
	// PasswordIs matches a string with the stored password.
	// If the stored password is hash, this function will apply to the
	// input before matching.
	PasswordIs(pass string) bool
}

// User of the API server
type User struct {
	ID       string    `db:"id,omitempty" json:"id"`
	Username string    `db:"username" json:"username"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password,omitempty" json:"-"`
	Name     string    `db:"name" json:"name"`
	MetaJSON string    `db:"meta_json" json:"-"`
	Token    string    `db:"token" json:"-"` // token for lost password request
	Created  time.Time `db:"created" json:"created"`
	Updated  time.Time `db:"updated" json:"updated"`
}

// userJSON is the struct for marshaling and unmarshaling
type userJSON struct {
	ID       string              `json:"id"`
	Username string              `json:"username"`
	Password string              `json:"password"`
	Email    string              `json:"email"`
	Name     string              `json:"name"`
	Meta     map[string][]string `json:"meta"`
	Created  time.Time           `json:"created"`
	Updated  time.Time           `json:"updated"`
}

// Meta read MetaJSON as map[string][]string
func (u User) Meta() (m map[string][]string) {
	m = make(map[string][]string)
	json.Unmarshal([]byte(u.MetaJSON), &m)
	return
}

// AddMeta adds Meta value
func (u *User) AddMeta(key, value string) {
	m := u.Meta()
	if _, ok := m[key]; !ok {
		m[key] = make([]string, 0, 1)
	}
	m[key] = append(m[key], value)

	b, _ := json.Marshal(m)
	u.MetaJSON = string(b)
	return
}

// MarshalJSON implements json.Marshaler
func (u User) MarshalJSON() ([]byte, error) {
	if u.MetaJSON == "" {
		u.MetaJSON = "{}"
	}
	vmap := userJSON{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Meta:     u.Meta(),
		Created:  u.Created,
		Updated:  u.Updated,
	}
	return json.Marshal(vmap)
}

// UnmarshalJSON implements json.Marshaler
func (u *User) UnmarshalJSON(data []byte) (err error) {

	val := userJSON{}
	if err = json.Unmarshal(data, &val); err != nil {
		return
	}

	// set all val to user
	u.ID = val.ID
	u.Username = val.Username
	u.Email = val.Email
	u.Name = val.Name
	u.Created = val.Created
	u.Updated = val.Updated

	// set password, if presents
	if val.Password != "" {
		u.SetPassword(val.Password)
	}

	// set Meta to MetaJSON
	b, err := json.Marshal(val.Meta)
	u.MetaJSON = string(b)

	return
}

// MarshalDB implement
func (u User) MarshalDB() (v interface{}, err error) {
	vmap := make(map[string]interface{})
	vmap["id"] = u.ID
	vmap["username"] = u.Username
	vmap["email"] = u.Email
	vmap["password"] = u.Password
	vmap["name"] = u.Name
	vmap["token"] = u.Token
	vmap["created"] = u.Created
	vmap["updated"] = u.Updated

	if u.MetaJSON == "" {
		vmap["meta_json"] = "{}"
	} else {
		vmap["meta_json"] = u.MetaJSON
	}

	v = vmap
	return
}

// PasswordIs matches the hash with database stored password
func (u *User) PasswordIs(pass string) bool {
	if u.Password == u.Hash(pass) {
		return true
	}
	return false
}

// SetPassword hashes the input and set to password field
func (u *User) SetPassword(pass string) {
	u.Password = u.Hash(pass)
}

// Hash provide the standard hashing for password
func (u *User) Hash(password string) string {
	h := md5.New()
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// UserDataID reads UserData field for AccessData / AuthorizeData
// then retrieve the ID string or return error
func UserDataID(UserData interface{}) (strID string, err error) {

	switch UserData.(type) {
	case *User:
		user := UserData.(*User)
		strID = user.ID
		return
	case map[string]interface{}:
		vmap := UserData.(map[string]interface{})
		if id, ok := vmap["id"]; !ok {
			err = fmt.Errorf(
				`.UserData["id"] not found (.UserData=%#v)`, vmap)
			return
		} else if strID, ok = id.(string); !ok {
			err = fmt.Errorf(
				`.UserData["id"] is not string (%#v)`, vmap)
			return
		}
	}

	err = fmt.Errorf(
		"unexpected .UserData type %#v", UserData)
	return
}
