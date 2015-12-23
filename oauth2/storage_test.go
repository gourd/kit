package oauth2_test

import (
	"fmt"
	"testing"

	"github.com/RangelReale/osin"
	"github.com/gourd/kit/oauth2"
	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
)

func TestStorage_AuthorizeData(t *testing.T) {

	// define test db
	getContext := func() context.Context {
		factory := store.NewFactory()
		factory.SetSource(store.DefaultSrc, defaultTestSrc())
		factory.Set(oauth2.KeyAccess, store.DefaultSrc, oauth2.AccessDataStoreProvider)
		factory.Set(oauth2.KeyAuth, store.DefaultSrc, oauth2.AuthorizeDataStoreProvider)
		factory.Set(oauth2.KeyClient, store.DefaultSrc, oauth2.ClientStoreProvider)
		factory.Set(oauth2.KeyUser, store.DefaultSrc, oauth2.UserStoreProvider)
		return store.WithFactory(context.Background(), factory)
	}

	// create dummy Client and user
	ctx := getContext()
	defer store.CloseAllIn(ctx)
	storage := &oauth2.Storage{}
	storage.SetContext(ctx)

	c, u := createStoreDummies(ctx, "password", "http://foobar.com/redirect")
	ad := dummyNewAuth(c, u)
	storage.SaveAuthorize(ad.ToOsin())

	// load the osin.AuthorizeData form store
	oad, err := storage.LoadAuthorize(ad.Code)
	if err != nil {
		t.Errorf("error: %#v", err.Error())
	}

	// Test if loaded Client equals to client in original one
	if want, have := c.GetId(), oad.Client.GetId(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := c.GetRedirectUri(), oad.Client.GetRedirectUri(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := c.GetSecret(), oad.Client.GetSecret(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// Test if UserData equals to original one
	if u1, u2 := ad.UserData.(*oauth2.User), oad.UserData.(*oauth2.User); true {
	} else if want, have := u1.ID, u2.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := u1.Email, u2.Email; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := u1.Name, u2.Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := u1.Password, u2.Password; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := u1.Created, u2.Created; want.Unix() != have.Unix() {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := u1.Updated, u2.Updated; want.Unix() != have.Unix() {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestStorage_AccessData(t *testing.T) {

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
		if access == nil {
			err = fmt.Errorf("unexpected nil *oauth2.AccessData")
			return
		}
		if oaccess == nil {
			err = fmt.Errorf("unexpected nil *osin.AccessData")
			return
		}

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
		if v1, v2 := access.RedirectURI, oaccess.RedirectUri; v1 != v2 {
			err = fmt.Errorf("RedirectUri mismatch.\n*oauth2.RedirectUri=%#v, *osin.RedirectUri=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.CreatedAt, oaccess.CreatedAt; v1.Unix() != v2.Unix() {
			err = fmt.Errorf("CreatedAt mismatch.\n*oauth2.CreatedAt=%#v, *osin.CreatedAt=%#v",
				v1, v2)
			return
		}
		return
	}

	// define test db
	getContext := func() context.Context {
		factory := store.NewFactory()
		factory.SetSource(store.DefaultSrc, defaultTestSrc())
		factory.Set(oauth2.KeyAccess, store.DefaultSrc, oauth2.AccessDataStoreProvider)
		factory.Set(oauth2.KeyAuth, store.DefaultSrc, oauth2.AuthorizeDataStoreProvider)
		factory.Set(oauth2.KeyClient, store.DefaultSrc, oauth2.ClientStoreProvider)
		factory.Set(oauth2.KeyUser, store.DefaultSrc, oauth2.UserStoreProvider)
		return store.WithFactory(context.Background(), factory)
	}

	// create dummy Client and user
	ctx := getContext()
	defer store.CloseAllIn(ctx)
	storage := &oauth2.Storage{}
	storage.SetContext(ctx)

	c, u := createStoreDummies(ctx, "password", "http://foobar.com/redirect")
	ad := dummyNewAuth(c, u)
	access1 := dummyNewAccess(c, u, ad, nil)

	storage.SaveAccess(access1.ToOsin())
	oaccess1, err := storage.LoadAccess(access1.AccessToken)
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	if err := accessMatch(access1, oaccess1); err != nil {
		t.Errorf("access1 != oaccess1, err = %#v", err.Error())
	}
	if err := authEqual(access1.AuthorizeData.ToOsin(), oaccess1.AuthorizeData); err != nil {
		t.Errorf("access1.AuthorizeData != oaccess1.AuthorizeData, err = %#v", err.Error())
		t.Logf("\naccess1.AuthorizeData=%#v\noaccess1.AuthorizeData=%#v",
			access1.AuthorizeData, oaccess1.AuthorizeData)
	}

	access2 := dummyNewAccess(c, u, ad, access1)
	if access2.AccessData == nil {
		t.Error("unexpected nil value")
	} else if access2.ToOsin().AccessData == nil {
		t.Error("unexpected nil value")
	}

	storage.SaveAccess(access2.ToOsin())
	oaccess2, err := storage.LoadAccess(access2.AccessToken)
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	if err := accessMatch(access2, oaccess2); err != nil {
		t.Errorf("access2 != oaccess2, err = %#v", err.Error())
	}
	if err := authEqual(access2.AuthorizeData.ToOsin(), oaccess2.AuthorizeData); err != nil {
		t.Errorf("access2.AuthorizeData != oaccess2.AuthorizeData, err = %#v", err.Error())
		t.Logf("\naccess2.AuthorizeData=%#v\noaccess2.AuthorizeData=%#v",
			access2.AuthorizeData, oaccess2.AuthorizeData)
	}
	if err := accessMatch(access1, oaccess2.AccessData); err != nil {
		t.Errorf("access1 != oaccess2.AccessData, err = %#v", err.Error())
	}

}

func TestStorage_AccessData_Refresh(t *testing.T) {

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
		if access == nil {
			err = fmt.Errorf("unexpected nil *oauth2.AccessData")
			return
		}
		if oaccess == nil {
			err = fmt.Errorf("unexpected nil *osin.AccessData")
			return
		}

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
		if v1, v2 := access.RedirectURI, oaccess.RedirectUri; v1 != v2 {
			err = fmt.Errorf("RedirectUri mismatch.\n*oauth2.RedirectUri=%#v, *osin.RedirectUri=%#v",
				v1, v2)
			return
		}
		if v1, v2 := access.CreatedAt, oaccess.CreatedAt; v1.Unix() != v2.Unix() {
			err = fmt.Errorf("CreatedAt mismatch.\n*oauth2.CreatedAt=%#v, *osin.CreatedAt=%#v",
				v1, v2)
			return
		}
		return
	}

	// define test db
	getContext := func() context.Context {
		factory := store.NewFactory()
		factory.SetSource(store.DefaultSrc, defaultTestSrc())
		factory.Set(oauth2.KeyAccess, store.DefaultSrc, oauth2.AccessDataStoreProvider)
		factory.Set(oauth2.KeyAuth, store.DefaultSrc, oauth2.AuthorizeDataStoreProvider)
		factory.Set(oauth2.KeyClient, store.DefaultSrc, oauth2.ClientStoreProvider)
		factory.Set(oauth2.KeyUser, store.DefaultSrc, oauth2.UserStoreProvider)
		return store.WithFactory(context.Background(), factory)
	}

	// create dummy Client and user
	ctx := getContext()
	defer store.CloseAllIn(ctx)
	storage := &oauth2.Storage{}
	storage.SetContext(ctx)

	c, u := createStoreDummies(ctx, "password", "http://foobar.com/redirect")
	ad := dummyNewAuth(c, u)
	access1 := dummyNewAccess(c, u, ad, nil)

	storage.SaveAccess(access1.ToOsin())
	oaccess1, err := storage.LoadRefresh(access1.RefreshToken)
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	if err := accessMatch(access1, oaccess1); err != nil {
		t.Errorf("access1 != oaccess1, err = %#v", err.Error())
	}
	if err := authEqual(access1.AuthorizeData.ToOsin(), oaccess1.AuthorizeData); err != nil {
		t.Errorf("access1.AuthorizeData != oaccess1.AuthorizeData, err = %#v", err.Error())
		t.Logf("\naccess1.AuthorizeData=%#v\noaccess1.AuthorizeData=%#v",
			access1.AuthorizeData, oaccess1.AuthorizeData)
	}

	access2 := dummyNewAccess(c, u, ad, access1)
	if access2.AccessData == nil {
		t.Error("unexpected nil value")
	} else if access2.ToOsin().AccessData == nil {
		t.Error("unexpected nil value")
	}

	storage.SaveAccess(access2.ToOsin())
	oaccess2, err := storage.LoadRefresh(access2.RefreshToken)
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	if err := accessMatch(access2, oaccess2); err != nil {
		t.Errorf("access2 != oaccess2, err = %#v", err.Error())
	}
	if err := authEqual(access2.AuthorizeData.ToOsin(), oaccess2.AuthorizeData); err != nil {
		t.Errorf("access2.AuthorizeData != oaccess2.AuthorizeData, err = %#v", err.Error())
		t.Logf("\naccess2.AuthorizeData=%#v\noaccess2.AuthorizeData=%#v",
			access2.AuthorizeData, oaccess2.AuthorizeData)
	}
	if err := accessMatch(access1, oaccess2.AccessData); err != nil {
		t.Errorf("access1 != oaccess2.AccessData, err = %#v", err.Error())
	}

}
