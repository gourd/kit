package upperio

import (
	"github.com/gourd/kit/store"
	"upper.io/db"
)

// Conn implements store.Conn
type Conn struct {
	db db.Database
}

// Raw implements store.Conn.Raw()
func (conn *Conn) Raw() interface{} {
	return conn.db
}

// Close implements store.Conn.Close()
func (conn *Conn) Close() {
	conn.db.Close()
}

// Source create store.Source from
func Source(adapter string, connURL db.ConnectionURL) store.Source {
	return func() (s store.Conn, err error) {
		database, err := db.Open(adapter, connURL)
		if err != nil {
			return
		}

		s = &Conn{db: database}
		return
	}
}
