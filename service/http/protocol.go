package httpservicse

import (
	"net/http"

	"github.com/gourd/kit/store"
)

// Request contains all common fields needed for a usual API request
type Request struct {

	// Request stores the raw *http.Request
	Request *http.Request

	// Query stores the parsed Query information
	Query store.Query

	// Previous stores, if any, previous entity information (mainly for update)
	Previous interface{}

	// Payload stores, if any, current request payload information
	Payload interface{}
}
