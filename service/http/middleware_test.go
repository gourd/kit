package http_test

import (
	"fmt"

	"github.com/go-kit/kit/endpoint"
	httpservice "github.com/gourd/kit/service/http"
	"golang.org/x/net/context"

	"testing"
)

func TestMiddleware(t *testing.T) {
	ms := &httpservice.Middlewares{}

	m1 := func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = inner(ctx, request)
			response = fmt.Sprintf("m1(%s)", response)
			return
		}
	}
	m2 := func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = inner(ctx, request)
			response = fmt.Sprintf("m2(%s)", response)
			return
		}
	}
	m3 := func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = inner(ctx, request)
			response = fmt.Sprintf("m3(%s)", response)
			return
		}
	}

	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response = fmt.Sprintf("ep(%s)", request)
		return
	}
	if resp, _ := ep(nil, "hello"); resp == nil {
		t.Errorf("unexpected nil")
	} else if want, have := "ep(hello)", resp; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	ms.Add(1, m1)
	ms.Add(5, m2)
	ms.Add(-1, m3)

	if mfinal := ms.Chain(); mfinal == nil {
		t.Errorf("Chain() returned nil")
	} else if resp, _ := mfinal(ep)(nil, "hello"); resp == nil {
		t.Errorf("unexpected nil")
	} else if want, have := "m3(m1(m2(ep(hello))))", resp; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}
