package oauth2_test

import (
	"github.com/gourd/kit/oauth2"

	"math/rand"
	"testing"
)

func dummyNewUser(password string) *oauth2.User {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	u := &oauth2.User{
		Username: randSeq(10),
	}
	u.Password = u.Hash(password)
	return u
}

func TestUser(t *testing.T) {
	var u oauth2.OAuth2User = &oauth2.User{}
	_ = u
}

func TestMeta(t *testing.T) {
	u := &oauth2.User{}
	u.MetaJSON = `{"hello": ["world 1", "world 2"]}`
	m := u.Meta()

	// inspect outer
	if want, have := 1, len(m); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}
	mHello, ok := m["hello"]
	if !ok {
		t.Errorf("unable to find %#v in meta", "hello")
		return
	}

	// inspect inner
	if want, have := 2, len(mHello); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}

	if want, have := "world 1", mHello[0]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "world 2", mHello[1]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestAddMeta(t *testing.T) {
	u := &oauth2.User{}
	u.AddMeta("hello", "world 1")
	u.AddMeta("hello", "world 2")
	m := u.Meta()

	// inspect outer
	if want, have := 1, len(m); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}
	mHello, ok := m["hello"]
	if !ok {
		t.Errorf("unable to find %#v in meta", "hello")
		return
	}

	// inspect inner
	if want, have := 2, len(mHello); want != have {
		t.Logf("result json: %#v", u.MetaJSON)
		t.Logf("result mHello: %#v", mHello)
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}

	if want, have := "world 1", mHello[0]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "world 2", mHello[1]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
