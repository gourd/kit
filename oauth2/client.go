//go:generate gourd gen store -type=Client -coll=oauth2_client $GOFILE
package oauth2

// Client implements the osin Client interface
type Client struct {
	Id          string      `db:"id,omitempty" json:"id"`
	Secret      string      `db:"secret" json:"-"`
	RedirectUri string      `db:"redirect_uri" json:"redirect_uri"`
	UserId      string      `db:"user_id" json:"user_id"`
	UserData    interface{} `db:"-" json:"-"`
}

func (c *Client) GetId() string {
	if c == nil {
		return ""
	}
	return c.Id
}

func (c *Client) GetSecret() string {
	if c == nil {
		return ""
	}
	return c.Secret
}

func (c *Client) GetRedirectUri() string {
	if c == nil {
		return ""
	}
	return c.RedirectUri
}

func (c *Client) GetUserData() interface{} {
	if c == nil {
		return nil
	}
	return c.UserData
}
