package oauth2_test

import (
	"github.com/gourd/kit/oauth2"

	"testing"
)

func testScopes() *oauth2.Scopes {
	return &oauth2.Scopes{"hello", "world"}
}

// test if BasicScopes implements Scopes
func TestScopesHas(t *testing.T) {
	s := testScopes()
	if !s.Has("hello") {
		t.Errorf("Cannot find existing scope")
	}
	if !s.Has("world") {
		t.Errorf("Cannot find existing scope")
	}
	if s.Has("foo") {
		t.Errorf("Found non-existing scope")
	}
}

// test if BasicScopes implements Scopes
func TestScopesHasAny(t *testing.T) {
	s := testScopes()
	if !s.HasAny("foo", "hello") {
		t.Errorf("Cannot find existing scope")
	}
	if s.HasAny("foo", "bar") {
		t.Errorf("Found non-existing scope")
	}
}
