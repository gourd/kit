package store_test

import (
	"fmt"
	"testing"

	"github.com/gourd/kit/store"
)

func TestFactory_source(t *testing.T) {
	dummySrc1 := func() (conn store.Conn, err error) {
		err = fmt.Errorf("hello dummySrc")
		return
	}

	factory := store.NewFactory()
	factory.SetSource(store.DefaultSrc, dummySrc1)
	dummySrc2 := factory.GetSource(store.DefaultSrc)

	if _, err1 := dummySrc1(); err1 == nil {
		t.Errorf("unexpected nil value")
	} else if _, err2 := dummySrc2(); err2 == nil {
		t.Errorf("unexpected nil value")
	} else if want, have := err1.Error(), err2.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestFactory_store(t *testing.T) {
	dummyPrvdr1 := func(sess interface{}) (s store.Store, err error) {
		err = fmt.Errorf("hello dummyPrvdr")
		return
	}

	type tempKey int
	var srcKey, storeKey tempKey = 0, 1

	factory := store.NewFactory()
	factory.Set(storeKey, srcKey, dummyPrvdr1)
	_, dummyPrvdr2 := factory.Get(storeKey)

	if _, err1 := dummyPrvdr1(nil); err1 == nil {
		t.Errorf("unexpected nil value")
	} else if _, err2 := dummyPrvdr2(nil); err2 == nil {
		t.Errorf("unexpected nil value")
	} else if want, have := err1.Error(), err2.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
