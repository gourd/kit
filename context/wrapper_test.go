package context_test

import (
	"testing"

	gcontext "github.com/gourd/kit/context"
	"golang.org/x/net/context"
)

func TestContextValue(t *testing.T) {
	// ensure the wrapper implements context.Context
	key := "foo"
	val := "bar"
	ctx0 := gcontext.New(nil)
	ctx1 := context.WithValue(ctx0, key, val)

	if ctx0 == ctx1 {
		t.Error("ctx1 should be a clone of ctx1")
	}

	if ctx0.Value(key) == val {
		t.Error("WithValue should not be %#v", val)
	}
}
