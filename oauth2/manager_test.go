package oauth2_test

import (
	"errors"
	"path"

	"golang.org/x/net/context"

	"github.com/gourd/kit/oauth2"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/pat"
	"github.com/gourd/kit/store"
)

// testRedirectErr helps capture the redirect URL
// in RedirectFunc
type testRedirectErr struct {
	msg      string
	redirect *url.URL
}

func (err testRedirectErr) Error() string {
	return err.msg
}

func (err testRedirectErr) Redirect() *url.URL {
	return err.redirect
}

// testNoRedirect implements RedirectFunc in *http.Request.
// It stops redirection and return an error containing
// the redirect URL
func testNoRedirect(req *http.Request, via []*http.Request) error {
	return testRedirectErr{"no redirect", req.URL}
}

// test the testRedirectErr type
func TestRedirectErr(t *testing.T) {
	redirect := &url.URL{}
	var err error = testRedirectErr{"hello", redirect}
	switch err.(type) {
	case testRedirectErr:
		// do nothing
	default:
		t.Errorf("type switch cannot identify the error raw type")
		return
	}
	if want, have := "hello", err.Error(); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
	if want, have := redirect, err.(testRedirectErr).Redirect(); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
}

// creates dummy client and user directly from the stores
func createDummies(password, redirect string) (*oauth2.Client, *oauth2.User) {

	type tempKey int
	const (
		testDB tempKey = iota
	)

	// define test db
	factory := store.NewFactory()
	factory.SetSource(testDB, defaultTestSrc())
	factory.Set(oauth2.KeyUser, testDB, oauth2.UserStoreProvider)
	factory.Set(oauth2.KeyClient, testDB, oauth2.ClientStoreProvider)
	factory.Set(oauth2.KeyAccess, testDB, oauth2.AccessDataStoreProvider)
	factory.Set(oauth2.KeyAuth, testDB, oauth2.AuthorizeDataStoreProvider)
	ctx := store.WithFactory(context.Background(), factory)
	defer store.CloseAllIn(ctx)

	// generate dummy user
	us, err := store.Get(ctx, oauth2.KeyUser)
	if err != nil {
		panic(err)
	}
	u := dummyNewUser(password)
	err = us.Create(store.NewConds(), u)
	if err != nil {
		panic(err)
	}

	// get related dummy client
	cs, err := store.Get(ctx, oauth2.KeyClient)
	if err != nil {
		panic(err)
	}
	c := dummyNewClient(redirect)
	c.UserId = u.Id
	err = cs.Create(store.NewConds(), c)
	if err != nil {
		panic(err)
	}

	return c, u
}

// getCodeRequest generates the http request to get code
func getCodeRequest(c *oauth2.Client, u *oauth2.User, password, authURL, redirect string) *http.Request {

	// login form request
	form := url.Values{}
	form.Add("user_id", u.Username)
	form.Add("password", password)
	log.Printf("form send: %s", form.Encode())

	// build the query string
	q := &url.Values{}
	q.Add("response_type", "code")
	q.Add("client_id", c.GetId())
	q.Add("redirect_uri", redirect)

	req, err := http.NewRequest("POST",
		authURL+"?"+q.Encode(),
		strings.NewReader(form.Encode()))
	if err != nil {
		panic(err) // not quite possible
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

// getCodeHTTP runs the getCodeRequest with actual
// HTTP client and parse the result as code, error
func getCodeHTTP(req *http.Request) (code string, err error) {

	log.Printf("Test retrieving code ====")

	// new http client to emulate user request
	hc := &http.Client{
		CheckRedirect: testNoRedirect,
	}

	_, rerr := hc.Do(req)
	if rerr == nil {
		err = errors.New("unexpected nil error, ecpecting testRedirectErr")
		return
	} else if _, ok := rerr.(*url.Error); !ok {
		err = fmt.Errorf("unexpected response error %#v", rerr)
		return
	}

	// examine error
	uerr := rerr.(*url.Error).Err
	switch uerr.(type) {
	case nil:
		err = errors.New("unexpected nil error, ecpecting testRedirectErr")
	case testRedirectErr:
		// do nothing
	default:
		err = fmt.Errorf("Failed run the request ??: %s", rerr.Error())
		return
	}

	// directly extract the code from the redirect url
	code = uerr.(testRedirectErr).Redirect().Query().Get("code")
	log.Printf("code: %#v", code)

	return
}

// getTokenRequest generates request which client app
// send to oauth2 server for the token
func getTokenRequest(c *oauth2.Client, code, tokenURL, redirect string) *http.Request {
	// build user request to token endpoint
	form := &url.Values{}
	form.Add("code", code)
	form.Add("client_id", c.GetId())
	form.Add("client_secret", c.Secret)
	form.Add("grant_type", "authorization_code")
	form.Add("redirect_uri", redirect)
	req, err := http.NewRequest("POST",
		tokenURL,
		strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	return req
}

// getTokenHTTP runs the getTokenRequest with actual
// HTTP client and parse the result as token, error
func getTokenHTTP(req *http.Request) (token string, err error) {

	log.Printf("Test retrieving token ====")

	// new http client to emulate user request
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		err = fmt.Errorf("Failed run the request: %s", err.Error())
	}

	// read token from token endpoint response (json)
	bodyDecoded := make(map[string]string)
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&bodyDecoded)

	log.Printf("Response Body: %#v", bodyDecoded)
	var ok bool
	if token, ok = bodyDecoded["access_token"]; !ok {
		err = fmt.Errorf(
			"Unable to parse access_token: %s", err.Error())
	}
	return
}

// getContentRequest generates a request to content endpoint
// of the OAuth2 server / OAuth2 guarded resource server
func getContentRequest(token, contentURL string) *http.Request {
	req, err := http.NewRequest("GET", contentURL, nil)
	if err != nil {
		panic(err)
	}

	// additional information
	req.Header.Add("Authority", token)
	return req
}

// getContentHTTP runs the getContentRequest with actual
// HTTP client and parse the result as body, error
func getContentHTTP(r *http.Request) (body string, err error) {
	// new http client to emulate user request
	hc := &http.Client{}
	resp, err := hc.Do(r)
	if err != nil {
		err = fmt.Errorf("Failed run the request: %s", err.Error())
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Failed to read body: %s", err.Error())
		return
	}

	body = string(raw)
	return
}

// example server web app
func testOAuth2Server(baseURL, msg string) http.Handler {

	rtr := pat.New()

	// oauth2 manager
	m := oauth2.NewManager()

	type tempKey int
	const (
		testDB tempKey = iota
	)

	// define store factory for storage
	factory := store.NewFactory()
	factory.SetSource(testDB, defaultTestSrc())
	factory.Set(oauth2.KeyUser, testDB, oauth2.UserStoreProvider)
	factory.Set(oauth2.KeyClient, testDB, oauth2.ClientStoreProvider)
	factory.Set(oauth2.KeyAccess, testDB, oauth2.AccessDataStoreProvider)
	factory.Set(oauth2.KeyAuth, testDB, oauth2.AuthorizeDataStoreProvider)

	// add oauth2 endpoints to router
	// ServeEndpoints bind OAuth2 endpoints to a given base path
	// Note: this is router specific and need to be generated somehow
	oauth2.RoutePat(rtr, baseURL, m.GetEndpoints(factory))

	// add a route the requires access
	rtr.Get("/content", func(w http.ResponseWriter, r *http.Request) {

		ctx := store.WithFactory(context.Background(), factory)
		ctx = oauth2.ReadTokenAccess(ctx, r)
		log.Printf("Dummy content page accessed")

		// obtain access
		a := oauth2.GetAccess(ctx)
		if a == nil {
			fmt.Fprint(w, "Unable to gain Access")
			return
		}

		// no news is good news
		fmt.Fprint(w, msg)
	})

	return rtr
}

// example client web app in the login
func testAppServer(path string) *pat.Router {
	rtr := pat.New()

	log.Printf("testAppServer(%#v)", path)

	// add dummy client reception of redirection
	rtr.Get(path, func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		enc := json.NewEncoder(w)
		enc.Encode(map[string]string{
			"code":  r.Form.Get("code"),
			"token": r.Form.Get("token"),
		})
	})

	return rtr
}

// TestOAuth2HTTP tests the stack with
// actual HTTP call against httptest.Server
// wrapped handlers
func TestOAuth2HTTP(t *testing.T) {

	// a dummy password for dummy user
	password := "password"
	message := "Success"

	// create test oauth2 server
	oauth2URL := "/oauth2"
	ts := httptest.NewServer(testOAuth2Server(oauth2URL, message))
	authEndpoint := ts.URL + path.Join(oauth2URL, "/authorize")
	tokenEndpoint := ts.URL + path.Join(oauth2URL, "/token")
	defer ts.Close()

	t.Logf("auth endpoint %#v", authEndpoint)

	// create test client server
	tcsbase := "/example_app/"
	tcspath := tcsbase + "code"
	tcs := httptest.NewServer(testAppServer(tcspath))
	defer tcs.Close()

	// create dummy oauth client and user
	c, u := createDummies(password, tcs.URL+tcsbase)
	code, err := getCodeHTTP(getCodeRequest(c, u, password, authEndpoint, tcs.URL+tcspath))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// retrieve token from token endpoint
	// get response from client web app redirect uri
	token, err := getTokenHTTP(getTokenRequest(c, code, tokenEndpoint, tcs.URL+tcspath))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// retrieve a testing content path
	body, err := func(token, contentURL string) (body string, err error) {

		log.Printf("Test accessing content with token ====")

		req, err := http.NewRequest("GET", contentURL, nil)
		req.Header.Add("Authority", token)

		// new http client to emulate user request
		hc := &http.Client{}
		resp, err := hc.Do(req)
		if err != nil {
			err = fmt.Errorf("Failed run the request: %s", err.Error())
			return
		}

		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("Failed to read body: %s", err.Error())
			return
		}

		body = string(raw)
		return
	}(token, ts.URL+"/content")

	// quit if error
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// final result
	if want, have := message, body; want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
	log.Printf("result: %#v", string(body))

}
