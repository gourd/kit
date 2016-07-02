package store_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gourd/kit/store"
)

// testConn implements store.Conn
type testConn struct {
	serial int
}

func (conn *testConn) String() string   { return fmt.Sprintf("%d", conn.serial) }
func (conn *testConn) Raw() interface{} { return nil }
func (conn *testConn) Close()           {}

// testSource implements store.Source
type testSource struct {
	serialLast int
}

func (src *testSource) Open() (store.Conn, error) {
	src.serialLast++
	return &testConn{serial: src.serialLast}, nil
}

// test store.PoolConn implements store.Conn
func TestPoolConn_storeConn(t *testing.T) {
	var conn store.Conn = &store.PoolConn{}
	_ = conn
}

// test store.SourcePool implements store.Source
func TestSourcePool_storeSource(t *testing.T) {
	var src store.Source = &store.SourcePool{}
	_ = src
}

// test normal blocking
func TestSourcePool_Open(t *testing.T) {
	src := &testSource{}
	pool := store.Pool(src, 10, 300*time.Second)
	grace := 10 * time.Millisecond

	getConn := func() <-chan store.Conn {
		out := make(chan store.Conn)
		go func() {
			conn, _ := pool.Open()
			out <- conn
			close(out)
		}()
		return out
	}

	// test getting connection within pool size
	timeout := false
OuterLoop:
	// should open 9 times without blocking
	for i := 0; i < 9; i++ {
		select {
		case <-getConn():
			break
		case <-time.After(grace):
			timeout = true
			t.Logf("time out on loop %d", i+1)
			break OuterLoop
		}
	}
	if want, have := false, timeout; want != have {
		t.Errorf("loop timeout before running 9 times")
	}

	// test getting the last connection in the pool
	var lastConn store.Conn
	select {
	case lastConn = <-getConn():
		break
	case <-time.After(grace):
		timeout = true
	}
	if want, have := false, timeout; want != have {
		t.Errorf("loop timeout before getting the last connection")
	}

	if lastConn == nil {
		t.Errorf("lastConn is nil")
		return
	}

	var currentConn *store.PoolConn
	setPoolConnPtr := func(currentConn **store.PoolConn, in <-chan store.Conn) {
		conn := <-in
		*currentConn, _ = conn.(*store.PoolConn)
	}

	// test getting connection after pool is dry off
	go setPoolConnPtr(&currentConn, getConn())
	time.Sleep(grace + 10*time.Millisecond) // wait long enough
	if currentConn != nil {
		t.Errorf("blocking failed. currentConn is not nil: %#v", currentConn)
	}

	// close a connection and see if the pool regains
	// the connection
	lastConn.Close()
	time.Sleep(30 * time.Millisecond) // wait for goroutine to do the work
	if currentConn == nil {
		t.Errorf("failed to connect after old connection closed")
	}
}

// test if connections can notify pool for their expiration
// which re-generates new connection to the pool
func TestSourcePool_connsReturning(t *testing.T) {
	src := &testSource{}
	expires := time.Millisecond
	pool := store.Pool(src, 1, expires)
	pool.Open()

	getConn := func() <-chan store.Conn {
		out := make(chan store.Conn)
		go func() {
			conn, _ := pool.Open()
			out <- conn
			close(out)
		}()
		return out
	}

	select {
	case <-getConn():
		// skip
	case <-time.After((expires * 15) / 10):
		// wait some more to see if get new connection
		t.Errorf("failed to get connection after conn expires")
	}
}
