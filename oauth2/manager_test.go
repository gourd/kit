package oauth2_test

import (
	"errors"
	"io"
	"path"

	"golang.org/x/net/context"

	"github.com/gourd/kit/oauth2"

	"encoding/json"
	"fmt"
	"io/ioutil"
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

type codeResponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func readCodeResponse(u string) (resp *codeResponse, err error) {
	redirect, err := url.Parse(u)
	if err != nil {
		err = fmt.Errorf("unexpected error: %s", err)
		return
	}
	q := redirect.Query()
	resp = &codeResponse{
		Code:  q.Get("code"),
		State: q.Get("state"),
	}
	return
}

type tokenRespnose struct {
	Token        string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type testContext struct {
	t            *testing.T
	client       *oauth2.Client
	user         *oauth2.User
	password     string
	redirectBase string
	redirectURL  string
	oauth2Base   string
	oauth2Path   string
	code         string
	token        string
	refresh      string
}

func (ctx *testContext) Code() string {
	return ctx.code
}

func (ctx *testContext) Token() string {
	return ctx.token
}

func (ctx *testContext) RefreshToken() string {
	return ctx.refresh
}

func (ctx *testContext) AuthEndpoint() string {
	return ctx.oauth2Base + path.Join(ctx.oauth2Path, "authorize")
}

func (ctx *testContext) TokenEndpoint() string {
	return ctx.oauth2Base + path.Join(ctx.oauth2Path, "token")
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
func createStoreDummies(ctx context.Context, password, redirect string) (*oauth2.Client, *oauth2.User) {

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
	c.UserID = u.ID
	err = cs.Create(store.NewConds(), c)
	if err != nil {
		panic(err)
	}

	return c, u
}

// getCodeRequest generates the http request to get code
func getCodeRequest(ctxTest *testContext) *http.Request {

	// login form request
	form := url.Values{}
	form.Add("user_id", ctxTest.user.Username)
	form.Add("password", ctxTest.password)
	ctxTest.t.Logf("form send: %s", form.Encode())

	// build the query string
	q := &url.Values{}
	q.Add("response_type", "code")
	q.Add("client_id", ctxTest.client.GetId())
	q.Add("redirect_uri", ctxTest.redirectURL)

	req, err := http.NewRequest("POST",
		ctxTest.AuthEndpoint()+"?"+q.Encode(),
		strings.NewReader(form.Encode()))
	if err != nil {
		panic(err) // not quite possible
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

// getCodeHTTP runs the getCodeRequest with actual
// HTTP client and parse the result as code, error
func getCodeHTTP(t *testing.T, req *http.Request) (code string, err error) {

	t.Logf("Test retrieving code ====")

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
	t.Logf("code: %#v", code)

	return
}

// getTokenRequest generates request which client app
// send to oauth2 server for the token
func getTokenRequest(ctxTest *testContext) *http.Request {
	// build user request to token endpoint
	form := &url.Values{}
	form.Add("code", ctxTest.code)
	form.Add("client_id", ctxTest.client.GetId())
	form.Add("client_secret", ctxTest.client.Secret)
	form.Add("grant_type", "authorization_code")
	form.Add("redirect_uri", ctxTest.redirectURL)
	req, err := http.NewRequest("POST",
		ctxTest.TokenEndpoint(),
		strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	return req
}

// getRefreshRequest generates request which client app
// send to oauth2 server for the token
func getRefreshRequest(ctxTest *testContext) *http.Request {
	// build user request to token endpoint
	form := &url.Values{}
	form.Add("client_id", ctxTest.client.GetId())
	form.Add("client_secret", ctxTest.client.Secret)
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", ctxTest.RefreshToken())
	req, err := http.NewRequest("POST",
		ctxTest.TokenEndpoint(),
		strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	return req
}

// getTokenHTTP runs the getTokenRequest with actual
// HTTP client and parse the result as token, error
func getTokenHTTP(t *testing.T, req *http.Request) (token, refresh string, err error) {

	t.Logf("Test retrieving token ====")

	// new http client to emulate user request
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		err = fmt.Errorf("Failed run the request: %s", err.Error())
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Failed read the response body: %s", err.Error())
	}

	t.Logf("raw response: %s", b)

	// read token from token endpoint response (json)
	bodyDecoded := make(map[string]interface{})
	json.Unmarshal(b, &bodyDecoded)

	t.Logf("Response Body: %#v", bodyDecoded)

	if v, ok := bodyDecoded["error"]; ok {
		errType := v.(string)
		if errDesc, ok := bodyDecoded["error_description"]; ok {
			// do nothing
			err = fmt.Errorf("oauth2 error: %s (%s)", errType, errDesc)
			return
		}
		err = fmt.Errorf("oauth2 error: %s", errType)
		return
	}

	if v, ok := bodyDecoded["access_token"]; !ok {
		err = fmt.Errorf("unable to find access_token in response")
		return
	} else if token, ok = v.(string); !ok {
		err = fmt.Errorf("access_token is not string")
		return
	}

	if v, ok := bodyDecoded["refresh_token"]; !ok {
		err = fmt.Errorf("unable to find refresh_token in response")
		return
	} else if refresh, ok = v.(string); !ok {
		err = fmt.Errorf("refresh_token is not string")
		return
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
func getContentHTTP(t *testing.T, r *http.Request) (body string, err error) {

	t.Logf("request URL: %#v", r.URL)

	// new http client to emulate user request
	hc := &http.Client{}
	resp, err := hc.Do(r)
	if err != nil && err != io.EOF {
		t.Logf("err 1: %#v", err)
		err = fmt.Errorf("Failed run the request: %s", err.Error())
		return
	}
	err = nil // in case EOF

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		t.Logf("err 2: %#v", err)
		err = fmt.Errorf("Failed to read body: %s", err.Error())
		return
	}
	err = nil // in case EOF

	body = string(raw)
	return
}

// example server web app
func testOAuth2Server(t *testing.T, baseURL, msg string) http.Handler {

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

	// router function
	rtrFunc := func(path string, methods []string, h http.Handler) error {
		for i := range methods {
			rtr.Add(methods[i], path, h)
		}
		return nil
	}

	// add oauth2 endpoints to router
	// ServeEndpoints bind OAuth2 endpoints to a given base path
	// Note: this is router specific and need to be generated somehow
	oauth2.Route(rtrFunc, baseURL, m.GetEndpoints(factory))

	// add a route the requires access
	rtr.Get("/content", func(w http.ResponseWriter, r *http.Request) {

		ctx := store.WithFactory(context.Background(), factory)
		ctx = oauth2.LoadTokenAccess(oauth2.UseToken(ctx, r))
		t.Logf("Dummy content page accessed")

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

// TestOAuth2HTTP tests the stack with
// actual HTTP call against httptest.Server
// wrapped handlers
func TestOAuth2HTTP(t *testing.T) {

	var err error

	// a dummy password for dummy user
	password := "password"
	message := "Success"

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

	testCtx := &testContext{
		password:     password,
		t:            t,
		redirectBase: "https://test.foobar/example_app/",
		redirectURL:  "https://test.foobar/example_app/code",
		oauth2Path:   "/oauth2",
	}

	// create test oauth2 server
	ts := httptest.NewServer(testOAuth2Server(t, testCtx.oauth2Path, message))
	defer ts.Close()
	testCtx.oauth2Base = ts.URL

	t.Logf("auth endpoint %#v", testCtx.AuthEndpoint())

	// create dummy client and user
	testCtx.client, testCtx.user = createStoreDummies(ctx, testCtx.password, testCtx.redirectBase)
	store.CloseAllIn(ctx)

	// create dummy oauth client and user
	testCtx.code, err = getCodeHTTP(t, getCodeRequest(testCtx))
	if err != nil {
		t.Error(err.Error())
		return
	}

	// retrieve token from token endpoint
	// get response from client web app redirect uri
	testCtx.token, testCtx.refresh, err = getTokenHTTP(t, getTokenRequest(testCtx))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// try to refresh token
	t.Logf(`refresh_token=%s token=%s msg="refresh token test"`, testCtx.refresh, testCtx.token)
	testCtx.token, testCtx.refresh, err = getTokenHTTP(t, getRefreshRequest(testCtx))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf(`refresh_token=%s token=%s msg="refresh token test success"`, testCtx.refresh, testCtx.token)

	// retrieve a testing content path
	body, err := getContentHTTP(t, getContentRequest(testCtx.token, ts.URL+"/content"))
	if err != nil {
		t.Logf("hello: %#v", err)
		t.Errorf(err.Error())
		return
	}

	// final result
	if want, have := message, body; want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
	t.Logf("result: %#v", string(body))

}
