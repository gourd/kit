package oauth2

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"github.com/RangelReale/osin"
	"github.com/gourd/kit/store"
)

// Endpoints contains http handler func of different endpoints
type Endpoints struct {
	Auth  http.HandlerFunc
	Token http.HandlerFunc
	Info  http.HandlerFunc
}

// NewManager returns a oauth2 manager with default configs
func NewManager() (m *Manager) {

	m = &Manager{}

	// provide storage to osin server
	// provide osin server to Manager
	m.InitOsin(DefaultOsinConfig())

	// set default login form handler
	// (only handles GET request of the authorize endpoint)
	m.SetLoginFormFunc(NewLoginFormFunc(DefaultLoginTpl))

	// set default login parser
	m.SetUserFunc(NewUserFunc("user_id"))

	return
}

// UserFunc reads the login form request and returns an OAuth2User
// for the reqeust. If there is error obtaining the user, an error
// is returned
type UserFunc func(r *http.Request, us store.Store) (u OAuth2User, err error)

// LoginFormFunc handles GET request of the authorize endpoint
// and displays a login form for user to login.
// The action parameter provides a pre-rendered URL to login
type LoginFormFunc func(w http.ResponseWriter, r *http.Request, action *url.URL) (err error)

// Manager handles oauth2 related request
// Also provide middleware for other http handler function
// to access scope related information
type Manager struct {
	storage       *Storage
	osinServer    *osin.Server
	loginFormFunc LoginFormFunc
	userFunc      UserFunc
}

// InitOsin set the OsinServer
func (m *Manager) InitOsin(cfg *osin.ServerConfig) *Manager {
	m.osinServer = osin.NewServer(cfg, m.storage)
	return m
}

// GetEndpoints generate endpoints http handers and return
func (m *Manager) GetEndpoints(factory store.Factory) *Endpoints {

	// handle login
	handleLogin := func(ar *osin.AuthorizeRequest, ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

		w.Header().Add("Content-Type", "text/html;charset=utf8")
		log.Printf("handleLogin")

		// parse POST input
		r.ParseForm()
		if r.Method == "POST" {

			var u OAuth2User
			var us store.Store

			// get and check password
			password := r.Form.Get("password")
			if password == "" {
				err = errors.New("empty password")
				return
			}

			// obtain user store
			us, err = store.Get(ctx, KeyUser)
			if err != nil {
				log.Printf("Error obtaining user store: %s", err.Error())
				err = errors.New("Internal Server Error")
				return
			}

			// get user by userFunc
			u, err = m.userFunc(r, us)
			if err != nil {
				serr := store.ExpandError(err)
				log.Printf("Error obtaining user: %#v", serr.ServerMsg)
				return
			}

			// if user is nil, user not found
			if u == nil {
				log.Printf("User not found")
				err = errors.New("user not found")
				return
			}

			// if password does not match
			if !u.PasswordIs(password) {
				log.Print("Incorrect password")
				err = errors.New("username or password incorrect")
			} else {
				log.Printf("Login success")
			}

			// return pointer of user object, allow it to be re-cast
			ar.UserData = u
			return
		}

		// no POST input or incorrect login, show form

		// build action query
		aq := url.Values{}
		aq.Add("response_type", string(ar.Type))
		aq.Add("client_id", ar.Client.GetId())
		aq.Add("state", ar.State)
		aq.Add("scope", ar.Scope)
		aq.Add("redirect_uri", ar.RedirectUri)

		// form action url
		aurl := r.URL
		aurl.RawQuery = aq.Encode()

		log.Printf("action URL: %#v", aurl)

		m.loginFormFunc(w, r, aurl)

		// end login handling sequence and wait for
		// user input from login form
		err = errors.New("need login")

		return
	}

	ep := Endpoints{}

	// authorize endpoint
	ep.Auth = func(w http.ResponseWriter, r *http.Request) {

		log.Printf("auth endpoint")

		// per connection based context.Context, with factory
		ctx := store.WithFactory(context.Background(), factory)
		defer store.CloseAllIn(ctx)

		srvr := m.osinServer
		resp := srvr.NewResponse()
		resp.Storage.(*Storage).SetContext(ctx)

		// handle authorize request with osin
		if ar := srvr.HandleAuthorizeRequest(resp, r); ar != nil {
			log.Printf("handle authorize request")
			if err := handleLogin(ar, ctx, w, r); err != nil {
				return
			}
			log.Printf("OAuth2 Authorize Request: User obtained: %#v", ar.UserData)
			ar.Authorized = true
			srvr.FinishAuthorizeRequest(resp, r, ar)
		}
		if resp.InternalError != nil {
			log.Printf("Internal Error: %s", resp.InternalError.Error())
		}
		log.Printf("OAuth2 Authorize Response: %#v", resp)
		osin.OutputJSON(resp, w, r)

	}

	// token endpoint
	ep.Token = func(w http.ResponseWriter, r *http.Request) {

		log.Printf("token endpoint")

		// per connection based context.Context, with factory
		ctx := store.WithFactory(context.Background(), factory)
		defer store.CloseAllIn(ctx)

		srvr := m.osinServer
		resp := srvr.NewResponse()
		resp.Storage.(*Storage).SetContext(ctx)

		if ar := srvr.HandleAccessRequest(resp, r); ar != nil {
			// TODO: handle authorization
			// check if the user has the permission to grant the scope
			log.Printf("Access successful")
			ar.Authorized = true
			srvr.FinishAccessRequest(resp, r, ar)
		} else if resp.InternalError != nil {
			log.Printf("Internal Error: %s", resp.InternalError.Error())
		}
		log.Printf("OAuth2 Token Response: %#v", resp)
		osin.OutputJSON(resp, w, r)

	}

	// information endpoint
	ep.Info = func(w http.ResponseWriter, r *http.Request) {

		// per connection based context.Context, with factory
		ctx := store.WithFactory(context.Background(), factory)
		defer store.CloseAllIn(ctx)

		log.Printf("information endpoint")
		srvr := m.osinServer

		resp := srvr.NewResponse()
		resp.Storage.(*Storage).SetContext(ctx)
		defer resp.Close()

		if ir := srvr.HandleInfoRequest(resp, r); ir != nil {
			srvr.FinishInfoRequest(resp, r, ir)
		}
		osin.OutputJSON(resp, w, r)

	}

	return &ep

}

// SetLoginFormFunc sets the handler to display login form
func (m *Manager) SetLoginFormFunc(f LoginFormFunc) {
	m.loginFormFunc = f
}

// SetUserFunc sets the parser for login request.
// Will be called when endpoint POST request
//
// Manager will then search user with `idField` equals to `id`.
// Then it will check User.HasPassword(`password`)
// (User should implement OAuth2User interface)
// to see if the password is correct
func (m *Manager) SetUserFunc(f UserFunc) {
	m.userFunc = f
}
