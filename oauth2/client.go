//go:generate gourd gen store -type=Client -coll=oauth2_client $GOFILE
package oauth2

// Client implements the osin Client interface
type Client struct {
	ID          string      `db:"id,omitempty" json:"id"`
	Secret      string      `db:"secret" json:"-"`
	RedirectURI string      `db:"redirect_uri" json:"redirect_uri"`
	UserID      string      `db:"user_id" json:"user_id"`
	UserData    interface{} `db:"-" json:"-"`
}

func (c *Client) GetId() string {
	if c == nil {
		return ""
	}
	return c.ID
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
	return c.RedirectURI
}

func (c *Client) GetUserData() interface{} {
	if c == nil {
		return nil
	}
	return c.UserData
}
