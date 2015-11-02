package store

// Result defines interface of an entity servie
// search result
type Result interface {

	// All fetches all results within the result set and dumps them into the
	// given pointer to slice of maps or structs
	All(el interface{}) (err error)

	// Raw returns the underlying database result variable
	Raw() (interface{}, error)

	// Close closes the result set
	Close() error
}
