package httpservice_test

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httpservice "github.com/gourd/kit/service/http"

	"testing"
)

func testDesc() httpservice.Desc {
	base, sing, plur := "/some/path", "ball", "balls"
	n := httpservice.NewNoun(sing, plur)
	p := httpservice.NewPaths(base, n, "someid")
	return httpservice.NewDesc(p)
}

func TestDesc_DecodeFunc(t *testing.T) {
	d := testDesc()
	d.SetDecodeFunc("hello", func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		request = "world"
		return
	})

	// test retrieve and use the decoder
	dec := d.GetDecodeFunc("hello")
	req, _ := dec(nil, nil)
	if want, have := "world", req; want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
}

func TestDesc_Middleware(t *testing.T) {
	d := testDesc()
	d.SetMiddleware("hello", func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			innerResp, err := inner(ctx, req)
			if err != nil {
				return
			}
			resp = fmt.Sprintf("hello %s", innerResp)
			return
		}
	})

	// test retrieve and use the middleware
	ep := d.GetMiddleware("hello")(func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		resp = req
		return
	})
	resp, err := ep(nil, "world")
	if err != nil {
		t.Errorf("error: %#v", err.Error())
	}
	if want, have := "hello world", resp; want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
}
