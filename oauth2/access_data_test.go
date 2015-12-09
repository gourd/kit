package oauth2_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	"github.com/gourd/kit/oauth2"
)

func dummyNewAccess(client *oauth2.Client, user *oauth2.User,
	ad *oauth2.AuthorizeData, prev *oauth2.AccessData) *oauth2.AccessData {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	access := &oauth2.AccessData{
		Id:            randSeq(10),
		ClientId:      client.Id,
		Client:        client,
		AuthorizeData: ad,
		AccessToken:   randSeq(10),
		RefreshToken:  randSeq(10),
		RedirectUri:   client.RedirectUri + "/" + randSeq(10),
		CreatedAt:     time.Now(),
		UserId:        user.Id,
		UserData:      user,
	}

	if prev != nil {
		access.AccessData = prev
	}

	return access
}

func TestAccess_ToOsin(t *testing.T) {

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

	authEqual := func(a1, a2 *osin.AuthorizeData) (err error) {
		if a1 == nil {
			err = fmt.Errorf("unexpected nil a1")
			return
		}
		if a2 == nil {
			err = fmt.Errorf("unexpected nil a2")
			return
		}

		if v1, v2 := a1.Code, a2.Code; v1 != v2 {
			err = fmt.Errorf("Code not equal. %#v != %#v", v1, v2)
			return
		}
		if v1, v2 := a1.ExpiresIn, a2.ExpiresIn; v1 != v2 {
			err = fmt.Errorf("Code not equal. %#v != %#v", v1, v2)
			return
		}
		return
	}

	accessMatch := func(access *oauth2.AccessData, oaccess *osin.AccessData) (err error) {
		if v1, v2 := access.AccessToken, oaccess.AccessToken; v1 != v2 {
			err = fmt.Errorf("AccessToken mismatch.\n*oauth2.AccessData=%#v, *osin.AccessData=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.RefreshToken, oaccess.RefreshToken; v1 != v2 {
			err = fmt.Errorf("RefreshToken mismatch.\n*oauth2.RefreshData=%#v, *osin.RefreshData=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.ExpiresIn, oaccess.ExpiresIn; v1 != v2 {
			err = fmt.Errorf("ExpiresIn mismatch.\n*oauth2.ExpiresIn=%#v, *osin.ExpiresIn=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.Scope, oaccess.Scope; v1 != v2 {
			err = fmt.Errorf("Scope mismatch.\n*oauth2.Scope=%#v, *osin.Scope=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.RedirectUri, oaccess.RedirectUri; v1 != v2 {
			err = fmt.Errorf("RedirectUri mismatch.\n*oauth2.RedirectUri=%#v, *osin.RedirectUri=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.CreatedAt, oaccess.CreatedAt; !v1.Equal(v2) {
			err = fmt.Errorf("CreatedAt mismatch.\n*oauth2.CreatedAt=%#v, *osin.CreatedAt=%#v",
				v1, v2)
			return
		}
		return
	}

	// create some dummies
	c, u := createDummies("password", "http://foobar.com/redirect")
	ad := dummyNewAuth(c, u)

	access1 := dummyNewAccess(c, u, ad, nil)
	oaccess1 := access1.ToOsin()
	if err := accessMatch(access1, oaccess1); err != nil {
		t.Errorf("access1 and oaccess1 mismatch, %#v", err.Error())
	}

	if want, have := access1.Client, oaccess1.Client; want != have {
		t.Errorf("\nexpected %#v\ngot      %#v", want, have)
	} else if want == nil {
		t.Errorf("unexpected nil value")
	}

	if want, have := access1.AuthorizeData.ToOsin(), oaccess1.AuthorizeData; true {
	} else if err := authEqual(want, have); err != nil {
		t.Errorf("want != have, err = %#v", err.Error())
	}

	access2 := dummyNewAccess(c, u, ad, access1)
	oaccess2 := access2.ToOsin()
	if err := accessMatch(access2, oaccess2); err != nil {
		t.Errorf("access2 and oaccess2 mismatch, %#v", err.Error())
	}

	if want, have := access2.Client, oaccess2.Client; want != have {
		t.Errorf("\nexpected %#v\ngot      %#v", want, have)
	} else if want == nil {
		t.Errorf("unexpected nil value")
	}

	if want, have := access2.AuthorizeData.ToOsin(), oaccess2.AuthorizeData; true {
	} else if err := authEqual(want, have); err != nil {
		t.Errorf("want != have, err = %#v", err.Error())
	}

	if want, have := access1, access2.AccessData.ToOsin(); true {
	} else if err := accessMatch(want, have); err != nil {
		t.Errorf("want != have, err = %#v", err.Error())
	}

}

func TestAccess_FromOsin(t *testing.T) {

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

	// create dummy osin.AccessData
	createOsinAccess := func() *osin.AccessData {
		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		randSeq := func(n int) string {
			b := make([]rune, n)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		}
		return &osin.AccessData{
			AccessToken:  randSeq(20),
			RefreshToken: randSeq(20),
			ExpiresIn:    rand.Int31(),
			Scope:        randSeq(20),
			RedirectUri:  randSeq(20),
			CreatedAt:    time.Now(),
		}
	}

	accessMatch := func(access *oauth2.AccessData, oaccess *osin.AccessData) (err error) {
		if v1, v2 := access.AccessToken, oaccess.AccessToken; v1 != v2 {
			err = fmt.Errorf("AccessToken mismatch.\n*oauth2.AccessData=%#v, *osin.AccessData=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.RefreshToken, oaccess.RefreshToken; v1 != v2 {
			err = fmt.Errorf("RefreshToken mismatch.\n*oauth2.RefreshData=%#v, *osin.RefreshData=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.ExpiresIn, oaccess.ExpiresIn; v1 != v2 {
			err = fmt.Errorf("ExpiresIn mismatch.\n*oauth2.ExpiresIn=%#v, *osin.ExpiresIn=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.Scope, oaccess.Scope; v1 != v2 {
			err = fmt.Errorf("Scope mismatch.\n*oauth2.Scope=%#v, *osin.Scope=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.RedirectUri, oaccess.RedirectUri; v1 != v2 {
			err = fmt.Errorf("RedirectUri mismatch.\n*oauth2.RedirectUri=%#v, *osin.RedirectUri=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.CreatedAt, oaccess.CreatedAt; !v1.Equal(v2) {
			err = fmt.Errorf("CreatedAt mismatch.\n*oauth2.CreatedAt=%#v, *osin.CreatedAt=%#v",
				v1, v2)
			return
		}
		return
	}

	c, u := createDummies("password", "http://foobar.com")
	oad := dummyNewAuth(c, u).ToOsin()
	oaccess1 := createOsinAccess()
	oaccess1.Client = c
	oaccess1.UserData = u
	oaccess1.AuthorizeData = oad

	access1 := &oauth2.AccessData{}
	access1.ReadOsin(oaccess1)
	if err := accessMatch(access1, oaccess1); err != nil {
		t.Errorf("access1 and oaccess1 mismatch, %#v", err.Error())
	}

	oaccess2 := createOsinAccess()
	oaccess2.Client = c
	oaccess2.UserData = u
	oaccess2.AuthorizeData = oad
	oaccess2.AccessData = oaccess1

	access2 := &oauth2.AccessData{}
	access2.ReadOsin(oaccess2)
	if err := accessMatch(access2, oaccess2); err != nil {
		t.Errorf("access2 and oaccess2 mismatch, %#v", err.Error())
	}

}
