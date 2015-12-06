package store_test

import (
	"fmt"
	"testing"

	"github.com/gourd/kit/store"
)

func TestDefs_source(t *testing.T) {
	dummySrc1 := func() (conn store.Conn, err error) {
		err = fmt.Errorf("hello dummySrc")
		return
	}

	defs := store.NewDefs()
	defs.SetSource(store.DefaultSrc, dummySrc1)
	dummySrc2 := defs.GetSource(store.DefaultSrc)

	if _, err1 := dummySrc1(); err1 == nil {
		t.Errorf("unexpected nil value")
	} else if _, err2 := dummySrc2(); err2 == nil {
		t.Errorf("unexpected nil value")
	} else if want, have := err1.Error(), err2.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestDefs_store(t *testing.T) {
	dummyPrvdr1 := func(sess interface{}) (s store.Store, err error) {
		err = fmt.Errorf("hello dummyPrvdr")
		return
	}

	type tempKey int
	var srcKey, storeKey tempKey = 0, 1

	defs := store.NewDefs()
	defs.Set(storeKey, srcKey, dummyPrvdr1)
	_, dummyPrvdr2 := defs.Get(storeKey)

	if _, err1 := dummyPrvdr1(nil); err1 == nil {
		t.Errorf("unexpected nil value")
	} else if _, err2 := dummyPrvdr2(nil); err2 == nil {
		t.Errorf("unexpected nil value")
	} else if want, have := err1.Error(), err2.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
