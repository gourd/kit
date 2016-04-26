package httpservice_test

import (
	httpservice "github.com/gourd/kit/service/http"

	"testing"
)

func TestNoun(t *testing.T) {
	sing, plur := "ball", "balls"
	n := httpservice.NewNoun(sing, plur)
	if want, have := sing, n.Singular(); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
	if want, have := plur, n.Plural(); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
}
