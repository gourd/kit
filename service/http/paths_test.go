package httpservice_test

import (
	httpservice "github.com/gourd/kit/service/http"

	"testing"
)

func TestPaths(t *testing.T) {
	base, sing, plur := "/some/path", "ball", "balls"
	n := httpservice.NewNoun(sing, plur)
	p := httpservice.NewPaths(base, n, "someid")

	if want, have := "/some/path/ball/someid", p.Singular(); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
	if want, have := "/some/path/balls", p.Plural(); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}

}
