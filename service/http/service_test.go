package httpservice_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/endpoint"
	"github.com/gourd/kit/context"
	httpservice "github.com/gourd/kit/service/http"
	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
)

func testServiceSuit(path, resultKey string) (s *httpservice.Service, mware endpoint.Middleware) {

	// dummy service
	s = httpservice.NewJSONService(path, func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response = map[string]interface{}{
			resultKey: request,
		}
		return
	})
	s.DecodeFunc = func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		request = "world"
		return
	}

	// dummy middleware
	mware = func(inner endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = inner(ctx, request)
			if err != nil {
				return
			}

			vmap := response.(map[string]interface{})
			vmap[resultKey] = fmt.Sprintf("hello %s", vmap[resultKey])
			return
		}
	}
	return
}

func TestService_Normal(t *testing.T) {
	resultKey := "result"

	// create handler with service
	s, mware := testServiceSuit("/foo/bar", resultKey)
	s.Middlewares.Add(httpservice.MWInner, mware)
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
	s.DecodeFunc = func(ctx context.Context, r *http.Request) (request interface{}, err error) {
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

func TestNewJSONService(t *testing.T) {

	str := `{"hello": "world"}`
	r1 := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(str)),
	}

	// create handler with service
	s := httpservice.NewJSONService("/foo/bar", func(ctx context.Context, request interface{}) (response interface{}, err error) {

		r := gourdctx.HTTPRequest(ctx)
		if want, have := r1, r; want != have {
			t.Errorf("expected %#v, got %#v", want, have)
		}

		// retrieve the request
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body error: %s", err)
			return
		}

		if want, have := str, string(b); want != have {
			err = fmt.Errorf(`expected "%s", got "%s"`, want, have)
			t.Error(err.Error())
		}
		response = map[string]string{
			"result": "success",
		}
		return
	})
	s.DecodeFunc = func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		return
	}
	h := s.Handler()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r1)

	b, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	t.Logf("result: %s", b)
}

func TestServiceSlice_Sort(t *testing.T) {
	slice := httpservice.ServiceSlice{
		{
			Weight: 5,
		},
		{
			Weight: -1,
		},
		{
			Weight: 2,
		},
		{
			Weight: 1,
		},
		{
			Weight: 0,
		},
		{
			Weight: 4,
		},
	}
	slice.Sort()

	expected := []int{-1, 0, 1, 2, 4, 5}
	for i, service := range slice {
		if want, have := expected[i], service.Weight; want != have {
			t.Errorf("[%d] exptected %d, got %d", i, want, have)
		}
	}
}

func TestServices_Each(t *testing.T) {
	services := httpservice.Services{
		"service1": {
			Weight: 5,
		},
		"service2": {
			Weight: -1,
		},
		"service3": {
			Weight: 2,
		},
		"service4": {
			Weight: 1,
		},
		"service5": {
			Weight: 0,
		},
		"service6": {
			Weight: 4,
		},
	}

	i := 0
	expected := []int{-1, 0, 1, 2, 4, 5}
	for service := range services.Each() {
		if want, have := expected[i], service.Weight; want != have {
			t.Errorf("[%d] exptected %d, got %d", i, want, have)
		}
		i++
	}
}

func TestServices_Route(t *testing.T) {
	resultKey := "result"
	servicePath := "/foo/bar"

	// dummy service to test
	var services httpservice.Services = make(map[string]*httpservice.Service)
	s, mware := testServiceSuit(servicePath, resultKey)
	s.Middlewares.Add(httpservice.MWInner, mware)
	services["example"] = s

	m := http.NewServeMux()

	// lazy implementation routerfunc for ServeMux
	// that simply skip method handling (don't try this at home)
	rtfn := func(m *http.ServeMux) httpservice.RouterFunc {
		return func(path string, methods []string, h http.Handler) error {
			m.Handle(path, h)
			return nil
		}
	}(m)

	// route dumy service
	services.Route(rtfn)

	// variables for decoding
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", servicePath, nil)
	m.ServeHTTP(w, req)
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

func TestServices_Patch(t *testing.T) {
	resultKey := "result"
	servicePath := "/foo/bar"

	// dummy service to test
	var services httpservice.Services = make(map[string]*httpservice.Service)
	s, mware := testServiceSuit(servicePath, resultKey)
	services["example"] = s

	m := http.NewServeMux()

	// lazy implementation routerfunc for ServeMux
	// that simply skip method handling (don't try this at home)
	rtfn := func(m *http.ServeMux) httpservice.RouterFunc {
		return func(path string, methods []string, h http.Handler) error {
			m.Handle(path, h)
			return nil
		}
	}(m)

	// add middleware with patch
	patch := func(services httpservice.Services) httpservice.Services {
		for name := range services {
			services[name].Middlewares.Add(httpservice.MWInner, mware)
		}
		return services
	}

	// route dumy service
	services.Patch(patch)
	services.Route(rtfn)

	// variables for decoding
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", servicePath, nil)
	m.ServeHTTP(w, req)
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
