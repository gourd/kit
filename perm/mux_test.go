package perm_test

import (
	"github.com/gourd/kit/perm"

	"golang.org/x/net/context"
	"testing"
)

func TestDefaultMux(t *testing.T) {
	// test if default mux implements mux
	var m perm.Mux = perm.NewMux()
	_ = m
}

func TestMuxFoundFunc(t *testing.T) {
	m := perm.NewMux()
	m.HandleFunc("access something", func(ctx context.Context, perm string, info ...interface{}) error {
		return nil
	})
	if err := m.Allow(nil, "access something"); err != nil {
		t.Errorf("Unexpected error. Failed to obtain handler for permission")
	}
}

func TestMuxFoundMux(t *testing.T) {

	// child mux
	m1 := perm.NewMux()
	m1.HandleFunc("access something", func(ctx context.Context, perm string, info ...interface{}) error {
		return nil
	})

	// parent mux
	m2 := perm.NewMux()
	m2.Handle("access something", m1)

	// test parent mux
	if err := m2.Allow(nil, "access something"); err != nil {
		t.Errorf("Unexpected error. Failed to obtain handler for permission")
	}
}

func TestMuxNotFound(t *testing.T) {
	m := perm.NewMux()
	err := m.Allow(nil, "access something")
	if err != perm.HandlerNotFound {
		t.Errorf("Error is not of expected type. Expecting perm.HandlerNotFound by get %#v", err)
	}
}
