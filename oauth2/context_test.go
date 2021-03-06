package oauth2_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gourd/kit/oauth2"
	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
)

// getCode use httptest.NewRecorder to interact with
// http.Handler of server to get auth code
func getCode(oauth2Srvr http.Handler, r *http.Request) (code string, err error) {

	// run request here
	w := httptest.NewRecorder()
	oauth2Srvr.ServeHTTP(w, r)
	w.Flush()

	// test status code
	if want, have := http.StatusFound, w.Code; want != have {
		err = fmt.Errorf("status expected %#v, got %#v", want, have)
		return
	}

	// read the location
	location := w.Header().Get("Location")
	if location == "" {
		bodyMsg := ""

		// try reading body
		b, bodyErr := ioutil.ReadAll(w.Body)
		if bodyErr != nil {
			bodyMsg = fmt.Sprintf("error reading body (%#v)", err.Error())
		} else {
			bodyMsg = fmt.Sprintf("body:   %#v", string(b))
		}

		// more details
		err = fmt.Errorf("no location found\n%s\n%s\n%s",
			fmt.Sprintf("status: %#v\n", w.Code),
			fmt.Sprintf("header: %#v\n", w.HeaderMap),
			bodyMsg)
		return
	}

	locURL, err := url.Parse(location)
	if err != nil {
		err = fmt.Errorf("error parsing location (%#v)", err.Error())
		return
	}
	code = locURL.Query().Get("code")
	if code == "" {
		err = errors.New("code not found")
	}
	return
}

// getToken use httptest.NewRecorder to interact with
// http.Handler of server to get token
func getToken(oauth2Srvr http.Handler, r *http.Request) (token string, err error) {

	// run request here
	w := httptest.NewRecorder()
	oauth2Srvr.ServeHTTP(w, r)
	w.Flush()

	// test status code
	if want, have := http.StatusOK, w.Code; want != have {
		err = fmt.Errorf("status expected %#v, got %#v", want, have)
		return
	}

	// read token from token endpoint response (json)
	bodyDecoded := make(map[string]string)
	body := w.Body.String()
	dec := json.NewDecoder(strings.NewReader(body))
	if err := dec.Decode(&bodyDecoded); err != nil {
		err = fmt.Errorf(
			"Unable to parse token response: %s\nbody: %#v",
			err.Error(), body)
	}

	var ok bool
	if token, ok = bodyDecoded["access_token"]; !ok {
		err = fmt.Errorf(
			"Unable to find access_token, body: %#v", body)
	} else if token == "" {
		err = errors.New("code not found")
	}
	return
}

// getContent use httptest.NewRecorder to interact with
// http.Handler of server to get content endpoint body
func getContent(srvr http.Handler, r *http.Request) (body string, err error) {

	// run request here
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)
	w.Flush()

	// test status code
	if want, have := http.StatusOK, w.Code; want != have {
		err = fmt.Errorf("status expected %#v, got %#v", want, have)
		return
	}

	// read body and convert to string
	raw, err := ioutil.ReadAll(w.Body)
	if err != nil {
		err = fmt.Errorf("Failed to read body: %s", err.Error())
		return
	}
	body = string(raw)
	return
}

func TestGetAccess_Session(t *testing.T) {

	var err error

	// test oauth2 server (router only)
	testCtx := &testContext{
		password:     "password",
		t:            t,
		redirectBase: "https://test.foobar/example_app/",
		redirectURL:  "https://test.foobar/example_app/code",
		oauth2Path:   "/oauth2/dummy",
	}

	message := "Success"
	oauth2Srvr := testOAuth2Server(t, testCtx.oauth2Path, message)
	contentURL := "/content"

	// test oauth2 client app (router only)
	//redirectURL := "/application/redirect"

	// test store context
	type tempKey int
	const (
		testDB tempKey = iota
	)
	factory := store.NewFactory()
	factory.SetSource(testDB, defaultTestSrc())
	factory.Set(oauth2.KeyUser, testDB, oauth2.UserStoreProvider)
	factory.Set(oauth2.KeyClient, testDB, oauth2.ClientStoreProvider)
	factory.Set(oauth2.KeyAccess, testDB, oauth2.AccessDataStoreProvider)
	factory.Set(oauth2.KeyAuth, testDB, oauth2.AuthorizeDataStoreProvider)
	ctx := store.WithFactory(context.Background(), factory)
	defer store.CloseAllIn(ctx)

	// create dummy oauth client and user
	testCtx.client, testCtx.user = createStoreDummies(ctx,
		testCtx.password, testCtx.redirectBase)

	// run the code request
	testCtx.code, err = getCode(oauth2Srvr, getCodeRequest(testCtx))
	if err != nil {
		t.Errorf("getCode error (%#v)", err.Error())
		return
	}
	t.Logf("code:  %#v", testCtx.code)

	// get oauth2 token
	testCtx.token, err = getToken(oauth2Srvr, getTokenRequest(testCtx))
	if err != nil {
		t.Errorf("getToken error (%s)", err.Error())
		return
	}
	t.Logf("token: %#v", testCtx.token)

	if want, have := (*oauth2.AccessData)(nil), oauth2.GetAccess(ctx); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// middleware routine: WithAccess set context with proper token passed
	// test getting AccessData from supposed context with AccessData
	r := getContentRequest(testCtx.token, contentURL)
	ctx = oauth2.LoadTokenAccess(oauth2.UseToken(ctx, r))
	access := oauth2.GetAccess(ctx)
	if access == nil {
		t.Errorf("expected *AccessData, got %#v", access)
		return
	}

	if want, have := "", access.ID; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if access.ClientID == "" {
		t.Errorf("access.ClientId expected to be not empty")
	}
	if want, have := testCtx.token, access.AccessToken; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if want, have := testCtx.user.ID, access.UserID; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if access.UserData == nil {
		t.Error("expect access.UserData not nil")
	} else if want, have := testCtx.user.ID, access.UserData.(*oauth2.User).ID; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if access.RefreshToken == "" {
		t.Errorf("access.RefreshToken expected to be not empty")
	}

}

/*
func TestGetRequestAccess(t *testing.T) {

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

	// get content endpoint
	body, err := getContent(oauth2Srvr, getContentRequest(token, contentURL))
	if err != nil {
		t.Errorf("getContent error (%#v)", err.Error())
		return
	}
	t.Logf("body:  %#v", body)

	// examine message
	if want, have := message, body; want != have {
		t.Errorf("expect: %#v, got: %#v", want, have)
	}

	// test getting access data from store,
	// emulating server environment
	access, err := getRequestAccess(token, contentURL)
	if err != nil {
		switch err.(type) {
		case *store.StoreError:
			serr := err.(*store.StoreError)
			t.Errorf("oauth2.GetRequestAccess StoreError (%#v)", serr.ServerMsg)
		default:
			t.Errorf("oauth2.GetRequest	Access error (%#v)", err.Error())
		}
	}

	if want, have := "", access.Id; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if access.ClientId == "" {
		t.Errorf("access.ClientId expected to be not empty")
	}
	if want, have := token, access.AccessToken; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if want, have := u.Id, access.UserId; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if access.UserData == nil {
		t.Error("expect access.UserData not nil")
	} else if want, have := u.Id, access.UserData.(*oauth2.User).Id; want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
	if access.RefreshToken == "" {
		t.Errorf("access.RefreshToken expected to be not empty")
	}

	t.Logf("access user: %#v", access.UserData)

}
*/

// TODO: test the refresh token routine
// TODO: implement and do the information endpoint
