package httpservice_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gourd/kit/service/http"
	"golang.org/x/net/context"
)

func TestProvideJSONDecoder(t *testing.T) {
	// test if httpservice.ProvideJSONDecoder implements
	// httptransport.RequestFunc
	var v httptransport.RequestFunc = httpservice.ProvideJSONDecoder
	_ = v
}

func TestDecoder(t *testing.T) {

	// mock request context
	body := ioutil.NopCloser(strings.NewReader(`{"hello": "world"}`))
	ctx := context.Background()
	r := &http.Request{Body: body}
	ctx = httpservice.ProvideJSONDecoder(ctx, r)

	// decode the context into entity struct
	entity := struct {
		Hello string `json:"hello"`
	}{}
	dec, ok := httpservice.DecoderFrom(ctx)

	if !ok {
		t.Errorf("expected ok to be true")
	}

	// test decoding
	if err := dec.Decode(&entity); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// test decoded result
	if want, have := "world", entity.Hello; want != have {
		t.Errorf("exptected %#v, got %#v", want, have)
	}
}

func TestDecoder_Nil(t *testing.T) {

	// mock request context
	ctx := context.Background()
	r := &http.Request{}
	ctx = httpservice.ProvideJSONDecoder(ctx, r)
	dec, ok := httpservice.DecoderFrom(ctx)

	if ok {
		t.Error("unexpected ok")
	}
	if dec != nil {
		t.Errorf("unexpected decoder value %#v", dec)
	}
}
