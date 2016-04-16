package oauth2

import (
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type privateKey int

const (
	tokenKey privateKey = iota
	accessKey
	osinAuthKey
)

// UseToken reads the token information from header ("Authority")
// and add to the context. Implements go-kit httptransport BeforeFunc
func UseToken(ctx context.Context, r *http.Request) context.Context {
	token := r.Header.Get("Authority")
	return context.WithValue(ctx, tokenKey, token)
}

// GetToken reads the token from context
func GetToken(ctx context.Context) (token string) {
	if v := ctx.Value(tokenKey); v == nil {
		return
	} else if str, ok := v.(string); ok {
		token = str
	}
	return
}

// WithAccess implements go-kit httptransport RequestFunc
// Adds the current HTTP Request to context.Context
func WithAccess(parent context.Context, ad *AccessData) context.Context {
	return context.WithValue(parent, accessKey, ad)
}

// GetAccess returns oauth2 AccessData stored in session
func GetAccess(ctx context.Context) (d *AccessData) {
	ptr := ctx.Value(accessKey)

	// TODO: use logger := log.NewContext(,sg)
	logger := msg
	logger.Log(
		"func", "GetAccess",
		"access", fmt.Sprintf("%#v", ptr))

	if ptr == nil {
		return
	}
	if ad, ok := ptr.(*AccessData); ok {
		d = ad
	}
	return
}

// withOsinAuthRequest implements go-kit httptransport RequestFunc.
// Adds an *osin.AuthorizeRequest to the context
func withOsinAuthRequest(parent context.Context, ar *osin.AuthorizeRequest) context.Context {
	return context.WithValue(parent, osinAuthKey, ar)
}

// getOsinAuthRequest retrive the *osin.AuthorizeRequest in
// context
func getOsinAuthRequest(ctx context.Context) (ar *osin.AuthorizeRequest) {
	ptr := ctx.Value(osinAuthKey)

	// TODO: use logger := log.NewContext(,sg)
	logger := msg
	logger.Log(
		"func", "getOsinAuthRequest",
		"access", fmt.Sprintf("%#v", ptr))
	if ptr == nil {
		return
	} else if v, ok := ptr.(*osin.AuthorizeRequest); ok {
		ar = v
	}
	return
}

// Middleware retrieves token from context with GetToken(),
// then set the AccessData to the context with WithAccess().
//
// Inner endpoint may retrieve the AccessData using GetAccess().
func Middleware(inner endpoint.Endpoint) endpoint.Endpoint {

	LoadTokenAccess := func(ctx context.Context) context.Context {
		storage := &Storage{}
		storage.SetContext(ctx)

		token := GetToken(ctx)
		if token == "" {
			return ctx
		}

		osinAccess, err := storage.LoadAccess(token)
		if err != nil {

			// TODO: use logger := log.NewContext(,sg)
			errLogger := errMsg
			errLogger.Log(
				"func", "LoadTokenAccess (Middleware)",
				"token", token,
				"message", "storage.LoadAccess failed to load access",
				"error", err.Error())
			return ctx
		}

		ad := &AccessData{}
		ad.ReadOsin(osinAccess)
		return WithAccess(ctx, ad)
	}

	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx = LoadTokenAccess(ctx)
		return inner(ctx, request)
	}
}

// LoadTokenAccess reads token information from header ("Authority")
// and, if AccessData found for the given token, add to context
func LoadTokenAccess(ctx context.Context) context.Context {

	storage := &Storage{}
	storage.SetContext(ctx)

	logger, errLogger := msg, errMsg

	token := GetToken(ctx)
	if token == "" {
		return ctx
	}

	osinAccess, err := storage.LoadAccess(token)
	if err != nil {
		errLogger.Log(
			"func", "LoadTokenAccess",
			"token", token,
			"message", "storage.LoadAccess failed to load access",
			"error", err.Error())
		return ctx
	}

	logger.Log(
		"func", "LoadTokenAccess",
		"token", token,
		"osinAccess.UserData", fmt.Sprintf("%#v", osinAccess.UserData))

	switch osinAccess.UserData.(type) {
	case *User:
	default:
		panic("hello")
	}
	ad := &AccessData{}
	ad.ReadOsin(osinAccess)
	return WithAccess(ctx, ad)
}
