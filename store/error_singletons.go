package store

import (
	"net/http"
)

// StatusFound singleton status in case entity found
var StatusFound = Error(http.StatusFound, "Success")

// ErrorNotFound singleton status in case not found
var ErrorNotFound = Error(http.StatusNotFound, "Not Found")

// ErrorForbidden singleton status in case permission denied
var ErrorForbidden = Error(http.StatusForbidden, "Permission Denied")

// ErrorInternal singleton status in case internal server error
var ErrorInternal = Error(http.StatusInternalServerError, "Internal Server Error")

// ErrorMethodNotAllowed singleton status in case HTTP method is not allowed to use
var ErrorMethodNotAllowed = Error(http.StatusMethodNotAllowed, "Method Not Allowed")

// ExpandError trys to cast the error into *StoreError
// or generate a new *ServicError based on it
func ExpandError(err error) *StoreError {

	if err == nil {
		return nil
	} else if serr, ok := err.(*StoreError); ok {
		return serr
	}

	return Error(http.StatusInternalServerError, err.Error())
}

// ParseError reads and parse a given status message.
// It reads Error type and unpack the status code and message.
// If it is not of Error type, it will return 400 internal server error.
// If it is nil, it will return 302 found
func ParseError(err error) (code int, msg string) {

	if err == nil {
		code = http.StatusFound
		msg = "Success"
		return
	}

	serr := ExpandError(err)
	code = serr.Code
	msg = serr.Error()
	return
}
