package perm_test

import (
	"github.com/gourd/kit/perm"

	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"testing"
)

func testEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, ok := perm.GetMuxOk(ctx)
		if !ok {
			err = errors.New("unable to get perm.Mux")
		}
		return
	}
}

func TestGetMux(t *testing.T) {
	ep := testEndpoint()
	_, err := ep(perm.WithMux(nil, perm.NewMux()), nil)
	if err != nil {
		t.Errorf("test error: %#v", err.Error())
	}
}

func TestGetFail(t *testing.T) {
	ep := testEndpoint()
	_, err := ep(perm.WithMux(nil, nil), nil)
	if err == nil {
		t.Error("unable to raise error when no perm.Mux in context")
	} else if want, have := "unable to get perm.Mux", err.Error(); want != have {
		t.Errorf("want: %#v, got: %#v", want, have)
	}
}

func TestMiddleware(t *testing.T) {
	ep := testEndpoint()
	var mw endpoint.Middleware = perm.UseMux(perm.NewMux())
	_, err := mw(ep)(nil, nil)
	if err != nil {
		t.Errorf("test error: %#v", err.Error())
	}
}
