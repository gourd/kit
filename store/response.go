package store

import (
	"encoding/json"
	"net/http"
)

// ExpandResponse converts a map[string]interface{} into proper
// JSON-marshalable response
func ExpandResponse(vmap map[string]interface{}) *Response {
	return &Response{
		vals: vmap,
	}
}

// NewResponse return a service response pointer
func NewResponse(key string, v interface{}) *Response {
	res := &Response{
		vals: make(map[string]interface{}),
	}
	res.Set(key, v)
	return res
}

// Response helps to wrap a
type Response struct {
	vals map[string]interface{}
}

// Set a value in the response
func (res *Response) Set(key string, v interface{}) {
	res.vals[key] = v
}

// Get a value in the repsonse. If there are no values
// associated with the
func (res *Response) Get(key string) interface{} {
	val, ok := res.vals[key]
	if !ok {
		return nil
	}
	return val
}

// MarshalJSON implements json Marshaler interface
func (res Response) MarshalJSON() ([]byte, error) {
	if status, ok := res.vals["status"]; !ok {
		res.vals["status"] = http.StatusOK
	} else {
		switch status.(type) {
		case int:
			// do nothing
		default:
			res.vals["status"] = http.StatusOK
		}
	}
	return json.Marshal(res.vals)
}
