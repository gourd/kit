package store_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/gourd/kit/store"
)

// mwareTestConn implements store.Conn
type mwareTestConn struct {
	raw    interface{}
	closed chan<- int
}

// Raw implements store.Conn
func (c mwareTestConn) Raw() interface{} {
	return c.raw
}

// Close implements store.Conn
func (c mwareTestConn) Close() {
	go func() {
		c.closed <- 0
	}()
}

func TestMiddleware(t *testing.T) {

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
	mware := func(msg interface{}, ch chan<- int) endpoint.Middleware {

		dummySrc := func() (conn store.Conn, err error) {
			conn = tConn{msg, ch}
			return
		}

		factory := store.NewFactory()
		factory.SetSource(srcKey, store.SourceFunc(dummySrc))
		factory.Set(key, srcKey, func(sess interface{}) (s store.Store, err error) {
			err = fmt.Errorf("%s", sess)
			return
		})

		return store.Middleware(factory)

	}(msg, ch)

	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {

		// get a store
		if _, err := store.Get(ctx, key); err == nil {
			t.Error("unexpected nil error")
		} else if want, have := msg, err.Error(); want != have {
			t.Errorf("expected %#v, got %#v", want, have)
		}

		return
	}

	ep = mware(ep)

	// test if store would close before timeout
	d, _ := time.ParseDuration("1s")
	timeout := time.After(d)
	ep(context.Background(), nil)

	select {
	case <-timeout:
		t.Error("mwareTestConn not closed before timeout")
	case <-ch:
		t.Log("mwareTestConn closed")
	}

}
