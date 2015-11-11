package http_test

import (
	"path"

	httpservice "github.com/gourd/kit/service/http"

	"testing"
)

func TestPaths(t *testing.T) {
	base, sing, plur := "/some/path", "ball", "balls"
	n := httpservice.NewNoun(sing, plur)
	p := httpservice.NewPaths(base, n, func(name string, noun httpservice.Noun) string {
		switch name {
		case "foo":
			return path.Join(noun.Plural(), "bar")
		case "hello":
			return path.Join(noun.Singular(), "world")
		}
		return ""
	})

	if want, have := "/some/path/balls/bar", p.Path("foo"); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
	if want, have := "/some/path/ball/world", p.Path("hello"); want != have {
		t.Errorf("expected: %#v, got: %#v", want, have)
	}
}
