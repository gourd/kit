package store_test

import (
	"github.com/gourd/kit/store"

	"net/http"
	"testing"
)

func TestProvideFunc(t *testing.T) {
	var f store.ProvideFunc = func(r *http.Request) (store.Store, error) {
		return nil, nil
	}
	var p store.Provider = f
	_ = p
	t.Log("ProvideFunc implements Provider")
}

func TestProviderStore(t *testing.T) {

	// define a database source
	store.Providers.DefineFunc(
		"dummy",
		func(r *http.Request) (store.Store, error) {
			return nil, nil
		},
	)

	// test creating the new database
	p, err := store.Providers.Get("dummy")
	if err != nil {
		t.Error(err.Error())
	}

	// test the provider
	r := &http.Request{}
	s, err := p.Store(r)
	if err != nil {
		t.Error(err.Error())
	} else if s != nil {
		t.Errorf(
			"Unexpected service provider result. Expecting nil but get %#v", s)
	}

	t.Log("Provider and Providers routine works")

}

func TestProviderMustStore(t *testing.T) {

	// define a database source
	store.Providers.DefineFunc(
		"dummy",
		func(r *http.Request) (store.Store, error) {
			return nil, nil
		},
	)
	r := &http.Request{}

	// test creating the new database
	s := store.Providers.MustStore(r, "dummy")
	if s != nil {
		t.Errorf(
			"Unexpected service provider result. Expecting nil but get %#v", s)
	}

	t.Log("Provider and Providers routine works")

}
