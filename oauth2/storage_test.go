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

	// creates dummy client and user directly from the stores
	createDummies := func(ctx context.Context, password, redirect string) (*oauth2.Client, *oauth2.User) {

		type tempKey int
		const (
			testDB tempKey = iota
		)

		// generate dummy user
		us, err := store.Get(ctx, oauth2.KeyUser)
		if err != nil {
			panic(err)
		}
		u := dummyNewUser(password)
		err = us.Create(store.NewConds(), u)
		if err != nil {
			panic(err)
		}

		// get related dummy client
		cs, err := store.Get(ctx, oauth2.KeyClient)
		if err != nil {
			panic(err)
		}
		c := dummyNewClient(redirect)
		c.UserId = u.Id
		err = cs.Create(store.NewConds(), c)
		if err != nil {
			panic(err)
		}

		return c, u
	}

	// create dummy Client and user
	ctx := getContext()
	defer store.CloseAllIn(ctx)
	storage := &oauth2.Storage{}
	storage.SetContext(ctx)

	c, u := createDummies(ctx, "password", "http://foobar.com/redirect")
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
