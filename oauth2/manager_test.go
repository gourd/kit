package oauth2_test

import (
	"errors"

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

	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/gourd/kit/store"
)

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
func testOauth2Dummies(password, redirect string) (*oauth2.Client, *oauth2.User) {
	r := &http.Request{}

	// generate dummy user
	us, err := store.Providers.Store(r, "User")
	if err != nil {
		panic(err)
	}
	u := dummyNewUser(password)
	err = us.Create(store.NewConds(), u)
	if err != nil {
		panic(err)
	}

	// get related dummy client
	cs, err := store.Providers.Store(r, "Client")
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

// handles redirection
func testNoRedirect(req *http.Request, via []*http.Request) error {
	log.Printf("redirect url: %#v", req.URL.Query().Get("code"))
	return testRedirectErr{"no redirect", req.URL}
}

// testGetCode request code from authorize endpoint
// with given redirect URL.
// It build user request to authorization endpoint
// get response from client web app redirect uri
func testGetCode(c *oauth2.Client, u *oauth2.User, password, authURL, redirect string) (code string, err error) {

	log.Printf("Test retrieving code ====")

	// login form
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
		err = fmt.Errorf("Failed to form new request: %s", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// new http client to emulate user request
	hc := &http.Client{
		CheckRedirect: testNoRedirect,
	}
	_, rerr := hc.Do(req)
	uerr := rerr.(*url.Error).Err

	// examine error
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

// example server web app
func testOAuth2ServerApp(msg string) http.Handler {

	rtr := pat.New()

	// oauth2 manager
	m := oauth2.NewManager()

	// add oauth2 endpoints to router
	// ServeEndpoints bind OAuth2 endpoints to a given base path
	// Note: this is router specific and need to be generated somehow
	oauth2.RoutePat(rtr, "/oauth", m.GetEndpoints())

	// add a route the requires access
	rtr.Get("/content", func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Dummy content page accessed")

		// obtain access
		a, err := oauth2.GetAccess(r)
		if err != nil {
			log.Printf("Dummy content: access error: %s", err.Error())
			fmt.Fprint(w, "Permission Denied")
			return
		}

		// test the access
		if a == nil {
			fmt.Fprint(w, "Unable to gain Access")
			return
		}

		// no news is good news
		fmt.Fprint(w, msg)
	})

	// create negroni middleware handler
	// with middlewares
	n := negroni.New()
	n.Use(negroni.Wrap(m.Middleware()))

	// use router in negroni
	n.UseHandler(rtr)

	return n
}

// example client web app in the login
func testOAuth2ClientApp(path string) *pat.Router {
	rtr := pat.New()

	log.Printf("testOAuth2ClientApp(%#v)", path)

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

func TestOAuth2(t *testing.T) {

	// create test oauth2 server
	ts := httptest.NewServer(testOAuth2ServerApp("Success"))
	defer ts.Close()

	// create test client server
	tcsbase := "/example_app/"
	tcspath := tcsbase + "code"
	tcs := httptest.NewServer(testOAuth2ClientApp(tcspath))
	defer tcs.Close()

	// a dummy password for dummy user
	password := "password"

	// create dummy oauth client and user
	c, u := testOauth2Dummies(password, tcs.URL+tcsbase)
	code, err := testGetCode(c, u, password, ts.URL+"/oauth/authorize", tcs.URL+tcspath)
	if err != nil {
		// quit if error
		t.Errorf(err.Error())
		return
	}

	// retrieve token from token endpoint
	// get response from client web app redirect uri
	token, err := func(c *oauth2.Client, code, redirect string) (token string, err error) {

		log.Printf("Test retrieving token ====")

		// build user request to token endpoint
		form := &url.Values{}
		form.Add("code", code)
		form.Add("client_id", c.GetId())
		form.Add("client_secret", c.Secret)
		form.Add("grant_type", "authorization_code")
		form.Add("redirect_uri", redirect)
		req, err := http.NewRequest("POST",
			ts.URL+"/oauth/token",
			strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Errorf("Failed to form new request: %s", err.Error())
		}

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

	}(c, code, tcs.URL+tcspath)

	// quit if error
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// retrieve a testing content path
	body, err := func(token string) (body string, err error) {

		log.Printf("Test accessing content with token ====")

		req, err := http.NewRequest("GET", ts.URL+"/content", nil)
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
	}(token)

	// quit if error
	if err != nil {
		t.Errorf(err.Error())
		return
	} else if body != "Success" {
		t.Errorf("Content Incorrect. Expecting \"Success\" but get \"%s\"", body)
	}

	// final result
	log.Printf("result: \"%s\"", body)

}
