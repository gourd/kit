package gourdctx_test

import (
	"net/http"
	"testing"

	gcontext "github.com/gorilla/context"
	gourdctx "github.com/gourd/kit/context"
	"golang.org/x/net/context"
)

func TestNew(t *testing.T) {
	// ensure the wrapper implements context.Context
	key := "foo"
	val := "bar"
	ctx0 := gourdctx.New(nil)
	ctx1 := context.WithValue(ctx0, key, val)

	if ctx0 == ctx1 {
		t.Error("ctx1 should be a clone of ctx1")
	}

	if ctx0.Value(key) == val {
		t.Errorf("WithValue should not be %#v", val)
	}
}

func TestHTTPRequest(t *testing.T) {
	r := &http.Request{}
	ctx0 := gourdctx.WithHTTPRequest(context.Background(), r)

	if want, have := r, gourdctx.HTTPRequest(ctx0); want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}
}

func TestGorilla(t *testing.T) {
	key := "foo"
	val0 := "bar 0"
	val1 := "bar 1"
	r := &http.Request{}
	ctx0 := gourdctx.WithGorilla(context.Background(), r)
	ctx1 := context.WithValue(ctx0, key, val1)
	gcontext.Set(r, key, val0)

	if ctx0 == ctx1 {
		t.Error("ctx1 should be a clone of ctx1")
	}

	if want, have := val0, ctx0.Value(key); want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}

	if want, have := val1, ctx1.Value(key); want != have {
		t.Errorf("expect %#v, got %#v", want, have)
	}

}

func TestWithID(t *testing.T) {
	id := "foobar"
	ctx0 := context.Background()
	ctx1 := gourdctx.WithID(ctx0, id)

	if want, have := "", gourdctx.GetID(ctx0); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := id, gourdctx.GetID(ctx1); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
