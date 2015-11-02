package context_test

import (
	"net/http"
	"testing"

	"github.com/gourd/kit/context"

	gkhttp "github.com/go-kit/kit/transport/http"
	gcontext "github.com/gorilla/context"
	glcontext "golang.org/x/net/context"
)

func TestUseGorilla(t *testing.T) {
	var fn gkhttp.RequestFunc = context.UseGorilla
	t.Log("context.UseGorilla implements http.RequestFunc")

	r := &http.Request{}
	key := "foo"
	val := "bar"

	// test set and get value
	gcontext.Set(r, key, val)
	ctx := fn(nil, r)

	res := ctx.Value(key)
	if res != val {
		t.Errorf("failed to set context value with gorilla/context. expected %#v, got %#v",
			val, res)
	}
}

func TestClearGorilla(t *testing.T) {
	r := &http.Request{}
	key := "foo"
	val := "bar"

	// endpoint to test
	ep1 := func(ctx glcontext.Context, request interface{}) (response interface{}, err error) {
		if r, ok := context.HTTPRequest(ctx); ok {
			gcontext.Set(r, key, val) // this is tested in TestUseGorilla
		}
		return
	}
	ep2 := context.ClearGorilla(ep1)

	// fake request to test the endpoint
	ep2(context.UseGorilla(nil, r), nil)
	if _, ok := gcontext.GetOk(r, key); ok {
		t.Error("still be able to get gorilla context value after clear")
	}
}
