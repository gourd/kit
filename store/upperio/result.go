package upperio

import (
	"github.com/gourd/kit/store"
	"upper.io/db"
)

func NewResult(fn func() (db.Result, error)) store.Result {
	return &Result{fn}
}

// Result implements store.Result
type Result struct {
	resultFunc func() (db.Result, error)
}

// All fetches all results within the result set and dumps them into the
// given pointer to slice of maps or structs
func (res *Result) All(el interface{}) (err error) {
	raw, err := res.raw()
	if err != nil {
		return
	}

	err = raw.All(el)
	if err != nil {
		serr := store.ErrorInternal
		serr.ServerMsg = err.Error()
		err = serr
	}
	return
}

// raw returns raw db.Result, or error
func (res *Result) raw() (raw db.Result, err error) {
	raw, err = res.resultFunc()
	if err != nil {
		serr := store.ErrorInternal
		serr.ServerMsg = err.Error()
		err = serr
		return
	}
	return
}

// Raw returns the underlying database result variable
func (res *Result) Raw() (interface{}, error) {
	return res.raw()
}

// Close closes the result set
func (res *Result) Close() error {
	return res.Close()
}
