package http_test

import (
	"fmt"
	"net/http"
	"path"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httpservice "github.com/gourd/kit/service/http"

	"testing"
)

func TestDesc_DecodeFunc(t *testing.T) {
	base, sing, plur := "/some/path", "ball", "balls"
	n := httpservice.NewNoun(sing, plur)
	p := httpservice.NewPaths(base, n, func(name string, noun httpservice.Noun) string {
		switch name {
		case "create":
			return noun.Plural()
		case "update":
			return path.Join(noun.Singular(), "{id}")
		}
		return ""
	})
	d := httpservice.NewDesc(p)
	d.SetDecodeFunc("hello", func(r *http.Request) (request interface{}, err error) {
		request = "world"
		return
	})

	// test retrieve and use the decoder
	dec := d.GetDecodeFunc("hello")
	req, _ := dec(nil)
	if want, have := "world", req; want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
}

func TestDesc_Middleware(t *testing.T) {
	base, sing, plur := "/some/path", "ball", "balls"
	n := httpservice.NewNoun(sing, plur)
	p := httpservice.NewPaths(base, n, func(name string, noun httpservice.Noun) string {
		switch name {
		case "create":
			return noun.Plural()
		case "update":
			return path.Join(noun.Singular(), "{id}")
		}
		return ""
	})
	d := httpservice.NewDesc(p)
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
