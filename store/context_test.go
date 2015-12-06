package store_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gourd/kit/store"
	"golang.org/x/net/context"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// tConn implements store.Conn
type tConn struct {
	raw    interface{}
	closed chan<- int
}

// Raw implements store.Conn
func (c tConn) Raw() interface{} {
	return c.raw
}

// Close implements store.Conn
func (c tConn) Close() {
	go func() {
		c.closed <- 0
	}()
}

func TestSetGet(t *testing.T) {

	type tempKey int

	const (
		srcKey tempKey = iota
		key
	)

	randString := func(n int) string {
		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return string(b)
	}

	msg := randString(20)
	ch := make(chan int)

	// prepare the context
	ctx := func(msg interface{}, ch chan<- int) context.Context {

		dummySrc := func() (conn store.Conn, err error) {
			conn = tConn{msg, ch}
			return
		}

		defs := store.NewDefs()
		defs.SetSource(srcKey, dummySrc)
		defs.Set(key, srcKey, func(sess interface{}) (s store.Store, err error) {
			err = fmt.Errorf("%s", sess)
			return
		})

		return store.WithStores(context.Background(), defs)

	}(msg, ch)

	// get a store
	if _, err := store.Get(ctx, key); err == nil {
		t.Error("unexpected nil error")
	} else if want, have := msg, err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// test if store would close before timeout
	d, _ := time.ParseDuration("1s")
	timeout := time.After(d)
	store.CloseAllIn(ctx)

	select {
	case <-timeout:
		t.Error("tConn not closed before timeout")
	case <-ch:
		t.Log("tConn closed")
	}

}
