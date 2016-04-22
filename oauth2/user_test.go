package oauth2_test

import (
	"encoding/json"

	"github.com/gourd/kit/oauth2"
	"github.com/gourd/kit/store"

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

func TestUnmarshalPassword(t *testing.T) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}
	pass := randSeq(20)

	jstr := "{\"username\":\"testing\",\"password\":\"" + pass + "\"}"
	u := &oauth2.User{}

	var unmr json.Unmarshaler = u
	_ = unmr // test if User implement json.Unmarshaler

	if err := json.Unmarshal([]byte(jstr), u); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	if want, have := u.Hash(pass), u.Password; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestSetGetPassword(t *testing.T) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}
	pass := randSeq(20)
	u := &oauth2.User{}
	u.SetPassword(pass)
	if !u.PasswordIs(pass) {
		t.Errorf("password is not hashed string of %#v, got %#v",
			pass, u.Password)
	}
}

func TestDBSavePassword(t *testing.T) {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}
	pass := randSeq(20)

	// test store context
	type tempKey int
	const (
		testDB tempKey = iota
	)

	db, err := defaultTestSrc().Open()
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	us, err := oauth2.UserStoreProvider(db.Raw())
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	u1 := &oauth2.User{Username: "TestDBSavePassword"}
	u1.Password = u1.Hash(pass)

	us.Create(nil, u1)

	users := []oauth2.User{}
	us.Search(store.NewQuery().AddCond("id", u1.ID)).All(&users)

	if want, have := 1, len(users); want != have {
		t.Errorf("expected users to have %d value, got %d instead (users=%#v)", want, have, users)
	}

	u2 := users[0]
	if want, have := u1.Password, u2.Password; !u2.PasswordIs(pass) {
		t.Errorf("password mismatch. expected %#v, got %#v", want, have)
	}
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

func TestMarshalJSON(t *testing.T) {
	u1, u2 := &oauth2.User{Name: "user 1"}, &oauth2.User{}
	u1.AddMeta("hello", "world 1")
	u1.AddMeta("hello", "world 2")

	b, err := json.Marshal(u1)
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}
	t.Logf("marshal result: %s", b)

	if err := json.Unmarshal(b, u2); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	t.Log("test unmarshal result")

	m := u2.Meta()
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
		t.Logf("result json: %#v", u2.MetaJSON)
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

func TestUnmarshalDB(t *testing.T) {
	u := &oauth2.User{Name: "user 1"}
	u.AddMeta("hello", "world 1")
	u.AddMeta("hello", "world 2")

	v, err := u.MarshalDB()
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	vmap, ok := v.(map[string]interface{})
	if !ok {
		t.Errorf("expected map[string]interface{}, got %#v", v)
		return
	}
	if _, ok := vmap["meta_json"]; !ok {
		t.Error("key me")
		return
	}

	metaJSON, ok := vmap["meta_json"].(string)
	if !ok {
		t.Errorf("expected string, got %#v", v)
		return
	}

	m := make(map[string][]string)
	if err := json.Unmarshal([]byte(metaJSON), &m); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

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
		t.Logf("result json: %#v", metaJSON)
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

func TestUnmarshalDB_EmptyMeta(t *testing.T) {
	u := &oauth2.User{Name: "user 1"}
	v, err := u.MarshalDB()
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	vmap, ok := v.(map[string]interface{})
	if !ok {
		t.Errorf("expected map[string]interface{}, got %#v", v)
		return
	}
	if _, ok := vmap["meta_json"]; !ok {
		t.Error("key me")
		return
	}

	metaJSON, ok := vmap["meta_json"].(string)
	if !ok {
		t.Errorf("expected string, got %#v", v)
		return
	}

	if want, have := "{}", metaJSON; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
