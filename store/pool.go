package store

import "time"

// PoolConn is a wrapper of Conn
// that also implments Conn
type PoolConn struct {
	pool    *SourcePool
	expires time.Time
	err     error
	dbConn  Conn
}

// Err returns the connection error of the wrapped source
func (conn *PoolConn) Err() error {
	return conn.err
}

// Expired return wether the connection has expired or not
func (conn *PoolConn) Expired() bool {
	return time.Now().After(conn.expires)
}

// Raw implements store.Conn.Raw()
func (conn *PoolConn) Raw() interface{} {
	return conn.dbConn.Raw()
}

// Close implements store.Conn.Close()
func (conn *PoolConn) Close() {
	go func() {
		// return itself to the pool if it is not expired
		if !conn.Expired() {
			pool := conn.pool
			conn.pool = nil // remove reference to original conn.pool
			pool.conns <- conn
		}
	}()
	conn.dbConn.Close()
}

// SourcePool helps to pool connection of any given source
type SourcePool struct {
	conns   chan *PoolConn
	src     Source
	expires time.Duration
}

// Open implements Source.Open()
func (pool *SourcePool) Open() (conn Conn, err error) {
	var poolConn *PoolConn
	for poolConn = <-pool.conns; poolConn.Expired(); poolConn = <-pool.conns {
		// get connection that is not expired
	}
	err = poolConn.err
	if err != nil {
		return
	}
	poolConn.pool = pool
	conn = poolConn
	return
}

// AddNewConn adds a new connection to the pool
func (pool *SourcePool) AddNewConn() {
	go func() {
		time.Sleep(pool.expires)
		pool.AddNewConn()
	}()

	dbConn, err := pool.src.Open()
	expires := time.Now().Add(pool.expires)
	if err != nil {
		pool.conns <- &PoolConn{
			err: err,
		}
	} else {
		pool.conns <- &PoolConn{
			dbConn:  dbConn,
			expires: expires,
		}
	}
}

// Pool wraps a source into a SourcePool
func Pool(src Source, size uint, expires time.Duration) *SourcePool {
	pool := &SourcePool{
		conns:   make(chan *PoolConn, size),
		src:     src,
		expires: expires,
	}

	// pump in new connections
	for i := uint(0); i < size; i++ {
		pool.AddNewConn()
	}

	return pool
}
