package context_test

import (
	"net/http"
	"testing"

	gourdctx "github.com/gourd/kit/context"

	httptransport "github.com/go-kit/kit/transport/http"
	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

func TestUseGorilla(t *testing.T) {
	var fn httptransport.RequestFunc = gourdctx.UseGorilla
	t.Log("context.UseGorilla implements http.RequestFunc")

	r := &http.Request{}
	key := "foo"
	val := "bar"

	// test set and get value
	gcontext.Set(r, key, val)
	if val2, ok := gcontext.GetOk(r, key); !ok {
		t.Errorf("failed to get value of key %#v", key)
	} else if want, have := val, val2; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// test getting value from the context
	ctx := fn(context.Background(), r)
	res := ctx.Value(key)
	if want, have := val, res; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestClearGorilla(t *testing.T) {
	r := &http.Request{}
	key := "foo"
	val := "bar"

	// endpoint to test
	ep1 := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		gcontext.Set(gourdctx.HTTPRequest(ctx), key, val) // this is tested in TestUseGorilla
		return
	}
	ep2 := gourdctx.ClearGorilla(ep1)

	// fake request to test the endpoint
	ep2(gourdctx.UseGorilla(nil, r), nil)
	if _, ok := gcontext.GetOk(r, key); ok {
		t.Error("still be able to get gorilla context value after clear")
	}
}

func TestUseID(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	ctx := gourdctx.UseID(context.Background(), r)
	var id string

	if id = gourdctx.GetRequestID(r); id == "" {
		t.Errorf("unexpected empty string")
		return
	}

	if want, have := id, gourdctx.GetID(ctx); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestUseID_Reuse(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	id := "hello"
	r.Header.Set("X-GOURD-ID", id)
	ctx := gourdctx.UseID(context.Background(), r)

	if want, have := id, gourdctx.GetRequestID(r); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := id, gourdctx.GetID(ctx); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}
