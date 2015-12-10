//go:generate gourd gen store -type=User -coll=user $GOFILE
package oauth2

import (
	"crypto/md5"
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
	Id       string    `db:"id,omitempty" json:"id"`
	Username string    `db:"username" json:"username"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"-"`
	Name     string    `db:"name" json:"name"`
	Created  time.Time `db:"created" json:"created"`
	Updated  time.Time `db:"updated" json:"updated"`
}

// PasswordIs matches the hash with database stored password
func (u *User) PasswordIs(pass string) bool {
	if u.Password == u.Hash(pass) {
		return true
	}
	return false
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
		strID = user.Id
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
