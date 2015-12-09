package oauth2_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	"github.com/gourd/kit/oauth2"
)

func dummyNewAuth(client *oauth2.Client, user *oauth2.User) *oauth2.AuthorizeData {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	ad := &oauth2.AuthorizeData{
		Id:          randSeq(10),
		ClientId:    client.Id,
		Client:      client,
		Code:        randSeq(10),
		ExpiresIn:   rand.Int31(),
		Scope:       randSeq(10),
		RedirectUri: client.RedirectUri + "/" + randSeq(10),
		State:       randSeq(10),
		CreatedAt:   time.Now(),
		UserId:      user.Id,
		UserData:    user,
	}
	return ad
}

func dummyNewOsinAuth(client *oauth2.Client, user *oauth2.User) *osin.AuthorizeData {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	ad := &osin.AuthorizeData{
		Client:      client,
		Code:        randSeq(10),
		ExpiresIn:   rand.Int31(),
		Scope:       randSeq(10),
		RedirectUri: client.RedirectUri + "/" + randSeq(10),
		State:       randSeq(10),
		CreatedAt:   time.Now(),
		UserData:    user,
	}
	return ad
}

func TestAuth_ToOsin(t *testing.T) {

	// creates dummy client and user directly from the stores
	createDummies := func(password, redirect string) (*oauth2.Client, *oauth2.User) {

		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		randSeq := func(n int) string {
			b := make([]rune, n)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		}

		u := dummyNewUser(password)
		c := dummyNewClient(redirect)
		u.Id = randSeq(20)
		c.Id = randSeq(20)
		return c, u
	}

	c, u := createDummies("password", "http://foobar.com/redirect")
	ad := dummyNewAuth(c, u)
	oad := ad.ToOsin()

	if want, have := ad.ClientId, oad.Client.GetId(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := ad.Client.GetId(), oad.Client.GetId(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := ad.Code, oad.Code; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := ad.ExpiresIn, oad.ExpiresIn; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := ad.Scope, oad.Scope; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := ad.RedirectUri, oad.RedirectUri; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := ad.State, oad.State; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	u1 := ad.UserData.(*oauth2.User)

	if u2, ok := oad.UserData.(*oauth2.User); !ok {
		t.Errorf(".UserData is not *oauth2.User, but %#v",
			oad.UserData)
	} else if want, have := u2.Id, u1.Id; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestAuth_ReadOsin(t *testing.T) {

	// creates dummy client and user directly from the stores
	createDummies := func(password, redirect string) (*oauth2.Client, *oauth2.User) {

		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		randSeq := func(n int) string {
			b := make([]rune, n)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		}

		u := dummyNewUser(password)
		c := dummyNewClient(redirect)
		u.Id = randSeq(20)
		c.Id = randSeq(20)
		return c, u
	}

	c, u := createDummies("password", "http://foobar.com/redirect")
	oad := dummyNewOsinAuth(c, u)
	ad := &oauth2.AuthorizeData{}

	ad.ReadOsin(oad)

	if want, have := oad.Client.GetId(), ad.ClientId; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := oad.Client.GetId(), ad.Client.GetId(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := oad.Code, ad.Code; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := oad.ExpiresIn, ad.ExpiresIn; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := oad.Scope, ad.Scope; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := oad.RedirectUri, ad.RedirectUri; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := oad.State, ad.State; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	u2 := ad.UserData.(*oauth2.User)

	if u1, ok := oad.UserData.(*oauth2.User); !ok {
		t.Errorf(".UserData is not *oauth2.User, but %#v",
			oad.UserData)
	} else if want, have := u2.Id, u1.Id; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}
