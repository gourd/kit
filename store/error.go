package store

import (
	"fmt"
)

// Error creates formated error message with entityError type
func Error(code int, msg string, v ...interface{}) *StoreError {
	var status int
	if code < 10000 {
		status = code
	} else {
		status = code / 100
	}

	renderedMsg := fmt.Sprintf(msg, v...)
	return &StoreError{
		status,
		code,
		renderedMsg,
		renderedMsg,
		"",
	}
}

// StoreError is for common service error message
type StoreError struct {

	// Status is the http status code
	Status int `json:"status"`

	// Code is service specific status code
	Code int `json:"code"`

	// ServerMsg is the server side log message
	ServerMsg string `json:"-"`

	// ClientMsg is the client side message which could
	// be displayed to their user directly
	ClientMsg string `json:"message"`

	// DeveloperMsg is the client side message which
	// should be of help to developer. Omit if empty
	DeveloperMsg string `json:"developer_message,omitempty"`
}

// Error implements the standard error type
// returns the client message
func (err *StoreError) Error() string {
	return err.ClientMsg
}

// String implements Stringer type which
// returns the client message
func (err *StoreError) String() string {
	return err.ClientMsg
}
