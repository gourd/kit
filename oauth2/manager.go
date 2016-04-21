package oauth2

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RangelReale/osin"
	"github.com/go-kit/kit/log"
	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
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
	m.SetLoginFormFunc(NewLoginFormFunc("user_id", DefaultLoginTpl))

	// set default login parser
	m.SetUserFunc(NewUserFunc("user_id"))

	return
}

// UserFunc reads the login form request and returns an OAuth2User
// for the reqeust. If there is error obtaining the user, an error
// is returned
type UserFunc func(r *http.Request, us store.Store) (u OAuth2User, err error)

// LoginFormContext represents the context of the login form rendering
type LoginFormContext struct {
	Context        context.Context
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	LoginErr       error
	ActionURL      *url.URL
	Logger         log.Logger
}

// LoginFormFunc handles GET request of the authorize endpoint
// and displays a login form for user to login.
// The action parameter provides a pre-rendered URL to login
type LoginFormFunc func(lctx *LoginFormContext) (err error)

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

func (m *Manager) showLoginForm(lctx *LoginFormContext, w http.ResponseWriter, r *http.Request) {

	logger := msg
	logger.Log(
		"func", "showLoginForm (Manager.GetEndpoints)")

	// build action query
	ar := getOsinAuthRequest(lctx.Context) // presume the context has *osin.AuthorizeRequest
	aq := url.Values{}
	aq.Add("response_type", string(ar.Type))
	aq.Add("client_id", ar.Client.GetId())
	aq.Add("state", ar.State)
	aq.Add("scope", ar.Scope)
	aq.Add("redirect_uri", ar.RedirectUri)

	// form action url
	aurl := r.URL
	aurl.RawQuery = aq.Encode()

	logger.Log(
		"func", "showLoginForm (Manager.GetEndpoints)",
		"action url", aurl)

	lctx.ActionURL = aurl

	w.Header().Add("Content-Type", "text/html;charset=utf8")
	if err := m.loginFormFunc(lctx); err != nil {
		serr := store.ExpandError(err)
		logger.Log(
			"func", "showLoginForm (Manager.GetEndpoints)",
			"action url", aurl,
			"error", serr.ServerMsg)
	}
}

// GetEndpoints generate endpoints http handers and return
func (m *Manager) GetEndpoints(factory store.Factory) *Endpoints {

	// try to login with given request login
	tryLogin := func(ctx context.Context, r *http.Request) (user OAuth2User, err error) {

		logger := msg
		logger.Log(
			"func", "tryLogin (Manager.GetEndpoints)")

		// parse POST input
		r.ParseForm()
		if r.Method == "POST" {

			var u OAuth2User
			var us store.Store

			// get and check password non-empty
			password := r.Form.Get("password")
			if password == "" {
				err = errors.New("empty password")
				return
			}

			// obtain user store
			us, err = store.Get(ctx, KeyUser)
			if err != nil {
				err = store.Error(
					http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError)).
					TellServer("error obtaining user store: %s", err.Error())
				return
			}

			// get user by userFunc
			u, err = m.userFunc(r, us)
			if err != nil {
				serr := store.ExpandError(err)
				if serr.Status == http.StatusNotFound {
					err = store.Error(http.StatusBadRequest, "user id or password incorrect").
						TellServer("user not found")
				} else {
					err = store.Error(
						http.StatusInternalServerError,
						http.StatusText(http.StatusInternalServerError)).
						TellServer("error obtaining user: %s", serr.ServerMsg)
				}
				return
			}

			// if user is nil, user not found
			if u == nil {
				err = store.Error(http.StatusBadRequest, "user not found")
				return
			}

			// if password does not match
			if !u.PasswordIs(password) {
				err = store.Error(http.StatusBadRequest, "user id or password incorrect").
					TellServer("incorrect password")
				return
			}

			// return pointer of user object, allow it to be re-cast
			logger.Log(
				"func", "tryLogin (Manager.GetEndpoints)",
				"message", "login success")
			user = u
			return
		}

		// no POST input or incorrect login, show form
		// end login handling sequence and wait for
		// user input from login form
		err = store.Error(http.StatusUnauthorized, "Require login").
			TellServer("no POST input")
		return
	}

	type ContextHandlerFunc func(ctx context.Context,
		w http.ResponseWriter, r *http.Request) *osin.Response

	// sessionContext takes a ContextHandlerFunc and returns
	// a http.HandlerFunc
	sessionContext := func(inner ContextHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// per connection based context.Context, with factory
			ctx := store.WithFactory(context.Background(), factory)
			defer store.CloseAllIn(ctx)
			if resp := inner(ctx, w, r); resp != nil {
				if resp.InternalError != nil {
					errLogger := errMsg
					errLogger.Log(
						"func", "sessionContext (Manager.GetEndpoints)",
						"error", resp.InternalError.Error())
				}
				osin.OutputJSON(resp, w, r)
			}
		}
	}

	ep := Endpoints{}

	// authorize endpoint
	ep.Auth = sessionContext(func(ctx context.Context,
		w http.ResponseWriter, r *http.Request) *osin.Response {

		logger := msg
		logger.Log(
			"endpoint", "auth")

		srvr := m.osinServer
		resp := srvr.NewResponse()
		resp.Storage.(*Storage).SetContext(ctx)

		// handle authorize request with osin
		if ar := srvr.HandleAuthorizeRequest(resp, r); ar != nil {
			logger.Log(
				"endpoint", "auth",
				"message", "handle authorize request")

			// TODO: maybe redirect to another URL for
			//       dedicated login form flow?
			var err error
			if ar.UserData, err = tryLogin(ctx, r); err != nil {
				serr := store.ExpandError(err)
				logger.Log(
					"endpoint", "auth",
					"message", "handle authorize request",
					"error", serr.ServerMsg)

				lctx := &LoginFormContext{
					Context:        withOsinAuthRequest(ctx, ar),
					LoginErr:       err,
					ResponseWriter: w,
					Request:        r,
					Logger:         logger,
				}
				m.showLoginForm(lctx, w, r)
				return nil
			}

			logger.Log(
				"endpoint", "auth",
				"message", "User obtained",
				"osin.AuthorizeData.UserData", fmt.Sprintf("%#v", ar.UserData))

			ar.Authorized = true
			srvr.FinishAuthorizeRequest(resp, r, ar)
		}

		logger.Log(
			"endpoint", "auth",
			"message", "User obtained",
			"response", fmt.Sprintf("%#v", resp))

		return resp
	})

	// token endpoint
	ep.Token = sessionContext(func(ctx context.Context,
		w http.ResponseWriter, r *http.Request) *osin.Response {

		logger := msg
		logger.Log(
			"endpoint", "token")

		srvr := m.osinServer
		resp := srvr.NewResponse()
		resp.Storage.(*Storage).SetContext(ctx)

		if ar := srvr.HandleAccessRequest(resp, r); ar != nil {
			// TODO: handle authorization
			// check if the user has the permission to grant the scope
			logger.Log(
				"endpoint", "token",
				"message", "access successful")
			ar.Authorized = true
			srvr.FinishAccessRequest(resp, r, ar)
		}

		logger.Log(
			"endpoint", "token",
			"response", fmt.Sprintf("%#v", resp))
		return resp
	})

	// information endpoint
	ep.Info = sessionContext(func(ctx context.Context,
		w http.ResponseWriter, r *http.Request) *osin.Response {

		logger := msg
		logger.Log(
			"endpoint", "information")

		srvr := m.osinServer

		resp := srvr.NewResponse()
		resp.Storage.(*Storage).SetContext(ctx)
		defer resp.Close()

		if ir := srvr.HandleInfoRequest(resp, r); ir != nil {
			srvr.FinishInfoRequest(resp, r, ir)
		}

		logger.Log(
			"endpoint", "information",
			"response", fmt.Sprintf("%#v", resp))
		return resp
	})

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
