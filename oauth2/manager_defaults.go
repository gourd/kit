package oauth2

import (
	"github.com/RangelReale/osin"
	"github.com/asaskevich/govalidator"
	"github.com/gourd/kit/store"

	"net/http"
	"text/template"
)

// DefaultStorage returns Storage that attachs to default stores
func DefaultStorage() (s *Storage) {
	s = &Storage{}
	return
}

// DefaultOsinConfig returns a preset config suitable
// for most generic oauth2 usage
func DefaultOsinConfig() (cfg *osin.ServerConfig) {
	cfg = osin.NewServerConfig()
	cfg.AllowGetAccessRequest = true
	cfg.AllowClientSecretInParams = true
	cfg.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
	}
	cfg.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{
		osin.CODE,
		osin.TOKEN,
	}
	return
}

// DefaultLoginTpl is the HTML template for login form by default
const DefaultLoginTpl = `
<!DOCTYPE html>
<html>
<head>
<title>Login</title>
<style>
body, html { margin: 0; font-size: 18pt; background-color: #EEE; }
#login-box { max-width: 100%; width: 400px; margin: 10% auto 0; box-shadow: 0 0 3px #777; background-color: #F9F9F9; }
#login-box h1 { font-size: 1.2em; margin: 0 0 0.5em; }
#login-box .content { margin: 0 20px; padding: 30px 0; text-align: center; }
#login-box .field { display: block; width: 88%; background-color: #FFF; }
#login-box .field { border: solid 1px #EEE; padding: 0.4em 1em; line-height: 1.3em; }
#login-box .actions { text-align: center; }
#login-box button { width: 100%; }
#error-message { background-color: rgba(255, 0, 0, 0.2); color: rgba(255, 0, 0, 0.8); padding: 0.5em 2em; }
</style>
</head>
<body>
	<div id="login-box"><div class="content">
		<h1>{{ .Title }}</h1>
		{{if .LoginErr }}<div id="error-message" class="message">{{ .LoginErr }}</div>{{end}}
		<form action="{{ .FormAction }}" method="POST">
			<div class="field-wrapper">
				<input name="user_id" type="text" class="field"
					placeholder="{{ .TextUserID }}" value="{{ .UserID }}" autofocus />
			</div>
			<div class="field-wrapper">
				<input name="password" type="password" class="field"
					placeholder="{{ .TextPassword }}" />
			</div>
			<div class="actions">
				<button type="submit">{{ .TextSubmit }}</button>
			</div>
		</form>
	</div></div>
</body>
</html>
`

// NewUserFunc creates the default parser of login HTTP request
func NewUserFunc(idName string) UserFunc {
	return func(r *http.Request, us store.Store) (ou OAuth2User, err error) {

		var c store.Conds

		id := r.Form.Get(idName)

		if id == "" {
			serr := store.Error(http.StatusBadRequest, "empty user identifier")
			err = serr
			return
		}

		// different condition based on the user_id field format
		if govalidator.IsEmail(id) {
			c = store.NewConds().Add("email", id)
		} else {
			c = store.NewConds().Add("username", id)
		}

		// get user from database
		u := us.AllocEntity()
		err = us.One(c, u)

		if err != nil {
			serr := store.ExpandError(err)
			if serr.Status != http.StatusNotFound {
				serr.TellServer("Error searching user %#v: %s", id, serr.ServerMsg)
				return
			}
			err = serr
			return
		}

		// if user does not exists
		if u == nil {
			serr := store.Error(http.StatusBadRequest, "Username or Password incorrect")
			serr.TellServer("Unknown user %#v attempt to login", id)
			err = serr
			return
		}

		// cast the user as OAuth2User
		// and do password check
		ou, ok := u.(OAuth2User)
		if !ok {
			serr := store.Error(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			serr.TellServer("User cannot be cast as OAuth2User")
			err = serr
			return
		}

		return
	}
}

// NewLoginFormFunc creates a LoginFormFunc from given template
func NewLoginFormFunc(idName, tpl string) LoginFormFunc {

	// compile template for login form
	loginTpl, err := template.New("loginForm").Parse(tpl)
	if err != nil {
		panic(err) // should not happen, simply panic
	}

	return func(lctx *LoginFormContext) (err error) {

		// TODO: pass the login error into showLoginForm context
		//       and display it to the visitor

		// template variables
		vars := map[string]interface{}{
			"Title":        "Login",
			"FormAction":   lctx.ActionURL,
			"UserID":       "",
			"TextUserID":   "Login ID or Email",
			"TextPassword": "Password",
			"TextSubmit":   "Login",
		}

		if lctx.Request.Method == "POST" && lctx.LoginErr != nil {
			vars["LoginErr"] = lctx.LoginErr.Error()
			vars["UserID"] = lctx.Request.Form.Get(idName)
		}

		// render the form with vars
		err = loginTpl.Execute(lctx.ResponseWriter, vars)
		if err != nil {
			serr := store.ExpandError(err)
			serr.TellServer("error executing login template: %#v", err.Error())
			err = serr
			return
		}

		return
	}
}
