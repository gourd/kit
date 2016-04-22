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

// Source is the upperio implementation of store.Source
type Source struct {
	adapter string
	connURL db.ConnectionURL
}

// Open implements store.Source
func (src *Source) Open() (s store.Conn, err error) {
	database, err := db.Open(src.adapter, src.connURL)
	if err != nil {
		return
	}

	s = &Conn{db: database}
	return
}

// NewSource create store.Source from
func NewSource(adapter string, connURL db.ConnectionURL) store.Source {
	return &Source{
		adapter: adapter,
		connURL: connURL,
	}
}
