package oauth2_test

import (
	"testing"

	"github.com/gourd/kit/store"
	"golang.org/x/net/context"

	"github.com/gourd/kit/oauth2"
)

func TestStorage(t *testing.T) {

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

	_ = ad

	//storage.SaveAuthorize(ad)

	// TODO
	// ----
	// 2. create dummy AuthorizeData, test load
	// 3. try to run data process in Token request,
	//    test if the AuthorizeData saved correctly with JSON.
	// 4. try to refresh token, test if the AuthorizeData and AccessData
	//    saved correctly with JSON
}
