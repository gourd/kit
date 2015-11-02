package upperio

import (
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"upper.io/db"
)

// Upper is the general registry to be used
// in database session management for upperio database
var defs map[string]Def

const upperCtxKey = "gourd/kit/store/upperio/"

func init() {
	defs = make(map[string]Def)
}

// Def contains the definition of a database source
// Has all the parameters needed by db.Database.Open()
type Def struct {
	Adapter string
	URL     db.ConnectionURL
}

// Define a database source with name
func Define(name, adapter string, conn db.ConnectionURL) {
	defs[name] = Def{
		Adapter: adapter,
		URL:     conn,
	}
}

// Open a database from existing definitions and error if there is problem
// or retrieve the previously openned database session
func Open(r *http.Request, name string) (d db.Database, err error) {

	// try getting from context
	if cv, ok := context.GetOk(r, upperCtxKey+name); ok {
		if d, ok = cv.(db.Database); ok {
			return
		}
	}

	// find definition
	if def, ok := defs[name]; ok {
		// connect
		d, err = db.Open(def.Adapter, def.URL)
		if err != nil {
			// remember the database in context
			context.Set(r, upperCtxKey+name, d)
		}
		return
	}

	// tell user that the definition doesn't exists
	err = fmt.Errorf(
		"Definition for upper.io source \"%s\" not exists", name)
	return
}

// Close down an existing database connection
func Close(r *http.Request, name string) error {
	var d db.Database

	// try getting from context
	if cv, ok := context.GetOk(r, upperCtxKey+name); ok {
		if d, ok = cv.(db.Database); ok {
			// disconnect
			return d.Close()
		}
	}

	// if connection doesn't exists, quit scilently
	return nil
}

// MustOpen equals to open except it return only the database and not error.
// It will panic when encountering error
func MustOpen(r *http.Request, name string) (d db.Database) {
	d, err := Open(r, name)
	if err != nil {
		panic(err.Error())
	}
	return
}
