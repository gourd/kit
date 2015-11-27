package oauth2_test

import (
	"testing"

	"github.com/gourd/kit/oauth2"
	"golang.org/x/net/context"
)

func TestGetAccess_Session(t *testing.T) {

	// test oauth2 server (router only)
	oauth2URL := "/oauth2/dummy"
	authURL := oauth2URL + "/authorize"
	tokenURL := oauth2URL + "/token"
	contentURL := "/content"
	message := "Success"
	oauth2Srvr := testOAuth2Server(oauth2URL, message)

	// test oauth2 client app (router only)
	redirectURL := "/application/redirect"
	password := "password"

	// create dummy oauth client and user
	c, u := createDummies(password, redirectURL)

	// run the code request
	code, err := getCode(oauth2Srvr, getCodeRequest(c, u, password, authURL, redirectURL))
	if err != nil {
		t.Errorf("getCode error (%#v)", err.Error())
		return
	}
	t.Logf("code:  %#v", code)

	// get oauth2 token
	token, err := getToken(oauth2Srvr, getTokenRequest(c, code, tokenURL, redirectURL))
	if err != nil {
		t.Errorf("getToken error (%#v)", err.Error())
		return
	}
	t.Logf("token: %#v", token)

	// test getting AccessData from empty context
	ctx := context.Background()
	if want, have := (*oauth2.AccessData)(nil), oauth2.GetAccess(ctx); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// middleware routine: WithAccess set context with proper token passed
	// test getting AccessData from supposed context with AccessData
	r := getContentRequest(token, contentURL)
	oauth2.NewManager().Middleware().ServeHTTP(nil, r)
	ctx = oauth2.WithAccess(ctx, r)
	access := oauth2.GetAccess(ctx)
	if access == nil {
		t.Errorf("expected *AccessData, got %#v", access)
	}

}
