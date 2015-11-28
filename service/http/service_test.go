package http_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/endpoint"
	httpservice "github.com/gourd/kit/service/http"
	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
)

func TestService_Normal(t *testing.T) {
	resultKey := "result"

	// create handler with service
	s := httpservice.NewJSONService("/foo/bar", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response = map[string]interface{}{
			resultKey: request,
		}
		return
	})
	s.Middlewares.Add(httpservice.MWInner, func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = inner(ctx, request)
			if err != nil {
				return
			}

			vmap := response.(map[string]interface{})
			vmap[resultKey] = fmt.Sprintf("hello %s", vmap[resultKey])
			return
		}
	})
	s.DecodeFunc = func(r *http.Request) (request interface{}, err error) {
		request = "world"
		return
	}
	h := s.Handler()

	// variables for decoding
	w := httptest.NewRecorder()
	h.ServeHTTP(w, nil)
	vmap := make(map[string]interface{})

	// try decoding
	dec := json.NewDecoder(w.Body)
	if err := dec.Decode(&vmap); err != nil {
		t.Errorf("error decoding response: %#v", err.Error())
	} else if result, ok := vmap[resultKey]; !ok {
		t.Errorf("got no %#v in vmap: %#v", resultKey, vmap)
	} else if want, have := "hello world", result; want != have {
		t.Errorf("expect: %#v, got: %#v", want, have)
	}
}

func TestService_Error(t *testing.T) {

	// create handler with service
	s := httpservice.NewJSONService("/foo/bar", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = store.Error(50123, "hello error")
		return
	})
	s.DecodeFunc = func(r *http.Request) (request interface{}, err error) {
		request = "hello world"
		return
	}
	h := s.Handler()

	// variables for decoding
	w := httptest.NewRecorder()
	h.ServeHTTP(w, nil)
	serr := &store.StoreError{}

	// try decoding
	dec := json.NewDecoder(w.Body)
	dec.Decode(serr)
	if err := dec.Decode(&serr); err != io.EOF {
		t.Errorf("error decoding response: %#v", err.Error())
	} else if want, have := "hello error", serr.ClientMsg; want != have {
		t.Errorf("expect: %#v, got: %#v", want, have)
	}
	t.Logf("err: %#v", serr)
}
